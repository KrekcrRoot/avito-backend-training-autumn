package restaurant

import (
	"avito-kitchen-service/internal/infrastructure/db/sqlc"
	"avito-kitchen-service/internal/transport/http/dto"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CreateOrderRequest struct {
	KitchenOrderID int64                    `json:"kitchen_order_id"`
	Items          []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) GetProducts(
	ctx context.Context,
	restaurant sqlc.Restaurant,
) ([]dto.Product, error) {

	url := strings.TrimRight(restaurant.BaseUrl, "/") + "/api/v1/products"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", restaurant.ApiKey)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"restaurant returned status %d",
			resp.StatusCode,
		)
	}

	var products []dto.Product

	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	return products, nil

}

func (c *Client) CreateOrder(
	ctx context.Context,
	restaurant sqlc.Restaurant,
	kitchenOrderID int64,
	items []CreateOrderItemRequest,
) (sqlc.Order, error) {

	requestBody := CreateOrderRequest{
		KitchenOrderID: kitchenOrderID,
		Items:          items,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return sqlc.Order{}, err
	}

	url := strings.TrimRight(restaurant.BaseUrl, "/") + "/api/v1/orders"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return sqlc.Order{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", restaurant.ApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return sqlc.Order{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return sqlc.Order{}, fmt.Errorf(
			"restaurant returned status %d",
			resp.StatusCode,
		)
	}

	var order sqlc.Order

	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return sqlc.Order{}, err
	}

	return order, nil
}
