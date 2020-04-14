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
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.Cipher;
import javax.crypto.KeyGenerator;
import javax.crypto.SecretKey;
import javax.crypto.spec.GCMParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Base64;
import java.util.Date;

//
// Authenticated Symmetric Encryption using Advanced Encryption Standard (AES) and Galois/Counter Mode (GCM) ciphers
// with no padding. Creates different random IVs and embed them in the returned encrypted data.
//
@Component
public class AESGCM implements Crypto {
    private static final int AES_KEY_SIZE = 256;
    private static final String ALGORITHM = "AES";
    private static final String CIPHER = "AES/GCM/NoPadding";
    private static final DateFormat formatter = new SimpleDateFormat("dd-MM-yyyy");

    // This size of the IV (in bytes) is normally (keysize / 8).
    // If the default keysize is 256, so the IV must be 32 bytes long.
    // Using a 16 character string here gives us 32 bytes when converted to a byte array.
    // For Galois/Counter Mode (GCM), in principle any IV size can be used as long as the IV doesn't ever repeat.
    // NIST however suggests that only an IV size of 12 bytes (96 bits) needs to be supported by implementations,
    // As other IV lengths will require additional calculations impairing performance.
    private static final int GCM_IV_LENGTH_BYTES = 12;

    // Size of authentication tags: the calculated tag will always be 16 bytes long, but the leftmost bytes can be used.
    // GCM is defined for the tag sizes 128, 120, 112, 104, or 96, 64 and 32.
    // Note that the security of GCM is strongly dependent on the tag size.
    // You should try and use a tag size of 64 bits at the very minimum, but in general a tag size of the full 128 bits
    // should be preferred.
    private static final int GCM_TAG_LENGTH_BYTES = 16; // 16 * 8 = 128 bits

    private final Config cfg;

    public AESGCM(Config cfg) {
        this.cfg = cfg;
    }

    // returns a new secret key
    @Override
    public String newKey() {
        KeyGenerator keyGenerator = null;
        try {
            keyGenerator = KeyGenerator.getInstance(ALGORITHM);
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        }
        keyGenerator.init(AES_KEY_SIZE);
        return Base64.getEncoder().encodeToString(keyGenerator.generateKey().getEncoded());
    }

    private SecretKey fromCharArray(char[] encodedKey) {
        byte[] decodedKey = Base64.getDecoder().decode(new String(encodedKey).trim());
        return new SecretKeySpec(decodedKey, 0, decodedKey.length, ALGORITHM);
    }

    // encrypts a plain text
    @Override
    public byte[] encrypt(String plaintext) {
        try {
            byte[] textBytes = plaintext.getBytes(StandardCharsets.UTF_8);

            SecureRandom secureRandom = new SecureRandom();
            byte[] iv = new byte[GCM_IV_LENGTH_BYTES];
            secureRandom.nextBytes(iv);

            Cipher cipher = Cipher.getInstance(CIPHER);
            GCMParameterSpec parameterSpec = new GCMParameterSpec(GCM_TAG_LENGTH_BYTES * 8, iv);
            cipher.init(Cipher.ENCRYPT_MODE, getKey(), parameterSpec);
            byte[] encryptedData = cipher.doFinal(textBytes);

            ByteBuffer byteBuffer = ByteBuffer.allocate(iv.length + encryptedData.length);
            byteBuffer.put(iv);
            byteBuffer.put(encryptedData);

            return byteBuffer.array();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    // decrypts a cipher
    @Override
    public byte[] decrypt(byte[] encryptedData, short keyIx) {
        try {
            ByteBuffer byteBuffer = ByteBuffer.wrap(encryptedData);

            byte[] iv = new byte[GCM_IV_LENGTH_BYTES];
            byteBuffer.get(iv);

            byte[] cipherBytes = new byte[byteBuffer.remaining()];
            byteBuffer.get(cipherBytes);

            Cipher cipher = Cipher.getInstance(CIPHER);
            GCMParameterSpec parameterSpec = new GCMParameterSpec(GCM_TAG_LENGTH_BYTES * 8, iv);
            cipher.init(Cipher.DECRYPT_MODE, getKey(keyIx), parameterSpec);
            return cipher.doFinal(cipherBytes);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private SecretKey getKey(int currentKey) {
        switch (currentKey) {
            case 1:
                return fromCharArray(cfg.getKey1());
            case 2:
                return fromCharArray(cfg.getKey2());
            default:
                throw new RuntimeException("Active key invalid. Has to be either A or B.");
        }
    }

    // the current key used for encryption
    private SecretKey getKey() {
        return getKey(getKeyIx());
    }

    // get key no 1 or 2 depending on expiry date
    @Override
    public short getKeyIx() {
        try {
            Date today = new Date();
            Date expiry = formatter.parse(cfg.getKeyExpiry());
            if (today.after(expiry)) {
                // return the non-active key
                return (cfg.getKeyDefault() == 1) ? (short)2 : 1;
            } else {
                // return active key
                return cfg.getKeyDefault();
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public short getDefaultKeyIx() {
        return cfg.getKeyDefault();
    }

    @Override
    public String getDefaultKeyExpiry() {
        return cfg.getKeyExpiry();
    }
}