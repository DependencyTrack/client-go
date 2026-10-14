package dtrack

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

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

// TestSetUpContainer_WithPostgres ensures that our testcontainer setup is able to start a
// Dependency-Track version >= 5 with a PostgreSQL container alongside it.
func TestSetUpContainer_WithPostgres(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		Version:     "5.2.0",
		UsePostgres: true,
	})

	about, err := client.About.Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, "5.2.0", about.Version)
	require.Equal(t, "Dependency-Track", about.Application)
}

type testContainerOptions struct {
	Version        string
	APIPermissions []string
	UsePostgres    bool
}

func setUpContainer(t *testing.T, options testContainerOptions) *Client {
	ctx := context.Background()

	version := "latest"
	if options.Version != "" {
		version = options.Version
	}

	// Latest tag points to v5, which requires a Postgres container.
	if version == "latest" {
		options.UsePostgres = true
	}

	image := fmt.Sprintf("dependencytrack/apiserver:%s", version)
	env := map[string]string{
		"JAVA_OPTIONS":                     "-Xmx1g",
		"SYSTEM_REQUIREMENT_CHECK_ENABLED": "false",
	}
	var networks []string
	var networkAliases map[string][]string
	var waitingFor wait.Strategy = wait.ForLog("Dependency-Track is ready")

	if options.UsePostgres {
		networkName := "dtrack-test-" + uuid.NewString()
		network, networkErr := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
			NetworkRequest: testcontainers.NetworkRequest{Name: networkName},
		})
		require.NoError(t, networkErr)
		t.Cleanup(func() {
			require.NoError(t, network.Remove(ctx))
		})

		postgres, postgresErr := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image: "postgres:18-alpine",
				Env: map[string]string{
					"POSTGRES_DB":       "dtrack",
					"POSTGRES_USER":     "dtrack",
					"POSTGRES_PASSWORD": "dtrack",
				},
				ExposedPorts: []string{"5432/tcp"},
				Networks:     []string{networkName},
				NetworkAliases: map[string][]string{
					networkName: {"postgres"},
				},
				WaitingFor: wait.ForListeningPort("5432/tcp"),
			},
			Started: true,
		})
		require.NoError(t, postgresErr)
		t.Cleanup(func() {
			require.NoError(t, postgres.Terminate(ctx))
		})

		networks = []string{networkName}
		networkAliases = map[string][]string{networkName: {"apiserver"}}
		env["DT_DATASOURCE_URL"] = "jdbc:postgresql://postgres:5432/dtrack"
		env["DT_DATASOURCE_USERNAME"] = "dtrack"
		env["DT_DATASOURCE_PASSWORD"] = "dtrack"
		waitingFor = wait.ForListeningPort("8080/tcp")
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:          image,
			Env:            env,
			ExposedPorts:   []string{"8080/tcp"},
			Networks:       networks,
			NetworkAliases: networkAliases,
			WaitingFor:     waitingFor,
		},
		Started: true,
	})
	if err != nil && container != nil {
		require.NoError(t, container.Terminate(ctx))
	}
	require.NoError(t, err)

	t.Cleanup(func() {
		err = container.Terminate(ctx)
		if err != nil {
			log.Fatalf("failed to terminate container: %v", err)
		}
	})

	// Select the API port explicitly instead of an arbitrary exposed port.
	apiURL, err := container.PortEndpoint(ctx, "8080/tcp", "http")
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

	client, err = NewClient(apiURL, WithAPIKey(apiKey.Key))
	require.NoError(t, err)

	return client
}
