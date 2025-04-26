package dtrack

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAclService(t *testing.T) {
	pageOptions := PageOptions{
		PageNumber: 1,
		PageSize:   10,
	}

	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionAccessManagement,
			PermissionPortfolioManagement,
		},
	})

	project, err := client.Project.Create(context.Background(), Project{
		Name: "Test_Acl_Project",
	})
	require.NoError(t, err)

	team, err := client.Team.Create(context.Background(), Team{
		Name: "Test_Acl_Team",
	})
	require.NoError(t, err)

	mappings, err := client.Acl.GetAllProjects(context.Background(), team.UUID, pageOptions)
	require.NoError(t, err)
	require.Equal(t, mappings.TotalCount, 0)
	require.Empty(t, mappings.Items)

	err = client.Acl.AddProjectMapping(context.Background(), AclMappingRequest{
		Team:    team.UUID,
		Project: project.UUID,
	})
	require.NoError(t, err)

	mappings, err = client.Acl.GetAllProjects(context.Background(), team.UUID, pageOptions)
	require.NoError(t, err)
	require.Equal(t, mappings.TotalCount, 1)
	require.Equal(t, len(mappings.Items), 1)
	require.Equal(t, mappings.Items[0], project)

	err = client.Acl.RemoveProjectMapping(context.Background(), team.UUID, project.UUID)
	require.NoError(t, err)

	mappings, err = client.Acl.GetAllProjects(context.Background(), team.UUID, pageOptions)
	require.NoError(t, err)
	require.Equal(t, mappings.TotalCount, 0)
	require.Empty(t, mappings.Items)
}
