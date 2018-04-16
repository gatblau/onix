package org.gatblau.onix.model;

import javax.persistence.*;
import java.io.Serializable;

@NamedQueries(value= {
    @NamedQuery(
        name = "itemType.deleteAll",
        query = "DELETE FROM ItemType i WHERE i.custom = true "
    ),
    @NamedQuery(
        name = "itemType.findAll",
        query = "SELECT i FROM ItemType i "
    )
})
@Entity()
@Table(name = "item_type")
public class ItemType implements Serializable {
    private static final long serialVersionUID = 1L;
    public static final String DELETE_ALL = "itemType.deleteAll";
    public static final String FIND_ALL = "itemType.findAll";

    @Id
    @GeneratedValue(strategy= GenerationType.IDENTITY)
    @Column(name = "id", updatable = false, nullable = false)
    private Integer id = null;

    @Column
    private String name;

    @Column
    private String description;

    @Column
    private boolean custom;

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
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
