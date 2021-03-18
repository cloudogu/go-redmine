package redmine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testIssueBodyJSON = `{
    "id": 1,
    "project": {
      "id": 1,
      "name": "example project1"
    },
    "tracker": {
      "id": 1,
      "name": "Bug"
    },
    "status": {
      "id": 1,
      "name": "New"
    },
    "priority": {
      "id": 2,
      "name": "Normal"
    },
    "author": {
      "id": 1,
      "name": "Redmine Admin"
    },
    "subject": "Something should be done",
    "description": "In this ticket an **important task** should be done1!\r\n\r\nGo ahead!\r\n\r\n` + "```bash\\r\\necho -n $PATH\\r\\n```" + `",
    "start_date": null,
    "due_date": null,
    "done_ratio": 0,
    "is_private": false,
    "estimated_hours": null,
    "total_estimated_hours": null,
    "spent_hours": 0,
    "total_spent_hours": 0,
    "created_on": "2021-02-23T14:20:48Z",
    "updated_on": "2021-02-23T14:39:02Z",
    "closed_on": null
  }`
const testIssueJSON = `{"issue":` + testIssueBodyJSON + `}`
const testIssuesJSON = `{"issues":[` + testIssueBodyJSON + `],"total_count":1,"offset":0,"limit":25}`
const projectID = 1

var testIssue = Issue{
	Id:          1,
	Subject:     "Something should be done",
	Description: "In this ticket an **important task** should be done1!\\r\\n\\r\\nGo ahead!\\r\\n\\r\\n` + \"```bash\\\\r\\\\necho -n $PATH\\\\r\\\\n```\" + `",
	ProjectId:   1,
	Project:     &IdName{Id: 1, Name: "example project1"},
	TrackerId:   0,
	Tracker:     &IdName{Id: 1, Name: "Bug"},
	ParentId:    0,
	Parent:      nil,
	StatusId:    1,
	Status:      nil,
	PriorityId:  2,
	Priority:    nil,
	Author:      nil,
	AssignedTo:  nil,
	Category:    nil,
	CategoryId:  0,
	CreatedOn:   "2021-02-23T14:20:48Z",
	UpdatedOn:   "2021-02-23T14:39:02Z",
	ClosedOn:    "",
}

func Test_getOneIssue(t *testing.T) {
	t.Run("should add auth token to issue GET request", func(t *testing.T) {
		actualCalledURL := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		_, err := getOneIssue(sut, 1, nil)

		require.NoError(t, err)
		assert.Equal(t, "/issues/1.json?key="+testAPIToken, actualCalledURL)
	})

	t.Run("should add basic auth to issue GET request", func(t *testing.T) {
		actualCalledURL := ""
		actualAuthUser := ""
		actualAuthPass := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			var ok bool
			actualAuthUser, actualAuthPass, ok = r.BasicAuth()
			assert.True(t, ok)
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		_, err := getOneIssue(sut, 1, nil)

		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.Equal(t, "/issues/1.json", actualCalledURL)
	})

	t.Run("should parse simple issue JSON without additional arguments", func(t *testing.T) {
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualHTTPMethod = r.Method
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		actual, err := getOneIssue(sut, 1, nil)

		// then
		require.NoError(t, err)
		assert.Equal(t, httpMethodGet, actualHTTPMethod)
		assert.Equal(t, 1, actual.Id)
		assert.Equal(t, "Something should be done", actual.Subject)
		assert.Equal(t, "In this ticket an **important task** should be done1!\r\n\r\nGo ahead!\r\n\r\n"+"```bash\r\necho -n $PATH\r\n```", actual.Description)
		assert.Equal(t, 0, actual.ProjectId)
		assert.Equal(t, IdName{Id: 1, Name: "example project1"}, *actual.Project)
		assert.Equal(t, 0, actual.TrackerId)
		assert.Equal(t, IdName{Id: 1, Name: "Bug"}, *actual.Tracker)
		assert.Equal(t, 0, actual.ParentId)
		assert.Nil(t, actual.Parent)
		assert.Equal(t, 0, actual.StatusId)
		assert.Equal(t, IdName{Id: 1, Name: "New"}, *actual.Status)
		assert.Equal(t, 0, actual.PriorityId)
		assert.Equal(t, IdName{Id: 2, Name: "Normal"}, *actual.Priority)
		assert.Equal(t, IdName{Id: 1, Name: "Redmine Admin"}, *actual.Author)
		assert.Equal(t, "2021-02-23T14:20:48Z", actual.CreatedOn)
		assert.Equal(t, "2021-02-23T14:39:02Z", actual.UpdatedOn)
		assert.Equal(t, "", actual.StartDate)
		assert.Equal(t, "", actual.DueDate)
		assert.Equal(t, "", actual.ClosedOn)
	})

	t.Run("should handle non-existing issues as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := getOneIssue(sut, 1, nil)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "issue (id: 1) was not found")
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := getOneIssue(sut, 1, nil)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "Something is not well\nAnother thing is also unacceptable")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := getOneIssue(sut, 1, nil)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_IssuesOf(t *testing.T) {
	t.Run("should add auth token to issue GET request", func(t *testing.T) {
		var actualCalledURLs []string
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualHTTPMethod = r.Method
			actualCalledURLs = append(actualCalledURLs, r.URL.String())
			_, _, ok := r.BasicAuth()
			assert.False(t, ok)
			_, _ = fmt.Fprintln(w, testIssuesJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		_, err := sut.IssuesOf(projectID)

		// then
		require.NoError(t, err)
		assert.Equal(t, httpMethodGet, actualHTTPMethod)
		assert.Len(t, actualCalledURLs, 2)
		assert.Contains(t, actualCalledURLs[0], "/issues.json?")
		assert.Contains(t, actualCalledURLs[0], "project_id=1")
		assert.Contains(t, actualCalledURLs[0], "key="+testAPIToken)
		assert.Contains(t, actualCalledURLs[0], "offset=0")
		assert.Contains(t, actualCalledURLs[1], "/issues.json?")
		assert.Contains(t, actualCalledURLs[1], "project_id=1")
		assert.Contains(t, actualCalledURLs[1], "key="+testAPIToken)
		assert.Contains(t, actualCalledURLs[1], "offset=1")
	})

	t.Run("should add basic auth to issue GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		var actualCalledURLs []string
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURLs = append(actualCalledURLs, r.URL.String())
			var ok bool
			actualAuthUser, actualAuthPass, ok = r.BasicAuth()
			assert.True(t, ok)
			_, _ = fmt.Fprintln(w, testIssuesJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		_, err := sut.IssuesOf(projectID)

		// then
		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.Len(t, actualCalledURLs, 2)
		assert.Contains(t, actualCalledURLs[0], "/issues.json?")
		assert.Contains(t, actualCalledURLs[0], "project_id=1")
		assert.NotContains(t, actualCalledURLs[0], testAPIToken)
		assert.NotContains(t, actualCalledURLs[0], "key=")
	})

	t.Run("should parse simple issue JSON without additional arguments", func(t *testing.T) {
		var actualCalledURLs []string
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURLs = append(actualCalledURLs, r.URL.String())

			if r.URL.Query().Get("offset") == "0" {
				_, _ = fmt.Fprintln(w, testIssuesJSON)
			} else {
				fakeOffsetResponse := `{"issues":[],"total_count":1,"offset":1,"limit":25}`
				_, _ = fmt.Fprintln(w, fakeOffsetResponse)
			}

		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		actualIssues, err := sut.IssuesOf(projectID)

		// then
		require.NoError(t, err)
		require.Len(t, actualIssues, 1)
		actual := actualIssues[0]
		assert.Equal(t, 1, actual.Id)
		assert.Equal(t, "Something should be done", actual.Subject)
		assert.Equal(t, "In this ticket an **important task** should be done1!\r\n\r\nGo ahead!\r\n\r\n"+"```bash\r\necho -n $PATH\r\n```", actual.Description)
		assert.Equal(t, 0, actual.ProjectId)
		assert.Equal(t, IdName{Id: 1, Name: "example project1"}, *actual.Project)
		assert.Equal(t, 0, actual.TrackerId)
		assert.Equal(t, IdName{Id: 1, Name: "Bug"}, *actual.Tracker)
		assert.Equal(t, 0, actual.ParentId)
		assert.Nil(t, actual.Parent)
		assert.Equal(t, 0, actual.StatusId)
		assert.Equal(t, IdName{Id: 1, Name: "New"}, *actual.Status)
		assert.Equal(t, 0, actual.PriorityId)
		assert.Equal(t, IdName{Id: 2, Name: "Normal"}, *actual.Priority)
		assert.Equal(t, IdName{Id: 1, Name: "Redmine Admin"}, *actual.Author)
		assert.Equal(t, "2021-02-23T14:20:48Z", actual.CreatedOn)
		assert.Equal(t, "2021-02-23T14:39:02Z", actual.UpdatedOn)
		assert.Equal(t, "", actual.StartDate)
		assert.Equal(t, "", actual.DueDate)
		assert.Equal(t, "", actual.ClosedOn)
	})
}

func TestClient_Issue(t *testing.T) {
	t.Run("should add auth token to issue GET request", func(t *testing.T) {
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualHTTPMethod = r.Method
			actualCalledURL = r.URL.String()
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		_, err := sut.Issue(1)

		// then
		require.NoError(t, err)
		assert.Equal(t, httpMethodGet, actualHTTPMethod)
		assert.Equal(t, "/issues/1.json?key="+testAPIToken, actualCalledURL)
	})

	t.Run("should add basic auth to issue GET request", func(t *testing.T) {
		actualCalledURL := ""
		actualAuthUser := ""
		actualAuthPass := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			var ok bool
			actualAuthUser, actualAuthPass, ok = r.BasicAuth()
			assert.True(t, ok)
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		_, err := sut.Issue(1)

		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.Equal(t, "/issues/1.json", actualCalledURL)
	})

	t.Run("should parse simple issue JSON without additional arguments", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		actual, err := sut.Issue(1)

		// then
		require.NoError(t, err)
		assert.Equal(t, 1, actual.Id)
		assert.Equal(t, "Something should be done", actual.Subject)
		assert.Equal(t, "In this ticket an **important task** should be done1!\r\n\r\nGo ahead!\r\n\r\n"+"```bash\r\necho -n $PATH\r\n```", actual.Description)
		assert.Equal(t, 0, actual.ProjectId)
		assert.Equal(t, IdName{Id: 1, Name: "example project1"}, *actual.Project)
		assert.Equal(t, 0, actual.TrackerId)
		assert.Equal(t, IdName{Id: 1, Name: "Bug"}, *actual.Tracker)
		assert.Equal(t, 0, actual.ParentId)
		assert.Nil(t, actual.Parent)
		assert.Equal(t, 0, actual.StatusId)
		assert.Equal(t, IdName{Id: 1, Name: "New"}, *actual.Status)
		assert.Equal(t, 0, actual.PriorityId)
		assert.Equal(t, IdName{Id: 2, Name: "Normal"}, *actual.Priority)
		assert.Equal(t, IdName{Id: 1, Name: "Redmine Admin"}, *actual.Author)
		assert.Equal(t, "2021-02-23T14:20:48Z", actual.CreatedOn)
		assert.Equal(t, "2021-02-23T14:39:02Z", actual.UpdatedOn)
		assert.Equal(t, "", actual.StartDate)
		assert.Equal(t, "", actual.DueDate)
		assert.Equal(t, "", actual.ClosedOn)
	})

	t.Run("should handle non-existing issues as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.Issue(1)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "issue (id: 1) was not found")
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.Issue(1)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "Something is not well\nAnother thing is also unacceptable")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.Issue(1)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_IssueWithArgs(t *testing.T) {
	actualCalledURL := ""
	actualHTTPMethod := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualHTTPMethod = r.Method
		actualCalledURL = r.URL.String()
		_, _ = fmt.Fprintln(w, testIssueJSON)
	}))
	defer ts.Close()

	sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()
	args := make(map[string]string, 1)
	args["leKey"] = "leValue"

	// when
	actual, err := sut.IssueWithArgs(1, args)

	// then
	require.NoError(t, err)
	assert.Equal(t, httpMethodGet, actualHTTPMethod)
	assert.Contains(t, actualCalledURL, "leKey=leValue")
	assert.Equal(t, 1, actual.Id)
	assert.Equal(t, "Something should be done", actual.Subject)
	assert.Equal(t, "In this ticket an **important task** should be done1!\r\n\r\nGo ahead!\r\n\r\n"+"```bash\r\necho -n $PATH\r\n```", actual.Description)
	assert.Equal(t, 0, actual.ProjectId)
	assert.Equal(t, IdName{Id: 1, Name: "example project1"}, *actual.Project)
	assert.Equal(t, 0, actual.TrackerId)
	assert.Equal(t, IdName{Id: 1, Name: "Bug"}, *actual.Tracker)
	assert.Equal(t, 0, actual.ParentId)
	assert.Nil(t, actual.Parent)
	assert.Equal(t, 0, actual.StatusId)
	assert.Equal(t, IdName{Id: 1, Name: "New"}, *actual.Status)
	assert.Equal(t, 0, actual.PriorityId)
	assert.Equal(t, IdName{Id: 2, Name: "Normal"}, *actual.Priority)
	assert.Equal(t, IdName{Id: 1, Name: "Redmine Admin"}, *actual.Author)
	assert.Equal(t, "2021-02-23T14:20:48Z", actual.CreatedOn)
	assert.Equal(t, "2021-02-23T14:39:02Z", actual.UpdatedOn)
	assert.Equal(t, "", actual.StartDate)
	assert.Equal(t, "", actual.DueDate)
	assert.Equal(t, "", actual.ClosedOn)
}

func TestClient_Issues(t *testing.T) {
	t.Run("should add auth token to issue GET request", func(t *testing.T) {
		var actualCalledURLs []string
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualHTTPMethod = r.Method
			actualCalledURLs = append(actualCalledURLs, r.URL.String())
			_, _, ok := r.BasicAuth()
			assert.False(t, ok)
			if r.URL.Query().Get("offset") == "0" {
				_, _ = fmt.Fprintln(w, testIssuesJSON)
			} else {
				fakeOffsetResponse := `{"issues":[],"total_count":1,"offset":1,"limit":25}`
				_, _ = fmt.Fprintln(w, fakeOffsetResponse)
			}
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		_, err := sut.Issues()

		// then
		require.NoError(t, err)
		assert.Equal(t, httpMethodGet, actualHTTPMethod)
		assert.Len(t, actualCalledURLs, 2)
		assert.Equal(t, "/issues.json?key="+testAPIToken+"&offset=0", actualCalledURLs[0])
		assert.Equal(t, "/issues.json?key="+testAPIToken+"&offset=1", actualCalledURLs[1])
	})

	t.Run("should add basic auth to issue GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		var actualCalledURLs []string
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURLs = append(actualCalledURLs, r.URL.String())
			var ok bool
			actualAuthUser, actualAuthPass, ok = r.BasicAuth()
			assert.True(t, ok)
			_, _ = fmt.Fprintln(w, testIssuesJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		_, err := sut.Issues()

		// then
		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.Len(t, actualCalledURLs, 2)
		assert.Contains(t, actualCalledURLs[0], "/issues.json")
		assert.NotContains(t, actualCalledURLs[0], testAPIToken)
		assert.NotContains(t, actualCalledURLs[0], "key=")
	})

	t.Run("should parse simple issue JSON without additional arguments", func(t *testing.T) {
		var actualCalledURLs []string
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURLs = append(actualCalledURLs, r.URL.String())

			if r.URL.Query().Get("offset") == "0" {
				_, _ = fmt.Fprintln(w, testIssuesJSON)
			} else {
				fakeOffsetResponse := `{"issues":[],"total_count":1,"offset":1,"limit":25}`
				_, _ = fmt.Fprintln(w, fakeOffsetResponse)
			}
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		actualIssues, err := sut.Issues()

		// then
		require.NoError(t, err)
		require.Len(t, actualIssues, 1)
		actual := actualIssues[0]
		assert.Equal(t, 1, actual.Id)
		assert.Equal(t, "Something should be done", actual.Subject)
		assert.Equal(t, "In this ticket an **important task** should be done1!\r\n\r\nGo ahead!\r\n\r\n"+"```bash\r\necho -n $PATH\r\n```", actual.Description)
		assert.Equal(t, 0, actual.ProjectId)
		assert.Equal(t, IdName{Id: 1, Name: "example project1"}, *actual.Project)
		assert.Equal(t, 0, actual.TrackerId)
		assert.Equal(t, IdName{Id: 1, Name: "Bug"}, *actual.Tracker)
		assert.Equal(t, 0, actual.ParentId)
		assert.Nil(t, actual.Parent)
		assert.Equal(t, 0, actual.StatusId)
		assert.Equal(t, IdName{Id: 1, Name: "New"}, *actual.Status)
		assert.Equal(t, 0, actual.PriorityId)
		assert.Equal(t, IdName{Id: 2, Name: "Normal"}, *actual.Priority)
		assert.Equal(t, IdName{Id: 1, Name: "Redmine Admin"}, *actual.Author)
		assert.Equal(t, "2021-02-23T14:20:48Z", actual.CreatedOn)
		assert.Equal(t, "2021-02-23T14:39:02Z", actual.UpdatedOn)
		assert.Equal(t, "", actual.StartDate)
		assert.Equal(t, "", actual.DueDate)
		assert.Equal(t, "", actual.ClosedOn)
	})
}

func TestClient_CreateIssue(t *testing.T) {
	t.Run("should add auth token to issue POST request", func(t *testing.T) {
		actualCalledURL := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, ok := r.BasicAuth()
			assert.False(t, ok)

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			actualCalledURL = r.URL.String()
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		_, err := sut.CreateIssue(testIssue)

		require.NoError(t, err)
		assert.Equal(t, "/issues.json?key="+testAPIToken, actualCalledURL)
	})

	t.Run("should add basic auth to issue POST request", func(t *testing.T) {
		actualCalledURL := ""
		actualAuthUser := ""
		actualAuthPass := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			var ok bool
			actualAuthUser, actualAuthPass, ok = r.BasicAuth()
			assert.True(t, ok)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		_, err := sut.CreateIssue(testIssue)

		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.Equal(t, "/issues.json", actualCalledURL)
		assert.NotContains(t, actualCalledURL, testAPIToken)
		assert.NotContains(t, actualCalledURL, "key=")
	})

	t.Run("should parse simple issue JSON without additional arguments", func(t *testing.T) {
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualHTTPMethod = r.Method
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		actual, err := sut.CreateIssue(testIssue)

		// then
		require.NoError(t, err)
		assert.Equal(t, httpMethodPost, actualHTTPMethod)
		assert.Equal(t, 1, actual.Id)
		assert.Equal(t, "Something should be done", actual.Subject)
		assert.Equal(t, "In this ticket an **important task** should be done1!\r\n\r\nGo ahead!\r\n\r\n"+"```bash\r\necho -n $PATH\r\n```", actual.Description)
		assert.Equal(t, 0, actual.ProjectId)
		assert.Equal(t, IdName{Id: 1, Name: "example project1"}, *actual.Project)
		assert.Equal(t, 0, actual.TrackerId)
		assert.Equal(t, IdName{Id: 1, Name: "Bug"}, *actual.Tracker)
		assert.Equal(t, 0, actual.ParentId)
		assert.Nil(t, actual.Parent)
		assert.Equal(t, 0, actual.StatusId)
		assert.Equal(t, IdName{Id: 1, Name: "New"}, *actual.Status)
		assert.Equal(t, 0, actual.PriorityId)
		assert.Equal(t, IdName{Id: 2, Name: "Normal"}, *actual.Priority)
		assert.Equal(t, IdName{Id: 1, Name: "Redmine Admin"}, *actual.Author)
		assert.Equal(t, "2021-02-23T14:20:48Z", actual.CreatedOn)
		assert.Equal(t, "2021-02-23T14:39:02Z", actual.UpdatedOn)
		assert.Equal(t, "", actual.StartDate)
		assert.Equal(t, "", actual.DueDate)
		assert.Equal(t, "", actual.ClosedOn)
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.CreateIssue(testIssue)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "Something is not well\nAnother thing is also unacceptable")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.CreateIssue(testIssue)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_DeleteIssue(t *testing.T) {
	t.Run("should add auth token to issue DELETE request", func(t *testing.T) {
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		err := sut.DeleteIssue(1)

		require.NoError(t, err)
		assert.Equal(t, httpMethodDelete, actualHTTPMethod)
		assert.Equal(t, "/issues/1.json?key="+testAPIToken, actualCalledURL)
	})

	t.Run("should add basic auth to issue DELETE request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualCalledURL := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			var ok bool
			actualAuthUser, actualAuthPass, ok = r.BasicAuth()
			assert.True(t, ok)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		err := sut.DeleteIssue(1)

		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.Equal(t, "/issues/1.json", actualCalledURL)
	})

	t.Run("should parse simple issue JSON without additional arguments", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		err := sut.DeleteIssue(1)

		// then
		require.NoError(t, err)
	})

	t.Run("should handle non-existing issues as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteIssue(1)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not delete issue (id: 1) because it was not found")
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteIssueCategory(1)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Something is not well\nAnother thing is also unacceptable")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteIssueCategory(1)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_UpdateIssue(t *testing.T) {
	t.Run("should add auth token to issue PUT request", func(t *testing.T) {
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualHTTPMethod = r.Method
			_, _, ok := r.BasicAuth()
			assert.False(t, ok)
			actualCalledURL = r.URL.String()
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		err := sut.UpdateIssue(testIssue)

		require.NoError(t, err)
		assert.Equal(t, httpMethodPut, actualHTTPMethod)
		assert.Equal(t, "/issues/1.json?key="+testAPIToken, actualCalledURL)
	})

	t.Run("should add basic auth to issue PUT request", func(t *testing.T) {
		actualCalledURL := ""
		actualAuthUser := ""
		actualAuthPass := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			var ok bool
			actualAuthUser, actualAuthPass, ok = r.BasicAuth()
			assert.True(t, ok)
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		err := sut.UpdateIssue(testIssue)

		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.Equal(t, "/issues/1.json", actualCalledURL)
	})

	t.Run("should parse simple issue JSON without additional arguments", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintln(w, testIssueJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		err := sut.UpdateIssue(testIssue)

		// then
		require.NoError(t, err)
	})

	t.Run("should handle non-existing issues as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateIssue(testIssue)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not update issue (id: 1) because it was not found")
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateIssue(testIssue)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Something is not well\nAnother thing is also unacceptable")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateIssue(testIssue)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestIssue_GetTitle(t *testing.T) {
	type fields struct {
		Id      int
		Subject string
		Tracker *IdName
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "simple issue", fields: fields{Id: 1, Subject: "subject", Tracker: &IdName{Id: 1, Name: "Bug"}}, want: "Bug #1: subject"},
		{name: "special chars", fields: fields{Id: 2, Subject: "Schönere UI", Tracker: &IdName{Id: 2, Name: "User Story"}}, want: "User Story #2: Schönere UI"},
		{name: "hashtags in subject", fields: fields{Id: 3, Subject: "Superseded by #2", Tracker: &IdName{Id: 2, Name: "User Story"}}, want: "User Story #3: Superseded by #2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issue := &Issue{
				Id:      tt.fields.Id,
				Subject: tt.fields.Subject,
				Tracker: tt.fields.Tracker,
			}
			if got := issue.GetTitle(); got != tt.want {
				t.Errorf("GetTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_argsToKeyValues(t *testing.T) {
	t.Run("should return nil slice", func(t *testing.T) {
		var args map[string]string

		actual := argsToKeyValues(args)

		assert.Nil(t, actual)
	})
	t.Run("should return single element slice", func(t *testing.T) {
		args := make(map[string]string)
		args["leKey"] = "leValue"

		actual := argsToKeyValues(args)

		require.NotNil(t, actual)
		expected := []keyValue{{key: "leKey", value: "leValue"}}
		assert.Equal(t, expected, actual)
	})
	t.Run("should return multi-element slice", func(t *testing.T) {
		args := make(map[string]string)
		args["leKey"] = "leValue"
		args["spaceship"] = "Heart of Gold"
		args["universe"] = "42"

		actual := argsToKeyValues(args)

		require.NotNil(t, actual)
		expected := []keyValue{
			{key: "leKey", value: "leValue"},
			{key: "spaceship", value: "Heart of Gold"},
			{key: "universe", value: "42"},
		}
		assert.ElementsMatch(t, expected, actual)
	})
}
