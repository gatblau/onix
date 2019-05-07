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

package org.gatblau.onix.security;

import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.oauth2.common.OAuth2AccessToken;

import java.util.*;

public class OIDCUserDetails implements UserDetails {

    private static final long serialVersionUID = 1L;

    private String userId;
    private String username;
    private OAuth2AccessToken token;
    private List<GrantedAuthority> authorities;

    public OIDCUserDetails(
            Map<String, String> authenticationInfo,
            Map<String, String> authorisationInfo,
            OAuth2AccessToken token
    ) {
        this.userId = authenticationInfo.get("sub");
        this.username = authenticationInfo.get("email");
        this.authorities = extractAuthorities(authorisationInfo);
        this.token = token;
    }

    private static List<GrantedAuthority> extractAuthorities(Map<String, String> authorisationInfo) {
        List<GrantedAuthority> authorities = new ArrayList<>();
        String role = authorisationInfo.get("role");
        if (role != null) {
            authorities.add(new SimpleGrantedAuthority("ROLE_" + role));
        }
        return authorities;
    }

    @Override
    public String getUsername() {
        return username;
    }

    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        return authorities;
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public OAuth2AccessToken getToken() {
        return token;
    }

    public void setToken(OAuth2AccessToken token) {
        this.token = token;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    @Override
    public String getPassword() {
        return null;
    }

    @Override
    public boolean isAccountNonExpired() {
        return true;
    }

    @Override
    public boolean isAccountNonLocked() {
        return true;
    }

    @Override
    public boolean isCredentialsNonExpired() {
        return true;
    }

    @Override
    public boolean isEnabled() {
        return true;
    }
}