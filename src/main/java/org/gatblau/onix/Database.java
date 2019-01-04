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
import org.springframework.stereotype.Service;

import java.sql.*;

@Service
class Database {
    private Connection conn;
    private PreparedStatement stmt;
    private String resultKey;

    @Autowired
    private DataSourceFactory ds;

    public Database() {
    }

    Connection createConnection() throws SQLException {
        conn = ds.instance().getConnection();
        return conn;
    }

    void prepare(String sql) throws SQLException {
        if (conn == null) {
            createConnection();
        }
        stmt = conn.prepareStatement(sql);
    }

    void setString(int parameterIndex, String value) throws SQLException {
        stmt.setString(parameterIndex, value);
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
        String result = null;
        ResultSet set = stmt.executeQuery();
        if (set.next()) {
            return set.getString(query_name);
        }
        throw new RuntimeException(String.format("Failed to execute query '%s'", query_name));
    }

    ResultSet executeQuery() throws SQLException {
        String result = null;
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
            stmt.close();
            conn.close();
            stmt = null;
            conn = null;
        }
        catch (Exception ex) {
            System.out.println("WARNING: failed to close database connection.");
            ex.printStackTrace();
        }
    }
}
