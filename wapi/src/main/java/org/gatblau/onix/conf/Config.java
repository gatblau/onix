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

import java.io.Serializable;

import static org.gatblau.onix.conf.Config.AuthMode.*;

/*
  Abstracts Spring configuration
 */
@Service
public class Config implements Serializable {

    public boolean isSwaggerEnabled() {
        return swaggerEnabled;
    }

    public void setSwaggerEnabled(boolean swaggerEnabled) {
        this.swaggerEnabled = swaggerEnabled;
    }

    public boolean isCsrfEnabled() {
        return csrfEnabled;
    }

    public boolean isSmtpSMTPS() {
        return smtpSMTPS;
    }

    public enum AuthMode {
        Basic,
        OIDC,
        None,
    }

    @Value("${wapi.swagger.enabled}")
    private boolean swaggerEnabled;

    @Value("${wapi.csrf.enabled}")
    private boolean csrfEnabled;

    @Value("${wapi.auth.mode}")
    private String authMode;
    
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

    @Value("${wapi.events.server.host}")
    private String eventsServerHost;

    @Value("${wapi.events.server.port}")
    private int eventsServerPort;

    @Value("${wapi.events.server.user}")
    private String eventsServerUser;

    @Value("${wapi.events.server.pwd}")
    private char[] eventsServerPwd;

    @Value("${wapi.events.enabled}")
    private boolean eventsEnabled;

    @Value("${wapi.events.client.retries}")
    private int eventsClientRetries;

    @Value("${wapi.events.client.backoffperiod}")
    private long eventsClientBackOffPeriod;
    
    @Value("${wapi.jwt.secret}")
    private String jwtSecret;

    @Value("${wapi.smtp.auth}")
    private boolean smtpAuth;
    
    @Value("${wapi.smtp.starttls.enable}")
    private boolean smtpStartTLS;

    @Value("${wapi.smtp.smtps.enabled}")
    private boolean smtpSMTPS;

    @Value("${wapi.smtp.host}")
    private String smtpHost;

    @Value("${wapi.smtp.port}")
    private int smtpPort;

    @Value("${wapi.smtp.from.user}")
    private String smtpFromUser;

    @Value("${wapi.smtp.from.pwd}")
    private String smtpFromPwd;

    @Value("${wapi.smtp.pwd.reset.uri}")
    private String smtpPwdResetURI;

    @Value("${wapi.smtp.pwd.setup.uri}")
    private String smtpPwdSetupURI;

    @Value("${wapi.smtp.enabled}")
    private boolean smtpEnabled;

    @Value("${wapi.smtp.pwd.reset.tokenexpiry}")
    private long smtpPwdResetTokenExpirySecs;

    @Value("${wapi.pwd.len}")
    private int pwdLen;
    
    @Value("${wapi.pwd.upper}")
    private int pwdUpper;
    
    @Value("${wapi.pwd.lower}")
    private int pwdLower;
    
    @Value("${wapi.pwd.digits}")
    private int pwdDigits;
    
    @Value("${wapi.pws.specialchars}")
    private int pwdSpecialChars;

    @Value("${database.server.url}")
    private String dbServerUrl;

    @Value("${database.name}")
    private String dbName;

    @Value("${spring.datasource.username}")
    private String dbUser;

    @Value("${spring.datasource.password}")
    private char[] dbPwd;

    @Value("${database.admin.pwd}")
    private char[] dbAdminPwd;

    public long getSmtpPwdResetTokenExpirySecs() {
        return smtpPwdResetTokenExpirySecs;
    }

    public boolean isSmtpEnabled() {
        return smtpEnabled;
    }

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

    public String getEventsServerHost() {
        return eventsServerHost;
    }

    public int getEventsServerPort() {
        return eventsServerPort;
    }

    public String getEventsServerUser() {
        return eventsServerUser;
    }

    public String getEventsServerPwd() {
        return new String(eventsServerPwd);
    }

    public boolean isEventsEnabled() {
        return eventsEnabled;
    }

    public int getEventsClientRetries() {
        return eventsClientRetries;
    }

    public long getEventsClientBackOffPeriod() {
        return eventsClientBackOffPeriod;
    }

    public AuthMode getAuthMode() {
        switch (authMode.toLowerCase()) {
            case "none":
                return None;
            case "oidc":
                return OIDC;
            case "basic":
            default:
                return Basic;
        }
    }

    public String getJwtSecret() {
        return jwtSecret;
    }

    public boolean isSmtpAuth() {
        return smtpAuth;
    }

    public boolean isSmtpStartTLS() {
        return smtpStartTLS;
    }

    public String getSmtpHost() {
        return smtpHost;
    }

    public int getSmtpPort() {
        return smtpPort;
    }

    public String getSmtpFromUser() {
        return smtpFromUser;
    }

    public String getSmtpFromPwd() { return smtpFromPwd; }

    public String getSmtpPwdResetURI() {
        return smtpPwdResetURI;
    }

    public String getSmtpPwdSetupURI() {
        return smtpPwdSetupURI;
    }

    public int getPwdLen() {
        return pwdLen;
    }

    public int getPwdUpper() {
        return pwdUpper;
    }

    public int getPwdLower() {
        return pwdLower;
    }

    public int getPwdSpecialChars() {
        return pwdSpecialChars;
    }

    public int getPwdDigits() {
        return pwdDigits;
    }

    public String getDbServerUrl() {
        return dbServerUrl;
    }

    public String getDbName() {
        return dbName;
    }

    public String getDbUser() {
        return dbUser;
    }

    public char[] getDbPwd() {
        return dbPwd;
    }

    public char[] getDbAdminPwd() {
        return dbAdminPwd;
    }
}
