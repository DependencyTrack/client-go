package dtrack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectPropertyService_GetAll(t *testing.T) {
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

	po := PageOptions{PageSize: 10}

	// Baseline: no properties yet
	properties, err := client.ProjectProperty.GetAll(context.Background(), project.UUID, po)
	require.NoError(t, err)
	// The api does not return TotalCount as it seems not to be paginated (although we can supply page options)
	//require.Equal(t, 0, properties.TotalCount)
	require.Empty(t, properties.Items)

	// Create a property
	_, err = client.ProjectProperty.Create(context.Background(), project.UUID, ProjectProperty{
		Group: "test-group",
		Name:  "test-property",
		Value: "test-value",
		Type:  "STRING",
	})
	require.NoError(t, err)

	// Verify it appears in GetAll
	properties, err = client.ProjectProperty.GetAll(context.Background(), project.UUID, po)
	require.NoError(t, err)
	// The api does not return TotalCount as it seems not to be paginated (although we can supply page options)
	//require.Equal(t, 1, properties.TotalCount)
	require.Len(t, properties.Items, 1)
	require.Equal(t, "test-group", properties.Items[0].Group)
	require.Equal(t, "test-property", properties.Items[0].Name)
	require.Equal(t, "test-value", properties.Items[0].Value)
	require.Equal(t, "STRING", properties.Items[0].Type)
}

func TestProjectPropertyService_Create(t *testing.T) {
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

	property, err := client.ProjectProperty.Create(context.Background(), project.UUID, ProjectProperty{
		Group:       "test-group",
		Name:        "test-property",
		Value:       "test-value",
		Type:        "STRING",
		Description: "a test property",
	})
	require.NoError(t, err)
	require.Equal(t, "test-group", property.Group)
	require.Equal(t, "test-property", property.Name)
	require.Equal(t, "test-value", property.Value)
	require.Equal(t, "STRING", property.Type)
	require.Equal(t, "a test property", property.Description)
}

func TestProjectPropertyService_Update(t *testing.T) {
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

	// Create a property first
	_, err = client.ProjectProperty.Create(context.Background(), project.UUID, ProjectProperty{
		Group: "test-group",
		Name:  "test-property",
		Value: "original-value",
		Type:  "STRING",
	})
	require.NoError(t, err)

	// Update the property
	updated, err := client.ProjectProperty.Update(context.Background(), project.UUID, ProjectProperty{
		Group: "test-group",
		Name:  "test-property",
		Value: "updated-value",
		Type:  "STRING",
	})
	require.NoError(t, err)
	require.Equal(t, "test-group", updated.Group)
	require.Equal(t, "test-property", updated.Name)
	require.Equal(t, "updated-value", updated.Value)
	require.Equal(t, "STRING", updated.Type)
}

func TestProjectPropertyService_Delete(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionPortfolioManagement,
		},
	})

	po := PageOptions{PageSize: 10}

	project, err := client.Project.Create(context.Background(), Project{
		Name:    "acme-app",
		Version: "1.0.0",
	})
	require.NoError(t, err)

	// Create a property
	_, err = client.ProjectProperty.Create(context.Background(), project.UUID, ProjectProperty{
		Group: "test-group",
		Name:  "test-property",
		Value: "test-value",
		Type:  "STRING",
	})
	require.NoError(t, err)

	// Verify presence
	properties, err := client.ProjectProperty.GetAll(context.Background(), project.UUID, po)
	require.NoError(t, err)
	require.Len(t, properties.Items, 1)

	// Delete the property
	err = client.ProjectProperty.Delete(context.Background(), project.UUID, "test-group", "test-property")
	require.NoError(t, err)

	// Verify absence
	properties, err = client.ProjectProperty.GetAll(context.Background(), project.UUID, po)
	require.NoError(t, err)
	require.Empty(t, properties.Items)
}
