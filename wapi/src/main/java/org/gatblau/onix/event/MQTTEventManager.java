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

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.eclipse.paho.client.mqttv3.*;
import org.eclipse.paho.client.mqttv3.persist.MemoryPersistence;
import org.gatblau.onix.conf.Config;
import org.gatblau.onix.data.ItemData;
import org.springframework.retry.support.RetryTemplate;
import org.springframework.stereotype.Service;

import javax.annotation.PostConstruct;
import java.util.UUID;

@Service
public class MQTTEventManager implements EventManager {
    private IMqttClient publisher;
    private final Config cfg;
    private final RetryTemplate retrySession;
    private final Logger log = LogManager.getLogger();
    private int retryCount;

    public MQTTEventManager(Config cfg, RetryTemplate retrySession) {
        this.cfg = cfg;
        this.retrySession = retrySession;
    }

    @PostConstruct
    private void init(){
        // only initialises the connectivity to the message broker if enabled in the configuration
        if (cfg.isEventsEnabled()) {
            log.atInfo().log("item change notifications are enabled");
            log.atInfo().log(String.format("attempting to connect to the message broker at %s", getConnectionURI()));
            // tries and acquire a session with the message broker
            // keeps retrying if it fails, based on the retry template
            publisher = retrySession.execute(context -> {
                IMqttClient client = null;
                try {
                    // create a new IMqttClient synchronous instance
                    String publisherId = UUID.randomUUID().toString();
                    client = new MqttClient(getConnectionURI(), publisherId, new MemoryPersistence());
                    // connect to the server
                    MqttConnectOptions options = new MqttConnectOptions();
                    options.setAutomaticReconnect(true);
                    options.setCleanSession(true);
                    options.setConnectionTimeout(10);
                    // if credentials are defined
                    if (cfg.getEventsServerUser() != null && !cfg.getEventsServerUser().isEmpty()) {
                        // set credentials
                        options.setUserName(cfg.getEventsServerUser());
                        options.setPassword(cfg.getEventsServerPwd().toCharArray());
                    }
                    client.connect(options);
                    // now the manager is ready to use
                    retryCount = 0;
                    log.atInfo().log(String.format("successfully connected to %s, publisher is ready to use", getConnectionURI()));
                } catch (Exception e) {
                    retryCount++;
                    log.atInfo().log(String.format("attempt %s - %s - retrying...", retryCount, e.getMessage()));
                    throw new RuntimeException(e);
                }
                return client;
            });
        } else {
            log.atInfo().log("item change notifications are disabled");
        }
    }

    private String getConnectionURI() {
        return String.format("tcp://%s:%s", cfg.getEventsServerHost(), cfg.getEventsServerPort());
    }

    @Override
    public boolean isReady() {
        return publisher != null && publisher.isConnected();
    }

    @Override
    public void notify(char notifyType, char changeType, ItemData item) {
        try {
            if (isReady()) {
                ItemChanged itemChanged = new ItemChanged(notifyType, changeType, item);
                // create the message
                MqttMessage msg = new MqttMessage(itemChanged.getBytes());
                // set exactly one semantics (loss is not acceptable and subscribers cannot handle duplicates)
                msg.setQos(2);
                // If true, each client that subscribes to a topic pattern that matches the topic of the retained message receives
                // the retained message immediately after they subscribe.
                // The broker stores only one retained message per topic.
                msg.setRetained(false);
                // publish the message
                publisher.publish(itemChanged.getTopicName(), msg);
                log.atDebug().log("message '%s' sent to topic '%s' ", itemChanged.toString(), itemChanged.getTopicName());
            } else {
                log.atWarn().log("broker service not ready to publish the message, the message will be discarded");
            }
        } catch (MqttException ex) {
            throw new RuntimeException(ex);
        }
    }
}


