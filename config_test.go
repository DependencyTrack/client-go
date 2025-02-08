package dtrack

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

const TEST_URL = "https://localhost"

func TestGetConfigProperty(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})

	property, err := client.Config.Get(context.Background(), "general", "base.url")
	require.NoError(t, err)
	require.Equal(t, property.GroupName, "general")
	require.Equal(t, property.Name, "base.url")
	require.Equal(t, property.Type, "URL")
	require.Equal(t, property.Description, "URL used to construct links back to Dependency-Track from external systems")
}

func TestUpdateConfigProperty(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})

	property, err := client.Config.Get(context.Background(), "general", "base.url")
	require.NoError(t, err)
	require.Empty(t, property.Value)

	property.Value = TEST_URL

	property, err = client.Config.Update(context.Background(), property)
	require.NoError(t, err)
	require.Equal(t, property.GroupName, "general")
	require.Equal(t, property.Name, "base.url")
	require.Equal(t, property.Type, "URL")
	require.Equal(t, property.Description, "URL used to construct links back to Dependency-Track from external systems")

	require.Equal(t, property.Value, TEST_URL)
}

func TestUpdateAllConfigProperty(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})

	baseUrl, err := client.Config.Get(context.Background(), "general", "base.url")
	require.NoError(t, err)
	badgeEnabled, err := client.Config.Get(context.Background(), "general", "badge.enabled")
	require.NoError(t, err)
	defaultLocale, err := client.Config.Get(context.Background(), "general", "default.locale")
	require.NoError(t, err)

	require.Empty(t, baseUrl.Value)
	require.Equal(t, badgeEnabled.Value, "false")
	require.Empty(t, defaultLocale.Value)

	baseUrl.Value = TEST_URL
	badgeEnabled.Value = "true"
	defaultLocale.Value = "de"

	cps, err := client.Config.UpdateAll(context.Background(), []ConfigProperty{baseUrl, badgeEnabled, defaultLocale})
	require.NoError(t, err)
	require.Equal(t, len(cps), 3)
	require.Equal(t, cps[0], baseUrl)
	require.Equal(t, cps[1], badgeEnabled)
	require.Equal(t, cps[2], defaultLocale)
}

func TestUnsetConfigProperty(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})

	baseUrl, err := client.Config.Get(context.Background(), "general", "base.url")
	require.NoError(t, err)
	baseUrl.Value = TEST_URL
	baseUrl, err = client.Config.Update(context.Background(), baseUrl)
	require.NoError(t, err)
	require.Equal(t, baseUrl.Value, TEST_URL)

	baseUrl.Value = ""
	baseUrl, err = client.Config.Update(context.Background(), baseUrl)
	require.NoError(t, err)
	require.Empty(t, baseUrl.Value)
}
