package redmine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testIssuePrioritiesJSON = `{"issue_priorities":[{"id":1,"name":"Low","is_default":false,"active":true},{"id":2,"name":"Normal","is_default":true,"active":true},{"id":3,"name":"High","is_default":false,"active":true},{"id":4,"name":"Urgent","is_default":false,"active":true},{"id":5,"name":"Immediate","is_default":false,"active":true}]}`

func TestClient_IssuePriorities(t *testing.T) {
	t.Run("should add auth token to issue GET request", func(t *testing.T) {
		actualCalledURL := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualCalledURL = r.URL.String()
			_, _ = fmt.Fprintln(w, testIssuePrioritiesJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		_, err := sut.IssuePriorities()

		require.NoError(t, err)
		assert.Equal(t, "/enumerations/issue_priorities.json?key="+testAPIToken, actualCalledURL)
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
			_, _ = fmt.Fprintln(w, testIssuePrioritiesJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthBasicAuth(authUser, authPassword).Build()

		_, err := sut.IssuePriorities()

		require.NoError(t, err)
		assert.Equal(t, authUser, actualAuthUser)
		assert.Equal(t, authPassword, actualAuthPass)
		assert.Equal(t, "/enumerations/issue_priorities.json", actualCalledURL)
	})

	t.Run("should parse simple issue JSON without additional arguments", func(t *testing.T) {
		actualHTTPMethod := ""
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualHTTPMethod = r.Method
			_, _ = fmt.Fprintln(w, testIssuePrioritiesJSON)
		}))
		defer ts.Close()

		sut, _ := NewClientBuilder().Endpoint(ts.URL).AuthAPIToken(testAPIToken).Build()

		// when
		actual, err := sut.IssuePriorities()

		// then
		require.NoError(t, err)
		assert.Equal(t, http.MethodGet, actualHTTPMethod)
		expectedIssuePriorities := []IssuePriority{
			{Id: 1, Name: "Low"},
			{Id: 2, Name: "Normal", IsDefault: true},
			{Id: 3, Name: "High"},
			{Id: 4, Name: "Urgent"},
			{Id: 5, Name: "Immediate"},
		}
		assert.Equal(t, expectedIssuePriorities, actual)
	})
}
