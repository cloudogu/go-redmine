package redmine

import (
	"encoding/json"
	"fmt"
	errors2 "github.com/pkg/errors"
	"net/http"
	"strings"
)

const entityEndpointNameProjects = "projects"

type projectRequest struct {
	Project Project `json:"project"`
}

type projectResult struct {
	Project Project `json:"project"`
}

type projectsResult struct {
	Projects []Project `json:"projects"`
}

// Project contains a Redmine API project object according Redmine 4.1 REST API.
//
// See also: https://www.redmine.org/projects/redmine/wiki/Rest_api
type Project struct {
	// Id uniquely identifies a project on technical level. This value will be generated on project creation and cannot
	// be changed. Id is mandatory for all project API calls except CreateProject()
	Id int `json:"id"`
	// ParentID may contain the Id of a parent project. If set, this project is then a child project of the parent project.
	// Projects can be unlimitedly nested.
	ParentID Id `json:"parent_id"`
	// Name contains a human readable project name.
	Name string `json:"name"`
	// Identifier used by the application for various things (eg. in URLs). It must be unique and cannot be composed of
	// only numbers. It must contain 1 to 100 characters of which only consist of lowercase latin characters, numbers,
	// hyphen (-) and underscore (_). Once the project is created, this identifier cannot be modified
	Identifier string `json:"identifier"`
	// Description contains a human readable project multiline description that appears on the project overview.
	Description string `json:"description"`
	// Homepage contains a URL to a project's website that appears on the project overview.
	Homepage string `json:"homepage"`
	// IsPublic controls who can view the project. If set to true the project can be viewed by all the users, including
	// those who are not members of the project. If set to false, only the project members have access to it, according to
	// their role.
	//
	// since Redmine 2.6.0
	IsPublic bool `json:"is_public"`
	// InheritMembers determines whether this project inherits members from a parent project. If set to true (and being a
	// nested project) all members from the parent project will apply also to this project.
	InheritMembers bool `json:"inherit_members"`
	// CreatedOn contains a timestamp of when the project was created.
	CreatedOn string `json:"created_on"`
	// UpdatedOn contains the timestamp of when the project was last updated.
	UpdatedOn string `json:"updated_on"`

	Status int `json:"status,omitempty"`
}

// Project returns a single project without additional fields.
func (c *Client) Project(id int) (*Project, error) {
	url := jsonResourceEndpointByID(c.endpoint, entityEndpointNameProjects, id)
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrapf(err, "error while creating GET request for project %d ", id)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrapf(err, "could not read project %d ", id)
	}
	defer res.Body.Close()

	var r projectResult
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project (id: %d) was not found", id)
	}

	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		return nil, errors2.Wrapf(decodeHTTPError(res), "error while reading project %d", id)
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.Project, nil
}

func (c *Client) Projects() ([]Project, error) {
	url := jsonResourceEndpoint(c.endpoint, entityEndpointNameProjects)
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for projects")
	}
	err = safelySetQueryParameters(req, c.getPaginationClauseParams())
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding pagination parameters to project request")
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrap(err, "could not read projects")
	}
	defer res.Body.Close()

	var r projectsResult
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		return nil, errors2.Wrap(decodeHTTPError(res), "error while reading projects")
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r.Projects, nil
}

func (c *Client) CreateProject(project Project) (*Project, error) {
	var ir projectRequest
	ir.Project = project
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}

	url := jsonResourceEndpoint(c.endpoint, entityEndpointNameProjects)
	req, err := c.authenticatedPost(url, strings.NewReader(string(s)))
	if err != nil {
		return nil, errors2.Wrapf(err, "error while creating POST request for project %s ", project.Identifier)
	}
	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrapf(err, "could not create project %s ", project.Identifier)
	}
	defer res.Body.Close()

	var r projectRequest
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusCreated}) {
		return nil, errors2.Wrapf(decodeHTTPError(res), "error while creating project %s", project.Identifier)
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.Project, nil
}

func (c *Client) UpdateProject(project Project) error {
	var ir projectRequest
	ir.Project = project
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}

	url := jsonResourceEndpointByID(c.endpoint, entityEndpointNameProjects, project.Id)
	req, err := c.authenticatedPut(url, strings.NewReader(string(s)))
	if err != nil {
		return errors2.Wrapf(err, "error while creating PUT request for project %d ", project.Id)
	}
	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return errors2.Wrapf(err, "could not update project %d ", project.Id)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not update project (id: %d) because it was not found", project.Id)
	}
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		return errors2.Wrapf(decodeHTTPError(res), "error while updating project %d", project.Id)
	}

	return nil
}

func (c *Client) DeleteProject(id int) error {
	url := jsonResourceEndpointByID(c.endpoint, entityEndpointNameProjects, id)
	req, err := c.authenticatedDelete(url, strings.NewReader(""))
	if err != nil {
		return errors2.Wrapf(err, "error while creating DELETE request for project %d ", id)
	}

	req.Header.Set(httpHeaderContentType, httpContentTypeApplicationJson)
	res, err := c.Do(req)
	if err != nil {
		return errors2.Wrapf(err, "could not delete project %d ", id)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not delete project (id: %d) because it was not found", id)
	}

	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		return errors2.Wrapf(decodeHTTPError(res), "error while deleting project %d", id)
	}

	return nil
}
