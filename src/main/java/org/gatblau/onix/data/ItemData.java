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

import org.json.simple.JSONObject;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

public class ItemData implements Serializable {
    private static final long serialVersionUID = 1L;

    private String key;
    private String name;
    private String description;
    private String itemType;
    private Short status;
    private JSONObject meta;
    private JSONObject attribute;
    private List<String> tag;
    private List<LinkData> fromLinks = new ArrayList<>();
    private List<LinkData> toLinks = new ArrayList<>();
    private String created;
    private String updated;
    private Integer version;

    public ItemData() {
    }

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

    public String getItemType() {
        return itemType;
    }

    public void setItemType(String itemType) {
        this.itemType = itemType;
    }

    public Short getStatus() {
        return status;
    }

    public void setStatus(Short status) {
        this.status = status;
    }

    public JSONObject getMeta() {
        return meta;
    }

    public void setMeta(JSONObject meta) {
        this.meta = meta;
    }

    public List<String> getTag() {
        return tag;
    }

    public void setTag(List<String> tag) {
        this.tag = tag;
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

    public List<LinkData> getFromLinks() {
        return fromLinks;
    }

    public void setFromLinks(List<LinkData> links) {
        this.fromLinks = links;
    }

    public List<LinkData> getToLinks() {
        return toLinks;
    }

    public void setToLinks(List<LinkData> links) {
        this.toLinks = links;
    }

    public JSONObject getAttribute() {
        return attribute;
    }

    public void setAttribute(JSONObject attribute) {
        this.attribute = attribute;
    }
}
