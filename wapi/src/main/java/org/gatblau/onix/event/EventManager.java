/*
    Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

    Contributors to this project, hereby assign copyright in their code to the
    project, to be licensed under the same terms as the rest of the code.
*/
package org.gatblau.onix.event;

import org.apache.qpid.jms.JmsConnectionFactory;
import org.gatblau.onix.conf.Config;
import org.gatblau.onix.data.ItemData;
import org.springframework.retry.support.RetryTemplate;
import org.springframework.stereotype.Service;

import javax.annotation.PostConstruct;
import javax.jms.*;
import java.time.ZonedDateTime;
import java.util.UUID;

/*
  Manages publication of change notification events to an AMQP message broker
 */
@Service
public class EventManager {
    private final Config cfg;
    private boolean ready;
    private final RetryTemplate retrySession;
    private final ConnectionFactory connFactory;
    private Session session;
    private Connection conn;

    public EventManager(Config cfg, RetryTemplate retrySession) {
        this.cfg = cfg;
        this.retrySession = retrySession;
        this.connFactory = new JmsConnectionFactory(String.format("amqp://%s:%s", cfg.getAmqpHost(), cfg.getAmqpPort()));
    }

    @PostConstruct
    private void init(){
        // tries and acquire a session with the message broker
        // keeps retrying if it fails, based on the retry template
        this.session = retrySession.execute(context -> {
            try {
                // create an amqp qpid 1.0 connection
                conn = connFactory.createConnection(cfg.getAmqpUser(), cfg.getAmqpPwd());
                // create a session
                session = conn.createSession(false, Session.AUTO_ACKNOWLEDGE);
                // now the manager is ready to use
                ready = true;
            } catch (Exception e) {
                System.out.println(String.format("WARNING: Failed to establish broker session: %s retrying...", e.getMessage()));
               throw new RuntimeException(e);
            }
            return session;
        });
    }

    /**
     * Send a change notification event
     * @param notifyType
     * @param item
     * @return
     * @throws JMSException
     */
    public boolean send(char notifyType, ItemData item) throws JMSException {
        if (ready) {
            ItemChanged itemChanged = new ItemChanged(notifyType, item);
            // create a sender
            Topic topic = session.createTopic(itemChanged.getTopicName());
            MessageProducer sender = session.createProducer(topic);
            // send the message
            sender.send(createMessage(itemChanged.toString()));
        } else {
            System.out.println("WARNING: broker service not ready");
        }
        return ready;
    }

    private Message createMessage(String text) throws JMSException {
        Message msg = null;
        if (ready) {
            msg = session.createTextMessage(text);
            // assigns a randomly generated GUID
            msg.setJMSMessageID(UUID.randomUUID().toString());
            msg.setJMSTimestamp(now());
        }
        return msg;
    }

    private long now() {
        ZonedDateTime zdt = ZonedDateTime.now();
        return zdt.toInstant().toEpochMilli();
    }
}
