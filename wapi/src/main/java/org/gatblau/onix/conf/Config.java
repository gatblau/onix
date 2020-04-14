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
package org.gatblau.onix.conf;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

/*
  Abstracts Spring configuration
 */
@Service
public class Config {
    @Value("${spring.datasource.username}")
    private String dbuser;

    @Value("${spring.datasource.password}")
    private String dbpwd;

    @Value("${spring.datasource.hikari.data-source-properties.cachePrepStmts}")
    private boolean cachePrepStmts;

    @Value("${spring.datasource.hikari.data-source-properties.prepStmtCacheSize}")
    private int prepStmtCacheSize;

    @Value("${spring.datasource.hikari.data-source-properties.prepStmtCacheSqlLimit}")
    private int prepStmtCacheSqlLimit;

    @Value("${spring.datasource.hikari.data-source-properties.useServerPrepStmts}")
    private boolean useServerPrepStmts;

    @Value("${spring.datasource.url}")
    private String connString;

    @Value("${wapi.ek.1}")
    private char[] key1;

    @Value("${wapi.ek.2}")
    private char[] key2;

    @Value("${wapi.ek.expiry.date}")
    private String keyExpiry;

    @Value("${wapi.ek.default}")
    private short keyDefault;

    @Value("${wapi.amqp.host}")
    private String amqpHost;

    @Value("${wapi.amqp.port}")
    private int amqpPort;

    @Value("${wapi.amqp.incomingWindowSize}")
    private int amqpIncomingWindowSize;

    @Value("${wapi.amqp.outgoingWindowSize}")
    private int amqpOutgoingWindowSize;

    public String getDbuser() {
        return dbuser;
    }

    public String getDbpwd() {
        return dbpwd;
    }

    public boolean isCachePrepStmts() {
        return cachePrepStmts;
    }

    public int getPrepStmtCacheSize() {
        return prepStmtCacheSize;
    }

    public int getPrepStmtCacheSqlLimit() {
        return prepStmtCacheSqlLimit;
    }

    public boolean isUseServerPrepStmts() {
        return useServerPrepStmts;
    }

    public String getConnString() {
        return connString;
    }

    public char[] getKey1() {
        return key1;
    }

    public char[] getKey2() {
        return key2;
    }

    public String getKeyExpiry() {
        return keyExpiry;
    }

    public short getKeyDefault() {
        return keyDefault;
    }

    public String getAmqpHost() {
        return amqpHost;
    }

    public int getAmqpPort() {
        return amqpPort;
    }

    public int getAmqpIncomingWindowSize() {
        return amqpIncomingWindowSize;
    }

    public int getAmqpOutgoingWindowSize() {
        return amqpOutgoingWindowSize;
    }
}
