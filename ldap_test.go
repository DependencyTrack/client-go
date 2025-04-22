package dtrack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdapMappings(t *testing.T) {
	// 1. Confirm absence of mappings
	// 2. Add mapping
	// 3. Confirm presence of mapping
	// 4. Remove mapping
	// 5. Confirm absence of mapping
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	team, err := client.Team.Create(context.Background(), Team{
		Name: "TestLdapMappings",
	})
	require.NoError(t, err)

	mappings, err := client.LDAP.GetTeamMappings(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Empty(t, mappings)

	mapping, err := client.LDAP.AddMapping(context.Background(), MappedLdapGroupRequest{
		Team:              team.UUID,
		DistinguishedName: "test.mapping.ldap.dependencytrack",
	})
	require.NoError(t, err)
	require.Equal(t, mapping.DistinguishedName, "test.mapping.ldap.dependencytrack")
	require.NotEmpty(t, mapping.UUID)

	mappings, err = client.LDAP.GetTeamMappings(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Equal(t, len(mappings), 1)
	require.Equal(t, mappings[0].DistinguishedName, mapping.DistinguishedName)
	require.Equal(t, mappings[0].UUID, mapping.UUID)

	err = client.LDAP.RemoveMapping(context.Background(), mapping.UUID)
	require.NoError(t, err)

	mappings, err = client.LDAP.GetTeamMappings(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Empty(t, mappings)
}

func TestLdapUsers(t *testing.T) {
	// 1. Confirm absence of users
	// 2. Create user
	// 3. Confirm presence of user
	// 4. Delete user
	// 5. Confirm absence of user
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	users, err := client.LDAP.GetUsers(context.Background(), PageOptions{
		PageSize: 10,
	})
	require.NoError(t, err)
	require.Equal(t, users.TotalCount, 0)
	require.Empty(t, users.Items)

	user, err := client.LDAP.CreateUser(context.Background(), LdapUser{
		Username:          "TestLdapUsers",
		DistinguishedName: "test.user.ldap.dependencytrack",
		Email:             "test@localhost",
	})
	require.NoError(t, err)
	require.Equal(t, user.Username, "TestLdapUsers")
	require.Equal(t, user.DistinguishedName, "Syncing...")
	require.Empty(t, user.Email)
	require.Empty(t, user.Permissions)
	require.Empty(t, user.Teams)

	users, err = client.LDAP.GetUsers(context.Background(), PageOptions{
		PageSize: 10,
	})
	require.NoError(t, err)
	require.Equal(t, users.TotalCount, 1)
	require.Equal(t, len(users.Items), 1)
	require.Equal(t, users.Items[0].DistinguishedName, user.DistinguishedName)
	require.Equal(t, users.Items[0].Email, user.Email)
	require.Equal(t, users.Items[0].Username, user.Username)

	err = client.LDAP.DeleteUser(context.Background(), user)
	require.NoError(t, err)

	users, err = client.LDAP.GetUsers(context.Background(), PageOptions{
		PageSize: 10,
	})
	require.NoError(t, err)
	require.Equal(t, users.TotalCount, 0)
	require.Empty(t, users.Items)
}
