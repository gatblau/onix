package features;

import org.apache.http.HttpHost;
import org.apache.tomcat.util.codec.binary.Base64;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.web.server.LocalServerPort;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.client.support.BasicAuthenticationInterceptor;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.web.client.RestTemplate;

import java.nio.charset.Charset;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@ContextConfiguration(classes= Config.class)
public class BaseTest implements InitializingBean {
    @Autowired
    protected Util util;

    @LocalServerPort
    protected int port;

    protected RestTemplate client;

    public void afterPropertiesSet()  {
        try {
            HttpHost host = new HttpHost("localhost", port, "http");
            client = new RestTemplate(new HttpComponentsClientHttpRequestFactoryBasicAuth(host));
            client.getInterceptors().add(new BasicAuthenticationInterceptor(util.adminUsername, util.adminPassword));
        }
        catch (Exception ex) {
        }
    }

    protected HttpEntity<String> getEntity() {
        return getEntity(null);
    }

    protected HttpEntity<String> getEntity(String payload) {
        HttpHeaders headers = new HttpHeaders();
        headers.add("Authorization", getAuthHeaderValue());
        headers.add("Content-Type", "application/json");
        headers.add("Accept", "application/json");
        if (payload != null) {
            return new HttpEntity<String>(payload, headers);
        }
        return new HttpEntity<String>(headers);
    }

    protected HttpEntity<String> getEntityFromKey(String payloadKey) {
        String payload = util.get(payloadKey);
        return getEntity(payload);
    }

    private String getAuthHeaderValue() {
        String auth = util.adminUsername + ":" + util.adminPassword;
        byte[] encodedAuth = Base64.encodeBase64(auth.getBytes(Charset.forName("US-ASCII")));
        return "Basic " + new String( encodedAuth );
    }
}