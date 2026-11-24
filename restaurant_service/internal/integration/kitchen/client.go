package kitchen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateOrderStatusRequest struct {
	KitchenOrderID int64  `json:"kitchen_order_id"`
	Status         string `json:"status"`
}

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) UpdateOrderStatus(
	ctx context.Context,
	kitchenBaseURL string,
	apiKey string,
	kitchenOrderID int64,
	status string,
) error {

	requestBody := UpdateOrderStatusRequest{
		KitchenOrderID: kitchenOrderID,
		Status:         status,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		kitchenBaseURL+"/integration/orders/status",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf(
			"kitchen returned status %d",
			response.StatusCode,
		)
	}

	return nil
}
