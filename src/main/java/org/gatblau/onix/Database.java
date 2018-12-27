package org.gatblau.onix;

import org.springframework.stereotype.Service;

import java.sql.Connection;
import java.sql.DriverManager;

@Service
public class Database {
    public Database() {
    }

    public Connection createConnection(){
        Connection conn = null;
        try {
            Class.forName("org.postgresql.Driver");
            conn = DriverManager.getConnection(getConnString(), getDbUser(), getDbPwd());
        }
        catch(Exception ex) {
            ex.printStackTrace();
        }
        return conn;
    }

    private String getConnString() {
        String connStr = System.getenv("ONIX_DB_CONN_STRING");
        if (connStr == null) {
            connStr = "jdbc:postgresql://localhost:5432/onix";
        }
        return connStr;
    }

    private String getDbUser(){
        String user = System.getenv("ONIX_DB_USER");
        if (user == null) {
            user = "onix";
        }
        return user;
    }

    private String getDbPwd(){
        String pwd = System.getenv("ONIX_DB_PWD");
        if (pwd == null) {
            pwd = "onix";
        }
        return pwd;
    }
}
