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

var (
	ErrRequestFailed    = errors.New("http request execution failed")
	ErrReadResponseBody = errors.New("failed to read response body")
	ErrApiReturnedError = errors.New("api returned an error status")
	ErrDecodeResponse   = errors.New("failed to decode api response")
)

type ContactOrigin string

const (
	OriginImported       ContactOrigin = "IMPORTED"
	OriginCreatedFromHub ContactOrigin = "CREATED_FROM_HUB"
	OriginCreatedByUser  ContactOrigin = "CREATED_BY_USER"
)

type ListContactsFilterRequest struct {
	Name               string          `json:"name,omitempty"`
	PhoneNumber        string          `json:"phonenumber,omitempty"`
	Instagram          string          `json:"instagram,omitempty"`
	Email              string          `json:"email,omitempty"`
	TagsID             []string        `json:"tagsId,omitempty"`
	TagsOperator       string          `json:"tagsOperator,omitempty"`
	ExceptTagsID       []string        `json:"exceptTagsId,omitempty"`
	ExceptTagsOperator string          `json:"exceptTagsOperator,omitempty"`
	Origins            []ContactOrigin `json:"origins,omitempty"`
	Status             string          `json:"status,omitempty"`
	Page               int             `json:"page,omitempty"`
	PageSize           int             `json:"pageSize,omitempty"`
	OrderBy            string          `json:"orderBy,omitempty"`
}

type Contact struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	CompanyID            string    `json:"companyId"`
	Name                 string    `json:"name"`
	NameWhatsapp         *string   `json:"nameWhatsapp"`
	NameInstagram        *string   `json:"nameInstagram"`
	NameMessenger        *string   `json:"nameMessenger"`
	PhoneNumber          string    `json:"phonenumber"`
	PhoneNumberFormatted string    `json:"phonenumberFormatted"`
	Email                string    `json:"email"`
	Instagram            *string   `json:"instagram"`
	MessengerID          *string   `json:"messengerId"`
	PictureFileID        *string   `json:"pictureFileId"`
	PictureURL           *string   `json:"pictureUrl"`
	Active               bool      `json:"active"`
	Annotation           *string   `json:"annotation"`
	TagsID               []string  `json:"tagsId"`
	PortfolioIDs         []string  `json:"portfolioIds"`
	Tags                 []any     `json:"tags"`
	Status               string    `json:"status"`
	Origin               string    `json:"origin"`
	CustomFieldValues    []any     `json:"customFieldValues"`
	Utm                  *any      `json:"utm"`
	OptInStatus          string    `json:"optInStatus"`
	OptInUpdatedAt       *string   `json:"optInUpdatedAt"`
	ImportedAt           *string   `json:"importedAt"`
	WsClientID           *string   `json:"wsClientId"`
	Metadata             *any      `json:"metadata"`
}

type ListContactsFilterResponse struct {
	LastTimestamp   string    `json:"lastTimestamp"`
	Items           []Contact `json:"items"`
	TotalItems      int       `json:"totalItems"`
	CountPages      int       `json:"countPages"`
	IsLastPage      bool      `json:"isLastPage"`
	Page            int       `json:"page"`
	PageSize        int       `json:"pageSize"`
	OrderBy         string    `json:"orderBy"`
	OrderByDesc     *bool     `json:"orderByDesc"`
	TimestampField  *string   `json:"timestampField"`
	TimestampFilter *string   `json:"timestampFilter"`
	NextPageToken   *string   `json:"nextPageToken"`
	Type            string    `json:"type"`
}

func (s *Client) ListContactsFilter(ctx context.Context, filter ListContactsFilterRequest) ([]Contact, error) {
	resp, err := s.request(ctx, filter, http.MethodPost, contactsFilterEndpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to list contacts: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to list contacts (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to list contacts (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload ListContactsFilterResponse
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding list contacts response: %w", ErrDecodeResponse, err)
	}

	return responsePayload.Items, nil
}

type CreateContactRequest struct {
	Name        string        `json:"name"`
	PhoneNumber string        `json:"phonenumber"`
	Email       string        `json:"email,omitempty"`
	Instagram   string        `json:"instagram,omitempty"`
	Annotation  string        `json:"annotation,omitempty"`
	Origin      ContactOrigin `json:"origin,omitempty"`
}

func (s *Client) CreateContact(ctx context.Context, contact CreateContactRequest) (*Contact, error) {
	resp, err := s.request(ctx, contact, http.MethodPost, createContactEndpoint)
	if err != nil {
		return nil, fmt.Errorf("%w: error making request to create contact: %w", ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("%w: failed to create contact (status: %d), read response body error: %w", ErrReadResponseBody, resp.StatusCode, readErr)
		}
		bodyErr := errors.New(string(bodyBytes))
		return nil, fmt.Errorf("%w: failed to create contact (status: %d): %w", ErrApiReturnedError, resp.StatusCode, bodyErr)
	}

	var responsePayload Contact
	if err := json.NewDecoder(resp.Body).Decode(&responsePayload); err != nil {
		return nil, fmt.Errorf("%w: error decoding create contact response: %w", ErrDecodeResponse, err)
	}

	return &responsePayload, nil
}
