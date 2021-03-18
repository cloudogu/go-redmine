package redmine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testIssueCategory1BodyJSON = `{"id":1,"project":{"id":1,"name":"Test Project"},"name":"Important Product"}`
const testIssueCategory2BodyJSON = `{"id":2,"project":{"id":1,"name":"Test Project"},"name":"Service \u0026 Maintenance","assigned_to":{"id":1,"name":"Redmine Admin"}}`
const testIssueCategoryJSON = `{ "issue_category":` + testIssueCategory1BodyJSON + "}"
const testIssueCategoriesJSON = `{"issue_categories":[` + testIssueCategory1BodyJSON + `],"total_count":1,"offset":0,"limit":25}`
const testProjectID = 1

var testIssueCategory1 = IssueCategory{
	Project:    IdName{Id: testProjectID, Name: "Test Project"},
	Id:         1,
	Name:       "Important Product",
	AssignedTo: IdName{},
}

func TestClient_IssueCategory(t *testing.T) {
	t.Run("should parse general issueCategory fields", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintln(w, testIssueCategoryJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		actual, err := sut.IssueCategory(1)

		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := &testIssueCategory1
		assert.Equal(t, expected, actual)
	})

	t.Run("should add basic auth to issueCategory GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			_, _ = fmt.Fprintln(w, testIssueCategoryJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		actual, err := sut.IssueCategory(1)

		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := &testIssueCategory1
		assert.Equal(t, expected, actual)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.True(t, actualBasicAuthOk)
		assert.Equal(t, httpMethodGet, actualHTTPMethod)
		assert.Equal(t, "/issue_categories/1.json", actualCalledURL)
	})

	t.Run("should add auth token to issueCategory GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			_, _ = fmt.Fprintln(w, testIssueCategoryJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.IssueCategory(1)

		// then
		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := &testIssueCategory1
		assert.Equal(t, expected, actual)
		assert.Empty(t, actualAuthUser)
		assert.Empty(t, actualAuthPass)
		assert.False(t, actualBasicAuthOk)
		assert.Equal(t, httpMethodGet, actualHTTPMethod)
		assert.Equal(t, "/issue_categories/1.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle non-existing issue_categories as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.IssueCategory(1)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "issue category (id: 1) was not found")
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.IssueCategory(1)

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
		actual, err := sut.IssueCategory(1)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_IssueCategories(t *testing.T) {
	t.Run("should add basic auth to issue categories GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			_, _ = fmt.Fprintln(w, testIssueCategoriesJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		// when
		actual, err := sut.IssueCategories(testProjectID)

		// then
		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := []IssueCategory{testIssueCategory1}
		assert.Equal(t, expected, actual)
		assert.Equal(t, expected, actual)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.True(t, actualBasicAuthOk)
		assert.Equal(t, httpMethodGet, actualHTTPMethod)
		assert.Equal(t, "/projects/1/issue_categories.json", actualCalledURL)
	})

	t.Run("should add auth token to issue categories GET request", func(t *testing.T) {
		actualAuthUser := ""
		actualAuthPass := ""
		actualBasicAuthOk := false
		actualCalledURL := ""
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			actualHTTPMethod = r.Method
			actualAuthUser, actualAuthPass, actualBasicAuthOk = r.BasicAuth()

			_, _ = fmt.Fprintln(w, testIssueCategoriesJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.IssueCategories(testProjectID)

		// then
		require.NoError(t, err)
		require.NotEmpty(t, actual)
		expected := []IssueCategory{testIssueCategory1}
		assert.Equal(t, expected, actual)
		assert.Empty(t, actualAuthUser)
		assert.Empty(t, actualAuthPass)
		assert.False(t, actualBasicAuthOk)
		assert.Equal(t, httpMethodGet, actualHTTPMethod)
		assert.Equal(t, "/projects/1/issue_categories.json?key=123456789", actualCalledURL)
	})

	t.Run("should handle HTTP 422 errors as error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errorAsJson := `{ "errors":[ "Something is not well", "Another thing is also unacceptable" ] }`
			http.Error(w, errorAsJson, http.StatusUnprocessableEntity)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(authToken).Build()

		// when
		actual, err := sut.IssueCategories(testProjectID)

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
		actual, err := sut.IssueCategories(testProjectID)

		// then
		require.Error(t, err)
		require.Empty(t, actual)
		assert.Contains(t, err.Error(), "HTTP 401 Unauthorized")
	})
}

func TestClient_CreateIssueCategory(t *testing.T) {

}

func TestClient_DeleteIssueCategory(t *testing.T) {

}

func TestClient_UpdateIssueCategory(t *testing.T) {

}
