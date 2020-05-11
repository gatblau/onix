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

import org.gatblau.onix.data.UserData;
import org.gatblau.onix.db.DbRepository;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.authentication.AuthenticationProvider;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import java.text.DateFormat;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/*
  authenticates the user based on credentials securely stored in the Onix database
 */
@Service
public class UserPwdAuthProvider implements AuthenticationProvider {
    private final DateFormat dateFormat = new SimpleDateFormat("dd-MM-yyyy HH:mm:ss Z");
    private final DbRepository db;
    private final PwdBasedEncryptor enc;
    public UserPwdAuthProvider(DbRepository db, PwdBasedEncryptor enc) {
        this.db = db;
        this.enc = enc;
    }

    @Override
    public Authentication authenticate(Authentication authentication) throws AuthenticationException {
        String username = (String)authentication.getPrincipal();
        String pwd = (String)authentication.getCredentials();

        // retrieve the user
        UserData user = db.getUser(username, new String[]{"ADMIN"});
        
        // if the user was not found
        if (user == null) {
            throw new UsernameNotFoundException(String.format("User details not found for username: %s", username));
        }

        // if there is a password expiration date set
        if (user.getExpires() != null) {
            try {
                // checks it has not expired
                Date now = new Date();
                Date expiry = dateFormat.parse(user.getExpires());
                // if the expiry date is in the past
                if (expiry.before(now)) {
                    throw new RuntimeException(String.format("password expired for user: %s", user.getName()));
                }
            } catch (ParseException pe) {
                throw new RuntimeException(String.format("cannot parse valuntil date in the user record: %s", pe.getMessage()));
            }
        }

        boolean authenticated = false;
        try {
            // check the user provided password matches the one stored in the database
            authenticated = enc.authenticate(pwd, user.getPwd(), user.getSalt());
        } catch (Exception ex) {
            // something went wrong with the encryption of the password in the database
            throw new RuntimeException(String.format("Failed to check credentials for %s", username), ex);
        }
        
        // if the user failed the authentication 
        if (!authenticated) {
            // issues a bad credentials exception
            throw new BadCredentialsException(String.format("Authentication failed for %s" + username));
        }

        AbstractAuthenticationToken token = null;
        try {
            // at this point, the user is authenticated so fetches the roles the user has been assigned
            List grantedAuthorities = getUserAuthorities(username);
            token = new UsernamePasswordAuthenticationToken(username, pwd, grantedAuthorities);
        } catch (Exception ex) {
            ex.printStackTrace();
        }
        
        // constructs an authorization object 
        return token;
    }

    private List<GrantedAuthority> getUserAuthorities(String username) {
        List<String> roles = db.getUserRolesInternal(username);
        List<GrantedAuthority> grantedAuthorities = new ArrayList<>();
        for (String userRole : roles) {
            grantedAuthorities.add(new SimpleGrantedAuthority(String.format("ROLE_%s", userRole)));
        }
        return grantedAuthorities;
    }

    @Override
    public boolean supports(Class<?> authenticationClass) {
        // the type of object returned by 
        return authenticationClass.equals(UsernamePasswordAuthenticationToken.class);
    }
}
