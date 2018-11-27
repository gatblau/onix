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

import javax.persistence.*;
import java.io.Serializable;

@NamedQueries(value= {
    @NamedQuery(
        name = "dimension.deleteAll",
        query = "DELETE FROM Dimension "
    ),
    @NamedQuery(
        name = "dimension.findByKey",
        query = "SELECT d " +
                "FROM Dimension d " +
                "WHERE d.key = :key " +
                "AND d.item.key = :itemKey "
    )
})
@Entity
public class Dimension implements Serializable {
    private static final long serialVersionUID = 1L;
    public static final String DELETE_ALL = "dimension.deleteAll";
    public static final String FIND_BY_KEY = "dimension.findByKey";
    public static final String PARAM_KEY = "key";
    public static final String PARAM_ITEM_KEY = "itemKey";

    @Id
    @GeneratedValue(strategy= GenerationType.IDENTITY)
    @Column(name = "id", updatable = false, nullable = false)
    private Long id = null;

    @Column
    private String key;

    @Column
    private String value;

    @ManyToOne(fetch=FetchType.LAZY)
    @JoinColumn(name="item_id")
    private Item item;

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

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public Item getItem() {
        return item;
    }

    public void setItem(Item item) {
        this.item = item;
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
            return getId().equals(((Dimension) that).getId());
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
