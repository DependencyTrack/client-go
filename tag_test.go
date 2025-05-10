package dtrack

import (
	"context"
	"testing"

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
