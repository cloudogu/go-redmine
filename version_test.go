package redmine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testVersionBodyJSON = `{"id":1,"project":{"id":1,"name":"Test Project"},"name":"Sprint 2021-06","description":"Target version for sprint 2021-06","status":"open","due_date":"2021-04-01","sharing":"descendants","wiki_page_title":"wikipage","created_on":"2021-03-18T14:55:25Z","updated_on":"2021-03-18T15:05:53Z"}`
const testVersionJSON = `{"version":` + testVersionBodyJSON + "}"
const testVersionsJSON = `{"versions":[` + testVersionBodyJSON + `],"total_count":1}`

var testVersion1 = Version{
	Id:          1,
	Project:     IdName{Id: testProjectID, Name: "Test Project"},
	Name:        "Sprint 2021-06",
	Description: "Target version for sprint 2021-06",
	Status:      "open",
	DueDate:     "2021-04-01",
	CreatedOn:   "2021-03-18T14:55:25Z",
	UpdatedOn:   "2021-03-18T15:05:53Z",
}

func TestClient_Version(t *testing.T) {
	t.Run("should parse general Version fields", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintln(w, testVersionJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		actual, err := sut.Version(1)

		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := &testVersion1
		assert.Equal(t, expected, actual)
	})

	t.Run("should add basic auth to Version GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			_, _ = fmt.Fprintln(w, testVersionJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		actual, err := sut.Version(1)

		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := &testVersion1
		assert.Equal(t, expected, actual)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.True(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodGet, actualHTTPMethod)
		assert.Equal(t, "/versions/1.json", actualCalledURL)
	})

	t.Run("should add auth token to Version GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			_, _ = fmt.Fprintln(w, testVersionJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.Version(1)

		// then
		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := &testVersion1
		assert.Equal(t, expected, actual)
		assert.Empty(t, actualAuthUser)
		assert.Empty(t, actualAuthPass)
		assert.False(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodGet, actualHTTPMethod)
		assert.Equal(t, "/versions/1.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle non-existing versions as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.Version(1)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "version (id: 1) was not found")
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.Version(1)

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
		actual, err := sut.Version(1)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_Versions(t *testing.T) {
	t.Run("should add basic auth to versions GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			_, _ = fmt.Fprintln(w, testVersionsJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		actual, err := sut.Versions(testProjectID)

		// then
		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := []Version{testVersion1}
		assert.Equal(t, expected, actual)
		assert.Equal(t, expected, actual)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.True(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodGet, actualHTTPMethod)
		assert.Equal(t, "/projects/1/versions.json", actualCalledURL)
	})

	t.Run("should add auth token to version GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			_, _ = fmt.Fprintln(w, testVersionsJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.Versions(testProjectID)

		// then
		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := []Version{testVersion1}
		assert.Equal(t, expected, actual)
		assert.Empty(t, actualAuthUser)
		assert.Empty(t, actualAuthPass)
		assert.False(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodGet, actualHTTPMethod)
		assert.Equal(t, "/projects/1/versions.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.Versions(testProjectID)

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
		actual, err := sut.Versions(testProjectID)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_CreateVersion(t *testing.T) {
	t.Run("should return without error on success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintln(w, testVersionJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualVersion, err := sut.CreateVersion(testVersion1)

		// then
		require.NoError(t, err)
		assert.Equal(t, testVersion1, *actualVersion)
	})

	t.Run("should add basic auth to version POST request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintln(w, testVersionJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		actualVersion, err := sut.CreateVersion(testVersion1)

		// then
		require.NoError(t, err)
		assert.Equal(t, testVersion1, *actualVersion)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.True(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodPost, actualHTTPMethod)
		assert.Equal(t, "/projects/1/versions.json", actualCalledURL)
	})

	t.Run("should add auth token to version POST request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprintln(w, testVersionJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualVersion, err := sut.CreateVersion(testVersion1)

		// then
		require.NoError(t, err)
		assert.Equal(t, testVersion1, *actualVersion)
		assert.Empty(t, actualAuthUser)
		assert.Empty(t, actualAuthPass)
		assert.False(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodPost, actualHTTPMethod)
		assert.Equal(t, "/projects/1/versions.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualVersion, err := sut.CreateVersion(testVersion1)

		// then
		require.Error(t, err)
		require.Empty(t, actualVersion)
		assert.Contains(t, err.Error(), "Something is not well\nAnother thing is also unacceptable")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actualVersion, err := sut.CreateVersion(testVersion1)

		// then
		require.Error(t, err)
		require.Empty(t, actualVersion)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_UpdateVersion(t *testing.T) {
	t.Run("should return without error on success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateVersion(testVersion1)

		// then
		require.NoError(t, err)
	})

	t.Run("should add basic auth to version PUT request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		err := sut.UpdateVersion(testVersion1)

		// then
		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.True(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodPut, actualHTTPMethod)
		assert.Equal(t, "/versions/1.json", actualCalledURL)
	})

	t.Run("should add auth token to version PUT request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateVersion(testVersion1)

		// then
		require.NoError(t, err)
		assert.Empty(t, actualAuthUser)
		assert.Empty(t, actualAuthPass)
		assert.False(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodPut, actualHTTPMethod)
		assert.Equal(t, "/versions/1.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle non-existing version as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateVersion(testVersion1)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not update version (id: 1)")
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.UpdateVersion(testVersion1)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_DeleteVersion(t *testing.T) {
	t.Run("should return without error on success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteVersion(1)

		// then
		require.NoError(t, err)
	})

	t.Run("should add basic auth to version DELETE request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		err := sut.DeleteVersion(1)

		// then
		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.True(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodDelete, actualHTTPMethod)
		assert.Equal(t, "/versions/1.json", actualCalledURL)
	})

	t.Run("should add auth token to version DELETE request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteVersion(1)

		// then
		require.NoError(t, err)
		assert.Empty(t, actualAuthUser)
		assert.Empty(t, actualAuthPass)
		assert.False(t, actualBasicAuthOk)
		assert.Equal(t, http.MethodDelete, actualHTTPMethod)
		assert.Equal(t, "/versions/1.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle non-existing versions as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteVersion(1)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not delete version (id: 1)")
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("should handle body-less HTTP responses as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusUnauthorized)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		err := sut.DeleteVersion(1)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}
