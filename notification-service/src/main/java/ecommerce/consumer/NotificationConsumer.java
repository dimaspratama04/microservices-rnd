package ecommerce.consumer;

import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;

import ecommerce.config.RabbitMQ;

@Component
public class NotificationConsumer {

    @RabbitListener(queues = RabbitMQ.NOTIFICATION_QUEUE)
    public void consume(String message) {
        System.out.println("Received: " + message);
    }
    

}
