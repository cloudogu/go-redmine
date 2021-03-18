package redmine

import (
	"encoding/json"
	"fmt"
	errors2 "github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

type issueCategoriesResult struct {
	IssueCategories []IssueCategory `json:"issue_categories"`
	TotalCount      int             `json:"total_count"`
}

type issueCategoryResult struct {
	IssueCategory IssueCategory `json:"issue_category"`
}

type issueCategoryRequest struct {
	IssueCategory IssueCategory `json:"issue_category"`
}

type IssueCategory struct {
	Id         int    `json:"id"`
	Project    IdName `json:"project"`
	Name       string `json:"name"`
	AssignedTo IdName `json:"assigned_to"`
}

func (c *Client) IssueCategories(projectId int) ([]IssueCategory, error) {
	url := jsonResourceEndpoint(c.endpoint, "projects/"+strconv.Itoa(projectId)+"/issue_categories")
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for issue_categories")
	}
	err = safelySetQueryParameters(req, c.getPaginationClauseParams())
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding pagination parameters to issue_categories")
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrap(err, "could not read issue_categories")
	}
	defer res.Body.Close()

	var r issueCategoriesResult
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		return nil, errors2.Wrap(decodeHTTPError(res), "error while reading issue_categories")
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r.IssueCategories, nil
}

func (c *Client) IssueCategory(id int) (*IssueCategory, error) {
	url := jsonResourceEndpointByID(c.endpoint, "issue_categories", id)
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrapf(err, "error while creating GET request for issue category %d ", id)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrapf(err, "could not read issue category %d ", id)
	}
	defer res.Body.Close()

	var r issueCategoryResult
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("issue category (id: %d) was not found", id)
	}

	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		return nil, errors2.Wrapf(decodeHTTPError(res), "error while reading issue category %d", id)
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.IssueCategory, nil
}

func (c *Client) CreateIssueCategory(issueCategory IssueCategory) (*IssueCategory, error) {
	var ir issueCategoryRequest
	ir.IssueCategory = issueCategory
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}

	url := jsonResourceEndpoint(c.endpoint, "issue_categories")
	req, err := c.authenticatedPost(url, strings.NewReader(string(s)))
	if err != nil {
		return nil, errors2.Wrapf(err, "error while creating POST request for issue category %s ", issueCategory.Name)
	}
	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrapf(err, "could not create issue category %s ", issueCategory.Name)
	}
	defer res.Body.Close()

	var r issueCategoryResult
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusCreated}) {
		return nil, errors2.Wrapf(decodeHTTPError(res), "error while creating issue category %s", issueCategory.Name)
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.IssueCategory, nil
}

func (c *Client) UpdateIssueCategory(issueCategory IssueCategory) error {
	var ir issueCategoryRequest
	ir.IssueCategory = issueCategory
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}

	url := jsonResourceEndpointByID(c.endpoint, "issue_categories", issueCategory.Id)
	req, err := c.authenticatedPut(url, strings.NewReader(string(s)))
	if err != nil {
		return errors2.Wrapf(err, "error while creating PUT request for issue category %d ", issueCategory.Id)
	}
	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return errors2.Wrapf(err, "could not update project %d ", issueCategory.Id)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not update issue category (id: %d) because it was not found", issueCategory.Id)
	}
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		return errors2.Wrapf(decodeHTTPError(res), "error while updating issue category %d", issueCategory.Id)
	}

	return nil
}

func (c *Client) DeleteIssueCategory(id int) error {
	url := jsonResourceEndpointByID(c.endpoint, "issue_categories", id)
	req, err := c.authenticatedDelete(url, strings.NewReader(""))
	if err != nil {
		return errors2.Wrapf(err, "error while creating DELETE request for issue category %d ", id)
	}

	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return errors2.Wrapf(err, "could not delete issue category %d ", id)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not delete issue category (id: %d) because it was not found", id)
	}

	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		return errors2.Wrapf(decodeHTTPError(res), "error while deleting issue category %d", id)
	}

	return nil
}
