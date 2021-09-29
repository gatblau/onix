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

    Code adapted from:
    http://blog.jerryorr.com/2012/05/secure-password-storage-lots-of-donts.html
*/
package org.gatblau.onix.security;

import org.gatblau.onix.data.UserData;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.PBEKeySpec;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.KeySpec;
import java.text.DateFormat;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Base64;
import java.util.Date;

@Service
public class PwdBasedEncryptor {
    private final DateFormat dateFormat = new SimpleDateFormat("dd-MM-yyyy HH:mm:ss Z");

    public boolean authenticate(String attemptedPassword, String encryptedPassword, String salt) throws NoSuchAlgorithmException, InvalidKeySpecException {
        // Encrypt the clear-text password using the same salt that was used to
        // encrypt the original password
        String encryptedAttemptedPassword = getEncryptedPwd(attemptedPassword, salt);

        // Authentication succeeds if encrypted password that the user entered
        // is equal to the stored hash
        return encryptedPassword.equals(encryptedAttemptedPassword);
    }

    public String generateSalt() throws NoSuchAlgorithmException {
        return Base64.getEncoder().encodeToString(generateSaltBytes());
    }

    public String getEncryptedPwd(String password, String salt) throws InvalidKeySpecException, NoSuchAlgorithmException {
        if (password == null) {
            return null;
        }
        return Base64.getEncoder().encodeToString(getEncryptedPwdBytes(password, Base64.getDecoder().decode(salt)));
    }

    private byte[] getEncryptedPwdBytes(String password, byte[] salt) throws NoSuchAlgorithmException, InvalidKeySpecException {
        // PBKDF2 with SHA-1 as the hashing algorithm. Note that the NIST
        // specifically names SHA-1 as an acceptable hashing algorithm for PBKDF2
        String algorithm = "PBKDF2WithHmacSHA1";
        // SHA-1 generates 160 bit hashes, so that's what makes sense here
        int derivedKeyLength = 160;
        // Pick an iteration count that works for you. The NIST recommends at
        // least 1,000 iterations:
        // http://csrc.nist.gov/publications/nistpubs/800-132/nist-sp800-132.pdf
        // iOS 4.x reportedly uses 10,000:
        // http://blog.crackpassword.com/2010/09/smartphone-forensics-cracking-blackberry-backup-passwords/
        int iterations = 20000;
        KeySpec spec = new PBEKeySpec(password.toCharArray(), salt, iterations, derivedKeyLength);
        SecretKeyFactory f = SecretKeyFactory.getInstance(algorithm);
        return f.generateSecret(spec).getEncoded();
    }

    private byte[] generateSaltBytes() throws NoSuchAlgorithmException {
        // VERY important to use SecureRandom instead of just Random
        SecureRandom random = SecureRandom.getInstance("SHA1PRNG");

        // Generate a 8 byte (64 bit) salt as recommended by RSA PKCS5
        byte[] salt = new byte[8];
        random.nextBytes(salt);

        return salt;
    }

    // authenticate user using username and password credentials
    public boolean authenticateUser(String username, String password, UserData user) {
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
            authenticated = authenticate(password, user.getPwd(), user.getSalt());
        } catch (Exception ex) {
            // something went wrong with the encryption of the password in the database
            throw new RuntimeException(String.format("Failed to check credentials for %s", username), ex);
        }
        return authenticated;
    }
}