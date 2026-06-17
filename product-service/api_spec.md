# Product Service API Specification

The Product Service manages the catalog of products in the ecommerce system.

**Base URL:** `http://localhost:8081/products`

---

## Endpoints

### 1. Get All Products
Retrieves a list of all products in the database.

- **URL:** `/`
- **Method:** `GET`
- **Response:**
    - **Status:** `200 OK`
    - **Body:**
      ```json
      {
        "message": "Products retrieved successfully",
        "data": [
          {
            "id": 1,
            "name": "Product Name",
            "price": 99.99
          }
        ]
      }
      ```
    - *Note:* If no products are found, the message will be "No products found" and data will be an empty list.

### 2. Get Product by ID
Retrieves details of a specific product.

- **URL:** `/{id}`
- **Method:** `GET`
- **Path Parameters:**
    - `id` (Long): Unique identifier of the product.
- **Response:**
    - **Status:** `200 OK`
    - **Body:**
      ```json
      {
        "message": "Product retrieved successfully",
        "data": {
          "id": 1,
          "name": "Product Name",
          "price": 99.99
        }
      }
      ```
    - **Status:** `404 Not Found`
    - **Body:**
      ```json
      {
        "message": "Product not found",
        "data": null
      }
      ```

### 3. Create Product
Creates a new product entry.

- **URL:** `/`
- **Method:** `POST`
- **Request Body:**
  ```json
  {
    "name": "New Product",
    "price": 49.99
  }
  ```
- **Response:**
    - **Status:** `201 Created`
    - **Body:**
      ```json
      {
        "message": "Product created successfully",
        "data": {
          "id": 2,
          "name": "New Product",
          "price": 49.99
        }
      }
      ```

### 4. Update Product
Updates an existing product's details.

- **URL:** `/{id}`
- **Method:** `PUT`
- **Path Parameters:**
    - `id` (Long): Unique identifier of the product to update.
- **Request Body:**
  ```json
  {
    "name": "Updated Product Name",
    "price": 59.99
  }
  ```
- **Response:**
    - **Status:** `200 OK`
    - **Body:**
      ```json
      {
        "message": "Product updated successfully",
        "data": {
          "id": 1,
          "name": "Updated Product Name",
          "price": 59.99
        }
      }
      ```
    - **Status:** `404 Not Found`
    - **Body:**
      ```json
      {
        "message": "Product not found",
        "data": null
      }
      ```

### 5. Delete Product
Removes a product from the database.

- **URL:** `/{id}`
- **Method:** `DELETE`
- **Path Parameters:**
    - `id` (Long): Unique identifier of the product to delete.
- **Response:**
    - **Status:** `200 OK`
    - **Body:**
      ```json
      {
        "message": "Product deleted successfully",
        "data": 1
      }
      ```
    - **Status:** `404 Not Found`
    - **Body:**
      ```json
      {
        "message": "Product not found",
        "data": null
      }
      ```

### 6. Health Check
Checks if the product service is available.

- **URL:** `/check`
- **Method:** `GET`
- **Response:**
    - **Status:** `200 OK`
    - **Body:**
      ```json
      {
        "message": "Products service is available",
        "data": null
      }
      ```

---

## Data Models

### Product
| Field | Type | Description |
|---|---|---|
| `id` | Long | Unique identifier |
| `name` | String | Name of the product |
| `price` | Double | Price of the product |

### ApiResponse
| Field | Type | Description |
|---|---|---|
| `message` | String | Status message |
| `data` | T | The payload of the response (can be an object, list, or null) |
