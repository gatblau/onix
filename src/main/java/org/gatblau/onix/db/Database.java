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

import org.gatblau.onix.FileUtil;
import org.gatblau.onix.scripts.ScriptSource;
import org.gatblau.onix.scripts.ScriptSourceFactory;
import org.json.simple.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Scope;
import org.springframework.context.annotation.ScopedProxyMode;
import org.springframework.stereotype.Service;
import org.springframework.web.context.WebApplicationContext;

import java.sql.*;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Service
@Scope(value = WebApplicationContext.SCOPE_REQUEST, proxyMode = ScopedProxyMode.TARGET_CLASS)
class Database {
    Logger log = LoggerFactory.getLogger(Database.class);

    private PreparedStatement stmt;

    private final ScriptSource script;

    private final DataSourceFactory ds;

    private final FileUtil file;

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

    @Value("${database.auto.deploy}")
    private boolean dbAutoDeploy;

    @Value("${database.auto.upgrade}")
    private boolean dbAutoUpgrade;

    private Version version;

    @Autowired
    public Database(ScriptSourceFactory selector, DataSourceFactory ds, FileUtil file) {
        // switches to the configured script source
        this.script = selector.source;
        this.ds = ds;
        this.file = file;
    }

    class Version {
        String app;
        String db;

        @Override
        public String toString() {
            if (app != null && db != null) {
                return String.format("App Version: '%s'; Db Version: '%s'.", app, db);
            } else {
                return super.toString();
            }
        }
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

    void createDb() throws SQLException {
        String ap = new String(dbAdminPwd);
        Map<String, String> vars = new HashMap<>();
        vars.put("<DB_NAME>", dbName);
        vars.put("<DB_USER>", dbUser);
        vars.put("<DB_PWD>", new String(dbPwd));
        // creates the database and db user as postgres user
        log.info(String.format("Creating database '%s' and user '%s'.", dbName, dbUser));
        runScriptFromResx(String.format("%s/postgres", dbServerUrl), "postgres", ap, "db/init/db_and_user.sql", vars);
        // creates the extensions in onix db as postgres user
        log.info(String.format("Creating extensions in database '%s'.", dbName));
        runScriptFromResx(String.format("%s/%s", dbServerUrl, dbName), "postgres", ap, "db/init/extensions.sql", null);
        log.info(String.format("Creating version control table in database '%s'.", dbName));
        runScriptFromResx(String.format("%s/%s", dbServerUrl, dbName), "postgres", ap, "db/init/version_table.sql", null);
    }

    private void runScriptFromResx(String dbServerUrl, String user, String pwd, String script, Map<String, String> vars) throws SQLException {
        Connection conn = DriverManager.getConnection(dbServerUrl, user, pwd);
        Statement stmt = conn.createStatement();
        final List<String> msg = Arrays.asList(file.getFile(script));
        if (vars != null) {
            vars.forEach((key, value) -> msg.set(0, msg.get(0).replace(key, value)));
        }
        stmt.execute(msg.get(0));
        stmt.close();
        conn.close();
    }

    private void runScriptFromString(String adminPwd, String script, String targetDb) throws SQLException {
        Connection conn = DriverManager.getConnection(String.format("%s/%s", dbServerUrl, targetDb), "postgres", adminPwd);
        Statement stmt = conn.createStatement();
        stmt.execute(script);
        stmt.close();
        conn.close();
    }

    private void deployScripts(Map<String, String> scripts, String adminPwd) {
        for (Map.Entry<String, String> script: scripts.entrySet()) {
            try {
                log.info(String.format("Executing script '%s'.", script.getKey()));
                runScriptFromString(adminPwd, script.getValue(), dbName);
            } catch (SQLException e) {
                throw new RuntimeException(String.format("Failed to apply script '%s'.", script.getKey()), e);
            }
        }
    }

    public void deployDb() throws SQLException {
        if (dbAutoDeploy) {
            // retrieve the relevant db scripts to be deployed before doing anything else
            Map<String, Map<String, String>> scripts = script.getDbScripts();

            // creates a local variable pwd that should go out of scope at the end of the scope and
            // be GC by the JVM
            String ap = new String(dbAdminPwd);

            // deploys the schemas first
            Map<String, String> schemas = scripts.get("schemas");
            deployScripts(schemas, ap);

            // deploys the functions
            Map<String, String> funcs = scripts.get("functions");
            deployScripts(funcs, ap);

            // updates the version table
            setVersion(script.getAppVersion(), script.appManifest.get("db").toString(), "Database automatically deployed by Onix", script.getSource());

            // resets the version in memory
            version = null;
        } else {
            throw new RuntimeException(String.format(
                    "Database does not exists.\n" +
                    "Database automatic deployment is disabled.\n" +
                    "Deploy database manually or enable auto deployment in the configuration."));
        }
    }

    /**
     * determines if the onix database exists
     * @return true if the onix database exists, otherwise false
     */
    public boolean exists() {
        boolean exists = false;
        Connection conn = null;
        try {
            conn = DriverManager.getConnection(String.format("%s/postgres", dbServerUrl), "postgres", new String(dbAdminPwd));
            Statement stmt = conn.createStatement();
            if (stmt.execute(String.format("SELECT 1 from pg_database WHERE datname='%s';", dbName))){
                ResultSet set = stmt.getResultSet();
                exists = set.next();
                if (exists) {
                    log.info(String.format("Database %s already exists.", dbName));
                } else {
                    log.info(String.format("Database %s does not exist.", dbName));
                }
            }
            stmt.close();
            conn.close();
        } catch (SQLException e) {
            throw new RuntimeException("Failed to check if database exists.", e);
        }
        return exists;
    }

    /**
     * gets the version information from the database
     * @return a varsion instance containing app and db versions
     */
    Version getVersion() {
        if (version == null) {
            version = new Version();
            Connection conn = null;
            try {
                conn = DriverManager.getConnection(String.format("%s/%s", dbServerUrl, dbName), "postgres", new String(dbAdminPwd));
                Statement stmt = conn.createStatement();
                if (stmt.execute(String.format("SELECT * from version ORDER BY time DESC LIMIT 1;", dbName))) {
                    ResultSet set = stmt.getResultSet();
                    if (set.next()) {
                        version.app = set.getString("application_version");
                        version.db = set.getString("database_version");
                    }
                }
                stmt.close();
                conn.close();
            } catch (SQLException e) {
                throw new RuntimeException("Failed to retrieve version from database", e);
            }
        }
        return version;
    }

    /**
     * determines whether the db should be upgraded based on the app version running
     * @return if 0, then should not upgrade
     * if 1, then it should upgrade db
     * if -1 then it should upgrade app, can't run on the current db version
     */
    public int shouldUpgrade() {
        log.info("Verifying database and application compatibility.");
        Version v = getVersion();
        JSONObject appManifest = script.getAppManifest();
        int dbv = Integer.parseInt(v.db); // the current db version
        int appdbv = Integer.parseInt(appManifest.get("db").toString()); // the db version required by the app
        if (dbv < appdbv){
            // version in the database is less than the one the app requires
            // should upgrade
            log.warn(String.format("Database version '%s' should be upgraded to version '%s' to meet the requirements of the application.", dbv, appdbv));
            return 1;
        }
        else if (dbv == appdbv) {
            // should not upgrade
            log.info(String.format("Database version '%s' is correct for the running application version '%s'.", v.db, v.app));
            return 0;
        }
        else {
            // app cannot run, db version is newer than required by the app!
            log.error(String.format(
                "The application cannot start. \n" +
                "The current database version (i.e. %s), is newer than the version required by the running application (i.e. %s).\n" +
                "The application should be upgraded for the current database version.", dbv, appdbv));
            return -1;
        }
    }

    /**
     * upgrades the database to the version indicated by the current application
     */
    public void upgrade() {
        throw new RuntimeException("Upgrade process not implemented.");
    }

    private void setVersion(String appVer, String dbVer, String desc, String scriptSrc) {
        log.info(String.format("Recording version of database installed: %s:%s.", appVer, dbVer));
        try {
            prepare("SELECT set_version(" +
                    "?::character varying," +
                    "?::character varying," +
                    "?::text," +
                    "?::character varying" +
                    ")");
            setString(1, appVer);
            setString(2, dbVer);
            setString(3, desc);
            setString(4,scriptSrc);
            executeQuery();
        } catch (Exception ex) {
            throw new RuntimeException("Failed to set version in database.", ex);
        }
    }

    /**
     * deletes the database.
     */
    void deleteDb() {
        try {
            // kills all existing connections and drops the database
            runScriptFromString(new String(dbAdminPwd),
                String.format(
                    "SELECT pid, pg_terminate_backend(pid) \n" +
                    "FROM pg_stat_activity \n" +
                    "WHERE datname = '%s' AND pid <> pg_backend_pid();\n" +
                    "DROP DATABASE IF EXISTS %s;\n" +
                    "DROP USER %s", dbName, dbName, dbUser), "postgres");
        } catch (Exception e) {
            log.warn(String.format("Failed to drop database '%s' after deployment failure: %s.", dbName, e.getMessage()));
        }
    }
}
