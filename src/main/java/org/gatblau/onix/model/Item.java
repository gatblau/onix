package org.gatblau.onix.model;

import javax.persistence.*;
import java.io.Serializable;
import java.util.Date;

@NamedQueries(value= {
    @NamedQuery(
        name = "item.findByTypeAndTag",
        query =   "SELECT i "
                + "FROM Item i "
                + "JOIN i.itemType t "
                + "WHERE t.id = :itemTypeId "
                + "AND i.tag = :tag "
    ),
    @NamedQuery(
        name = "item.findByKey",
        query =   "SELECT i "
                + "FROM Item i "
                + "WHERE i.key = :key "
    ),
//    @NamedQuery(
//        name = "node.findLinkedNodesByTypeAndTag",
//        query = "SELECT endNode " +
//                "FROM link " +
//                "INNER JOIN node endNode " +
//                "ON endNode.id = link.end_node_id " +
//                "AND endNode.tag LIKE '%:tag%' " +
//                "AND endNode.node_type_id = :nodeTypeId " +
//                "INNER JOIN node startNode " +
//                "ON startNode.id = link.start_node_id " +
//                "AND startNode.id = :nodeId " +
//                "ORDER BY startNode.id ASC"
//    ),
    @NamedQuery(
        name = "item.deleteAll",
        query = "DELETE FROM Item"
    ),
})
@Entity
public class Item implements Serializable {
    private static final long serialVersionUID = 1L;

    public static final String FIND_BY_TYPE_AND_TAG = "item.findByTypeAndTag";
    public static final String FIND_LINKED_NODES_BY_TYPE_AND_TAG = "item.findLinkedNodesByTypeAndTag";
    public static final String FIND_BY_KEY = "item.findByKey";
    public static final String DELETE_ALL = "item.deleteAll";

    public static final String PARAM_ITEM_TYPE_ID = "itemTypeId";
    public static final String PARAM_TAG = "tag";
    public static final String PARAM_KEY = "key";

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

//    @Column
//    private String meta;

    @Column
    private Date created;

    @Column
    private Date updated;

    @Column
    private String tag;

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

//    public String getMeta() {
//        return meta;
//    }
//
//    public void setMeta(String meta) {
//        this.meta = meta;
//    }

    public String getTag(){
        return tag;
    }

    public void setTag(String tag){
        this.tag = tag;
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

    public ItemType getItemType() {
        return itemType;
    }

    public void setItemType(ItemType itemType) {
        this.itemType = itemType;
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
