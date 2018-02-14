package org.gatblau.onix;

import java.io.Serializable;

public class Info implements Serializable {
    private static final long serialVersionUID = 1L;

    private String description;
    private String version;

    public Info() {
    }

    public Info(String description, String version) {
        this.description = description;
        this.version = version;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }
}
