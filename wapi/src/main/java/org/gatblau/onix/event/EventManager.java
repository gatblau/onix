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

import com.swiftmq.amqp.AMQPContext;
import com.swiftmq.amqp.v100.client.*;
import com.swiftmq.amqp.v100.messaging.AMQPMessage;
import org.gatblau.onix.conf.Config;
import org.springframework.stereotype.Service;

/*
  Manages publication of change notification events to an AMQP message broker
 */
@Service
public class EventManager {
    // initialises thread pool and tracing facilities - context must always be created in Client mode
    private AMQPContext ctx = new AMQPContext(AMQPContext.CLIENT);
//    private final Connection connection;
//    private final Session session;
    private final Config cfg;

    public EventManager(Config cfg) throws ConnectionClosedException, SessionHandshakeException {
        this.cfg = cfg;
        // creates a connection without SASL
//        this.connection = new Connection(ctx, cfg.getAmqpHost(), cfg.getAmqpPort(), false);
//        this.session = connection.createSession(cfg.getAmqpIncomingWindowSize(), cfg.getAmqpOutgoingWindowSize());
    }

    /**
     * Send a change notification event
     * @param target queue / exchange name
     * @throws AMQPException
     */
    public void send(String target) throws AMQPException {
//        AMQPMessage msg = new AMQPMessage();
//        Producer p = session.createProducer(target, QoS.EXACTLY_ONCE);
//        p.send(msg); // Always asynchronously
//        p.close(); // Settlement completed after this call returns
    }
}
