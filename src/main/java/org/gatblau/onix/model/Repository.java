package org.gatblau.onix.model;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.json.simple.JSONObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import javax.persistence.EntityManager;
import javax.persistence.NoResultException;
import javax.persistence.Query;
import javax.persistence.TypedQuery;
import java.io.IOException;
import java.time.ZonedDateTime;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

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

        List<LinkedHashMap<String, String>> dims =  (List<LinkedHashMap<String, String>>)json.get("dimensions");
        if (dims != null && action.equals("CREATED")) {
            for (LinkedHashMap<String, String> dim : dims) {
                Map.Entry<String, String> entry = dim.entrySet().iterator().next();
                Dimension d = new Dimension();
                d.setItem(item);
                d.setKey((String)entry.getKey());
                d.setValue((String)entry.getValue());
                em.persist(d);
            }
        }
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
     * Removes all transactional data in the database.
     */
    @Transactional
    public void clear() {
        if (em != null) {
            em.createNamedQuery(ItemType.DELETE_ALL).executeUpdate();
            em.createNamedQuery(Dimension.DELETE_ALL).executeUpdate();
            em.createNamedQuery(Link.DELETE_ALL).executeUpdate();
            em.createNamedQuery(Item.DELETE_ALL).executeUpdate();
        }
    }

    @Transactional
    public void deleteItemTypes() {
        if (em != null) {
            em.createNamedQuery(ItemType.DELETE_ALL).executeUpdate();
        }
    }

    @Transactional
    public String createOrUpdateLink(String fromItemKey, String toItemKey, JSONObject json) throws IOException {
        String action = "UPDATED";
        TypedQuery<Link> query = em.createNamedQuery(Link.FIND_BY_KEYS, Link.class);
        query.setParameter(Link.KEY_START_ITEM, fromItemKey);
        query.setParameter(Link.KEY_END_ITEM, toItemKey);
        Link link = null;
        ZonedDateTime time = ZonedDateTime.now();
        try {
            link = query.getSingleResult();
        }
        catch (NoResultException e) {
            link = new Link();
            link.setCreated(time);
            action = "CREATED";

            TypedQuery<Item> startItemQuery = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
            startItemQuery.setParameter(Item.PARAM_KEY, fromItemKey);
            Item startItem = null;
            try {
                startItem = startItemQuery.getSingleResult();
            }
            catch (NoResultException nre) {
                throw new RuntimeException("Could create link to start configuration item with key '" + fromItemKey + "' as it does not exist.");
            }

            TypedQuery<Item> endItemQuery = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
            endItemQuery.setParameter(Item.PARAM_KEY, toItemKey);
            Item endItem = null;
            try {
                endItem = endItemQuery.getSingleResult();
            }
            catch (NoResultException nre) {
                throw new RuntimeException("Could not create link to end configuration item with key '" + toItemKey + "' as it does not exist.");
            }

            link.setStartItem(startItem);
            link.setEndItem(endItem);

            link.setKey(String.format("%s_%s", startItem.getKey(), endItem.getKey()));
        }
        link.setDescription((String)json.get("description"));
        link.setMeta(mapper.readTree((String)json.get("meta")));
        link.setTag((String)json.get("tag"));
        link.setRole((String)json.get("role"));

        link.setUpdated(time);

        em.persist(link);

        return action;
    }

    @Transactional
    public void deleteLink(String fromItemKey, String toItemKey) {
        if (em != null) {
            TypedQuery<Link> query = em.createNamedQuery(Link.FIND_BY_KEYS, Link.class);
            query.setParameter(Link.KEY_START_ITEM, fromItemKey);
            query.setParameter(Link.KEY_END_ITEM, toItemKey);
            Link link = null;
            try {
                query.getSingleResult();
            }
            catch (NoResultException nre){
                // do nothing so link stays null
            }
            if (link != null) {
                em.remove(link);
            }
        }
    }
}
