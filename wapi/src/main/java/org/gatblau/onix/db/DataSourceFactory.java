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

package org.gatblau.onix.db;

import com.zaxxer.hikari.HikariDataSource;
import org.gatblau.onix.conf.Config;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.sql.Connection;
import java.sql.SQLException;

@Service
public class DataSourceFactory {
    private HikariDataSource ds;
    private Connection conn;

    private final Config cfg;

    public DataSourceFactory(Config cfg) {
        this.cfg = cfg;
    }

    private HikariDataSource instance() {
        if (ds == null) {
            ds = new HikariDataSource();
            System.out.println(String.format("JDBC ==> Setting JDBC URL to: '%s'", cfg.getConnString()));
            ds.setJdbcUrl(cfg.getConnString());
            ds.setUsername(cfg.getDbuser());
            ds.setPassword(cfg.getDbpwd());
            ds.setPoolName("onix-connection-pool");
            ds.addDataSourceProperty("cachePrepStmts", cfg.isCachePrepStmts());
            ds.addDataSourceProperty("prepStmtCacheSize", cfg.getPrepStmtCacheSize());
            ds.addDataSourceProperty("prepStmtCacheSqlLimit", cfg.getPrepStmtCacheSqlLimit());
            ds.addDataSourceProperty("useServerPrepStmts", cfg.isUseServerPrepStmts());
        }
        return ds;
    }

    public Connection getConn() {
        try {
            if (conn == null || conn.isClosed()) {
                conn = instance().getConnection();
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