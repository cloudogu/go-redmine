package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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
	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(id) + ".json?" + c.apiKeyParameter())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectResult
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project (id: %d) was not found", id)
	}
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
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
	return &r.Project, nil
}

func isHTTPStatusSuccessful(httpStatus int, acceptedStatuses []int) bool {
	for _, acceptedStatus := range acceptedStatuses {
		if httpStatus == acceptedStatus {
			return true
		}
	}

	return false
}

func (c *Client) Projects() ([]Project, error) {
	parameters := c.concatParameters(c.apiKeyParameter(), c.getPaginationClause())
	res, err := c.Get(c.endpoint + "/projects.json?" + parameters)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectsResult
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
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
	return r.Projects, nil
}

func (c *Client) CreateProject(project Project) (*Project, error) {
	var ir projectRequest
	ir.Project = project
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}

	parameters := c.concatParameters(c.apiKeyParameter())
	req, err := http.NewRequest("POST", c.endpoint+"/projects.json?"+parameters, strings.NewReader(string(s)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r projectRequest
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
	return &r.Project, nil
}

func (c *Client) UpdateProject(project Project) error {
	var ir projectRequest
	ir.Project = project
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}

	parameters := c.concatParameters(c.apiKeyParameter())
	req, err := http.NewRequest("PUT", c.endpoint+"/projects/"+strconv.Itoa(project.Id)+".json?"+parameters, strings.NewReader(string(s)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not update project (id: %d) because it was not found", project.Id)
	}
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		decoder := json.NewDecoder(res.Body)
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	if err != nil {
		return err
	}
	return err
}

func (c *Client) DeleteProject(id int) error {
	parameters := c.concatParameters(c.apiKeyParameter())
	req, err := http.NewRequest("DELETE", c.endpoint+"/projects/"+strconv.Itoa(id)+".json?"+parameters, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("could not delete project (id %d) because it was not found", id)
	}

	decoder := json.NewDecoder(res.Body)
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK, http.StatusNoContent}) {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}
