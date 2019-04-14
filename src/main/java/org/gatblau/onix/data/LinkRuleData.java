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

package org.gatblau.onix.data;

import java.io.Serializable;

public class LinkRuleData implements Serializable {
    private static final long serialVersionUID = 1L;

    private String key;
    private String name;
    private String description;
    private String linkTypeKey;
    private String startItemTypeKey;
    private String endItemTypeKey;
    private String created;
    private String updated;
    private Integer version;
    private String changedBy;

    public String getKey() {
        return key;
    }

    public void setKey(String key) {
        this.key = key;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getLinkTypeKey() {
        return linkTypeKey;
    }

    public void setLinkTypeKey(String linkTypeKey) {
        this.linkTypeKey = linkTypeKey;
    }

    public String getStartItemTypeKey() {
        return startItemTypeKey;
    }

    public void setStartItemTypeKey(String startItemTypeKey) {
        this.startItemTypeKey = startItemTypeKey;
    }

    public String getEndItemTypeKey() {
        return endItemTypeKey;
    }

    public void setEndItemTypeKey(String endItemTypeKey) {
        this.endItemTypeKey = endItemTypeKey;
    }

    public String getCreated() {
        return created;
    }

    public void setCreated(String created) {
        this.created = created;
    }

    public String getUpdated() {
        return updated;
    }

    public void setUpdated(String updated) {
        this.updated = updated;
    }

    public Integer getVersion() {
        return version;
    }

    public void setVersion(Integer version) {
        this.version = version;
    }

    public String getChangedBy() {
        return changedBy;
    }

    public void setChangedBy(String changedBy) {
        this.changedBy = changedBy;
    }
}
