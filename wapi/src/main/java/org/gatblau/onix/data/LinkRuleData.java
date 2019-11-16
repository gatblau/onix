/*
Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

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

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.io.Serializable;

@ApiModel(
    description = "A rule determining which items can be connected using a specific link."
)
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

    @ApiModelProperty(
        position = 0,
        required = true,
        value = "The natural key uniquely identifying this model.",
        example = "model_01"
    )
    public String getKey() {
        return key;
    }

    public void setKey(String key) {
        this.key = key;
    }

    @ApiModelProperty(
        position = 1,
        required = true,
        value = "The rule name.",
        example = "Item Type X to Item Type Y rule."
    )
    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    @ApiModelProperty(
        position = 2,
        required = false,
        value = "The description of the link rule.",
        example = "This rule allow linking items of type X."
    )
    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    @ApiModelProperty(
        position = 3,
        required = true,
        value = "The natural key identifying the type of link this rule is for.",
        example = "test_link"
    )
    public String getLinkTypeKey() {
        return linkTypeKey;
    }

    public void setLinkTypeKey(String linkTypeKey) {
        this.linkTypeKey = linkTypeKey;
    }

    @ApiModelProperty(
        position = 4,
        required = true,
        value = "The item type from which the link should depart.",
        example = "item_type_A"
    )
    public String getStartItemTypeKey() {
        return startItemTypeKey;
    }

    public void setStartItemTypeKey(String startItemTypeKey) {
        this.startItemTypeKey = startItemTypeKey;
    }

    @ApiModelProperty(
        position = 5,
        required = true,
        value = "The item type to which the link should arrive.",
        example = "item_type_B"
    )
    public String getEndItemTypeKey() {
        return endItemTypeKey;
    }

    public void setEndItemTypeKey(String endItemTypeKey) {
        this.endItemTypeKey = endItemTypeKey;
    }

    @ApiModelProperty(
        position = 6,
        required = false,
        value = "Date and time on which the rule was created.",
        example = "01-02-2017 09:16:38"
    )
    public String getCreated() {
        return created;
    }

    public void setCreated(String created) {
        this.created = created;
    }

    @ApiModelProperty(
        position = 7,
        required = false,
        value = "Date and time on which the rule was updated.",
        example = "01-03-2017 14:32:57"
    )
    public String getUpdated() {
        return updated;
    }

    public void setUpdated(String updated) {
        this.updated = updated;
    }

    @ApiModelProperty(
        position = 8,
        required = false,
        value = "The version number for this rule.",
        example = "7"
    )
    public Integer getVersion() {
        return version;
    }

    public void setVersion(Integer version) {
        this.version = version;
    }

    @ApiModelProperty(
        position = 9,
        required = false,
        value = "The user who made the last change.",
        example = "admin"
    )
    public String getChangedBy() {
        return changedBy;
    }

    public void setChangedBy(String changedBy) {
        this.changedBy = changedBy;
    }
}
