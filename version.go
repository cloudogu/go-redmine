package redmine

import (
	"encoding/json"
	"fmt"
	errors2 "github.com/pkg/errors"
	"net/http"
	"strings"
)

const entityEndpointNameVersions = "versions"

type versionRequest struct {
	Version Version `json:"version"`
}

type versionResult struct {
	Version Version `json:"version"`
}

type versionsResult struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	Id           int            `json:"id"`
	Project      IdName         `json:"project"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Status       string         `json:"status"`
	DueDate      string         `json:"due_date"`
	CreatedOn    string         `json:"created_on"`
	UpdatedOn    string         `json:"updated_on"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
}

func (c *Client) Version(id int) (*Version, error) {
	url := jsonResourceEndpointByID(c.endpoint, entityEndpointNameVersions, id)
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrapf(err, "error while creating GET request for version %d ", id)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrapf(err, "could not read version %d ", id)
	}
	defer res.Body.Close()

	var r versionResult
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("version (id: %d) was not found", id)
	}

	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		return nil, errors2.Wrapf(decodeHTTPError(res), "error while reading version %d", id)
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.Version, nil
}

func (c *Client) Versions(projectId int) ([]Version, error) {
	compoundEndpointName := fmt.Sprintf("%s/%d/%s", entityEndpointNameProjects, projectId, entityEndpointNameVersions)
	url := jsonResourceEndpoint(c.endpoint, compoundEndpointName)
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for versions")
	}
	err = safelySetQueryParameters(req, c.getPaginationClauseParams())
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding pagination parameters to versions")
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrap(err, "could not read issue_categories")
	}
	defer res.Body.Close()

	var r versionsResult
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		return nil, errors2.Wrap(decodeHTTPError(res), "error while reading versions")
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r.Versions, nil
}

func (c *Client) CreateVersion(version Version) (*Version, error) {
	var ir versionRequest
	ir.Version = version
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}

	compoundEndpointName := fmt.Sprintf("%s/%d/%s", entityEndpointNameProjects, version.Project.Id, entityEndpointNameVersions)
	url := jsonResourceEndpoint(c.endpoint, compoundEndpointName)
	req, err := c.authenticatedPost(url, strings.NewReader(string(s)))
	if err != nil {
		return nil, errors2.Wrapf(err, "error while creating POST request for version %s ", version.Name)
	}
	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrapf(err, "could not create version %s ", version.Name)
	}
	defer res.Body.Close()

	var r versionRequest
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusCreated}) {
		return nil, errors2.Wrapf(decodeHTTPError(res), "error while creating version %s", version.Name)
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.Version, nil
}

func (c *Client) UpdateVersion(version Version) error {
	var ir versionRequest
	ir.Version = version
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}

	url := jsonResourceEndpointByID(c.endpoint, entityEndpointNameVersions, version.Id)
	req, err := c.authenticatedPut(url, strings.NewReader(string(s)))
	if err != nil {
		return errors2.Wrapf(err, "error while creating PUT request for version %d ", version.Id)
	}
	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return errors2.Wrapf(err, "could not update version %d ", version.Id)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not update version (id: %d) because it was not found", version.Id)
	}
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		return errors2.Wrapf(decodeHTTPError(res), "error while updating version %d", version.Id)
	}

	return nil
}

func (c *Client) DeleteVersion(id int) error {
	url := jsonResourceEndpointByID(c.endpoint, entityEndpointNameVersions, id)
	req, err := c.authenticatedDelete(url, strings.NewReader(""))
	if err != nil {
		return errors2.Wrapf(err, "error while creating DELETE request for version %d ", id)
	}

	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return errors2.Wrapf(err, "could not delete version %d ", id)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not delete version (id: %d) because it was not found", id)
	}

	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		return errors2.Wrapf(decodeHTTPError(res), "error while deleting version %d", id)
	}

	return nil
}
