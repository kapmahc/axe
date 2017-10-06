package com.github.kapmahc.axe;

import org.springframework.amqp.core.Queue;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class AmqpConfig {

    @Bean(name = "emailsQueue")
    Queue queue() {
        return new Queue("emails", true);
    }


}
