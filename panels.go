package evo_crm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Panel struct {
	ID               string    `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CompanyID        string    `json:"companyId"`
	Archived         bool      `json:"archived"`
	Scope            string    `json:"scope"`
	DepartmentIDs    []string  `json:"departmentIds"`
	UserID           string    `json:"userId"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	ThumbnailID      *string   `json:"thumbnailId"`
	ThumbnailFile    *any      `json:"thumbnailFile"`
	Key              string    `json:"key"`
	OverdueCardCount int       `json:"overdueCardCount"`
	StepTitles       []string  `json:"stepTitles"`
	Tags             []any     `json:"tags"`
	Steps            []any     `json:"steps"`
}

type ListPanelsResponse struct {
	Items          []Panel `json:"items"`
	TotalItems     int     `json:"totalItems"`
	TotalPages     int     `json:"totalPages"`
	HasMorePages   bool    `json:"hasMorePages"`
	PageNumber     int     `json:"pageNumber"`
	PageSize       int     `json:"pageSize"`
	OrderBy        string  `json:"orderBy"`
	OrderDirection string  `json:"orderDirection"`
}

func (s *Client) ListPanels(ctx context.Context) (*ListPanelsResponse, error) {
	resp, err := s.request(ctx, nil, http.MethodGet, listPanelsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to list panels: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to list panels (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to list panels (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload ListPanelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding list panels response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}

type PanelCard struct {
	ID                 string    `json:"id"`
	Active             bool      `json:"active"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	Archived           bool      `json:"archived"`
	CompanyID          string    `json:"companyId"`
	PanelID            string    `json:"panelId"`
	PanelTitle         *string   `json:"panelTitle"`
	StepID             string    `json:"stepId"`
	StepTitle          *string   `json:"stepTitle"`
	StepPrimaryColor   *string   `json:"stepPrimaryColor"`
	StepSecondaryColor *string   `json:"stepSecondaryColor"`
	StepPhase          *string   `json:"stepPhase"`
	Position           float64   `json:"position"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Key                string    `json:"key"`
	Number             int       `json:"number"`
	DueDate            *string   `json:"dueDate"`
	IsOverdue          bool      `json:"isOverdue"`
	IsInDueDate        bool      `json:"isInDueDate"`
	ResponsibleUserID  string    `json:"responsibleUserId"`
	TagIDs             []string  `json:"tagIds"`
	ContactIDs         []string  `json:"contactIds"`
	SessionID          *string   `json:"sessionId"`
	MonetaryAmount     float64   `json:"monetaryAmount"`
	SearchableText     string    `json:"searchableText"`
	CustomFields       any       `json:"customFields"`
	ResponsibleUser    *any      `json:"responsibleUser"`
	Contacts           []any     `json:"contacts"`
	Notes              []any     `json:"notes"`
	Logs               []any     `json:"logs"`
	Metadata           *any      `json:"metadata"`
}

type ListPanelCardsFilterRequest struct {
	PanelID           string `query:"panelId"`
	StepID            string `query:"stepId"`
	ContactID         string `query:"contactId"`
	ResponsibleUserID string `query:"responsibleUserId"`
	TextFilter        string `query:"textFilter"`
	IncludeArchived   bool   `query:"includeArchived"`
	IncludeDetails    bool   `query:"includeDetails"`
	PageNumber        int    `query:"pageNumber"`
	PageSize          int    `query:"pageSize"`
	OrderBy           string `query:"orderBy"`
	OrderDirection    string `query:"orderDirection"`
}

type ListPanelCardsResponse struct {
	Items          []PanelCard `json:"items"`
	TotalItems     int         `json:"totalItems"`
	TotalPages     int         `json:"totalPages"`
	HasMorePages   bool        `json:"hasMorePages"`
	PageNumber     int         `json:"pageNumber"`
	PageSize       int         `json:"pageSize"`
	OrderBy        string      `json:"orderBy"`
	OrderDirection string      `json:"orderDirection"`
}

func (s *Client) ListPanelCardsFilter(ctx context.Context, filter ListPanelCardsFilterRequest) (*ListPanelCardsResponse, error) {
	if filter.PanelID == "" {
		return nil, errors.New("panelId is required for ListPanelCardsFilter")
	}

	queryString, err := StructToQueryString(filter)
	if err != nil {
		return nil, fmt.Errorf("error creating query string from filter: %w", err)
	}

	fullPath := listPanelCardsEndpoint
	if queryString != "" {
		fullPath += "?" + queryString
	}

	resp, err := s.request(ctx, nil, http.MethodGet, fullPath)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to list panel cards: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to list panel cards (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to list panel cards (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload ListPanelCardsResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding list panel cards response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}

type PanelIncludeDetails string

const (
	IncludeSteps          PanelIncludeDetails = "Steps"
	IncludeStepsFields    PanelIncludeDetails = "StepsFields"
	IncludeStepsCardCount PanelIncludeDetails = "StepsCardCount"
	IncludeTags           PanelIncludeDetails = "Tags"
	IncludeCards          PanelIncludeDetails = "Cards"
)

func (s *Client) GetPanelByID(ctx context.Context, panelId string, includeDetails []PanelIncludeDetails) (*Panel, error) {
	if panelId == "" {
		return nil, errors.New("panelId is required for GetPanelByID")
	}

	endpointPath := fmt.Sprintf(getPanelByIDEndpoint, panelId)
	fullPath := endpointPath

	if len(includeDetails) > 0 {
		queryParams := url.Values{}
		for _, detail := range includeDetails {
			queryParams.Add("includeDetails", string(detail))
		}
		fullPath += "?" + queryParams.Encode()
	}

	resp, err := s.request(ctx, nil, http.MethodGet, fullPath)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to get panel by id: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to get panel by id (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to get panel by id (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload Panel
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding get panel by id response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}

type CreatePanelCardRequest struct {
	StepID            string         `json:"stepId"`
	Title             string         `json:"title"`
	Description       string         `json:"description,omitempty"`
	Position          float64        `json:"position,omitempty"`
	ResponsibleUserID string         `json:"responsibleUserId,omitempty"`
	TagIDs            []string       `json:"tagIds,omitempty"`
	ContactIDs        []string       `json:"contactIds,omitempty"`
	MonetaryAmount    float64        `json:"monetaryAmount,omitempty"`
	Metadata          map[string]any `json:"metadata,omitempty"`
	CustomFields      map[string]any `json:"customFields,omitempty"`
}

func (s *Client) CreatePanelCard(ctx context.Context, card CreatePanelCardRequest) (*PanelCard, error) {
	if card.StepID == "" {
		return nil, errors.New("stepId is required for CreatePanelCard")
	}
	if card.Title == "" {
		return nil, errors.New("title is required for CreatePanelCard")
	}

	resp, err := s.request(ctx, card, http.MethodPost, createPanelCardEndpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to create panel card: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to create panel card (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to create panel card (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload PanelCard
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding create panel card response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}
