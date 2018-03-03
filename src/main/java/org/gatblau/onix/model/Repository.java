package org.gatblau.onix.model;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.json.simple.JSONObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import javax.persistence.EntityManager;
import javax.persistence.NoResultException;
import javax.persistence.TypedQuery;
import java.io.IOException;
import java.time.ZonedDateTime;
import java.util.List;

@Service
public class Repository {
    private ObjectMapper mapper = new ObjectMapper();

    @Autowired
    private EntityManager em;

    public Repository() {
    }

    /***
     *
     * @param node
     * @return
     */
    public long createItem(Item node) {
        Item n = em.merge(node);
        return n.getId();
    }

    @Transactional
    public String createOrUpdateItem(String key, JSONObject json) throws IOException {
        String action = "UPDATED";
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
        query.setParameter(Item.PARAM_KEY, key);
        Item item = null;
        ZonedDateTime time = ZonedDateTime.now();
        try {
            item = query.getSingleResult();
        }
        catch (NoResultException e) {
            item = new Item();
            item.setCreated(time);
            action = "CREATED";
        }
        ItemType itemType = em.getReference(ItemType.class, Long.parseLong(json.get("itemTypeId").toString()));
        item.setKey(key);
        item.setItemType(itemType);
        item.setDescription((String)json.get("description"));
        item.setTag((String)json.get("tag"));
        item.setName((String)json.get("name"));
        item.setDeployed((boolean)json.get("deployed"));
        item.setUpdated(time);
        item.setMeta(mapper.readTree((String)json.get("meta")));
        em.persist(item);
        return action;
    }

    /***
     * Find all nodes of a particular type and with a specific tags.
     * @param itemTypeId
     * @param tag the tags used to filter the result of the search.
     * @return
     */
    public List<Item> findItemsByType(long itemTypeId, String tag) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE_AND_TAG, Item.class);
        query.setParameter(Item.PARAM_ITEM_TYPE_ID, itemTypeId);
        query.setParameter(Item.PARAM_TAG, tag);
        return query.getResultList();
    }

    /***
     * Finds links that depart from a specified node.
     * @param itemId the Id of the node from where links depart.
     * @return
     */
    public List<Link> findLinksFromItem(long itemId) {
        TypedQuery<Link> query = em.createNamedQuery(Link.FIND_FROM_ITEM_ID, Link.class);
        query.setParameter(Link.PARAM_ITEM_ID, itemId);
        return query.getResultList();
    }

    /***
     * Find links that arrive to a specified node.
     * @param itemId the Id of the node where links arrive.
     * @return
     */
    public List<Link> findLinksToItem(long itemId) {
        TypedQuery<Link> query = em.createNamedQuery(Link.FIND_TO_ITEM_ID, Link.class);
        query.setParameter(Link.PARAM_ITEM_ID, itemId);
        return query.getResultList();
    }

    /***
     * Finds nodes of a specified type with specified tags which are linked to the passed-in node.
     * @param itemId the id of the node used to find linked nodes
     * @param itemTypeId the type of linked nodes to find
     * @param tag the tags associated to the linked nodes used as a filter for the search.
     * @return
     */
    public List<Item> findLinkedItemsByTypeAndTag(long itemId, int itemTypeId, String tag) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_LINKED_NODES_BY_TYPE_AND_TAG, Item.class);
        query.setParameter(Link.PARAM_ITEM_ID, itemId);
        query.setParameter(Item.PARAM_TAG, tag);
        query.setParameter(Item.PARAM_ITEM_TYPE_ID, itemTypeId);
        return query.getResultList();
    }

    /***
     * Removes all transactional data in the database.
     */
    @Transactional
    public void clear() {
        if (em != null) {
            em.createNamedQuery(Link.DELETE_ALL).executeUpdate();
            em.createNamedQuery(Item.DELETE_ALL).executeUpdate();
        }
    }

    public void deleteItemTypes() {
        if (em != null) {
            em.createNamedQuery(ItemType.DELETE_ALL).executeUpdate();
        }
    }
}
