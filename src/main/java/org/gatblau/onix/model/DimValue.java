package org.gatblau.onix.model;

import javax.persistence.*;
import java.io.Serializable;

@Entity
public class DimValue implements Serializable {
    private static final long serialVersionUID = 1L;

    @Id
    @GeneratedValue(strategy= GenerationType.IDENTITY)
    @Column(name = "id", updatable = false, nullable = false)
    private Long id = null;

    @Column
    private String value;

    @ManyToOne
    @JoinColumn(name="dim_type_id")
    private DimType dimType;

    @ManyToOne(fetch=FetchType.LAZY)
    @JoinColumn(name="item_id")
    private Item item;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public DimType getDimType() {
        return dimType;
    }

    public void setDimType(DimType dimType) {
        this.dimType = dimType;
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
            return getId().equals(((DimValue) that).getId());
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
