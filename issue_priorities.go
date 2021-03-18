package redmine

import (
	"encoding/json"
	errors2 "github.com/pkg/errors"
	"net/http"
)

type issuePrioritiesResult struct {
	IssuePriorities []IssuePriority `json:"issue_priorities"`
}

type IssuePriority struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

func (c *Client) IssuePriorities() ([]IssuePriority, error) {
	url := jsonResourceEndpoint(c.endpoint, "enumerations/issue_priorities")
	req, err := c.authenticatedGet(url)
	if err != nil {
		return nil, errors2.Wrap(err, "error while creating GET request for issue priorities")
	}
	err = safelySetQueryParameters(req, c.getPaginationClauseParams())
	if err != nil {
		return nil, errors2.Wrap(err, "error while adding additional parameters to issue priorities request")
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, errors2.Wrap(err, "error while reading issue priorities response")
	}
	defer res.Body.Close()

	var r issuePrioritiesResult
	if !isHTTPStatusSuccessful(res.StatusCode, []int{http.StatusOK}) {
		return nil, errors2.Wrapf(decodeHTTPError(res), "issue request returned non-successfully, URL: %s", req.URL.String())
	}

	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r.IssuePriorities, nil
}
