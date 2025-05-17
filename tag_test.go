package dtrack

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	po := PageOptions{
		PageSize: 10,
	}
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionTagManagement,
			PermissionViewPortfolio,
		},
	})

	tags, err := client.Tag.GetAll(context.Background(), po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, tags.TotalCount, 0)
	require.Empty(t, tags.Items)

	err = client.Tag.Create(context.Background(), []string{"test_foo", "test_bar"})
	require.NoError(t, err)

	tags, err = client.Tag.GetAll(context.Background(), po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, tags.TotalCount, 2)
	require.Equal(t, tags.Items[0].Name, "test_bar")
	require.Equal(t, tags.Items[1].Name, "test_foo")
	require.Equal(t, tags.Items[0].NotificationRuleCount, 0)
	require.Equal(t, tags.Items[0].PolicyCount, 0)
	require.Equal(t, tags.Items[0].ProjectCount, 0)

	err = client.Tag.Delete(context.Background(), []string{"test_bar", "test_foo"})
	require.NoError(t, err)

	tags, err = client.Tag.GetAll(context.Background(), po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, tags.TotalCount, 0)
	require.Empty(t, tags.Items)
}

func TestTagCounts(t *testing.T) {
	po := PageOptions{PageSize: 10}
	tag := Tag{Name: "test_tag"}
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionTagManagement,
			PermissionViewPortfolio,
			PermissionPortfolioManagement,
			PermissionPolicyManagement,
		},
	})

	tags, err := client.Tag.GetAll(context.Background(), po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, tags.TotalCount, 0)
	require.Empty(t, tags.Items)

	err = client.Tag.Create(context.Background(), []string{tag.Name})
	require.NoError(t, err)

	// Setup Project
	project, err := client.Project.Create(context.Background(), Project{
		Name: "test_project",
		Tags: []Tag{tag},
	})
	require.NoError(t, err)
	require.Equal(t, len(project.Tags), 1)

	// Setup Policy
	policy, err := client.Policy.Create(context.Background(), Policy{
		Name:           "test_policy",
		Operator:       PolicyOperatorAll,
		ViolationState: PolicyViolationStateFail,
	})
	require.NoError(t, err)
	_, err = client.Policy.AddTag(context.Background(), policy.UUID, tag.Name)
	require.NoError(t, err)

	// Count assertions
	tags, err = client.Tag.GetAll(context.Background(), po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, tags.TotalCount, 1)
	require.Equal(t, tags.Items[0].NotificationRuleCount, int64(0))
	require.Equal(t, tags.Items[0].PolicyCount, int64(1))
	require.Equal(t, tags.Items[0].ProjectCount, int64(1))
}

func TestTagProject(t *testing.T) {
	po := PageOptions{PageSize: 10}
	tag := Tag{Name: "test_tag_project"}
	projectName := "test_project"
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionTagManagement,
			PermissionPortfolioManagement,
			PermissionViewPortfolio,
		},
	})

	// Setup
	err := client.Tag.Create(context.Background(), []string{tag.Name})
	require.NoError(t, err)

	project, err := client.Project.Create(context.Background(), Project{Name: projectName})
	require.NoError(t, err)
	require.Equal(t, project.Name, projectName)

	// Baseline
	projects, err := client.Tag.GetProjects(context.Background(), tag.Name, po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, projects.TotalCount, 0)
	require.Empty(t, projects.Items)

	// Tag
	err = client.Tag.TagProjects(context.Background(), tag.Name, []uuid.UUID{project.UUID})
	require.NoError(t, err)

	// Check Presence
	projects, err = client.Tag.GetProjects(context.Background(), tag.Name, po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, projects.TotalCount, 1)
	require.Equal(t, projects.Items[0].UUID, project.UUID)
	require.Equal(t, projects.Items[0].Name, project.Name)
	require.Equal(t, projects.Items[0].Version, project.Version)

	// Untag
	err = client.Tag.UntagProjects(context.Background(), tag.Name, []uuid.UUID{project.UUID})
	require.NoError(t, err)

	// Check Absence
	projects, err = client.Tag.GetProjects(context.Background(), tag.Name, po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, projects.TotalCount, 0)
	require.Empty(t, projects.Items)
}
