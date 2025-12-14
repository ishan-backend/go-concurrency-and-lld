package wmg

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	BaseURL         = "http://127.0.0.1:9200"
	MaxWorkers      = 10 // Concurrent order processors
	MaxRetries      = 3  // For 500 errors
	PollInterval    = 50 * time.Millisecond
	MaxPollAttempts = 100 // Timeout for async operations
	BatchSize       = 5   // Orders to take at once
)

type Client struct {
	ingestion  *IngestionAPI
	legacy     *LegacyGenerationAPI
	submission *SubmissionAPI
	wg         sync.WaitGroup
	mu         sync.Mutex
	processed  int
}

func NewClient(baseURL string) *Client {
	return &Client{
		ingestion:  NewIngestionAPI(baseURL),
		legacy:     NewLegacyGenerationAPI(baseURL),
		submission: NewSubmissionAPI(baseURL),
	}
}

func (c *Client) Run() {
	log.Println("Client started")

	// Print initial status
	if status, err := c.legacy.Status(); err == nil {
		log.Printf("Initial status: %+v", status)
	}

	// Create worker pool
	ordersChan := make(chan Order, MaxWorkers*2)

	// Start workers
	for i := 0; i < MaxWorkers; i++ {
		c.wg.Add(1)
		go c.worker(i, ordersChan)
	}

	// Feed orders to workers
	c.feedOrders(ordersChan)

	// Wait for all workers to complete
	c.wg.Wait()

	// Print final status
	if status, err := c.legacy.Status(); err == nil {
		log.Printf("Final status: %+v", status)
	}

	log.Printf("Client finished. Processed %d orders", c.processed)
}

func (c *Client) feedOrders(ordersChan chan<- Order) {
	defer close(ordersChan)

	for {
		// Peek to see if there are orders
		peeked, err := c.retryOperation(func() (interface{}, error) {
			return c.ingestion.PeekOrders(BatchSize)
		})

		if err != nil {
			log.Printf("Error peeking orders: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		orders := peeked.([]Order)
		if len(orders) == 0 {
			log.Println("No more orders in queue")
			break
		}

		// Take the orders
		refs := make([]string, len(orders))
		for i, o := range orders {
			refs[i] = o.ReferenceNumber
		}

		taken, err := c.retryOperation(func() (interface{}, error) {
			return c.ingestion.TakeOrders(refs)
		})

		if err != nil {
			log.Printf("Error taking orders: %v", err)
			continue
		}

		// Send to workers
		for _, order := range taken.([]Order) {
			ordersChan <- order
		}
	}
}

func (c *Client) worker(id int, orders <-chan Order) {
	defer c.wg.Done()

	for order := range orders {
		log.Printf("Worker %d processing order %s", id, order.ReferenceNumber)

		if err := c.processOrder(order); err != nil {
			log.Printf("Worker %d error processing order %s: %v", id, order.ReferenceNumber, err)
		} else {
			c.mu.Lock()
			c.processed++
			c.mu.Unlock()
		}
	}
}

func (c *Client) processOrder(order Order) error {
	// Step 1: Queue asset generation
	err := c.retryOperationSimple(func() error {
		return c.legacy.QueueGenerateAssets(order.ReferenceNumber)
	})
	if err != nil {
		return fmt.Errorf("failed to queue assets: %w", err)
	}

	// Step 2: Wait for assets to complete and retrieve them
	assets, assetsFailed, err := c.waitForAssets(order.ReferenceNumber)
	if err != nil {
		return fmt.Errorf("failed to get assets: %w", err)
	}

	// If assets failed, submit failed order without metadata
	if assetsFailed {
		log.Printf("Assets failed for order %s, submitting failed order", order.ReferenceNumber)
		return c.retryOperationSimple(func() error {
			return c.submission.SubmitFailedOrder(FailedOrder{
				Order:      order,
				FullAssets: nil,
				Metadata:   nil,
			})
		})
	}

	// Step 3: Get asset details
	assetDetails, err := c.getAssetDetails(order.ReferenceNumber, assets)
	if err != nil {
		return fmt.Errorf("failed to get asset details: %w", err)
	}

	// Build full assets
	fullAssets := c.buildFullAssets(assets, assetDetails)

	// Step 4: Queue metadata generation
	err = c.retryOperationSimple(func() error {
		return c.legacy.QueueGenerateMetadata(order.ReferenceNumber)
	})
	if err != nil {
		return fmt.Errorf("failed to queue metadata: %w", err)
	}

	// Step 5: Wait for metadata
	metadata, metadataFailed, err := c.waitForMetadata(order.ReferenceNumber)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	// If metadata failed, submit failed order with assets but no metadata
	if metadataFailed {
		log.Printf("Metadata failed for order %s, submitting failed order with assets", order.ReferenceNumber)
		return c.retryOperationSimple(func() error {
			return c.submission.SubmitFailedOrder(FailedOrder{
				Order:      order,
				FullAssets: fullAssets,
				Metadata:   nil,
			})
		})
	}

	// Step 6: Submit shippable order
	log.Printf("Order %s completed successfully", order.ReferenceNumber)
	return c.retryOperationSimple(func() error {
		return c.submission.SubmitShippableOrder(ShippableOrder{
			Order:      order,
			FullAssets: fullAssets,
			Metadata:   metadata,
		})
	})
}

func (c *Client) waitForAssets(refNum string) ([]Asset, bool, error) {
	for attempt := 0; attempt < MaxPollAttempts; attempt++ {
		result, err := c.retryOperation(func() (interface{}, error) {
			return c.legacy.GetAssets(refNum)
		})

		if err != nil {
			return nil, false, err
		}

		resp := result.(AssetsResponse)

		switch resp.Status {
		case "SUCCESS":
			return resp.Assets, false, nil
		case "FAILED":
			return nil, true, nil
		case "IN_PROGRESS", "QUEUED":
			time.Sleep(PollInterval)
			continue
		default:
			return nil, false, fmt.Errorf("unknown assets status: %s", resp.Status)
		}
	}

	return nil, false, fmt.Errorf("asset generation timeout")
}

func (c *Client) waitForMetadata(refNum string) (*Metadata, bool, error) {
	for attempt := 0; attempt < MaxPollAttempts; attempt++ {
		status, err := c.retryOperation(func() (interface{}, error) {
			return c.legacy.GetMetadataStatus(refNum)
		})

		if err != nil {
			return nil, false, err
		}

		statusStr := status.(string)

		switch statusStr {
		case "SUCCESS":
			// Now get the actual metadata
			result, err := c.retryOperation(func() (interface{}, error) {
				return c.legacy.GetMetadata([]string{refNum})
			})
			if err != nil {
				return nil, false, err
			}
			metadatas := result.([]Metadata)
			if len(metadatas) == 0 {
				return nil, false, fmt.Errorf("no metadata returned")
			}
			return &metadatas[0], false, nil

		case "FAILED":
			return nil, true, nil

		case "IN_PROGRESS", "QUEUED":
			time.Sleep(PollInterval)
			continue

		default:
			return nil, false, fmt.Errorf("unknown metadata status: %s", statusStr)
		}
	}

	return nil, false, fmt.Errorf("metadata generation timeout")
}

func (c *Client) getAssetDetails(refNum string, assets []Asset) ([]AssetDetail, error) {
	assetIDs := make([]int, len(assets))
	for i, a := range assets {
		assetIDs[i] = a.AssetID
	}

	result, err := c.retryOperation(func() (interface{}, error) {
		return c.legacy.GetAssetsDetails(refNum, assetIDs)
	})

	if err != nil {
		return nil, err
	}

	return result.([]AssetDetail), nil
}

func (c *Client) buildFullAssets(assets []Asset, details []AssetDetail) []FullAsset {
	detailMap := make(map[int]AssetDetail)
	for _, d := range details {
		detailMap[d.AssetID] = d
	}

	fullAssets := make([]FullAsset, 0, len(assets))
	for _, asset := range assets {
		if detail, ok := detailMap[asset.AssetID]; ok {
			fullAssets = append(fullAssets, FullAsset{
				Asset:       asset,
				AssetDetail: detail,
			})
		}
	}

	return fullAssets
}

// Retry wrapper for operations returning values
func (c *Client) retryOperation(op func() (interface{}, error)) (interface{}, error) {
	var lastErr error
	for attempt := 0; attempt < MaxRetries; attempt++ {
		result, err := op()
		if err == nil {
			return result, nil
		}

		// Only retry on 500 errors
		if !is500Error(err) {
			return nil, err
		}

		lastErr = err
		if attempt < MaxRetries-1 {
			backoff := time.Duration(attempt+1) * 50 * time.Millisecond
			time.Sleep(backoff)
		}
	}
	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// Retry wrapper for operations returning only error
func (c *Client) retryOperationSimple(op func() error) error {
	var lastErr error
	for attempt := 0; attempt < MaxRetries; attempt++ {
		err := op()
		if err == nil {
			return nil
		}

		if !is500Error(err) {
			return err
		}

		lastErr = err
		if attempt < MaxRetries-1 {
			backoff := time.Duration(attempt+1) * 50 * time.Millisecond
			time.Sleep(backoff)
		}
	}
	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

func is500Error(err error) bool {
	// Check if error message contains "500" or "Internal Server Error"
	errStr := err.Error()
	return containsAny(errStr, "500", "Internal Server Error")
}

func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

func main() {
	client := NewClient(BaseURL)
	client.Run()
}
