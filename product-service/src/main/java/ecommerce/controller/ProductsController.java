package ecommerce.controller;

import ecommerce.dto.ApiResponse;
import ecommerce.model.Product;
import ecommerce.service.ProductService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/products")
public class ProductsController {

    @Autowired
    private ProductService productService;

    @GetMapping
    public ResponseEntity<ApiResponse<List<Product>>> getAllProducts() {
        List<Product> products = productService.getAll();
        String message = products.isEmpty() ? "No products found" : "Products retrieved successfully";
        ApiResponse<List<Product>> response = ApiResponse.<List<Product>>builder()
                .message(message)
                .data(products)
                .build();
        return ResponseEntity.ok(response);
    }

    @GetMapping("/{id}")
    public ResponseEntity<ApiResponse<Product>> getProductById(@PathVariable Long id) {
        Product product = productService.getById(id);
        if (product == null) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND)
                    .body(ApiResponse.<Product>builder()
                            .message("Product not found")
                            .build());
        }
        return ResponseEntity.ok(ApiResponse.<Product>builder()
                .message("Product retrieved successfully")
                .data(product)
                .build());
    }

    @PostMapping
    public ResponseEntity<ApiResponse<Product>> createProduct(@RequestBody Product product) {
        Product savedProduct = productService.save(product);
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(ApiResponse.<Product>builder()
                        .message("Product created successfully")
                        .data(savedProduct)
                        .build());
    }

    @PutMapping("/{id}")
    public ResponseEntity<ApiResponse<Product>> updateProduct(@PathVariable Long id, @RequestBody Product productDetails) {
        Product product = productService.getById(id);
        if (product == null) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND)
                    .body(ApiResponse.<Product>builder()
                            .message("Product not found")
                            .build());
        }
        
        product.setName(productDetails.getName());
        product.setPrice(productDetails.getPrice());
        Product updatedProduct = productService.save(product);
        return ResponseEntity.ok(ApiResponse.<Product>builder()
                .message("Product updated successfully")
                .data(updatedProduct)
                .build());
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<ApiResponse<Long>> deleteProduct(@PathVariable Long id) {
        Product product = productService.getById(id);
        if (product == null) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND)
                    .body(ApiResponse.<Long>builder()
                            .message("Product not found")
                            .build());
        }
        productService.delete(id);
        return ResponseEntity.ok(ApiResponse.<Long>builder()
                .message("Product deleted successfully")
                .data(id)
                .build());
    }

    @GetMapping("/check")
    public ResponseEntity<ApiResponse<String>> checkProducts() {
        return ResponseEntity.ok(ApiResponse.<String>builder()
                .message("Products service is available")
                .build());
    }
}
