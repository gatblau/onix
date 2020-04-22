/*
Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org

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

import org.gatblau.onix.conf.Config;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.authentication.AuthenticationProvider;
import org.springframework.security.config.annotation.authentication.builders.AuthenticationManagerBuilder;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.builders.WebSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.client.OAuth2RestTemplate;
import org.springframework.security.oauth2.client.filter.OAuth2ClientContextFilter;
import org.springframework.security.web.authentication.LoginUrlAuthenticationEntryPoint;
import org.springframework.security.web.authentication.preauth.AbstractPreAuthenticatedProcessingFilter;
import org.springframework.security.web.authentication.www.BasicAuthenticationFilter;
import org.springframework.security.web.csrf.CookieCsrfTokenRepository;

@Configuration
@EnableWebSecurity
public class WebSecurityConfig extends WebSecurityConfigurerAdapter {
    private final OAuth2RestTemplate restTemplate;
    private final Config cfg;
    private final AuthenticationProvider dbAuthProvider;

    public WebSecurityConfig(OAuth2RestTemplate restTemplate, UserPwdAuthProvider dbAuthProvider, Config cfg) {
        this.restTemplate = restTemplate;
        this.dbAuthProvider = dbAuthProvider;
        this.cfg = cfg;
    }

    @Autowired
    public void configureGlobal(AuthenticationManagerBuilder auth) throws Exception {
        if (cfg.getAuthMode().equals(Config.AuthMode.Basic)) {
            auth.authenticationProvider(dbAuthProvider);
        }
    }

    @Override
    protected void configure(HttpSecurity http) throws Exception {
        if (cfg.getAuthMode().equals(Config.AuthMode.None)) {
            http.csrf().csrfTokenRepository(CookieCsrfTokenRepository.withHttpOnlyFalse()).and()
                    .authorizeRequests()
                    .antMatchers("/**").permitAll().and()
                    .sessionManagement().sessionCreationPolicy(SessionCreationPolicy.STATELESS);
        } else if (cfg.getAuthMode().equals(Config.AuthMode.Basic)) {
            http.authorizeRequests()
                    .antMatchers("/").permitAll()
                    .antMatchers("/live").permitAll()
                    .antMatchers("/ready").permitAll()
                    .anyRequest().authenticated()
                    .and().httpBasic();
        } else if (cfg.getAuthMode().equals(Config.AuthMode.OIDC)) {
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
                String.format("Incorrect AUTH_MODE value '%s': expected one of 'none', 'basic' or 'oidc'.", cfg.getAuthMode()));
        }

        if (cfg.isCsrfEnabled()) {
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
