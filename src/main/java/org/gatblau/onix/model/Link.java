package org.gatblau.onix.model;

import javax.persistence.*;
import java.io.Serializable;
import java.util.Date;

@NamedQueries(value= {
//    @NamedQuery(
//        name = "link.findFromNodeId",
//        query = "SELECT l FROM Link l "
//                + "JOIN Node l.startNode s "
//                + "WHERE s.id = :nodeId "
//    ),
//    @NamedQuery(
//        name = "link.findToNodeId",
//        query = "SELECT l FROM Link l "
//                + "JOIN Node l.endNode e "
//                + "WHERE e.id = :nodeId "
//    ),
    @NamedQuery(
        name = "link.deleteAll",
        query = "DELETE FROM Link "
    ),
})
@Entity
public class Link implements Serializable {
    private static final long serialVersionUID = 1L;

    public static final String FIND_FROM_ITEM_ID = "link.findFromItemId";
    public static final String FIND_TO_ITEM_ID = "link.findToItemId";
    public static final String DELETE_ALL = "link.deleteAll";

    public static final String PARAM_ITEM_ID = "itemId";

    @Id
    @GeneratedValue(strategy= GenerationType.IDENTITY)
    @Column(name = "id", updatable = false, nullable = false)
    private Long id = null;

    @ManyToOne(fetch= FetchType.LAZY)
    @JoinColumn(name="start_item_id")
    private Item startItem;

    @ManyToOne(fetch= FetchType.LAZY)
    @JoinColumn(name="end_item_id")
    private Item endItem;

    @Column
    private String key;

    @Column
    private String tag;

    @Column
    private String description;

    @Column
    private String meta;

    @Column
    private Date created;

    @Column
    private Date updated;

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

    public String getTag() {
        return tag;
    }

    public void setTag(String tag) {
        this.tag = tag;
    }

    public Item getStartItem() {
        return startItem;
    }

    public void setStartItem(Item startNode) {
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

    public String getMeta() {
        return meta;
    }

    public void setMeta(String meta) {
        this.meta = meta;
    }

    public Date getCreated() {
        return created;
    }

    public void setCreated(Date created) {
        this.created = created;
    }

    public Date getUpdated() {
        return updated;
    }

    public void setUpdated(Date updated) {
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
