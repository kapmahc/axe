package com.github.kapmahc.axe.nut.helper;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.rabbit.annotation.RabbitHandler;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.core.env.Environment;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import javax.annotation.Resource;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

@Component("nut.emailJobReceiver")
@RabbitListener(queues = "emails")
public class EmailJobReceiver {
    @RabbitHandler
    public void process(String message) throws IOException {
        TypeReference<HashMap<String, String>> ref = new TypeReference<HashMap<String, String>>() {
        };
        Map<String, String> args = mapper.readValue(message, ref);
        String to = args.get("to");
        String subject = args.get("subject");
        String body = args.get("body");
        if (env.acceptsProfiles("production")) {
            // TODO send mail
        } else {
            logger.debug("send email to {}\n{}\n{}", to, subject, body);
        }
    }

    @PostConstruct
    void init() {
        mapper = new ObjectMapper();
    }


    @Resource
    Environment env;
    private ObjectMapper mapper;
    private final static Logger logger = LoggerFactory.getLogger(EmailJobReceiver.class);
}
