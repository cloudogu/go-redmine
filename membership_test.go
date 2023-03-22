package redmine

import (
	"testing"
)

func TestClient_CreateMembershipByProjectID(t *testing.T) {
	client := NewClient("https://ecosystem.cloudogu.com/redmine", "3fee4c85d2af5fde4873d909b53d79b5a20ee809")
	memberEntry := MembershipDTO{}
	memberEntry.Membership.RoleIds = append(memberEntry.Membership.RoleIds, 4)
	memberEntry.Membership.UserId = 299
	member, err := client.CreateMembershipByProjectID(memberEntry, 125)
	if err != nil {
		t.Fatalf(err.Error())
	}
	println(member.Id)
}
