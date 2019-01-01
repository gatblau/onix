package org.gatblau.onix;

import com.zaxxer.hikari.HikariDataSource;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public class DataSourceFactory {
    private HikariDataSource ds;

    @Value("${spring.datasource.username}")
    private String dbuser;

    @Value("${spring.datasource.password}")
    private String dbpwd;

    @Value("${spring.datasource.url}")
    private String connString;

    public HikariDataSource instance() {
        if (ds == null) {
            ds = new HikariDataSource();
            ds = new HikariDataSource();
            ds.setJdbcUrl(connString);
            ds.setUsername(dbuser);
            ds.setPassword(dbpwd);
            ds.setDataSourceClassName("org.postgresql.ds.PGSimpleDataSource");
        }
        return ds;
    }
}
