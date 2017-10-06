package com.github.kapmahc.axe.nut.helper;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.core.AmqpTemplate;
import org.springframework.amqp.core.Queue;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import javax.annotation.Resource;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

@Component("nut.emailJobSender")
public class EmailJobSender {

    public void send(String to, String subject, String body) throws IOException {
        Map<String, String> map = new HashMap<>();
        map.put("to", to);
        map.put("subject", subject);
        map.put("body", body);
        amqpTemplate.convertAndSend(
                queue.getName(),
                mapper.writeValueAsString(map)
        );
    }


    @PostConstruct
    void init() {
        mapper = new ObjectMapper();
    }

    @Resource
    AmqpTemplate amqpTemplate;
    @Resource(name = "emailsQueue")
    Queue queue;

    private ObjectMapper mapper;
    private final static Logger logger = LoggerFactory.getLogger(EmailJobSender.class);
}
