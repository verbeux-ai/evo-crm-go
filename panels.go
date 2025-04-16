package evo_crm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
