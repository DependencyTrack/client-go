package dtrack

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateAPIKey(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	team, err := client.Team.Create(context.Background(), Team{
		Name: "GenerateAPIKey",
	})
	require.NoError(t, err)

	key, err := client.Team.GenerateAPIKey(context.Background(), team.UUID)
	require.NoError(t, err)

	keys, err := client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Equal(t, len(keys), 1)
	require.Equal(t, keys[0].Key, key)
}

func TestDeleteAPIKey(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	team, err := client.Team.Create(context.Background(), Team{
		Name: "DeleteAPIKey",
	})
	require.NoError(t, err)

	key, err := client.Team.GenerateAPIKey(context.Background(), team.UUID)
	require.NoError(t, err)

	err = client.Team.DeleteAPIKey(context.Background(), key)
	require.NoError(t, err)

	keys, err := client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Empty(t, keys)
}

func TestUpdateAPIKeyComment(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	team, err := client.Team.Create(context.Background(), Team{
		Name: "UpdateAPIKeyComment",
	})
	require.NoError(t, err)

	key, err := client.Team.GenerateAPIKey(context.Background(), team.UUID)
	require.NoError(t, err)

	comment, err := client.Team.UpdateAPIKeyComment(context.Background(), key, "test-comment")
	require.NoError(t, err)
	require.Equal(t, comment, "test-comment")

	keys, err := client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Equal(t, len(keys), 1)
	require.Equal(t, keys[0].Key, key)
	require.Equal(t, keys[0].Comment, "test-comment")
}
