package redmine

import (
	"testing"
)

func TestClient_DeactivatedUser(t *testing.T) {
	c := NewClient("https://ecosystem.cloudogu.com/redmine", "3fee4c85d2af5fde4873d909b53d79b5a20ee809")
	statue := Status{}
	statue.User.Status = 3
	err := c.DeactivatedUser(statue, 304)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
