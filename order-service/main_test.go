package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"order-service/domain"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&domain.OrderModel{})
	return db
}

func TestOrderCRUD(t *testing.T) {
	os.Setenv("GO_ENV", "test")
	db := setupTestDB()

	// Mock Product Service
	productSvc := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/products/1" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": 1, "name": "Test Product"}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer productSvc.Close()
	os.Setenv("PRODUCT_SERVICE_URL", productSvc.URL)

	app := SetupApp(db)

	// 1. Create - Success
	newOrder := domain.OrderModel{ProductID: 1, Quantity: 2, Total: 100.0}
	body, _ := json.Marshal(newOrder)
	req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Order created successfully", result["message"])

	data := result["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["product_id"])
	assert.Equal(t, "Pending", data["status"])

	// 1.1 Create - Fail (Product Not Found)
	newOrderFail := domain.OrderModel{ProductID: 99, Quantity: 1, Total: 50.0}
	bodyFail, _ := json.Marshal(newOrderFail)
	reqFail := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(bodyFail))
	reqFail.Header.Set("Content-Type", "application/json")
	respFail, _ := app.Test(reqFail)

	assert.Equal(t, 400, respFail.StatusCode)
	var errResult map[string]interface{}
	json.NewDecoder(respFail.Body).Decode(&errResult)
	assert.Equal(t, "product not found or product service unavailable", errResult["error"])

	// 2. Read All
	req = httptest.NewRequest("GET", "/orders", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	var readAllResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&readAllResult)
	assert.Equal(t, "Orders retrieved successfully", readAllResult["message"])
	assert.NotEmpty(t, readAllResult["data"])

	// 3. Read One
	req = httptest.NewRequest("GET", "/orders/1", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	// 4. Update
	updatePayload := map[string]interface{}{"status": "Paid"}
	body, _ = json.Marshal(updatePayload)
	req = httptest.NewRequest("PUT", "/orders/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	var updateResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&updateResult)
	assert.Equal(t, "Order updated successfully", updateResult["message"])

	updatedData := updateResult["data"].(map[string]interface{})
	assert.Equal(t, "Paid", updatedData["status"])

	// 5. Delete
	req = httptest.NewRequest("DELETE", "/orders/1", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	var deleteResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&deleteResult)
	assert.Equal(t, "Order deleted successfully", deleteResult["message"])
	assert.Equal(t, "1", deleteResult["id"])

	// 6. Verify Delete
	req = httptest.NewRequest("GET", "/orders/1", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, 404, resp.StatusCode)
}
