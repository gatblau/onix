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

import org.apache.commons.codec.binary.Base64;
import org.gatblau.onix.FileUtil;
import org.gatblau.onix.conf.Info;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;

/**
 * read database scripts from a remote git repository exposed via http
 */
@Service
public class RemoteSource extends ScriptSource {

    private String scriptsUrl;

    public RemoteSource(@Value("${database.scripts}") String scriptsUrl,
                        FileUtil file,
                        Info info) {
        super(file, info);
        this.scriptsUrl = scriptsUrl;
    }

    protected String getAppManifestURI(){
        return String.format("%s/app/%s.json", getScriptsURI(), getAppVersion());
    }

    protected String getDbManifestURI(String dbVersion) {
        return String.format("%s/db/install/%s/db.json", getScriptsURI(), dbVersion);
    }

    protected String getScriptsURI() {
        // if the url has a commit variable the replace it with the commit number of the app
        if (scriptsUrl.contains("<app_commit>")) {
            scriptsUrl = scriptsUrl.replace("<app_commit>", getCommit());
        }
        return scriptsUrl;
    }

    protected String getContent(String urlString, String username, String password) {
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

    @Override
    public String getSource() {
        return String.format("url: %s", scriptsUrl);
    }
}
