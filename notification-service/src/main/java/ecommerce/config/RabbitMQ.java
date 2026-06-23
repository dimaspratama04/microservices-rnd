package ecommerce.config;

import org.springframework.amqp.core.Queue;
import org.springframework.amqp.rabbit.annotation.RabbitHandler;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitMQ {
    public static final String NOTIFICATION_QUEUE = "notification.created";

    @Bean
    public Queue notificationWorker(){
        return new Queue(NOTIFICATION_QUEUE,true);
    }

}