package org.gatblau.onix.model;

import javax.persistence.*;
import java.io.Serializable;

@NamedQueries(value= {
    @NamedQuery(
        name = "itemType.deleteAll",
        query = "DELETE FROM ItemType i WHERE i.custom = true "
    )
})
@Entity()
@Table(name = "item_type")
public class ItemType implements Serializable {
    private static final long serialVersionUID = 1L;
    public static final String DELETE_ALL = "itemType.deleteAll";

    @Id
    @GeneratedValue(strategy= GenerationType.IDENTITY)
    @Column(name = "id", updatable = false, nullable = false)
    private Long id = null;

    @Column
    private String name;

    @Column
    private String description;

    @Column
    private boolean custom;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
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
