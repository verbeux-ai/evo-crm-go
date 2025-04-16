package evo_crm_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	evo_crm "github.com/verbeux-ai/evo-crm-go"
)

func TestListContactsFilter(t *testing.T) {
	ctx := context.Background()

	filter := evo_crm.ListContactsFilterRequest{
		PageSize: 5,
		Page:     1,
	}

	contacts, err := client.ListContactsFilter(ctx, filter)

	require.NoError(t, err)
	require.NotNil(t, contacts)
}

func TestCreateContact(t *testing.T) {
	ctx := context.Background()

	uniqueSuffix := time.Now().Format("20060102150405")
	testPhone := fmt.Sprintf("+55|%s", uniqueSuffix)
	if len(testPhone) > 15 {
		testPhone = testPhone[:15]
	}

	contactData := evo_crm.CreateContactRequest{
		Name:        fmt.Sprintf("Test Contact %s", uniqueSuffix),
		PhoneNumber: testPhone,
		Email:       fmt.Sprintf("test-%s@example.com", uniqueSuffix),
		Instagram:   fmt.Sprintf("testgram_%s", uniqueSuffix),
		Annotation:  "Integration test contact created",
		Origin:      evo_crm.OriginCreatedByUser,
	}

	createdContact, err := client.CreateContact(ctx, contactData)

	require.NoError(t, err)
	require.NotNil(t, createdContact)
	require.NotEmpty(t, createdContact.ID)
	require.Equal(t, contactData.Name, createdContact.Name)
	require.Equal(t, contactData.PhoneNumber, createdContact.PhoneNumber)
	require.Equal(t, contactData.Email, createdContact.Email)
	require.Equal(t, contactData.Instagram, *createdContact.Instagram)
	require.Equal(t, string(contactData.Origin), createdContact.Origin)
}
