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

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.authentication.builders.AuthenticationManagerBuilder;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.builders.WebSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.client.OAuth2RestTemplate;
import org.springframework.security.oauth2.client.filter.OAuth2ClientContextFilter;
import org.springframework.security.web.authentication.LoginUrlAuthenticationEntryPoint;
import org.springframework.security.web.authentication.preauth.AbstractPreAuthenticatedProcessingFilter;
import org.springframework.security.web.csrf.CookieCsrfTokenRepository;

@Configuration
@EnableWebSecurity
public class WebSecurityConfig extends WebSecurityConfigurerAdapter {
    @Autowired
    private OAuth2RestTemplate restTemplate;

    @Autowired
    private OnixBasicAuthEntryPoint authenticationEntryPoint;

    @Value("${wapi.auth.mode}")
    private String authMode;

    @Value("${wapi.admin.user}")
    private String adminUsername;

    @Value("${wapi.admin.pwd}")
    private String adminPassword;

    @Value("${wapi.reader.user}")
    private String readerUsername;

    @Value("${wapi.reader.pwd}")
    private String readerPassword;

    @Value("${wapi.writer.user}")
    private String writerUsername;

    @Value("${wapi.writer.pwd}")
    private String writerPassword;

    @Value("${WAPI_CSRF_ENABLED:false}")
    private boolean csrfEnabled;

    @Autowired
    public void configureGlobal(AuthenticationManagerBuilder auth) throws Exception {
        if (authMode.equals("basic")) {
            auth.inMemoryAuthentication()
                .withUser(readerUsername).password(passwordEncoder().encode(readerPassword)).roles("READER").and()
                .withUser(writerUsername).password(passwordEncoder().encode(writerPassword)).roles("WRITER").and()
                .withUser(adminUsername).password(passwordEncoder().encode(adminPassword)).roles("ADMIN");
        }
    }

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        if (authMode.equals("basic")) {
            // Basic Access Authentication: the request should contain a header field of the form Authorization: Basic <credentials>
            //  where credentials is the base64 encoding of username and password joined by a single colon (e.g. user:password)
            http
                .authorizeRequests()
                .antMatchers("/").permitAll()
                .antMatchers("/live").permitAll()
                .antMatchers("/ready").permitAll()
                .antMatchers(HttpMethod.GET, "(item|link|data|tag|linktype|itemtype|model|enckey)/**").hasAnyRole("READER", "WRITER")
                .antMatchers(HttpMethod.PUT, "(item|link|tag|data)/**").hasRole("WRITER")
                .antMatchers(HttpMethod.POST, "(tag)/**").hasRole("WRITER")
                .antMatchers(HttpMethod.DELETE, "(item|link|tag)/**").hasRole("WRITER")
                .antMatchers("/**").hasRole("ADMIN")
                // http basic authentication
                .and().httpBasic()
                // 401 challenge response configuration
                .authenticationEntryPoint(authenticationEntryPoint)
                // No need session.
                .and().sessionManagement().sessionCreationPolicy(SessionCreationPolicy.STATELESS);
        }
        else if (authMode.equals("none")) {
            http.csrf().csrfTokenRepository(CookieCsrfTokenRepository.withHttpOnlyFalse()).and()
                .authorizeRequests()
                .antMatchers("/**").permitAll().and()
                .sessionManagement().sessionCreationPolicy(SessionCreationPolicy.STATELESS);
        }
        else if (authMode.equals("oidc")) {
            // OpenId Connect authentication and OAuth 2.0 authorisation
            // Implements both authorisation code and password grant types (flows) to support
            //   Web UI and Native or Automated client implementations
            http
                .csrf().csrfTokenRepository(CookieCsrfTokenRepository.withHttpOnlyFalse()).and()
                // required by the grant_type=authorization_code
                .addFilterAfter(new OAuth2ClientContextFilter(), AbstractPreAuthenticatedProcessingFilter.class)
                .addFilterAfter(oidcFilter(), OAuth2ClientContextFilter.class)
                .httpBasic()
                    .authenticationEntryPoint(new LoginUrlAuthenticationEntryPoint("/oidc-login"))
                .and()
                    .authorizeRequests()
                    // permits access to endpoints foe liveliness and readyness probes
                    .antMatchers("/").permitAll()
                    .antMatchers("/live").permitAll()
                    .antMatchers("/ready").permitAll()
                    .anyRequest()
                        .authenticated()
                // required by the grant_type=password
                .and()
                    .oauth2ResourceServer()
                        .jwt()
                            // adds a Jwt Token Converter to extract the 'roles' claim from the token and turn it into
                            // authorities for authorisation
                            .jwtAuthenticationConverter(new JwtTokenClaimsConverter())
                ;
        }
        else {
            throw new RuntimeException(
                String.format("Incorrect AUTH_MODE value '%s': expected one of 'none', 'basic' or 'oidc'.", authMode));
        }

        if (csrfEnabled) {
            // persists the CSRF token in a cookie named "XSRF-TOKEN"
            http.csrf().csrfTokenRepository(CookieCsrfTokenRepository.withHttpOnlyFalse());
        } else {
            // disables CSRF
            http.csrf().disable();
        }
    }

    @Bean
    public OIDCFilter oidcFilter() {
        final OIDCFilter filter = new OIDCFilter("/oidc-login");
        filter.setRestTemplate(restTemplate);
        return filter;
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Override
    public void configure(WebSecurity web) throws Exception {
        web.ignoring().antMatchers("/resources/**");
    }
}
