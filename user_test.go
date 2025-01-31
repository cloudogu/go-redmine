package redmine

import (
	"testing"
)

const REDMINE_TEST_ENDPOINT = "https://placeholder.com/redmine"

func TestClient_SetUserStatus(t *testing.T) {
	c := NewClient(REDMINE_TEST_ENDPOINT, "")
	statue := Status{}
	statue.User.Status = 3
	err := c.SetUserStatus(statue, 304)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestClient_GetAllUser(t *testing.T) {
	c := NewClient(REDMINE_TEST_ENDPOINT, "")
	statue := Status{}
	statue.User.Status = 3
	Users, err := c.AllUsers()
	if err != nil {
		t.Fatalf(err.Error())
	}
	print(len(Users))
}

func TestClient_GetTotalCount(t *testing.T) {
	c := NewClient(REDMINE_TEST_ENDPOINT, "")
	statue := Status{}
	statue.User.Status = 3
	num, err := c.totalCount()
	if err != nil {
		t.Fatalf(err.Error())
	}
	print(num)
}
