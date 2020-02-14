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

import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;

//
// Interface for symmetric encryption algorithms
//
@Component
public interface Crypto {
    // returns a new secret key
    String newKey();

    // get the secret key from an UTF8 encoded string
    SecretKey fromString(String key);

    // encrypts a plain text
    byte[] encrypt(String plaintext, SecretKey key);

    // decrypts a cipher
    byte[] decrypt(byte[] encryptedData, SecretKey key);
}
