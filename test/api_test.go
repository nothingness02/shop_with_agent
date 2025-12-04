package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080/api/v2"

// ============= Shop Tests =============

func TestCreateShop(t *testing.T) {
	req := map[string]interface{}{
		"name":        "Test Shop",
		"description": "This is a test shop",
		"owner_id":    1,
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/shops", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create shop: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Create Shop Response: %s", string(data))
}

func TestListShops(t *testing.T) {
	resp, err := http.Get(baseURL + "/shops")
	if err != nil {
		t.Fatalf("Failed to list shops: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("List Shops Response: %s", string(data))
}

func TestGetShop(t *testing.T) {
	shopID := 1
	resp, err := http.Get(fmt.Sprintf("%s/shops/%d", baseURL, shopID))
	if err != nil {
		t.Fatalf("Failed to get shop: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Get Shop Response: %s", string(data))
}

func TestUpdateShop(t *testing.T) {
	shopID := 1
	req := map[string]interface{}{
		"name":        "Updated Shop",
		"description": "Updated description",
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/shops/%d", baseURL, shopID), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("Failed to update shop: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Update Shop Response: %s", string(data))
}

func TestDeleteShop(t *testing.T) {
	shopID := 1
	httpReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/shops/%d", baseURL, shopID), nil)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("Failed to delete shop: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Delete Shop Response: %s", string(data))
}

// ============= Product Tests =============

func TestCreateProduct(t *testing.T) {
	shopID := 1
	req := map[string]interface{}{
		"name":        "Test Product",
		"description": "This is a test product",
		"price":       99.99,
		"stock":       100,
		"product_img": "https://example.com/product.jpg",
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post(fmt.Sprintf("%s/shops/%d/products", baseURL, shopID), "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Create Product Response: %s", string(data))
}

func TestListProducts(t *testing.T) {
	shopID := 1
	resp, err := http.Get(fmt.Sprintf("%s/shops/%d/products", baseURL, shopID))
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("List Products Response: %s", string(data))
}

func TestGetProduct(t *testing.T) {
	productID := 1
	resp, err := http.Get(fmt.Sprintf("%s/products/%d", baseURL, productID))
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Get Product Response: %s", string(data))
}

func TestSearchProductByName(t *testing.T) {
	shopID := 1
	productName := "iPhone"
	resp, err := http.Get(fmt.Sprintf("%s/shops/%d/products/search?name=%s", baseURL, shopID, productName))
	if err != nil {
		t.Fatalf("Failed to search product by name: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Search Product Response: %s", string(data))
}

func TestUpdateProduct(t *testing.T) {
	productID := 1
	req := map[string]interface{}{
		"name":        "Updated Product",
		"description": "Updated product description",
		"price":       149.99,
		"stock":       50,
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/products/%d", baseURL, productID), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Update Product Response: %s", string(data))
}

func TestDeleteProduct(t *testing.T) {
	productID := 1
	httpReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/products/%d", baseURL, productID), nil)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Delete Product Response: %s", string(data))
}

// ============= Order Tests =============

func TestCreateOrder(t *testing.T) {
	req := map[string]interface{}{
		"user_id": 101,
		"items": []map[string]interface{}{
			{
				"product_id":   1,
				"product_name": "iPhone 15",
				"price":        999.99,
				"quantity":     1,
			},
			{
				"product_id":   3,
				"product_name": "AirPods Pro",
				"price":        249.99,
				"quantity":     2,
			},
		},
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post(baseURL+"/orders", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Create Order Response: %s", string(data))
}

func TestListOrders(t *testing.T) {
	resp, err := http.Get(baseURL + "/orders")
	if err != nil {
		t.Fatalf("Failed to list orders: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("List Orders Response: %s", string(data))
}

func TestGetOrder(t *testing.T) {
	orderID := 1
	resp, err := http.Get(fmt.Sprintf("%s/orders/%d", baseURL, orderID))
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Get Order Response: %s", string(data))
}

func TestUpdateOrderStatus(t *testing.T) {
	orderID := 1
	req := map[string]interface{}{
		"status": "completed",
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/orders/%d/status", baseURL, orderID), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("Failed to update order status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Update Order Status Response: %s", string(data))
}

func TestDeleteOrder(t *testing.T) {
	orderID := 1
	httpReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/orders/%d", baseURL, orderID), nil)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("Failed to delete order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Delete Order Response: %s", string(data))
}

// ============= Batch Tests =============

func TestBatchDeleteProducts(t *testing.T) {
	req := map[string]interface{}{
		"ids": []uint{2, 3, 4},
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodDelete, baseURL+"/products", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("Failed to batch delete products: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Batch Delete Products Response: %s", string(data))
}

func TestBatchDeleteShops(t *testing.T) {
	req := map[string]interface{}{
		"ids": []uint{2, 3},
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodDelete, baseURL+"/shops", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		t.Fatalf("Failed to batch delete shops: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)
	t.Logf("Batch Delete Shops Response: %s", string(data))
}
