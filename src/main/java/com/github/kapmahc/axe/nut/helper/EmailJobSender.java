package com.github.kapmahc.axe.nut.helper;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.kapmahc.axe.TaskReceiver;
import org.springframework.amqp.core.MessageBuilder;
import org.springframework.amqp.core.MessageProperties;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import javax.annotation.Resource;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

@Component("nut.emailJobSender")
public class EmailJobSender {

    public void send(String to, String subject, String body) throws IOException {
        Map<String, String> map = new HashMap<>();
        map.put("to", to);
        map.put("subject", subject);
        map.put("body", body);
        MessageProperties mp = new MessageProperties();
        mp.setMessageId(UUID.randomUUID().toString());
        mp.setType(TYPE);
        rabbitTemplate.send(
                MessageBuilder
                        .withBody(mapper.writeValueAsBytes(map))
                        .andProperties(mp)
                        .build()
        );
    }


    @PostConstruct
    void init() {
        taskReceiver.register(TYPE, emailJobReceiver);
        mapper = new ObjectMapper();
    }

    @Resource
    RabbitTemplate rabbitTemplate;
    @Resource
    EmailJobReceiver emailJobReceiver;
    @Resource
    TaskReceiver taskReceiver;
    private ObjectMapper mapper;


    private final static String TYPE = "send-email";

}
