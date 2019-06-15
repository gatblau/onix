package org.gatblau.onix;

import org.apache.commons.codec.binary.Base64;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@Service
public class ScriptSource {
    @Value("${database.scripts}")
    private String scriptsUrl;

    private Pattern versionPattern = Pattern.compile("\\d*\\.\\d*\\.\\d*");

    private JSONParser jsonParser = new JSONParser();

    @Autowired
    private Info info;

    public ScriptSource() {
    }

    private String getManifestUrl(){
        return String.format("%s/%s.json", scriptsUrl, getVersion());
    }

    private String getVersion(){
        Matcher matcher = versionPattern.matcher(info.getVersion());
        if (matcher.find()) {
            return matcher.group(0);
        }
        throw new RuntimeException("Can't find version.");
    }

    private JSONObject getJSON(String urlString, String username, String password){
        try {
            URL url = new URL(urlString);
            HttpURLConnection conn = (HttpURLConnection) url.openConnection();
            conn.setRequestMethod("GET");
            conn.setConnectTimeout(10000);
            conn.setReadTimeout(10000);
            // add authorization header?
            if (username != null && password != null) {
                String userCredentials = String.format("%s:%s", username, password);
                String basicAuth = "Basic " + new String(Base64.encodeBase64(userCredentials.getBytes()));
                conn.setRequestProperty ("Authorization", basicAuth);
            }
            BufferedReader in = new BufferedReader(new InputStreamReader(conn.getInputStream()));
            String inputLine;
            StringBuffer content = new StringBuffer();
            while ((inputLine = in.readLine()) != null) {
                content.append(inputLine);
            }
            in.close();
            conn.disconnect();
            String response = content.toString();
            return (JSONObject)jsonParser.parse(content.toString());

        } catch (Exception e) {
           throw new RuntimeException("Can't retrieve manifest.", e);
        }
    }

    public JSONObject getManifest() {
        return getJSON(getManifestUrl(), null, null);
    }
}
