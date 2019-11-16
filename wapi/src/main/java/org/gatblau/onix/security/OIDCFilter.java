/*
Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

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

package org.gatblau.onix.security;

import com.auth0.jwk.Jwk;
import com.auth0.jwk.JwkProvider;
import com.auth0.jwk.UrlJwkProvider;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.jwt.Jwt;
import org.springframework.security.jwt.JwtHelper;
import org.springframework.security.jwt.crypto.sign.RsaVerifier;
import org.springframework.security.oauth2.client.OAuth2RestOperations;
import org.springframework.security.oauth2.client.OAuth2RestTemplate;
import org.springframework.security.oauth2.common.OAuth2AccessToken;
import org.springframework.security.oauth2.common.exceptions.OAuth2Exception;
import org.springframework.security.web.authentication.AbstractAuthenticationProcessingFilter;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.net.URL;
import java.security.interfaces.RSAPublicKey;
import java.util.Date;
import java.util.Map;

/*
    Authenticates against the specified OpenId server
 */
public class OIDCFilter extends AbstractAuthenticationProcessingFilter {
    @Value("${oidc.clientid}")
    private String clientId;

    @Value("${oidc.issuer}")
    private String issuer;

    @Value("${oidc.jwkUrl}")
    private String jwkUrl;

    public OAuth2RestOperations restTemplate;

    public OIDCFilter(String defaultFilterProcessesUrl) {
        super(defaultFilterProcessesUrl);
        setAuthenticationManager(new NoopAuthenticationManager());
    }

    @Override
    public Authentication attemptAuthentication(HttpServletRequest request, HttpServletResponse response) throws AuthenticationException {
        checkConfig("OIDC_ISSUER", issuer);
        checkConfig("OIDC_CLIENT_ID", clientId);
        checkConfig("OIDC_JWKURL", jwkUrl);
        OAuth2AccessToken accessToken;
        try {
            accessToken = restTemplate.getAccessToken();
        } catch (final OAuth2Exception e) {
            throw new BadCredentialsException("Could not obtain access token", e);
        }
        try {
            // process the openid id_token
            final String idToken = accessToken.getAdditionalInformation().get("id_token").toString();
            String kid = JwtHelper.headers(idToken).get("kid");
            final RsaVerifier verifier = verifier(kid);

            final Jwt decodedIdToken = JwtHelper.decodeAndVerify(idToken, verifier);
            final Map<String, String> authenticationInfo = new ObjectMapper().readValue(decodedIdToken.getClaims(), Map.class);
            verifyAuthenticationClaims(authenticationInfo);

            //process the oauth 2.0 access_token
            final Jwt decodedAccessToken = JwtHelper.decodeAndVerify(accessToken.getValue(), verifier);
            final Map<String, String> authorisationInfo = new ObjectMapper().readValue(decodedAccessToken.getClaims(), Map.class);
            verifyAuthorisationClaims(authorisationInfo);

            final OIDCUserDetails user = new OIDCUserDetails(authenticationInfo, authorisationInfo, accessToken);
            return new UsernamePasswordAuthenticationToken(user, null, user.getAuthorities());
        } catch (final Exception e) {
            throw new BadCredentialsException("Could not obtain user details from token", e);
        }
    }

    public void verifyAuthenticationClaims(Map claims) {
        int exp = (int) claims.get("exp");
        Date expireDate = new Date(exp * 1000L);
        Date now = new Date();
        if (expireDate.before(now) || !claims.get("iss").equals(issuer) || !claims.get("aud").equals(clientId)) {
            throw new RuntimeException("Invalid claims");
        }
    }

    public void verifyAuthorisationClaims(Map claims) {
        int exp = (int) claims.get("exp");
        Date expireDate = new Date(exp * 1000L);
        Date now = new Date();
        if (
            expireDate.before(now) ||
            !claims.get("iss").equals(issuer) ||
            claims.get("roles") == null ||
            claims.get("sub") == null
        ) {
            throw new RuntimeException("Invalid claims");
        }
    }

    private RsaVerifier verifier(String kid) throws Exception {
        JwkProvider provider = new UrlJwkProvider(new URL(jwkUrl));
        Jwk jwk = provider.get(kid);
        return new RsaVerifier((RSAPublicKey) jwk.getPublicKey());
    }

    public void setRestTemplate(OAuth2RestTemplate restTemplate2) {
        restTemplate = restTemplate2;
    }

    private static class NoopAuthenticationManager implements AuthenticationManager {
        @Override
        public Authentication authenticate(Authentication authentication) throws AuthenticationException {
            throw new UnsupportedOperationException("No authentication should be done with this AuthenticationManager");
        }
    }

    private static void checkConfig(String variableName, String value) {
        if (value == null || value.trim().length() == 0) {
            throw new RuntimeException(String.format("Variable '%s' not specified.", variableName));
        }
    }
}

