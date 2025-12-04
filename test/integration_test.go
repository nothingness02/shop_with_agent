package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL2 = "http://localhost:8080/api/v2"

// APIClient provides methods to interact with the API
type APIClient struct {
	client  *http.Client
	baseURL string
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: baseURL,
	}
}

// DoRequest makes an HTTP request and returns the response body
func (c *APIClient) DoRequest(method, endpoint string, body interface{}) ([]byte, int, error) {
	url := c.baseURL + endpoint

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return data, resp.StatusCode, nil
}

// ============= Integration Tests =============

// TestFullShopWorkflow tests complete shop creation, update, and deletion
func (c *APIClient) TestFullShopWorkflow() error {
	fmt.Println("\n========== Testing Shop Workflow ==========")

	// 1. Create a shop
	fmt.Println("\n[1] Creating a new shop...")
	createShopReq := map[string]interface{}{
		"name":        "TechMart Store",
		"description": "A complete tech store",
		"owner_id":    1,
	}

	resp, status, err := c.DoRequest("POST", "/shops", createShopReq)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to create shop: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Shop created\nResponse: %s\n", string(resp))

	// Parse shop ID from response
	var shopResp map[string]interface{}
	if err := json.Unmarshal(resp, &shopResp); err != nil {
		return fmt.Errorf("failed to parse shop response: %v", err)
	}
	shopID := int(shopResp["ID"].(float64))
	fmt.Printf("Shop ID: %d\n", shopID)

	// 2. Get the shop
	fmt.Println("\n[2] Getting shop details...")
	resp, status, err = c.DoRequest("GET", fmt.Sprintf("/shops/%d", shopID), nil)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to get shop: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Shop retrieved\n")

	// 3. Update the shop
	fmt.Println("\n[3] Updating shop...")
	updateShopReq := map[string]interface{}{
		"name":        "TechMart Store Premium",
		"description": "Premium tech store with better inventory",
	}
	resp, status, err = c.DoRequest("PATCH", fmt.Sprintf("/shops/%d", shopID), updateShopReq)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to update shop: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Shop updated\n")

	return nil
}

// TestFullProductWorkflow tests product creation, update, search, and deletion
func (c *APIClient) TestFullProductWorkflow() error {
	fmt.Println("\n========== Testing Product Workflow ==========")

	shopID := 1

	// 1. Create products
	fmt.Println("\n[1] Creating products...")
	products := []map[string]interface{}{
		{
			"name":        "Laptop Pro",
			"description": "High-performance laptop",
			"price":       1299.99,
			"stock":       20,
			"product_img": "https://example.com/laptop.jpg",
		},
		{
			"name":        "Wireless Mouse",
			"description": "Ergonomic wireless mouse",
			"price":       39.99,
			"stock":       100,
			"product_img": "https://example.com/mouse.jpg",
		},
		{
			"name":        "USB Hub",
			"description": "7-port USB hub",
			"price":       49.99,
			"stock":       50,
			"product_img": "https://example.com/hub.jpg",
		},
	}

	productIDs := []int{}
	for i, product := range products {
		resp, status, err := c.DoRequest("POST", fmt.Sprintf("/shops/%d/products", shopID), product)
		if err != nil || status != http.StatusOK {
			return fmt.Errorf("failed to create product %d: %v (status: %d)", i+1, err, status)
		}

		var productResp map[string]interface{}
		if err := json.Unmarshal(resp, &productResp); err != nil {
			return fmt.Errorf("failed to parse product response: %v", err)
		}
		productID := int(productResp["ID"].(float64))
		productIDs = append(productIDs, productID)
		fmt.Printf("✓ Product %d created (ID: %d)\n", i+1, productID)
	}

	// 2. List products
	fmt.Println("\n[2] Listing all products in shop...")
	_, status, err := c.DoRequest("GET", fmt.Sprintf("/shops/%d/products", shopID), nil)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to list products: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Products listed\n")

	// 3. Search for product by name
	fmt.Println("\n[3] Searching for 'Laptop' product...")
	_, status, err = c.DoRequest("GET", fmt.Sprintf("/shops/%d/products/search?name=Laptop", shopID), nil)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to search product: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Product search successful\n")

	// 4. Get specific product
	fmt.Println("\n[4] Getting specific product details...")
	_, status, err = c.DoRequest("GET", fmt.Sprintf("/products/%d", productIDs[0]), nil)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to get product: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Product details retrieved\n")

	// 5. Update product
	fmt.Println("\n[5] Updating product...")
	updateProductReq := map[string]interface{}{
		"name":        "Laptop Pro Max",
		"description": "Ultra high-performance laptop",
		"price":       1499.99,
		"stock":       15,
	}
	_, status, err = c.DoRequest("PATCH", fmt.Sprintf("/products/%d", productIDs[0]), updateProductReq)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to update product: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Product updated\n")

	// 6. Delete single product
	fmt.Println("\n[6] Deleting single product...")
	_, status, err = c.DoRequest("DELETE", fmt.Sprintf("/products/%d", productIDs[2]), nil)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to delete product: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Product deleted\n")

	return nil
}

// TestFullOrderWorkflow tests order creation, status update, and deletion
func (c *APIClient) TestFullOrderWorkflow() error {
	fmt.Println("\n========== Testing Order Workflow ==========")

	// 1. Create order
	fmt.Println("\n[1] Creating order...")
	createOrderReq := map[string]interface{}{
		"user_id": 101,
		"items": []map[string]interface{}{
			{
				"product_id":   1,
				"product_name": "Laptop Pro",
				"product_img":  "https://example.com/laptop.jpg",
				"price":        1299.99,
				"quantity":     1,
			},
			{
				"product_id":   2,
				"product_name": "Wireless Mouse",
				"product_img":  "https://example.com/mouse.jpg",
				"price":        39.99,
				"quantity":     2,
			},
		},
	}

	resp, status, err := c.DoRequest("POST", "/orders", createOrderReq)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to create order: %v (status: %d)", err, status)
	}

	var orderResp map[string]interface{}
	if err := json.Unmarshal(resp, &orderResp); err != nil {
		return fmt.Errorf("failed to parse order response: %v", err)
	}
	orderID := int(orderResp["ID"].(float64))
	fmt.Printf("✓ Order created (ID: %d)\n", orderID)

	// 2. Get order
	fmt.Println("\n[2] Getting order details...")
	resp, status, err = c.DoRequest("GET", fmt.Sprintf("/orders/%d", orderID), nil)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to get order: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Order retrieved\n")

	// 3. Update order status
	fmt.Println("\n[3] Updating order status...")
	updateOrderReq := map[string]interface{}{
		"status": "shipped",
	}
	resp, status, err = c.DoRequest("PATCH", fmt.Sprintf("/orders/%d/status", orderID), updateOrderReq)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to update order status: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Order status updated to 'shipped'\n")

	// 4. List orders
	fmt.Println("\n[4] Listing orders...")
	resp, status, err = c.DoRequest("GET", "/orders", nil)
	if err != nil || status != http.StatusOK {
		return fmt.Errorf("failed to list orders: %v (status: %d)", err, status)
	}
	fmt.Printf("✓ Orders listed\n")

	return nil
}

// RunAllTests executes all workflow tests
func (c *APIClient) RunAllTests() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║     API Integration Tests Started        ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	if err := c.TestFullShopWorkflow(); err != nil {
		fmt.Printf("✗ Shop workflow failed: %v\n", err)
	}

	if err := c.TestFullProductWorkflow(); err != nil {
		fmt.Printf("✗ Product workflow failed: %v\n", err)
	}

	if err := c.TestFullOrderWorkflow(); err != nil {
		fmt.Printf("✗ Order workflow failed: %v\n", err)
	}

	fmt.Println("\n╔══════════════════════════════════════════╗")
	fmt.Println("║     All Tests Completed                  ║")
	fmt.Println("╚══════════════════════════════════════════╝")
}

// Main entry point for testing
func main() {
	client := NewAPIClient(baseURL2)
	client.RunAllTests()
}
