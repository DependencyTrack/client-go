package dtrack

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAboutService_Get(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{})

	about, err := client.About.Get(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, about)

	require.NotEmpty(t, about.Timestamp)
	require.NotEmpty(t, about.Version)
	require.NotEqual(t, uuid.Nil, about.UUID)
	require.NotEqual(t, uuid.Nil, about.SystemUUID)
	require.Equal(t, "Dependency-Track", about.Application)

	require.NotEmpty(t, about.Framework.Timestamp)
	require.NotEmpty(t, about.Framework.Version)
	require.NotEqual(t, uuid.Nil, about.Framework.UUID)
	require.Equal(t, "Alpine", about.Framework.Name)
}

type testContainerOptions struct {
	Version        string
	APIPermissions []string
}

func setUpContainer(t *testing.T, options testContainerOptions) *Client {
	ctx := context.Background()

	version := "latest"
	if options.Version != "" {
		version = options.Version
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: fmt.Sprintf("dependencytrack/apiserver:%s", version),
			Env: map[string]string{
				"JAVA_OPTIONS":                     "-Xmx1g",
				"SYSTEM_REQUIREMENT_CHECK_ENABLED": "false",
			},
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForLog("Dependency-Track is ready"),
		},
		Started: true,
	})

	t.Cleanup(func() {
		err = container.Terminate(ctx)
		if err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	})
	require.NoError(t, err)

	apiURL, err := container.Endpoint(ctx, "http")
	require.NoError(t, err)

	client, err := NewClient(apiURL)
	require.NoError(t, err)

	err = client.User.ForceChangePassword(ctx, "admin", "admin", "test")
	require.NoError(t, err)

	bearerToken, err := client.User.Login(ctx, "admin", "test")
	require.NoError(t, err)

	client, err = NewClient(apiURL, WithBearerToken(bearerToken))
	require.NoError(t, err)

	team, err := client.Team.Create(ctx, Team{Name: "test"})
	require.NoError(t, err)

	for _, permissionName := range options.APIPermissions {
		_, err = client.Permission.AddPermissionToTeam(ctx, Permission{Name: permissionName}, team.UUID)
		require.NoError(t, err)
	}

	apiKey, err := client.Team.GenerateAPIKey(ctx, team.UUID)
	require.NoError(t, err)

	client, err = NewClient(apiURL, WithAPIKey(apiKey))
	require.NoError(t, err)

	return client
}
