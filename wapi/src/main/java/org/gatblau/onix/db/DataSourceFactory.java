/*
Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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

package org.gatblau.onix.db;

import com.zaxxer.hikari.HikariDataSource;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.sql.Connection;
import java.sql.SQLException;

@Service
public class DataSourceFactory {
    private HikariDataSource ds;
    private Connection conn;

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

    private HikariDataSource instance() {
        if (ds == null) {
            ds = new HikariDataSource();
            System.out.println(String.format("JDBC ==> Setting JDBC URL to: '%s'", connString));
            ds.setJdbcUrl(connString);
            ds.setUsername(dbuser);
            ds.setPassword(dbpwd);
            ds.setPoolName("onix-connection-pool");
            ds.addDataSourceProperty("cachePrepStmts", cachePrepStmts);
            ds.addDataSourceProperty("prepStmtCacheSize", prepStmtCacheSize);
            ds.addDataSourceProperty("prepStmtCacheSqlLimit", prepStmtCacheSqlLimit);
            ds.addDataSourceProperty("useServerPrepStmts", useServerPrepStmts);
        }
        return ds;
    }

    public Connection getConn() {
        try {
            if (conn == null || conn.isClosed()) {
                try {
                    conn = instance().getConnection();
                } catch (SQLException e) {
                    throw new RuntimeException(e);
                }
            }
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
        return conn;
    }

    public void closeConn() {
        try {
            if (conn != null || !conn.isClosed()) {
                conn.close();
            }
        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
    }

}