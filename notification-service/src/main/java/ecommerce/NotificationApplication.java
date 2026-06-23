package ecommerce;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import ecommerce.config.RabbitMQ;

@SpringBootApplication
public class NotificationApplication {

	public static void main(String[] args) {
		SpringApplication.run(NotificationApplication.class, args);
		System.out.println(RabbitMQ.NOTIFICATION_QUEUE);
	}

}
