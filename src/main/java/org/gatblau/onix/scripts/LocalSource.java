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

package org.gatblau.onix.scripts;

import org.gatblau.onix.FileUtil;
import org.gatblau.onix.conf.Info;
import org.springframework.stereotype.Service;

/**
 * read database scripts embedded in the application binary file
 */
@Service
public class LocalSource extends ScriptSource {

    public LocalSource(FileUtil file, Info info) {
        super(file, info);
    }

    @Override
    protected String getAppManifestURI() {
        return String.format("app/%s.json", getAppVersion());
    }

    @Override
    protected String getDbManifestURI(String dbVersion) {
        return String.format("db/install/%s/db.json", dbVersion);
    }

    @Override
    protected String getDbScriptURI(String dbVersion, String scriptName) {
        return String.format("db/install/%s/%s", dbVersion, scriptName);
    }

    @Override
    protected String getContent(String uriString, String username, String password) {
        try {
            return file.getFile(uriString);
        } catch (Exception e) {
            throw new RuntimeException("Can't retrieve content.", e);
        }
    }

    @Override
    public String getSource() {
        return String.format("image: %s", file.getFile("version"));
    }
}
