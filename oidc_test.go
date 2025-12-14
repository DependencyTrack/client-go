package dtrack

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOIDCGroup(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})
	// Check absence
	{
		groups, err := client.OIDC.GetAllGroups(ctx)
		require.NoError(t, err)
		require.Empty(t, groups)
	}
	// Create Group
	group, err := client.OIDC.CreateGroup(ctx, "Test_Group")
	require.NoError(t, err)
	require.Equal(t, group.Name, "Test_Group")
	require.NotZero(t, group.UUID)

	// Check presence
	{
		groups, err := client.OIDC.GetAllGroups(ctx)
		require.NoError(t, err)
		require.Equal(t, len(groups), 1)
		require.Equal(t, groups[0], group)
	}

	// Update Group
	{
		updated, err := client.OIDC.UpdateGroup(ctx, OIDCGroup{
			UUID: group.UUID,
			Name: "Updated_Test_Group",
		})
		require.NoError(t, err)
		require.Equal(t, updated.UUID, group.UUID)
		require.Equal(t, updated.Name, "Updated_Test_Group")
	}

	// Check updated
	{
		groups, err := client.OIDC.GetAllGroups(ctx)
		require.NoError(t, err)
		require.Equal(t, len(groups), 1)
		require.Equal(t, groups[0], OIDCGroup{
			UUID: group.UUID,
			Name: "Updated_Test_Group",
		})
	}

	// Delete Group
	{
		err := client.OIDC.DeleteGroup(ctx, group.UUID)
		require.NoError(t, err)
	}

	// Check absence
	{
		groups, err := client.OIDC.GetAllGroups(ctx)
		require.NoError(t, err)
		require.Empty(t, groups)
	}
}

func TestOIDCTeamMappings(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	// Add Team
	team, err := client.Team.Create(ctx, Team{
		Name: "Team_Name",
	})
	require.NoError(t, err)

	// Add OIDC Group
	group, err := client.OIDC.CreateGroup(ctx, "Group_Name")
	require.NoError(t, err)

	// Check absence
	{
		teams, err := client.OIDC.GetAllTeamsOf(ctx, group)
		require.NoError(t, err)
		require.Empty(t, teams)
	}

	// Add Mapping
	mapping, err := client.OIDC.AddTeamMapping(ctx, OIDCMappingRequest{
		Team:  team.UUID,
		Group: group.UUID,
	})
	require.NoError(t, err)

	// Check presence
	{
		teams, err := client.OIDC.GetAllTeamsOf(ctx, group)
		require.NoError(t, err)
		require.Equal(t, len(teams), 1)
		require.Equal(t, teams[0].UUID, team.UUID)
	}

	// Delete using mapping ID
	{
		err := client.OIDC.RemoveTeamMapping(ctx, mapping.UUID)
		require.NoError(t, err)
	}

	// Check absence
	{
		teams, err := client.OIDC.GetAllTeamsOf(ctx, group)
		require.NoError(t, err)
		require.Empty(t, teams)
	}

	// Add Mapping
	_, err = client.OIDC.AddTeamMapping(ctx, OIDCMappingRequest{
		Team:  team.UUID,
		Group: group.UUID,
	})
	require.NoError(t, err)

	// Check presence
	{
		teams, err := client.OIDC.GetAllTeamsOf(ctx, group)
		require.NoError(t, err)
		require.Equal(t, len(teams), 1)
		require.Equal(t, teams[0].UUID, team.UUID)
	}

	// Delete using Team ID, Group ID
	{
		err := client.OIDC.RemoveTeamMapping2(ctx, group.UUID, team.UUID)
		require.NoError(t, err)
	}

	// Check absence
	{
		teams, err := client.OIDC.GetAllTeamsOf(ctx, group)
		require.NoError(t, err)
		require.Empty(t, teams)
	}
}

func TestOIDCUsers(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	// Check absence
	{
		users, err := client.OIDC.GetAllUsers(ctx)
		require.NoError(t, err)
		require.Empty(t, users.Items)
		require.Zero(t, users.TotalCount)
	}

	// Create User
	user, err := client.OIDC.CreateUser(ctx, OIDCUser{
		Username: "Username",
	})
	require.NoError(t, err)
	fmt.Printf("%+v", user)
	require.Equal(t, user.Username, "Username")

	// Check presence
	{
		users, err := client.OIDC.GetAllUsers(ctx)
		require.NoError(t, err)
		require.Equal(t, users.TotalCount, 1)
		require.Equal(t, users.Items[0], user)
	}

	// Delete User
	{
		err := client.OIDC.DeleteUser(ctx, user)
		require.NoError(t, err)
	}

	// Check absence
	{
		users, err := client.OIDC.GetAllUsers(ctx)
		require.NoError(t, err)
		require.Empty(t, users.Items)
		require.Zero(t, users.TotalCount)
	}
}
