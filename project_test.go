package redmine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testProjectJSON = `{
    "id": 1,
    "name": "example project",
    "identifier": "exampleproject",
    "description": "This is an example project.",
    "homepage": "http://github.com/cloudogu/go-redmine",
    "status": 1,
    "is_public": true,
    "inherit_members": true,
    "trackers": [
      {
        "id": 1,
        "name": "Bug"
      },
      {
        "id": 2,
        "name": "Feature"
      }
    ],
    "enabled_modules": [
      {
        "id": 71,
        "name": "issue_tracking"
      },
			{
        "id": 73,
        "name": "wiki"
      }
    ],
    "created_on": "2021-02-19T16:51:03Z",
    "updated_on": "2021-02-19T16:51:25Z"
  }`
const simpleProjectJSON = `{ "project":` + testProjectJSON + "}"
const simpleProjectsJSON = `{ "projects":[` + testProjectJSON + `],"total_count":1,"offset":0,"limit":25}`

var testProject = Project{
	Id:             1,
	ParentID:       Id{},
	Name:           "example project",
	Identifier:     "exampleproject",
	Description:    "This is an example project.",
	Homepage:       "http://github.com/cloudogu/go-redmine",
	IsPublic:       true,
	InheritMembers: true,
	CreatedOn:      "2021-02-19T16:51:03Z",
	UpdatedOn:      "2021-02-19T16:51:25Z",
}

const (
	authUser     = "leUser"
	authPassword = "Passwort1! äöü+ß"
	authToken    = "123456789"
)

func TestClient_Project(t *testing.T) {
	t.Run("should parse general project fields, and ignore module names and trackers from http response", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintln(w, simpleProjectJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		actualProject, err := sut.Project(1)

		require.NoError(t, err)
		require.NotEmpty(t, actualProject)
		expectedProject := &Project{
			Id:             1,
			ParentID:       Id{},
			Name:           "example project",
			Identifier:     "exampleproject",
			Description:    "This is an example project.",
			Homepage:       "http://github.com/cloudogu/go-redmine",
			IsPublic:       true,
			InheritMembers: true,
			CreatedOn:      "2021-02-19T16:51:03Z",
			UpdatedOn:      "2021-02-19T16:51:25Z",
		}
		assert.Equal(t, expectedProject, actualProject)
	})

	t.Run("should add basic auth to project GET request", func(t *testing.T) {
		actualCalledURL := ""

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			user, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, authUser, user)
			assert.Equal(t, authPassword, password)
			_, _ = fmt.Fprintln(w, simpleProjectJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		actualProject, err := sut.Project(1)

		require.NoError(t, err)
		require.NotEmpty(t, actualProject)
		expectedProject := &Project{
			Id:             1,
			ParentID:       Id{},
			Name:           "example project",
			Identifier:     "exampleproject",
			Description:    "This is an example project.",
			Homepage:       "http://github.com/cloudogu/go-redmine",
			IsPublic:       true,
			InheritMembers: true,
			CreatedOn:      "2021-02-19T16:51:03Z",
			UpdatedOn:      "2021-02-19T16:51:25Z",
		}
		assert.Equal(t, expectedProject, actualProject)
		assert.Equal(t, "/projects/1.json", actualCalledURL)
	})

	t.Run("should add auth token to project GET request", func(t *testing.T) {
		actualCalledURL := ""

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			_, _ = fmt.Fprintln(w, simpleProjectJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProject, err := sut.Project(1)

		// then
		require.NoError(t, err)
		require.NotEmpty(t, actualProject)
		expectedProject := &Project{
			Id:             1,
			ParentID:       Id{},
			Name:           "example project",
			Identifier:     "exampleproject",
			Description:    "This is an example project.",
			Homepage:       "http://github.com/cloudogu/go-redmine",
			IsPublic:       true,
			InheritMembers: true,
			CreatedOn:      "2021-02-19T16:51:03Z",
			UpdatedOn:      "2021-02-19T16:51:25Z",
		}
		assert.Equal(t, expectedProject, actualProject)
		assert.Equal(t, "/projects/1.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle non-existing projects as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProject, err := sut.Project(1)

		// then
		require.Error(t, err)
		require.Empty(t, actualProject)
		assert.Contains(t, err.Error(), "project (id: 1) was not found")
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProject, err := sut.Project(1)

		// then
		require.Error(t, err)
		require.Empty(t, actualProject)
		assert.Contains(t, err.Error(), "Something is not well\nAnother thing is also unacceptable")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProject, err := sut.Project(1)

		// then
		require.Error(t, err)
		require.Empty(t, actualProject)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_Projects(t *testing.T) {
	t.Run("should add basic auth to project GET request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, authUser, user)
			assert.Equal(t, authPassword, password)
			_, _ = fmt.Fprintln(w, simpleProjectsJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		actualProjects, err := sut.Projects()

		// then
		require.NoError(t, err)
		require.NotEmpty(t, actualProjects)
		expectedProject := []Project{
			{
				Id:             1,
				ParentID:       Id{},
				Name:           "example project",
				Identifier:     "exampleproject",
				Description:    "This is an example project.",
				Homepage:       "http://github.com/cloudogu/go-redmine",
				IsPublic:       true,
				InheritMembers: true,
				CreatedOn:      "2021-02-19T16:51:03Z",
				UpdatedOn:      "2021-02-19T16:51:25Z",
			},
		}
		assert.Equal(t, expectedProject, actualProjects)
	})

	t.Run("should add auth token to project GET request", func(t *testing.T) {
		actualCalledURL := ""

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			_, _ = fmt.Fprintln(w, simpleProjectsJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProjects, err := sut.Projects()

		// then
		require.NoError(t, err)
		require.NotEmpty(t, actualProjects)
		expectedProject := []Project{
			testProject,
		}
		assert.Equal(t, expectedProject, actualProjects)
		assert.Equal(t, "/projects.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProjects, err := sut.Projects()

		// then
		require.Error(t, err)
		require.Empty(t, actualProjects)
		assert.Contains(t, err.Error(), "Something is not well\nAnother thing is also unacceptable")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProjects, err := sut.Projects()

		// then
		require.Error(t, err)
		require.Empty(t, actualProjects)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_CreateProject(t *testing.T) {
	t.Run("should return without error on success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintln(w, simpleProjectJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProject, err := sut.CreateProject(testProject)

		// then
		require.NoError(t, err)
		assert.Equal(t, testProject, *actualProject)
	})

	t.Run("should add basic auth to project POST request", func(t *testing.T) {
		actualCalledURL := ""

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			user, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, authUser, user)
			assert.Equal(t, authPassword, password)

			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintln(w, simpleProjectJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		actualProject, err := sut.CreateProject(testProject)

		// then
		require.NoError(t, err)
		assert.Equal(t, testProject, *actualProject)
		assert.Equal(t, "/projects.json", actualCalledURL)
	})

	t.Run("should add auth token to project POST request", func(t *testing.T) {
		actualCalledURL := ""

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintln(w, simpleProjectJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProject, err := sut.CreateProject(testProject)

		// then
		require.NoError(t, err)
		assert.Equal(t, testProject, *actualProject)
		assert.Equal(t, "/projects.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProject, err := sut.CreateProject(testProject)

		// then
		require.Error(t, err)
		require.Empty(t, actualProject)
		assert.Contains(t, err.Error(), "Something is not well\nAnother thing is also unacceptable")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualProject, err := sut.CreateProject(testProject)

		// then
		require.Error(t, err)
		require.Empty(t, actualProject)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_UpdateProject(t *testing.T) {
	t.Run("should return without error on success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateProject(testProject)

		// then
		require.NoError(t, err)
	})

	t.Run("should add basic auth to project PUT request", func(t *testing.T) {
		actualCalledURL := ""

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			user, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, authUser, user)
			assert.Equal(t, authPassword, password)

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		err := sut.UpdateProject(testProject)

		// then
		require.NoError(t, err)
		assert.Equal(t, "/projects/1.json", actualCalledURL)
	})

	t.Run("should add auth token to project PUT request", func(t *testing.T) {
		actualCalledURL := ""

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateProject(testProject)

		// then
		require.NoError(t, err)
		assert.Equal(t, "/projects/1.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle non-existing projects as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateProject(testProject)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not update project (id: 1)")
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateProject(testProject)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_DeleteProject(t *testing.T) {
	t.Run("should return without error on success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteProject(1)

		// then
		require.NoError(t, err)
	})

	t.Run("should add basic auth to project DELETE request", func(t *testing.T) {
		actualCalledURL := ""

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			user, password, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, authUser, user)
			assert.Equal(t, authPassword, password)

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		err := sut.DeleteProject(1)

		// then
		require.NoError(t, err)
		assert.Equal(t, "/projects/1.json", actualCalledURL)
	})

	t.Run("should add auth token to project DELETE request", func(t *testing.T) {
		actualCalledURL := ""

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteProject(1)

		// then
		require.NoError(t, err)
		assert.Equal(t, "/projects/1.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle non-existing projects as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteProject(1)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not delete project (id: 1)")
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteProject(1)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func Test_jsonResourceEndpointByID(t *testing.T) {
	type args struct {
		baseURL      string
		resourceName string
		entityID     int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should return simple JSON endpoint", args{"http://1.2.3.4", "projects", 42}, "http://1.2.3.4/projects/42.json"},
		{"should return complex JSON endpoint ", args{"https://domain.ex-ample.com:3000/redmine", "projects", 1}, "https://domain.ex-ample.com:3000/redmine/projects/1.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonResourceEndpointByID(tt.args.baseURL, tt.args.resourceName, tt.args.entityID); got != tt.want {
				t.Errorf("jsonResourceEndpointByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jsonResourceEndpoint(t *testing.T) {
	type args struct {
		baseURL      string
		resourceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should return simple JSON endpoint", args{"http://1.2.3.4", "projects"}, "http://1.2.3.4/projects.json"},
		{"should return complex JSON endpoint ", args{"https://domain.ex-ample.com:3000/redmine", "projects"}, "https://domain.ex-ample.com:3000/redmine/projects.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonResourceEndpoint(tt.args.baseURL, tt.args.resourceName); got != tt.want {
				t.Errorf("jsonResourceEndpointByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
