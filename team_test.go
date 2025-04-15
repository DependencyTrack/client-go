package dtrack

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAPIKey_v4_12(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		Version: "4.12.7",
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	team, err := client.Team.Create(context.Background(), Team{
		Name: "GenerateAPIKey_v4_12",
	})
	require.NoError(t, err)

	key, err := client.Team.GenerateAPIKey(context.Background(), team.UUID)
	require.NoError(t, err)

	keys, err := client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Equal(t, len(keys), 1)
	require.Equal(t, keys[0].Key, key.Key)
	require.Equal(t, keys[0].MaskedKey, key.MaskedKey)
	require.Equal(t, keys[0].Comment, "")
	require.Equal(t, key.Comment, "")
	require.Equal(t, keys[0].Created, key.Created)
	require.Equal(t, len(key.Key), 36)
	require.Equal(t, len(key.MaskedKey), 36)
}

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
	require.Equal(t, keys[0].PublicId, key.PublicId)
	require.Equal(t, keys[0].Comment, key.Comment)
	require.Equal(t, keys[0].Created, key.Created)
	require.Equal(t, keys[0].MaskedKey, key.MaskedKey)
	require.Equal(t, keys[0].Key, "")
	require.Equal(t, keys[0].Legacy, false)
	require.Equal(t, len(keys[0].PublicId), 8)
	require.Equal(t, keys[0].MaskedKey, "odt_"+key.PublicId+strings.Repeat("*", 32))
}

func TestDeleteAPIKey_v4_12(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		Version: "4.12.7",
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	team, err := client.Team.Create(context.Background(), Team{
		Name: "DeleteAPIKey_v4_12",
	})
	require.NoError(t, err)

	key, err := client.Team.GenerateAPIKey(context.Background(), team.UUID)
	require.NoError(t, err)

	keys, err := client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Equal(t, len(keys), 1)

	err = client.Team.DeleteAPIKey(context.Background(), key.Key)
	require.NoError(t, err)

	keys, err = client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Empty(t, keys)
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

	keys, err := client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Equal(t, len(keys), 1)

	err = client.Team.DeleteAPIKey(context.Background(), key.PublicId)
	require.NoError(t, err)

	keys, err = client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Empty(t, keys)
}

func TestUpdateAPIKeyComment_v4_12(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		Version: "4.12.7",
		APIPermissions: []string{
			PermissionAccessManagement,
		},
	})

	team, err := client.Team.Create(context.Background(), Team{
		Name: "UpdateAPIKeyComment_v4_12",
	})
	require.NoError(t, err)

	key, err := client.Team.GenerateAPIKey(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Equal(t, key.Comment, "")

	comment, err := client.Team.UpdateAPIKeyComment(context.Background(), key.Key, "test-comment")
	require.NoError(t, err)
	require.Equal(t, comment, "test-comment")

	keys, err := client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Equal(t, len(keys), 1)
	require.Equal(t, keys[0].Key, key.Key)
	require.Equal(t, keys[0].Comment, "test-comment")
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
	require.Equal(t, key.Comment, "")

	comment, err := client.Team.UpdateAPIKeyComment(context.Background(), key.PublicId, "test-comment")
	require.NoError(t, err)
	require.Equal(t, comment, "test-comment")

	keys, err := client.Team.GetAPIKeys(context.Background(), team.UUID)
	require.NoError(t, err)
	require.Equal(t, len(keys), 1)
	require.Equal(t, keys[0].PublicId, key.PublicId)
	require.Equal(t, keys[0].Comment, "test-comment")
}
