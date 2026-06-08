package ecommerce;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/products")
public class ProductController {

    @Autowired
    private ProductRepository productRepository;

    @GetMapping("/check")
    public ResponseEntity<Object> checkProduct() {
        Map<String, Object> response = new HashMap<>();
        response.put("message", "Product service is available");
        return ResponseEntity.ok(response);
    }

    @GetMapping
    public ResponseEntity<Object> getAllProducts() {
        List<Product> products = productRepository.findAll();
        Map<String, Object> response = new HashMap<>();
        if (products.isEmpty()) {
            response.put("message", "No products found");
            response.put("data", products);
            return ResponseEntity.ok(response);
        }
        response.put("message", "Products retrieved successfully");
        response.put("data", products);
        return ResponseEntity.ok(response);
    }

    @GetMapping("/{id}")
    public ResponseEntity<Object> getProductById(@PathVariable Long id) {
        return productRepository.findById(id)
                .<ResponseEntity<Object>>map(ResponseEntity::ok)
                .orElseGet(() -> {
                    Map<String, Object> response = new HashMap<>();
                    response.put("error", "Product not found");
                    return ResponseEntity.status(HttpStatus.NOT_FOUND).body(response);
                });
    }

    @PostMapping
    public ResponseEntity<Object> createProduct(@RequestBody Product product) {
        Product savedProduct = productRepository.save(product);
        Map<String, Object> response = new HashMap<>();
        response.put("message", "Product created successfully");
        response.put("data", savedProduct);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    @PutMapping("/{id}")
    public ResponseEntity<Object> updateProduct(@PathVariable Long id, @RequestBody Product productDetails) {
        return productRepository.findById(id)
                .map(product -> {
                    product.setName(productDetails.getName());
                    product.setPrice(productDetails.getPrice());
                    Product updatedProduct = productRepository.save(product);
                    Map<String, Object> response = new HashMap<>();
                    response.put("message", "Product updated successfully");
                    response.put("data", updatedProduct);
                    return ResponseEntity.ok((Object) response);
                })
                .orElseGet(() -> {
                    Map<String, Object> response = new HashMap<>();
                    response.put("error", "Product not found");
                    return ResponseEntity.status(HttpStatus.NOT_FOUND).body(response);
                });
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Object> deleteProduct(@PathVariable Long id) {
        return productRepository.findById(id)
                .map(product -> {
                    productRepository.delete(product);
                    Map<String, Object> response = new HashMap<>();
                    response.put("message", "Product deleted successfully");
                    response.put("id", id);
                    return ResponseEntity.ok((Object) response);
                })
                .orElseGet(() -> {
                    Map<String, Object> response = new HashMap<>();
                    response.put("error", "Product not found");
                    return ResponseEntity.status(HttpStatus.NOT_FOUND).body(response);
                });
    }
}
