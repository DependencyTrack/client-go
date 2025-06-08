package dtrack

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestComponentLifecycle(t *testing.T) {
	po := PageOptions{PageSize: 10}
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionPortfolioManagement,
			PermissionSystemConfiguration,
			PermissionViewPortfolio,
		},
	})

	project, err := client.Project.Create(context.Background(), Project{
		Name:    "TestComponentLifecycleProject",
		Version: "1.0.0",
	})
	require.NoError(t, err)
	require.Equal(t, project.Name, "TestComponentLifecycleProject")

	// Check absence
	{
		components, err := client.Component.GetAll(context.Background(), project.UUID, po, ComponentFilterOptions{})
		require.NoError(t, err)
		require.Equal(t, components.TotalCount, 0)
		require.Empty(t, components.Items)
	}

	// Create Component
	component, err := client.Component.Create(context.Background(), project.UUID, Component{
		Name:       "Component-Name",
		Version:    "1.2.3",
		Classifier: "APPLICATION",
	})
	require.NoError(t, err)
	require.Equal(t, component.Name, "Component-Name")

	// Check presence
	{
		components, err := client.Component.GetAll(context.Background(), project.UUID, po, ComponentFilterOptions{})
		require.NoError(t, err)
		require.Equal(t, components.TotalCount, 1)
		require.Equal(t, len(components.Items), 1)
		require.Equal(t, components.Items[0].UUID, component.UUID)
	}

	// Update component
	{
		component.Name = component.Name + "-With-Change"
		newComponent, err := client.Component.Update(context.Background(), component)
		require.NoError(t, err)
		require.Equal(t, newComponent.UUID, component.UUID)
		require.Equal(t, newComponent.Name, component.Name)
	}

	// Check values
	{
		singleComponent, err := client.Component.Get(context.Background(), component.UUID)
		require.NoError(t, err)
		require.Equal(t, singleComponent.UUID, component.UUID)
		require.Equal(t, singleComponent.Name, "Component-Name-With-Change")
	}

	// Delete
	{
		err := client.Component.Delete(context.Background(), component.UUID)
		require.NoError(t, err)
	}

	// Check absence
	{
		components, err := client.Component.GetAll(context.Background(), project.UUID, po, ComponentFilterOptions{})
		require.NoError(t, err)
		require.Equal(t, components.TotalCount, 0)
		require.Empty(t, components.Items)
	}
}

func TestComponentPropertyLifecycle(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionPortfolioManagement,
			PermissionViewPortfolio,
		},
	})

	project, err := client.Project.Create(context.Background(), Project{
		Name:    "TestComponentPropertyLifecycleProject",
		Version: "1.0.0",
	})
	require.NoError(t, err)
	require.Equal(t, project.Name, "TestComponentPropertyLifecycleProject")

	component, err := client.Component.Create(context.Background(), project.UUID, Component{
		Name:       "Component-Name",
		Version:    "1.2.3",
		Classifier: "APPLICATION",
	})
	require.NoError(t, err)
	require.Equal(t, component.Name, "Component-Name")

	// Check absence
	{
		components, err := client.Component.GetProperties(context.Background(), component.UUID)
		require.NoError(t, err)
		require.Empty(t, components)
	}

	// Create
	property, err := client.Component.CreateProperty(context.Background(), component.UUID, ComponentProperty{
		Group:       "Property-Group",
		Name:        "Property-Name",
		Value:       "Property-Value",
		Type:        "STRING",
		Description: "Property-Description",
	})
	require.NoError(t, err)
	require.Equal(t, property.Group, "Property-Group")
	require.Equal(t, property.Name, "Property-Name")
	require.Equal(t, property.Value, "Property-Value")
	require.Equal(t, property.Type, "STRING")
	require.Equal(t, property.Description, "Property-Description")

	// Check presence
	{
		properties, err := client.Component.GetProperties(context.Background(), component.UUID)
		require.NoError(t, err)
		require.Equal(t, len(properties), 1)
		require.Equal(t, properties[0], property)
	}

	// Delete
	{
		err := client.Component.DeleteProperty(context.Background(), component.UUID, property.UUID)
		require.NoError(t, err)
	}

	// Check absence
	{
		components, err := client.Component.GetProperties(context.Background(), component.UUID)
		require.NoError(t, err)
		require.Empty(t, components)
	}
}

func TestComponentLocate(t *testing.T) {
	po := PageOptions{PageSize: 10}
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionPortfolioManagement,
			PermissionViewPortfolio,
		},
	})

	project, err := client.Project.Create(context.Background(), Project{
		Name:    "TestLocateComponentProject",
		Version: "1.0.0",
	})
	require.NoError(t, err)
	require.Equal(t, project.Name, "TestLocateComponentProject")

	component, err := client.Component.Create(context.Background(), project.UUID, Component{
		Name:       "Component-Name",
		Version:    "1.2.3",
		Classifier: "APPLICATION",
		MD5:        "0123456789abcdef0123456789abcdef",
	})
	require.NoError(t, err)
	require.Equal(t, component.Name, "Component-Name")
	require.Equal(t, component.MD5, "0123456789abcdef0123456789abcdef")

	// Hash
	{
		components, err := client.Component.GetByHash(context.Background(), "0123456789abcdef0123456789abcdef", po, SortOptions{})
		require.NoError(t, err)
		require.Equal(t, components.TotalCount, 1)

		require.Nil(t, components.Items[0].Project.Tags)
		require.Nil(t, components.Items[0].Project.Properties)
		components.Items[0].Project.Tags = []Tag{}
		components.Items[0].Project.Properties = []ProjectProperty{}

		require.Equal(t, components.Items[0], component)
	}

	// Identity
	{
		components, err := client.Component.GetByIdentity(context.Background(), po, SortOptions{}, ComponentIdentityQueryOptions{
			Name: "Component-Name",
		})
		require.NoError(t, err)
		require.Equal(t, components.TotalCount, 1)

		require.Nil(t, components.Items[0].Project.Tags)
		require.Nil(t, components.Items[0].Project.Properties)
		components.Items[0].Project.Tags = []Tag{}
		components.Items[0].Project.Properties = []ProjectProperty{}

		require.Equal(t, components.Items[0], component)
	}
}

func TestComponentIdentifyInternal(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionSystemConfiguration,
		},
	})
	err := client.Component.IdentifyInternal(context.Background())
	require.NoError(t, err)
}
