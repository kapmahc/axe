package com.github.kapmahc.axe;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.core.MessageProperties;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

@Component
public class TaskReceiver {
    public interface Handler {
        void Do(String id, byte[] body) throws IOException;
    }

    public void receiveMessage(Message message) {
        MessageProperties mp = message.getMessageProperties();
        String id = mp.getMessageId();
        String type = mp.getType();
        logger.info("receive {}@{}", id, type);
        Handler hnd = handlerMap.get(type);
        if (hnd == null) {
            logger.error("can't find job handler {}", type);
            return;
        }
        try {
            hnd.Do(id, message.getBody());
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public void register(String type, Handler handler) {
        if (handlerMap.containsKey(type)) {
            logger.warn("already have job handler {}, will override it", type);
        }
        handlerMap.put(type, handler);
    }

    @PostConstruct
    void init() {
        handlerMap = new HashMap<>();
    }

    private Map<String, Handler> handlerMap;
    private final static Logger logger = LoggerFactory.getLogger(TaskReceiver.class);

}
