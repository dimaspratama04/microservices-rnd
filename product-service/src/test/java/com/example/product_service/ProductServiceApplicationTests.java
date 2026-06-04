package com.example.product_service;

import com.fasterxml.jackson.databind.ObjectMapper;
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
		Product product = new Product("Laptop", 1200.0);

		// Create
		String json = mockMvc.perform(post("/products")
				.contentType(MediaType.APPLICATION_JSON)
				.content(objectMapper.writeValueAsString(product)))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$.name", is("Laptop")))
				.andReturn().getResponse().getContentAsString();

		Product createdProduct = objectMapper.readValue(json, Product.class);
		Long id = createdProduct.getId();

		// Read All
		mockMvc.perform(get("/products"))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$", hasSize(greaterThanOrEqualTo(1))));

		// Read One
		mockMvc.perform(get("/products/" + id))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$.name", is("Laptop")));

		// Update
		createdProduct.setName("Gaming Laptop");
		mockMvc.perform(put("/products/" + id)
				.contentType(MediaType.APPLICATION_JSON)
				.content(objectMapper.writeValueAsString(createdProduct)))
				.andExpect(status().isOk())
				.andExpect(jsonPath("$.name", is("Gaming Laptop")));

		// Delete
		mockMvc.perform(delete("/products/" + id))
				.andExpect(status().isOk());

		// Verify Delete
		mockMvc.perform(get("/products/" + id))
				.andExpect(status().isNotFound());
	}

}
