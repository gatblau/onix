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

import org.apache.commons.codec.binary.Base64;
import org.json.simple.JSONArray;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.*;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/*
  retrieves database deployment scripts from an online repository
  the resolution of the scripts version is as follows:
  - Lookup application version in the app resources
  - Retrieve online app manifest for application version
  - Get the db version from the app manifest
  - Retrieve db scripts using db version
 */
@Service
public class ScriptSource {
    Logger log = LoggerFactory.getLogger(ScriptSource.class);

    @Value("${database.scripts}")
    String scriptsUrl;

    private Pattern versionPattern = Pattern.compile("\\d*\\.\\d*\\.\\d*");

    private JSONParser jsonParser = new JSONParser();

    private JSONObject appManifest;

    @Autowired
    private Info info;

    public ScriptSource() {
    }

    private String getAppManifestUrl(){
        return String.format("%s/%s.json", scriptsUrl, getAppVersion());
    }

    private String getDbManifestUrl(String dbVersion) {
        return String.format("%s/%s/manifest.json", scriptsUrl, dbVersion);
    }

    private String getDbScriptUrl(String dbVersion, String scriptName) {
        return String.format("%s/%s/%s", scriptsUrl, dbVersion, scriptName);
    }

    String getAppVersion(){
        Matcher matcher = versionPattern.matcher(info.getVersion());
        if (matcher.find()) {
            return matcher.group(0);
        }
        throw new RuntimeException("Can't find version.");
    }

    private String getContent(String urlString, String username, String password) {
        try {
            URL url = new URL(urlString);
            HttpURLConnection conn = (HttpURLConnection) url.openConnection();
            conn.setRequestMethod("GET");
            conn.setConnectTimeout(10000);
            conn.setReadTimeout(10000);
            // try disable caching
            conn.setRequestProperty("Cache-Control", "no-cache");
            conn.setRequestProperty("Pragma", "no-cache");
            conn.setUseCaches(false);
            // add authorization header?
            if (username != null && password != null) {
                String userCredentials = String.format("%s:%s", username, password);
                String basicAuth = "Basic " + new String(Base64.encodeBase64(userCredentials.getBytes()));
                conn.setRequestProperty("Authorization", basicAuth);
            }
            BufferedReader in = new BufferedReader(new InputStreamReader(conn.getInputStream()));
            String inputLine;
            StringBuffer content = new StringBuffer();
            while ((inputLine = in.readLine()) != null) {
                content.append(inputLine).append("\n");
            }
            in.close();
            conn.disconnect();
            return content.toString();
        } catch (Exception e) {
            throw new RuntimeException("Can't retrieve content.", e);
        }
    }

    private JSONObject getJSON(String urlString, String username, String password){
        String content = getContent(urlString, username, password);
        JSONObject json = null;
        try {
            json = (JSONObject)jsonParser.parse(content.toString());
        } catch (ParseException e) {
            throw new RuntimeException("Cant parse JSON.", e);
        }
        return json;
    }

    JSONObject getAppManifest() {
        if (appManifest == null) {
            appManifest = getJSON(getAppManifestUrl(), null, null);
        }
        return appManifest;
    }

    private JSONObject getDbManifest() {
        // fetches the app version manifest
        JSONObject appManifest = getAppManifest();
        // gets the db version is supposed to apply to the app version
        String dbVersion = appManifest.get("db").toString();
        // works out the url of the db manifest
        String dbManifestUrl = getDbManifestUrl(dbVersion);
        // fetches the db manifest from the remote url
        return getJSON(dbManifestUrl, null, null);
    }

    private Map<String, String> getScript(JSONObject dbManifest, String type) {
        // uses a LinkedHashMap instead of a HashMap so as to keep keys in order
        // then applying scripts occur in the intended sequence as per manifest order
        Map<String, String> map = new LinkedHashMap<>();
        JSONArray scripts = (JSONArray)dbManifest.get(type);
        for (Object obj: scripts) {
            JSONObject schema = (JSONObject)obj;
            String scriptName = (String)schema.get("file");
            String dbScriptUrl = getDbScriptUrl(dbManifest.get("release").toString(), scriptName);
            String dbScript = getContent(dbScriptUrl, null, null);
            map.put(scriptName, dbScript);
        }
        return map;
    }

    public Map<String, Map<String, String>> getDbScripts() {
        log.info(String.format("Fetching database deployment scripts from %s.", scriptsUrl));
        Map<String, Map<String, String>> scripts = new HashMap<>();
        JSONObject dbManifest = getDbManifest();
        scripts.put("schemas", getScript(dbManifest, "schemas"));
        scripts.put("functions", getScript(dbManifest, "functions"));
        return scripts;
    }
}
