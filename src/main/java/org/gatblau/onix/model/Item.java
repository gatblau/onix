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
import org.hibernate.annotations.Type;

import javax.persistence.*;
import java.io.Serializable;
import java.time.ZonedDateTime;
import java.util.List;

@NamedQueries(value= {
    @NamedQuery(
        name = "item.findAllByDateDesc",
        query =   "SELECT i "
                + "FROM Item i "
                + "ORDER BY i.updated DESC "
    ),
    @NamedQuery(
        name = "item.findByKey",
        query =   "SELECT i "
                + "FROM Item i "
                + "WHERE i.key = :key "
    ),
    @NamedQuery(
        name = "item.findByTag",
        query =   "SELECT i "
                + "FROM Item i "
                + "WHERE i.tag LIKE :tag "
                + "ORDER BY i.updated DESC "
    ),
    @NamedQuery(
        name = "item.findByType",
        query =   "SELECT i "
                + "FROM Item i "
                + "WHERE i.itemType.key = :itemTypeKey "
                + "ORDER BY i.updated DESC "
    ),
    @NamedQuery(
            name = "item.findByDate",
            query =   "SELECT i "
                + "FROM Item i "
                + "WHERE i.updated >= :fromDate "
                + "AND i.updated <= :toDate "
                + "ORDER BY i.updated DESC "
    ),
    @NamedQuery(
        name = "item.findByTypeAndTag",
        query =   "SELECT i "
                + "FROM Item i "
                + "WHERE i.itemType.key = :itemTypeKey "
                + "AND i.tag LIKE :tag "
                + "ORDER BY i.updated DESC "
    ),
    @NamedQuery(
        name = "item.findByTypeAndDate",
        query =   "SELECT i "
                + "FROM Item i "
                + "WHERE i.itemType.key = :itemTypeKey "
                + "AND i.updated >= :fromDate "
                + "AND i.updated <= :toDate "
                + "ORDER BY i.updated DESC "
    ),
    @NamedQuery(
        name = "item.findByTypeTagAndDate",
        query =   "SELECT i "
            + "FROM Item i "
            + "WHERE i.itemType.key = :itemTypeKey "
            + "AND i.tag LIKE :tag "
            + "AND i.updated >= :fromDate "
            + "AND i.updated <= :toDate "
            + "ORDER BY i.updated DESC "
    ),
    @NamedQuery(
        name = "item.deleteAll",
        query = "DELETE FROM Item "
    )
})
@Entity
public class Item implements Serializable {
    private static final long serialVersionUID = 1L;

    public static final String FIND_ALL_BY_DATE_DESC = "item.findAllByDateDesc";
    public static final String FIND_BY_KEY = "item.findByKey";
    public static final String FIND_BY_TYPE = "item.findByType";
    public static final String FIND_BY_TAG = "item.findByTag";
    public static final String FIND_BY_DATE = "item.findByDate";
    public static final String FIND_BY_TYPE_AND_TAG = "item.findByTypeAndTag";
    public static final String FIND_BY_TYPE_AND_DATE = "item.findByTypeAndDate";
    public static final String FIND_BY_TYPE_TAG_AND_DATE = "item.findByTypeTagAndDate";
    public static final String DELETE_ALL = "item.deleteAll";

    public static final String PARAM_ITEM_TYPE_KEY = "itemTypeKey";
    public static final String PARAM_TAG = "tag";
    public static final String PARAM_KEY = "key";
    public static final String PARAM_FROM_DATE = "fromDate";
    public static final String PARAM_TO_DATE = "toDate";

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", updatable = false, nullable = false)
    private Long id = null;

    @ManyToOne
    @JoinColumn(name="item_type_id")
    private ItemType itemType;

    @Column
    private String key;

    @Column
    private String name;

    @Column
    private String description;

    @SuppressWarnings("JpaAttributeTypeInspection")
    @Column(name = "meta", nullable = true)
    @Convert(converter = JSONBConverter.class)
    private JsonNode meta;

    @Column(columnDefinition= "TIMESTAMP WITH TIME ZONE")
    @Type(type="java.time.ZonedDateTime")
    private ZonedDateTime created;

    @Column(columnDefinition= "TIMESTAMP WITH TIME ZONE")
    @Type(type="java.time.ZonedDateTime")
    private ZonedDateTime updated;

    @Column
    private String tag;

    @Version
    @Column
    private int version;

    @Column
    private Short status;

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

    public JsonNode getMeta() {
        return meta;
    }

    public void setMeta(JsonNode meta) {
        this.meta = meta;
    }

    public String getTag(){
        return tag;
    }

    public void setTag(String tag){
        this.tag = tag;
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

    public ItemType getItemType() {
        return itemType;
    }

    public void setItemType(ItemType itemType) {
        this.itemType = itemType;
    }

    public void setStatus(Short status) {
        this.status = status;
    }

    public Short getStatus() {
        return this.status;
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
            return getId().equals(((Item) that).getId());
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