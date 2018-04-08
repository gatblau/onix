package org.gatblau.onix.data;

import com.fasterxml.jackson.databind.JsonNode;

import java.io.Serializable;

public class LinkData implements Serializable {
    private static final long serialVersionUID = 1L;

    private String key;
    private String role;
    private String tag;
    private String description;
    private JsonNode meta;
    private LinkedItemData item;

    public String getKey() {
        return key;
    }

    public void setKey(String key) {
        this.key = key;
    }

    public String getRole() {
        return role;
    }

    public void setRole(String role) {
        this.role = role;
    }

    public String getTag() {
        return tag;
    }

    public void setTag(String tag) {
        this.tag = tag;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public JsonNode getMeta() {
        return meta;
    }

    public void setMeta(JsonNode meta) {
        this.meta = meta;
    }

    public LinkedItemData getItem() {
        return item;
    }

    public void setItem(LinkedItemData item) {
        this.item = item;
    }
}