package wmg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ============= Types =============

type Order struct {
	ReferenceNumber string `json:"referenceNumber"`
	// Add other fields as needed
}

type Asset struct {
	AssetID         int    `json:"assetId"`
	ReferenceNumber string `json:"referenceNumber"`
	// Add other fields as needed
}

type AssetDetail struct {
	AssetID int    `json:"assetId"`
	Detail  string `json:"detail"`
	// Add other fields as needed
}

type FullAsset struct {
	Asset       Asset       `json:"asset"`
	AssetDetail AssetDetail `json:"assetDetail"`
}

type Metadata struct {
	ReferenceNumber string `json:"referenceNumber"`
	Data            string `json:"data"`
	// Add other fields as needed
}

type ShippableOrder struct {
	Order      Order       `json:"order"`
	FullAssets []FullAsset `json:"fullAssets"`
	Metadata   *Metadata   `json:"metadata"`
}

type FailedOrder struct {
	Order      Order       `json:"order"`
	FullAssets []FullAsset `json:"fullAssets,omitempty"`
	Metadata   *Metadata   `json:"metadata,omitempty"`
}

type AssetsResponse struct {
	Status string  `json:"status"`
	Assets []Asset `json:"assets"`
}

type StatusResponse struct {
	Takeable            int `json:"takeable"`
	Taken               int `json:"taken"`
	AssetsGenQueued     int `json:"assetsGenQueued"`
	AssetsGenComplete   int `json:"assetsGenComplete"`
	MetadataGenQueued   int `json:"metadataGenQueued"`
	MetadataGenComplete int `json:"metadataGenComplete"`
	SubmittedSuccess    int `json:"submittedSuccess"`
	SubmittedFailed     int `json:"submittedFailed"`
	Errored             int `json:"errored"`
}

// ============= Ingestion API =============

type IngestionAPI struct {
	baseURL string
	client  *http.Client
}

func NewIngestionAPI(baseURL string) *IngestionAPI {
	return &IngestionAPI{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (api *IngestionAPI) PeekOrders(limit int) ([]Order, error) {
	url := fmt.Sprintf("%s/ingestion/peek?limit=%d", api.baseURL, limit)
	resp, err := api.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("peek orders failed: %d - %s", resp.StatusCode, string(body))
	}

	var orders []Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (api *IngestionAPI) TakeOrders(references []string) ([]Order, error) {
	url := fmt.Sprintf("%s/ingestion/take", api.baseURL)

	body, err := json.Marshal(map[string][]string{
		"references": references,
	})
	if err != nil {
		return nil, err
	}

	resp, err := api.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("take orders failed: %d - %s", resp.StatusCode, string(respBody))
	}

	var orders []Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// ============= Legacy Generation API =============

type LegacyGenerationAPI struct {
	baseURL string
	client  *http.Client
}

func NewLegacyGenerationAPI(baseURL string) *LegacyGenerationAPI {
	return &LegacyGenerationAPI{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (api *LegacyGenerationAPI) QueueGenerateAssets(referenceNumber string) error {
	url := fmt.Sprintf("%s/legacy/assets/queue", api.baseURL)

	body, err := json.Marshal(map[string]string{
		"referenceNumber": referenceNumber,
	})
	if err != nil {
		return err
	}

	resp, err := api.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("queue assets failed: %d - %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (api *LegacyGenerationAPI) GetAssets(referenceNumber string) (AssetsResponse, error) {
	url := fmt.Sprintf("%s/legacy/assets?referenceNumber=%s", api.baseURL, referenceNumber)
	resp, err := api.client.Get(url)
	if err != nil {
		return AssetsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return AssetsResponse{}, fmt.Errorf("get assets failed: %d - %s", resp.StatusCode, string(body))
	}

	var result AssetsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return AssetsResponse{}, err
	}

	return result, nil
}

func (api *LegacyGenerationAPI) GetAssetsDetails(referenceNumber string, assetIDs []int) ([]AssetDetail, error) {
	// Convert asset IDs to comma-separated string
	idStrs := make([]string, len(assetIDs))
	for i, id := range assetIDs {
		idStrs[i] = fmt.Sprintf("%d", id)
	}
	idsParam := strings.Join(idStrs, ",")

	url := fmt.Sprintf("%s/legacy/assets/details?referenceNumber=%s&assetIds=%s",
		api.baseURL, referenceNumber, idsParam)

	resp, err := api.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get asset details failed: %d - %s", resp.StatusCode, string(body))
	}

	var details []AssetDetail
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, err
	}

	return details, nil
}

func (api *LegacyGenerationAPI) QueueGenerateMetadata(referenceNumber string) error {
	url := fmt.Sprintf("%s/legacy/metadata/queue", api.baseURL)

	body, err := json.Marshal(map[string]string{
		"referenceNumber": referenceNumber,
	})
	if err != nil {
		return err
	}

	resp, err := api.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("queue metadata failed: %d - %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (api *LegacyGenerationAPI) GetMetadataStatus(referenceNumber string) (string, error) {
	url := fmt.Sprintf("%s/legacy/metadata/status?referenceNumber=%s", api.baseURL, referenceNumber)
	resp, err := api.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("get metadata status failed: %d - %s", resp.StatusCode, string(body))
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result["status"], nil
}

func (api *LegacyGenerationAPI) GetMetadata(references []string) ([]Metadata, error) {
	refsParam := strings.Join(references, ",")
	url := fmt.Sprintf("%s/legacy/metadata?references=%s", api.baseURL, refsParam)

	resp, err := api.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get metadata failed: %d - %s", resp.StatusCode, string(body))
	}

	var metadatas []Metadata
	if err := json.NewDecoder(resp.Body).Decode(&metadatas); err != nil {
		return nil, err
	}

	return metadatas, nil
}

func (api *LegacyGenerationAPI) Status() (StatusResponse, error) {
	url := fmt.Sprintf("%s/status", api.baseURL)
	resp, err := api.client.Get(url)
	if err != nil {
		return StatusResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return StatusResponse{}, fmt.Errorf("get status failed: %d - %s", resp.StatusCode, string(body))
	}

	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return StatusResponse{}, err
	}

	return status, nil
}

// ============= Submission API =============

type SubmissionAPI struct {
	baseURL string
	client  *http.Client
}

func NewSubmissionAPI(baseURL string) *SubmissionAPI {
	return &SubmissionAPI{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (api *SubmissionAPI) SubmitShippableOrder(order ShippableOrder) error {
	url := fmt.Sprintf("%s/submission/shippable", api.baseURL)

	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	resp, err := api.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("submit shippable order failed: %d - %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (api *SubmissionAPI) SubmitFailedOrder(order FailedOrder) error {
	url := fmt.Sprintf("%s/submission/failed", api.baseURL)

	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	resp, err := api.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("submit failed order failed: %d - %s", resp.StatusCode, string(respBody))
	}

	return nil
}
