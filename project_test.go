package dtrack

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProjectService_Clone(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionPortfolioManagement,
		},
	})

	project, err := client.Project.Create(context.Background(), Project{
		Name:    "acme-app",
		Version: "1.0.0",
	})
	require.NoError(t, err)

	token, err := client.Project.Clone(context.Background(), ProjectCloneRequest{
		ProjectUUID: project.UUID,
		Version:     "2.0.0",
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestProjectService_Clone_v4_10(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		Version: "4.10.1",
		APIPermissions: []string{
			PermissionPortfolioManagement,
		},
	})

	project, err := client.Project.Create(context.Background(), Project{
		Name:    "acme-app",
		Version: "1.0.0",
	})
	require.NoError(t, err)

	token, err := client.Project.Clone(context.Background(), ProjectCloneRequest{
		ProjectUUID: project.UUID,
		Version:     "2.0.0",
	})
	require.NoError(t, err)
	require.Empty(t, token)
}
