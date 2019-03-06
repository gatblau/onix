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

package org.gatblau.onix;

import com.zaxxer.hikari.HikariDataSource;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Component
public class DataSourceFactory {
    private HikariDataSource ds;

    @Value("${spring.datasource.username}")
    private String dbuser;

    @Value("${spring.datasource.password}")
    private String dbpwd;

    @Value("${spring.datasource.url}")
    private String connString;

    public HikariDataSource instance() {
        if (ds == null) {
            ds = new HikariDataSource();
            System.out.println(String.format("JDBC ==> Setting JDBC URL to: '%s'", connString));
            ds.setJdbcUrl(connString);
            ds.setUsername(dbuser);
            ds.setPassword(dbpwd);
        }
        return ds;
    }
}