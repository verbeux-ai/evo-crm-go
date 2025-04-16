package evo_crm_test

import (
	"context"
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	evo_crm "github.com.br/verbeux-ai/evo-crm-go"
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

func TestListPanelCardsFilter(t *testing.T) {
	ctx := context.Background()

	panelIdForTest := os.Getenv("EVO_CRM_PANEL_ID_FOR_TEST")
	require.NotEmpty(t, panelIdForTest, "Environment variable EVO_CRM_PANEL_ID_FOR_TEST must be set")

	filter := evo_crm.ListPanelCardsFilterRequest{
		PanelID:    panelIdForTest,
		PageSize:   5,
		PageNumber: 1,
		ContactID:  "39df1926-092a-4140-8f2c-c75db6e87e8e",
	}

	listResponse, err := client.ListPanelCardsFilter(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, listResponse)
	require.NotNil(t, listResponse.Items)
	require.GreaterOrEqual(t, listResponse.TotalItems, 0)
	require.GreaterOrEqual(t, len(listResponse.Items), 0)
}

func TestGetPanelByID(t *testing.T) {
	ctx := context.Background()

	panelIdForTest := os.Getenv("EVO_CRM_PANEL_ID_FOR_TEST")
	require.NotEmpty(t, panelIdForTest, "Environment variable EVO_CRM_PANEL_ID_FOR_TEST must be set")

	includeDetails := []evo_crm.PanelIncludeDetails{
		evo_crm.IncludeSteps,
		evo_crm.IncludeTags,
	}

	panel, err := client.GetPanelByID(ctx, panelIdForTest, includeDetails)

	require.NoError(t, err)
	require.NotNil(t, panel)
	require.Equal(t, panelIdForTest, panel.ID)
	require.NotEmpty(t, panel.Title)

	if slices.Contains(includeDetails, evo_crm.IncludeSteps) {
		require.NotNil(t, panel.Steps)
	}
	if slices.Contains(includeDetails, evo_crm.IncludeTags) {
		require.NotNil(t, panel.Tags)
	}
}

func TestCreatePanelCard(t *testing.T) {
	ctx := context.Background()

	stepIdForTest := os.Getenv("EVO_CRM_STEP_ID_FOR_TEST")
	require.NotEmpty(t, stepIdForTest, "Environment variable EVO_CRM_STEP_ID_FOR_TEST must be set")

	uniqueSuffix := time.Now().Format("20060102150405")
	cardData := evo_crm.CreatePanelCardRequest{
		StepID:         stepIdForTest,
		Title:          fmt.Sprintf("Test Card %s", uniqueSuffix),
		Description:    "Card created via integration test",
		MonetaryAmount: 123.45,
	}

	createdCard, err := client.CreatePanelCard(ctx, cardData)

	require.NoError(t, err)
	require.NotNil(t, createdCard)
	require.NotEmpty(t, createdCard.ID)
	require.Equal(t, cardData.StepID, createdCard.StepID)
	require.Equal(t, cardData.Title, createdCard.Title)
	require.Equal(t, cardData.Description, createdCard.Description)
	require.Equal(t, cardData.MonetaryAmount, createdCard.MonetaryAmount)
}
