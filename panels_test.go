package evo_crm_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListPanels(t *testing.T) {
	ctx := context.Background()

	listResponse, err := client.ListPanels(ctx)

	require.NoError(t, err)
	require.NotNil(t, listResponse)
	require.NotNil(t, listResponse.Items)
	require.GreaterOrEqual(t, listResponse.TotalItems, 0)
	require.GreaterOrEqual(t, len(listResponse.Items), 0)
}
