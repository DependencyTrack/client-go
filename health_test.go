package dtrack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{})

	health, err := client.Health.Get(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, health)

	require.Equal(t, "UP", health.Status)
	require.Equal(t, 1, len(health.Checks))
	require.Equal(t, "database", health.Checks[0].Name)
	require.Equal(t, "UP", health.Checks[0].Status)
	require.NotNil(t, health.Checks[0].Data)
	require.Equal(t, "UP", health.Checks[0].Data.(map[string]any)["nontx_connection_pool"])
	require.Equal(t, "UP", health.Checks[0].Data.(map[string]any)["tx_connection_pool"])
}
