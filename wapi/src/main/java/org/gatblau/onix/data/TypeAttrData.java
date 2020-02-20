package org.gatblau.onix.data;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.io.Serializable;

@ApiModel(
    description = "Defines an attribute of a configuration item type."
)
public class TypeAttrData implements Serializable {
    private static final long serialVersionUID = 1L;

    private String key;
    private String name;
    private String description;
    private String type;
    private String defValue;
    private Boolean managed;
    private Boolean required;
    private String regex;
    private String itemTypeKey;
    private String linkTypeKey;
    private String created;
    private String updated;
    private Integer version;
    private String changedBy;

    @ApiModelProperty(
        position = 0,
        required = true,
        value = "The natural key that uniquely identifies this attribute of a configuration item type.",
        example = "test_item_type_attribute"
    )
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

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getDefValue() {
        return defValue;
    }

    public void setDefValue(String defValue) {
        this.defValue = defValue;
    }

    public Boolean getManaged() {
        return managed;
    }

    public void setManaged(Boolean managed) {
        this.managed = managed;
    }

    public Boolean getRequired() {
        return required;
    }

    public void setRequired(Boolean required) {
        this.required = required;
    }

    public String getRegex() {
        return regex;
    }

    public void setRegex(String regex) {
        this.regex = regex;
    }

    public String getItemTypeKey() {
        return itemTypeKey;
    }

    public void setItemTypeKey(String itemTypeKey) {
        this.itemTypeKey = itemTypeKey;
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

    public String getLinkTypeKey() {
        return linkTypeKey;
    }

    public void setLinkTypeKey(String linkTypeKey) {
        this.linkTypeKey = linkTypeKey;
    }
}
