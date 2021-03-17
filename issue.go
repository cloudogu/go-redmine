package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	errors2 "github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

type issueRequest struct {
	Issue Issue `json:"issue"`
}

type issueResult struct {
	Issue Issue `json:"issue"`
}

type issuesResult struct {
	Issues     []Issue `json:"issues"`
	TotalCount uint    `json:"total_count"`
	Offset     uint    `json:"offset"`
	Limit      uint    `json:"limit"`
}

type JournalDetails struct {
	Property string `json:"property"`
	Name     string `json:"name"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}
type Journal struct {
	Id        int              `json:"id"`
	User      *IdName          `json:"user"`
	Notes     string           `json:"notes"`
	CreatedOn string           `json:"created_on"`
	Details   []JournalDetails `json:"details"`
}

type Issue struct {
	Id           int            `json:"id"`
	Subject      string         `json:"subject"`
	Description  string         `json:"description"`
	ProjectId    int            `json:"project_id"`
	Project      *IdName        `json:"project"`
	TrackerId    int            `json:"tracker_id"`
	Tracker      *IdName        `json:"tracker"`
	ParentId     int            `json:"parent_issue_id,omitempty"`
	Parent       *Id            `json:"parent"`
	StatusId     int            `json:"status_id"`
	Status       *IdName        `json:"status"`
	PriorityId   int            `json:"priority_id,omitempty"`
	Priority     *IdName        `json:"priority"`
	Author       *IdName        `json:"author"`
	FixedVersion *IdName        `json:"fixed_version"`
	AssignedTo   *IdName        `json:"assigned_to"`
	Category     *IdName        `json:"category"`
	CategoryId   int            `json:"category_id"`
	Notes        string         `json:"notes"`
	StatusDate   string         `json:"status_date"`
	CreatedOn    string         `json:"created_on"`
	UpdatedOn    string         `json:"updated_on"`
	StartDate    string         `json:"start_date"`
	DueDate      string         `json:"due_date"`
	ClosedOn     string         `json:"closed_on"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
	Uploads      []*Upload      `json:"uploads"`
	DoneRatio    float32        `json:"done_ratio"`
	Journals     []*Journal     `json:"journals"`
}

type IssueFilter struct {
	ProjectId    string
	SubprojectId string
	TrackerId    string
	StatusId     string
	AssignedToId string
	UpdatedOn    string
	ExtraFilters map[string]string
}

type CustomField struct {
	Id          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Multiple    bool        `json:"multiple"`
	Value       interface{} `json:"value"`
}

func (c *Client) IssuesOf(projectId int) ([]Issue, error) {
	url := jsonResourceEndpoint(c.endpoint, "issues")
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for issues")
	}
	err = safelySetQueryParameter(req, "project_id", strconv.Itoa(projectId))
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding additional parameters to issue request")
	}
	err = safelySetQueryParameters(req, c.getPaginationClauseParams())
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding additional parameters to issue request")
	}

	issues, err := getPagedIssuesForRequest(c, req)
	if err != nil {
		return nil, errors2.Wrapf(err, "error while reading issues for project %d", projectId)
	}

	return issues, nil
}

func (c *Client) Issue(id int) (*Issue, error) {
	return getOneIssue(c, id, nil)
}

func (c *Client) IssueWithArgs(id int, args map[string]string) (*Issue, error) {
	return getOneIssue(c, id, args)
}

func (c *Client) IssuesByQuery(queryId int) ([]Issue, error) {
	url := jsonResourceEndpoint(c.endpoint, "issues")
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for issues")
	}
	err = safelySetQueryParameters(req, c.getPaginationClauseParams())
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding additional parameters to issue request")
	}
	err = safelySetQueryParameter(req, "query_id", strconv.Itoa(queryId))
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding query_id parameter to issue request")
	}

	issues, err := getPagedIssuesForRequest(c, req)
	if err != nil {
		return nil, errors2.Wrapf(err, "error while reading issues for query id %d", queryId)
	}

	return issues, nil
}

// IssuesByFilter filters issues applying the f criteria
func (c *Client) IssuesByFilter(f *IssueFilter) ([]Issue, error) {
	url := jsonResourceEndpoint(c.endpoint, "issues")
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for issues")
	}
	err = safelySetQueryParameters(req, c.getPaginationClauseParams())
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding additional parameters to issue request")
	}
	filterClauses := strings.Split(getIssueFilterClause(f), "&")
	for _, clause := range filterClauses {
		kv := strings.Split(clause, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("could not properly split issue filter %s", clause)
		}
		err = safelySetQueryParameter(req, kv[0], kv[1])
		if err != nil {
			return nil, errors2.Wrap(err, "error while adding query_id parameter to issue request")
		}
	}

	issues, err := getPagedIssuesForRequest(c, req)
	if err != nil {
		return nil, errors2.Wrapf(err, "error while reading issues by filter %v", f)
	}

	return issues, nil
}

func (c *Client) Issues() ([]Issue, error) {
	url := jsonResourceEndpoint(c.endpoint, "issues")
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for issues")
	}
	err = safelySetQueryParameters(req, c.getPaginationClauseParams())
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding additional parameters to issue request")
	}

	issues, err := getPagedIssuesForRequest(c, req)
	if err != nil {
		return nil, errors2.Wrap(err, "error while reading issues")
	}

	return issues, nil
}

func (c *Client) CreateIssue(issue Issue) (*Issue, error) {
	url := jsonResourceEndpoint(c.endpoint, "issues")

	var ir issueRequest
	ir.Issue = issue
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}
	req, err := c.authenticatedPost(url, strings.NewReader(string(s)))
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating PUT request for issue")
	}
	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r issueRequest
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusCreated}) {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return &r.Issue, nil
}

func (c *Client) UpdateIssue(issue Issue) error {
	url := jsonResourceEndpointByID(c.endpoint, "issues", issue.Id)

	var ir issueRequest
	ir.Issue = issue
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := c.authenticatedPut(url, strings.NewReader(string(s)))
	if err != nil {
		return errors2.Wrap(err, "error while creating PUT request for issue")
	}
	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not update issue (id: %d) because it was not found", issue.Id)
	}
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		return errors2.Wrapf(decodeHTTPError(res), "error while deleting issue %d", issue.Id)
	}

	return nil
}

func (c *Client) DeleteIssue(id int) error {
	url := jsonResourceEndpointByID(c.endpoint, "issues", id)
	req, err := c.authenticatedDelete(url, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not delete issue (id: %d) because it was not found", id)
	}

	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		return errors2.Wrapf(decodeHTTPError(res), "error while deleting issue %d", id)
	}

	return nil
}

func (issue *Issue) GetTitle() string {
	return fmt.Sprintf("%s #%d: %s", issue.Tracker.Name, issue.Id, issue.Subject)
}

// MarshalJSON marshals issue to JSON.
// This overrides the default MarshalJSON() to reset parent issue.
func (issue Issue) MarshalJSON() ([]byte, error) {
	type Issue2 Issue

	// To reset parent issue, set empty string to "parent_issue_id"
	var parentIssueID *string
	if issue.Parent == nil {
		// reset parent issue
		id := ""
		parentIssueID = &id
	} else if issue.ParentId > 0 {
		// set parent issue
		id := strconv.Itoa(issue.ParentId)
		parentIssueID = &id
	}

	return json.Marshal(&struct {
		Issue2
		ParentId *string `json:"parent_issue_id,omitempty"`
	}{
		Issue2:   Issue2(issue),
		ParentId: parentIssueID,
	})
}

func getIssueFilterClause(filter *IssueFilter) string {
	if filter == nil {
		return ""
	}
	clause := ""
	if filter.ProjectId != "" {
		clause = clause + fmt.Sprintf("&project_id=%v", filter.ProjectId)
	}
	if filter.SubprojectId != "" {
		clause = clause + fmt.Sprintf("&subproject_id=%v", filter.SubprojectId)
	}
	if filter.TrackerId != "" {
		clause = clause + fmt.Sprintf("&tracker_id=%v", filter.TrackerId)
	}
	if filter.StatusId != "" {
		clause = clause + fmt.Sprintf("&status_id=%v", filter.StatusId)
	}
	if filter.AssignedToId != "" {
		clause = clause + fmt.Sprintf("&assigned_to_id=%v", filter.AssignedToId)
	}
	if filter.UpdatedOn != "" {
		clause = clause + fmt.Sprintf("&updated_on=%v", filter.UpdatedOn)
	}

	if filter.ExtraFilters != nil {
		extraFilter := make([]string, 0)
		for key, value := range filter.ExtraFilters {
			extraFilter = append(extraFilter, fmt.Sprintf("%s=%s", key, value))
		}
		clause = clause + "&" + strings.Join(extraFilter[:], "&")
	}

	return clause
}

func getOneIssue(c *Client, id int, args map[string]string) (*Issue, error) {
	kvs := argsToKeyValues(args)

	url := jsonResourceEndpointByID(c.endpoint, "issues", id)
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for issue")
	}
	err = safelySetQueryParameters(req, kvs)
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding additional parameters to issue request")
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrapf(err, "could not read issue %d ", id)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("issue (id: %d) was not found", id)
	}

	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		return nil, errors2.Wrapf(decodeHTTPError(res), "error while reading issue %d", id)
	}

	var r issueResult
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.Issue, nil
}

func argsToKeyValues(args map[string]string) []keyValue {
	if args == nil {
		return nil
	}

	kvs := []keyValue{}
	for key, value := range args {
		kv := keyValue{key: key, value: value}
		kvs = append(kvs, kv)
	}
	return kvs
}

func getPagedIssuesForRequest(c *Client, req *http.Request) ([]Issue, error) {
	completed := false
	var issues []Issue

	for !completed {
		r, err := getOffsetIssueForRequest(c, req, len(issues))

		if err != nil {
			return nil, err
		}

		if r.TotalCount == uint(len(issues)) {
			completed = true
		}

		issues = append(issues, r.Issues...)
	}

	return issues, nil
}

func getOffsetIssueForRequest(c *Client, req *http.Request, offset int) (*issuesResult, error) {
	err := safelySetQueryParameter(req, "offset", strconv.Itoa(offset))
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding additional parameters to issue request")
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrap(err, "error while reading issue response")
	}
	defer res.Body.Close()

	var r issuesResult
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		return nil, errors2.Wrapf(decodeHTTPError(res), "issue request returned non-successfully, URL: %s", req.URL.String())
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
