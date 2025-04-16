package evo_crm_test

import (
	"log"
	"os"
	"testing"

	evo_crm "github.com.br/verbeux-ai/evo-crm-go"
	"github.com/joho/godotenv"
)

var client *evo_crm.Client

func TestMain(m *testing.M) {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("Warning: could not load .env file. Ensure EVO_CRM_TOKEN is set via environment.")
	}

	apiToken := os.Getenv("EVO_CRM_TOKEN")
	if apiToken == "" {
		log.Fatal("Error: EVO_CRM_TOKEN environment variable not set.")
	}

	client = evo_crm.NewClient(
		evo_crm.WithToken(apiToken),
	)

	exitCode := m.Run()
	os.Exit(exitCode)
}
