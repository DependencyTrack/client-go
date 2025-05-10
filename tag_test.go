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

	err = client.Tag.Delete(context.Background(), []string{"test_bar", "test_foo"})
	require.NoError(t, err)

	tags, err = client.Tag.GetAll(context.Background(), po, SortOptions{})
	require.NoError(t, err)
	require.Equal(t, tags.TotalCount, 0)
	require.Empty(t, tags.Items)
}
