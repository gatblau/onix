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

package org.gatblau.onix;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Scope;
import org.springframework.context.annotation.ScopedProxyMode;
import org.springframework.stereotype.Service;
import org.springframework.web.context.WebApplicationContext;

import java.io.File;
import java.io.IOException;
import java.sql.*;
import java.sql.Date;
import java.util.*;

@Service
@Scope(value = WebApplicationContext.SCOPE_REQUEST, proxyMode = ScopedProxyMode.TARGET_CLASS)
class Database {
    private PreparedStatement stmt;

    @Autowired
    private DataSourceFactory ds;

    @Value("${database.server.url}")
    private String dbServerUrl;

    @Value("${database.name}")
    private String dbName;

    @Value("${spring.datasource.username}")
    private String dbUser;

    @Value("${spring.datasource.password}")
    private String dbPwd;

    public Database() {
    }

    void prepare(String sql) throws SQLException {
        stmt = ds.getConn().prepareStatement(sql);
    }

    void setString(int parameterIndex, String value) throws SQLException {
        stmt.setString(parameterIndex, value);
    }

    void setArray(int parameterIndex, String[] value) throws SQLException {
        stmt.setArray(parameterIndex, ds.getConn().createArrayOf("varchar", value));
    }

    void setInt(int parameterIndex, Integer value) throws SQLException {
        stmt.setInt(parameterIndex, value);
    }

    void setShort(int parameterIndex, Short value) throws SQLException {
        stmt.setShort(parameterIndex, value);
    }

    void setDate(int parameterIndex, Date value) throws SQLException {
        stmt.setDate(parameterIndex, value);
    }

    void setObject(int parameterIndex, Object value) throws SQLException {
        stmt.setObject(parameterIndex, value);
    }

    void setObjectRange(int fromIndex, int toIndex, Object value) throws SQLException {
        for (int i = fromIndex; i < toIndex + 1; i++) {
            setObject(i, null);
        }
    }

    String executeQueryAndRetrieveStatus(String query_name) throws SQLException {
        ResultSet set = stmt.executeQuery();
        if (set.next()) {
            return set.getString(query_name);
        }
        throw new RuntimeException(String.format("Failed to execute query '%s'", query_name));
    }

    ResultSet executeQuery() throws SQLException {
        return stmt.executeQuery();
    }

    ResultSet executeQuerySingleRow() throws SQLException {
        String result = null;
        ResultSet set = stmt.executeQuery();
        if (set.next()) {
            return set;
        }
        throw new RuntimeException("No results found.");
    }

    boolean execute() throws SQLException {
        return stmt.execute();
    }

    void close() {
        try {
            if (stmt != null) {
                stmt.close();
                stmt = null;
            }
            ds.closeConn();
        }
        catch (Exception ex) {
            System.out.println("WARNING: failed to close database statement.");
            ex.printStackTrace();
        }
    }

    void createDb(String adminPwd) throws SQLException {
        Map<String, String> vars = new HashMap<>();
        vars.put("<DB_NAME>", dbName);
        vars.put("<DB_USER>", dbUser);
        vars.put("<DB_PWD>", dbPwd);
        // creates the database and db user as postgres user
        runScriptFromResx(String.format("%s/postgres", dbServerUrl), "postgres", adminPwd, "db/1_create_db_user.sql", vars);
        // creates the extensions in onix db as postgres user
        runScriptFromResx(String.format("%s/%s", dbServerUrl, dbName), "postgres", adminPwd, "db/2_create_ext.sql", null);
    }

    private void runScriptFromResx(String dbServerUrl, String user, String pwd, String script, Map<String, String> vars) throws SQLException {
        Connection conn = DriverManager.getConnection(dbServerUrl, user, pwd);
        Statement stmt = conn.createStatement();
        final List<String> msg = Arrays.asList(getFile(script));
        if (vars != null) {
            vars.forEach((key, value) -> msg.set(0, msg.get(0).replace(key, value)));
        }
        stmt.execute(msg.get(0));
        stmt.close();
        conn.close();
    }

    private void runScriptFromString(String adminPwd, String script) throws SQLException {
        Connection conn = DriverManager.getConnection(String.format("%s/onix", dbServerUrl), "postgres", adminPwd);
        Statement stmt = conn.createStatement();
        stmt.execute(script);
        stmt.close();
        conn.close();
    }

    void deployScripts(Map<String, String> scripts, String adminPwd) {
        for (Map.Entry<String, String> script: scripts.entrySet()) {
            try {
                runScriptFromString(adminPwd, script.getValue());
            } catch (SQLException e) {
                throw new RuntimeException(String.format("Failed to apply %s script.", script.getKey()), e);
            }
        }
    }

    private String getFile(String fileName) {
        StringBuilder result = new StringBuilder("");
        //Get file from resources folder
        ClassLoader classLoader = getClass().getClassLoader();
        File file = new File(classLoader.getResource(fileName).getFile());
        try (Scanner scanner = new Scanner(file)) {
            while (scanner.hasNextLine()) {
                String line = scanner.nextLine();
                result.append(line).append("\n");
            }

        } catch (IOException e) {
            e.printStackTrace();
        }
        return result.toString();

    }
}
