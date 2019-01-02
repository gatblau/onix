package org.gatblau.onix.data;

import java.io.Serializable;

public class ItemTypeData implements Serializable {
    private static final long serialVersionUID = 1L;

    private String key;
    private String name;
    private String description;

    public ItemTypeData() {
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
}