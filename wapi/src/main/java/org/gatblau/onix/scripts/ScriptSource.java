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

package org.gatblau.onix.scripts;

import org.gatblau.onix.FileUtil;
import org.gatblau.onix.conf.Info;
import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;

import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * represents a source of database deployment / upgrade scripts
 * it is extended by a local and remote sources for binary resource and github storage respectively
 * the resolution of the scripts version is as follows:
 * - Lookup application version in the app resources
 * - Retrieve local/remote app manifest for application version
 * - Get the db version from the app manifest
 * - Retrieve local/remote db manifest for db version
 * - Retrieve local/remote db scripts using db version
 */
public class ScriptSource {
    Logger log = LoggerFactory.getLogger(ScriptSource.class);

    protected final FileUtil file;
    private Pattern versionPattern = Pattern.compile("\\d*\\.\\d*\\.\\d*");
    private Pattern commitPattern = Pattern.compile("-([^-]+)-");

    JSONParser jsonParser = new JSONParser();

    public JSONObject appManifest;

    private final Info info;

    @Autowired
    public ScriptSource(
            FileUtil file,
            Info info) {
        this.file = file;
        this.info = info;
    }

    public String getAppVersion(){
        Matcher matcher = versionPattern.matcher(info.getVersion());
        if (matcher.find()) {
            return matcher.group(0);
        }
        throw new RuntimeException("Can't find version.");
    }

    protected String getAppManifestURI(){
        throw new UnsupportedOperationException();
    };

    protected String getDbManifestURI(String dbVersion) {
        throw new UnsupportedOperationException();
    }

    protected String getDbScriptURI(String dbVersion, String scriptName) {
        return String.format("%s/db/install/%s/%s", getScriptsURI(), dbVersion, scriptName);
    }

    protected String getScriptsURI() {
        throw new UnsupportedOperationException();
    }

    protected JSONObject getJSON(String urlString, String username, String password){
        String content = getContent(urlString, username, password);
        JSONObject json = null;
        try {
            json = (JSONObject)jsonParser.parse(content.toString());
        } catch (ParseException e) {
            throw new RuntimeException("Can't parse JSON.", e);
        }
        return json;
    }

    protected String getContent(String uriString, String username, String password) {
        throw new UnsupportedOperationException();
    }

    private Map<String, String> getScript(JSONObject dbManifest, String type) {
        // uses a LinkedHashMap instead of a HashMap so as to keep keys in order
        // then applying scripts occur in the intended sequence as per manifest order
        Map<String, String> map = new LinkedHashMap<>();
        JSONArray scripts = (JSONArray)dbManifest.get(type);
        for (Object obj: scripts) {
            JSONObject schema = (JSONObject)obj;
            String scriptName = (String)schema.get("file");
            String dbScriptUrl = getDbScriptURI(dbManifest.get("release").toString(), scriptName);
            String dbScript = getContent(dbScriptUrl, null, null);
            map.put(scriptName, dbScript);
        }
        return map;
    }

    private JSONObject getDbManifest() {
        // gets the db version is supposed to apply to the app version
        String dbVersion = getAppManifest().get("db").toString();
        // works out the url of the db manifest
        String dbManifestUrl = getDbManifestURI(dbVersion);
        // fetches the db manifest from the remote url
        return getJSON(dbManifestUrl, null, null);
    }

    protected String getCommit(){
        Matcher matcher = commitPattern.matcher(info.getVersion());
        if (matcher.find()) {
            String v = matcher.group(0);
            return v.substring(1, v.length()-1);
        }
        throw new RuntimeException("Can't find version.");
    }

    public Map<String, Map<String, String>> getDbScripts() {
        log.info(String.format("Fetching database deployment scripts from %s.", getSource()));
        Map<String, Map<String, String>> scripts = new HashMap<>();
        JSONObject dbManifest = getDbManifest();
        scripts.put("schemas", getScript(dbManifest, "schemas"));
        scripts.put("functions", getScript(dbManifest, "functions"));
        scripts.put("upgrade", getScript(dbManifest, "upgrade"));
        return scripts;
    }

    public String getSource(){
        throw new UnsupportedOperationException();
    };

    public JSONObject getAppManifest(){
        if (appManifest == null) {
            appManifest = getJSON(getAppManifestURI(), null, null);
        }
        return appManifest;
    }
}
