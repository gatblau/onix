/*
Onix CMDB - Copyright (c) 2018 by www.gatblau.org

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

package org.gatblau.onix.model;

import com.fasterxml.jackson.databind.JsonNode;

import javax.persistence.*;
import java.io.Serializable;
import java.time.ZonedDateTime;
import java.util.Date;

@NamedQueries(value= {
    @NamedQuery(
        name = "link.deleteAll",
        query = "DELETE FROM Link "
    ),
    @NamedQuery(
        name = "link.findByKey",
        query =   "SELECT L "
            + "FROM Link L "
            + "WHERE L.key = :key "
    ),
    @NamedQuery(
        name = "link.findFromItem",
        query = "SELECT L "
            + "FROM Link L "
            + "WHERE L.startItem.id = :itemId"
    ),
    @NamedQuery(
        name = "link.findToItem",
        query = "SELECT L "
            + "FROM Link L "
            + "WHERE L.endItem.id = :itemId"
    )
})
@Entity
public class Link implements Serializable {
    private static final long serialVersionUID = 1L;

    public static final String FIND_FROM_ITEM = "link.findFromItem";
    public static final String FIND_TO_ITEM = "link.findToItem";
    public static final String DELETE_ALL = "link.deleteAll";

    public static final String FIND_BY_KEY = "link.findByKey";

    public static final String KEY_LINK = "key";
    public static final String KEY_ITEM_ID = "itemId";

    @Id
    @GeneratedValue(strategy= GenerationType.IDENTITY)
    @Column(name = "id", updatable = false, nullable = false)
    private Long id = null;

    @Column
    private String key;

    @Column
    private String role;

    @ManyToOne(fetch= FetchType.LAZY)
    @JoinColumn(name="start_item_id")
    private Item startItem;

    @ManyToOne(fetch= FetchType.LAZY)
    @JoinColumn(name="end_item_id")
    private Item endItem;

    @Column
    private String tag;

    @Column
    private String description;

    @SuppressWarnings("JpaAttributeTypeInspection")
    @Column(name = "meta", nullable = true)
    @Convert(converter = JSONBConverter.class)
    private JsonNode meta;

    @Column
    private ZonedDateTime created;

    @Column
    private ZonedDateTime updated;

    @Version
    @Column
    private int version;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

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

    public Item getStartItem() {
        return startItem;
    }

    public void setStartItem(Item startItem) {
        this.startItem = startItem;
    }

    public Item getEndItem() {
        return endItem;
    }

    public void setEndItem(Item endItem) {
        this.endItem = endItem;
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

    public ZonedDateTime getCreated() {
        return created;
    }

    public void setCreated(ZonedDateTime created) {
        this.created = created;
    }

    public ZonedDateTime getUpdated() {
        return updated;
    }

    public void setUpdated(ZonedDateTime updated) {
        this.updated = updated;
    }

    public int getVersion() {
        return version;
    }

    protected void setVersion(int version) {
        this.version = version;
    }

    @Override
    public boolean equals(Object that) {
        if (this == that) {
            return true;
        }
        if (that == null) {
            return false;
        }
        if (getClass() != that.getClass()) {
            return false;
        }
        if (getId() != null) {
            return getId().equals(((Link) that).getId());
        }
        return super.equals(that);
    }

    @Override
    public int hashCode() {
        if (getId() != null) {
            return getId().hashCode();
        }
        return super.hashCode();
    }
}
