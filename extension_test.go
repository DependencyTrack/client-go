package dtrack

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtensionService(t *testing.T) {
	ctx := context.Background()
	client := setUpContainer(t, testContainerOptions{
		Version:     "5.2.0",
		UsePostgres: true,
		APIPermissions: []string{
			PermissionAccessManagement,
			PermissionSystemConfiguration,
		},
	})

	requireAPIError := func(t *testing.T, err error, statusCode int) {
		t.Helper()
		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, statusCode, apiErr.StatusCode)
	}

	team, err := client.Team.Create(ctx, Team{Name: "extension-no-permissions"})
	require.NoError(t, err)
	apiKey, err := client.Team.GenerateAPIKey(ctx, team.UUID)
	require.NoError(t, err)
	unprivileged, err := NewClient(client.baseURL.String(), WithAPIKey(apiKey.Key))
	require.NoError(t, err)

	t.Run("GetAllExtensionPoints", func(t *testing.T) {
		points, err := client.Extension.GetAllExtensionPoints(ctx)
		require.NoError(t, err)
		require.Contains(t, points, ExtensionPoint{Name: "vuln-analyzer"})

		_, err = unprivileged.Extension.GetAllExtensionPoints(ctx)
		requireAPIError(t, err, http.StatusForbidden)
	})

	t.Run("GetAllExtensions", func(t *testing.T) {
		extensions, err := client.Extension.GetAllExtensions(ctx, "vuln-analyzer")
		require.NoError(t, err)
		require.Contains(t, extensions, Extension{
			Name: "oss-index", DisplayName: "Sonatype OSS Index", Configurable: true, Testable: true,
		})

		_, err = client.Extension.GetAllExtensions(ctx, "does-not-exist")
		requireAPIError(t, err, http.StatusNotFound)
	})

	t.Run("GetConfigSchema", func(t *testing.T) {
		schema, err := client.Extension.GetConfigSchema(ctx, "vuln-analyzer", "internal")
		require.NoError(t, err)
		var parsed struct {
			Schema string `json:"$schema"`
		}
		require.NoError(t, json.Unmarshal(schema, &parsed))
		require.NotEmpty(t, parsed.Schema)

		schema, err = client.Extension.GetConfigSchema(ctx, "notification-publisher", "console")
		require.NoError(t, err)
		require.Nil(t, schema)
	})

	t.Run("GetConfig", func(t *testing.T) {
		config, err := client.Extension.GetConfig(ctx, "vuln-analyzer", "internal")
		require.NoError(t, err)
		require.JSONEq(t, `{"enabled":true}`, string(config))

		_, err = client.Extension.GetConfig(ctx, "vuln-analyzer", "does-not-exist")
		requireAPIError(t, err, http.StatusNotFound)
	})

	// The update endpoint has no JSON success body; these calls fail with HTTP 406
	// if the client does not accept application/problem+json.
	t.Run("UpdateConfig", func(t *testing.T) {
		original, err := client.Extension.GetConfig(ctx, "vuln-analyzer", "internal")
		require.NoError(t, err)

		// Unchanged configuration (HTTP 304).
		require.NoError(t, client.Extension.UpdateConfig(ctx, "vuln-analyzer", "internal", original))

		require.NoError(t, client.Extension.UpdateConfig(ctx, "vuln-analyzer", "internal", json.RawMessage(`{"enabled":false}`)))
		updated, err := client.Extension.GetConfig(ctx, "vuln-analyzer", "internal")
		require.NoError(t, err)
		require.JSONEq(t, `{"enabled":false}`, string(updated))

		err = client.Extension.UpdateConfig(ctx, "vuln-analyzer", "internal", json.RawMessage(`{"enabled":"invalid"}`))
		requireAPIError(t, err, http.StatusBadRequest)
		require.ErrorContains(t, err, `"instance_location":"/enabled"`)

		err = unprivileged.Extension.UpdateConfig(ctx, "vuln-analyzer", "internal", original)
		requireAPIError(t, err, http.StatusForbidden)
	})

	t.Run("Test", func(t *testing.T) {
		config, err := client.Extension.GetConfig(ctx, "vuln-analyzer", "oss-index")
		require.NoError(t, err)

		result, err := client.Extension.Test(ctx, "vuln-analyzer", "oss-index", config)
		require.NoError(t, err)
		require.NotEmpty(t, result.Checks)
		for _, check := range result.Checks {
			require.NotEmpty(t, check.Name)
			require.Contains(t, []ExtensionTestCheckStatus{
				ExtensionTestCheckStatusPassed, ExtensionTestCheckStatusFailed, ExtensionTestCheckStatusSkipped,
			}, check.Status)
		}

		// Configurable extensions require a configuration.
		_, err = client.Extension.Test(ctx, "vuln-analyzer", "oss-index", nil)
		requireAPIError(t, err, http.StatusBadRequest)
	})
}

func TestExtensionTestResult_Failed(t *testing.T) {
	require.False(t, ExtensionTestResult{Checks: []ExtensionTestCheck{
		{Status: ExtensionTestCheckStatusPassed}, {Status: ExtensionTestCheckStatusSkipped},
	}}.Failed())
	require.True(t, ExtensionTestResult{Checks: []ExtensionTestCheck{
		{Status: ExtensionTestCheckStatusPassed}, {Status: ExtensionTestCheckStatusFailed},
	}}.Failed())
}
