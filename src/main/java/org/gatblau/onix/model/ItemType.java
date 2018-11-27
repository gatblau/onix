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

import org.hibernate.annotations.Type;

import javax.persistence.*;
import java.io.Serializable;
import java.time.ZonedDateTime;

@NamedQueries(value= {
    @NamedQuery(
        name = "itemType.deleteAll",
        query = "DELETE FROM ItemType i WHERE i.custom = true "
    ),
    @NamedQuery(
        name = "itemType.findAll",
        query = "SELECT i FROM ItemType i "
    ),
    @NamedQuery(
        name = "itemType.findByKey",
        query = "SELECT i " +
                "FROM ItemType i " +
                "WHERE i.key = :key "
    )
})
@Entity()
@Table(name = "item_type")
public class ItemType implements Serializable {
    private static final long serialVersionUID = 1L;
    public static final String DELETE_ALL = "itemType.deleteAll";
    public static final String FIND_ALL = "itemType.findAll";
    public static final String FIND_BY_KEY = "itemType.findByKey";
    public static final String PARAM_KEY = "key";

    @Id
    @GeneratedValue(strategy= GenerationType.IDENTITY)
    @Column(name = "id", updatable = false, nullable = false)
    private Integer id = null;

    @Column
    private String key;

    @Column
    private String name;

    @Column
    private String description;

    @Column
    private boolean custom;

    @Column(columnDefinition= "TIMESTAMP WITH TIME ZONE")
    @Type(type="java.time.ZonedDateTime")
    private ZonedDateTime created;

    @Column(columnDefinition= "TIMESTAMP WITH TIME ZONE")
    @Type(type="java.time.ZonedDateTime")
    private ZonedDateTime updated;

    @Version
    @Column
    private int version;

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
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

    public boolean isCustom() {
        return custom;
    }

    public void setCustom(boolean custom) {
        this.custom = custom;
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
            return getId().equals(((ItemType) that).getId());
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
