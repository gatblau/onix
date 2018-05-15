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
                "WHERE d.key = :key "
    )
})
@Entity
public class Dimension implements Serializable {
    private static final long serialVersionUID = 1L;
    public static final String DELETE_ALL = "dimension.deleteAll";
    public static final String FIND_BY_KEY = "dimension.findByKey";
    public static final String PARAM_KEY = "key";

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
