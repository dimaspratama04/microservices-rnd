package ecommerce;

import ecommerce.model.Product;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.math.BigDecimal;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;
import static org.hamcrest.Matchers.*;

@SpringBootTest
@AutoConfigureMockMvc
class ProductServiceApplicationTests {

	@Autowired
	private MockMvc mockMvc;

	@Autowired
	private ObjectMapper objectMapper;

	@Test
	void contextLoads() {
	}

	@Test
	void testProductCRUD() throws Exception {
		Product product = new Product();
		product.setName("Laptop");		
		product.setPrice(new BigDecimal("999.99"));

		// Create
		String responseJson = mockMvc.perform(post("/products")
				.contentType(MediaType.APPLICATION_JSON)
				.content(objectMapper.writeValueAsString(product)))
				.andExpect(status().isCreated())
				.andExpect(jsonPath("$.message", is("Product created successfully")))
				.andExpect(jsonPath("$.data.name", is("Laptop")))
				.andReturn().getResponse().getContentAsString();

		JsonNode rootNode = objectMapper.readTree(responseJson);
		Long id = rootNode.path("data").path("id").asLong();

		// Read All
		mockMvc.perform(get("/products"))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$.data", hasSize(greaterThanOrEqualTo(1))));

		// Read One
		mockMvc.perform(get("/products/" + id))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$.data.name", is("Laptop")));

		// Update
		product.setName("Gaming Laptop");
		mockMvc.perform(put("/products/" + id)
				.contentType(MediaType.APPLICATION_JSON)
				.content(objectMapper.writeValueAsString(product)))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$.message", is("Product updated successfully")))
				.andExpect(jsonPath("$.data.name", is("Gaming Laptop")));

		// Delete
		mockMvc.perform(delete("/products/" + id))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$.message", is("Product deleted successfully")))
				.andExpect(jsonPath("$.data", is(id.intValue())));

		// Verify Delete
		mockMvc.perform(get("/products/" + id))
				.andExpect(status().isNotFound())
				.andExpect(jsonPath("$.message", is("Product not found")));
	}

}
