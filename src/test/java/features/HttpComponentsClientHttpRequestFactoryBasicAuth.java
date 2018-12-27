package features;

import org.apache.http.HttpHost;
import org.apache.http.client.AuthCache;
import org.apache.http.client.protocol.HttpClientContext;
import org.apache.http.impl.auth.BasicScheme;
import org.apache.http.impl.client.BasicAuthCache;
import org.apache.http.protocol.BasicHttpContext;
import org.apache.http.protocol.HttpContext;
import org.springframework.http.HttpMethod;
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory;

import java.net.URI;

/*
    The RestTemplate requires an HTTP request factory that supports Basic Authentication.
    To support basic authentication, subclasses HttpComponentsClientHttpRequestFactory and
    overrides the createHttpContext method.
 */
public class HttpComponentsClientHttpRequestFactoryBasicAuth extends HttpComponentsClientHttpRequestFactory {

    HttpHost host;

    public HttpComponentsClientHttpRequestFactoryBasicAuth(HttpHost host) {
        super();
        this.host = host;
    }

    /*
        the creation of the HttpContext containing the basic authentication support built in.
        the authentication info is cached and the process of setting up this authentication cache
        is very manual and unintuitive.
     */
    protected HttpContext createHttpContext(HttpMethod httpMethod, URI uri) {
        AuthCache authCache = new BasicAuthCache();
        authCache.put(host, new BasicScheme());
        BasicHttpContext localcontext = new BasicHttpContext();
        localcontext.setAttribute(HttpClientContext.AUTH_CACHE, authCache);
        return localcontext;
    }
}
