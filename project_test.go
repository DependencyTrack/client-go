package dtrack

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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

func TestProjectService_CreateWithCollection(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionPortfolioManagement,
			PermissionTagManagement,
		},
	})

	collectionDirectChildrenProject, err := client.Project.Create(context.Background(), Project{
		Name:            "acme-app",
		Version:         "1.0.0",
		Active:          true,
		CollectionLogic: &CollectionLogicAggregateDirectChildren,
	})
	require.NoError(t, err)

	tag := "weewoo"
	err = client.Tag.Create(context.Background(), []string{tag})
	require.NoError(t, err)

	collectionTags, err := client.Project.Create(context.Background(), Project{
		Name:            "acme-app-2",
		Version:         "1.0.0",
		Active:          true,
		CollectionLogic: &CollectionLogicAggregateDirectChildrenWithTag,
		CollectionTag:   &Tag{Name: tag},
	})
	require.NoError(t, err)

	cases := []struct {
		name          string
		childProjects []Project
		projectName   string
		parent        uuid.UUID
		expErr        error
	}{
		{
			name: "aggregate_collection_logic_direct_children_single_child",
			childProjects: []Project{
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					ParentRef: &ParentRef{
						UUID: collectionDirectChildrenProject.UUID,
					},
				},
			},
			parent: collectionDirectChildrenProject.UUID,
		},
		{
			name: "aggregate_collection_logic_direct_children_multiple_children",
			childProjects: []Project{
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					ParentRef: &ParentRef{
						UUID: collectionDirectChildrenProject.UUID,
					},
				},
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					ParentRef: &ParentRef{
						UUID: collectionDirectChildrenProject.UUID,
					},
				},
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					ParentRef: &ParentRef{
						UUID: collectionDirectChildrenProject.UUID,
					},
				},
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					ParentRef: &ParentRef{
						UUID: collectionDirectChildrenProject.UUID,
					},
				},
			},
			parent: collectionDirectChildrenProject.UUID,
		},
		{
			name: "aggregate_collection_logic_tags_single_child",
			childProjects: []Project{
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					Tags: []Tag{{Name: tag}},
					ParentRef: &ParentRef{
						UUID: collectionTags.UUID,
					},
				},
			},
			parent: collectionTags.UUID,
		},
		{
			name: "aggregate_collection_logic_tags_ multiple_children",
			childProjects: []Project{
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					Tags: []Tag{{Name: tag}},
					ParentRef: &ParentRef{
						UUID: collectionTags.UUID,
					},
				},
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					Tags: []Tag{{Name: tag}},
					ParentRef: &ParentRef{
						UUID: collectionTags.UUID,
					},
				},
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					Tags: []Tag{{Name: tag}},
					ParentRef: &ParentRef{
						UUID: collectionTags.UUID,
					},
				},
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					Tags: []Tag{{Name: tag}},
					ParentRef: &ParentRef{
						UUID: collectionTags.UUID,
					},
				},
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					Tags: []Tag{{Name: tag}},
					ParentRef: &ParentRef{
						UUID: collectionTags.UUID,
					},
				},
				{
					Name: fmt.Sprintf("child-%d", rand.Intn(100000)), //nolint:gosec
					Tags: []Tag{{Name: tag}},
					ParentRef: &ParentRef{
						UUID: collectionTags.UUID,
					},
				},
			},
			parent: collectionTags.UUID,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			for _, project := range tc.childProjects {
				childProject, err := client.Project.Create(context.Background(), project)
				require.Equal(t, tc.expErr, err)
				require.NotEmpty(t, childProject.ParentRef)
				require.Equal(t, tc.parent, childProject.ParentRef.UUID)
			}
		})
	}
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

func TestProjectService_Latest(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		Version: "4.12.7",
		APIPermissions: []string{
			PermissionPortfolioManagement,
			PermissionViewPortfolio,
		},
	})
	name := "acme-app"
	project, err := client.Project.Create(context.Background(), Project{
		Name:     name,
		Version:  "1.0.0",
		IsLatest: OptionalBoolOf(true),
	})
	require.NoError(t, err)
	latest, err := client.Project.Latest(context.Background(), name)

	require.NoError(t, err)
	require.Equal(t, project.Version, latest.Version)

	token, err := client.Project.Clone(context.Background(), ProjectCloneRequest{
		ProjectUUID:     project.UUID,
		Version:         "2.0.0",
		MakeCloneLatest: OptionalBoolOf(true),
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	for {
		processing, err := client.Event.IsBeingProcessed(context.Background(), token)
		require.NoError(t, err)
		if !processing {
			break
		}
	}

	latest, err = client.Project.Latest(context.Background(), name)

	require.NoError(t, err)
	require.Equal(t, "2.0.0", latest.Version)
}
