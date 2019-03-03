package org.gatblau.onix.data;

import java.io.Serializable;

public class SnapshotData implements Serializable {
    private static final long serialVersionUID = 1L;

    private String label;
    private String rootItemKey;
    private String name;
    private String description;
    private String created;
    private String updated;
    private Integer version;
    private String changedBy;

    public SnapshotData() {
    }

    public String getLabel() {
        return label;
    }

    public void setLabel(String label) {
        this.label = label;
    }

    public String getRootItemKey() {
        return rootItemKey;
    }

    public void setRootItemKey(String rootItemKey) {
        this.rootItemKey = rootItemKey;
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
