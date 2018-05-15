package org.gatblau.onix;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.gatblau.onix.data.ItemData;
import org.gatblau.onix.data.LinkData;
import org.gatblau.onix.data.LinkedItemData;
import org.gatblau.onix.model.Dimension;
import org.gatblau.onix.model.Item;
import org.gatblau.onix.model.ItemType;
import org.gatblau.onix.model.Link;
import org.json.simple.JSONObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import javax.persistence.EntityManager;
import javax.persistence.NoResultException;
import javax.persistence.TypedQuery;
import java.io.IOException;
import java.time.ZonedDateTime;
import java.util.*;
import java.util.function.Consumer;

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
        ItemType itemType = em.getReference(ItemType.class, Integer.parseInt(json.get("itemTypeId").toString()));
        item.setKey(key);
        item.setItemType(itemType);
        item.setDescription((String)json.get("description"));
        item.setTag((String)json.get("tag"));
        item.setName((String)json.get("name"));
        item.setStatus(Short.parseShort(json.get("status").toString()));
        item.setUpdated(time);
        item.setMeta(mapper.valueToTree(json.get("meta")));

        em.persist(item);

        LinkedHashMap<String, String> dims =  (LinkedHashMap<String, String>)json.get("dimensions");
        if (dims != null) {
            if (action.equals("CREATED")) {
                Iterator<Map.Entry<String, String>> iterator = dims.entrySet().iterator();
                while (iterator.hasNext()) {
                    Map.Entry<String, String> entry = iterator.next();
                    Dimension d = new Dimension();
                    d.setItem(item);
                    d.setKey((String) entry.getKey());
                    d.setValue((String) entry.getValue());
                    em.persist(d);
                }
            }
            if (action.equals("UPDATED")) {
                Iterator<Map.Entry<String, String>> iterator = dims.entrySet().iterator();
                while (iterator.hasNext()) {
                    Map.Entry<String, String> entry = iterator.next();
                    TypedQuery<Dimension> dimQuery = em.createNamedQuery(Dimension.FIND_BY_KEY, Dimension.class);
                    dimQuery.setParameter(Dimension.PARAM_KEY, entry.getKey());
                    Dimension d;
                    try {
                        d = dimQuery.getSingleResult();
                    }
                    catch (NoResultException nre) {
                        // if the dimension does not exist, then creates one
                        d = new Dimension();
                    }
                    d.setKey(entry.getKey());
                    d.setValue(entry.getValue());
                    d.setItem(item);
                    em.persist(d);
                }
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
    public List<Item> findItemsByTypeAndTag(long itemTypeId, String tag) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE_AND_TAG, Item.class);
        query.setParameter(Item.PARAM_ITEM_TYPE_ID, itemTypeId);
        query.setParameter(Item.PARAM_TAG, "%" + tag + "%");
        return query.getResultList();
    }

    /***
     * Removes all transactional data in the database.
     */
    @Transactional
    public void clear() {
        if (em != null) {
            em.createNamedQuery(Dimension.DELETE_ALL).executeUpdate();
            em.createNamedQuery(Link.DELETE_ALL).executeUpdate();
            em.createNamedQuery(Item.DELETE_ALL).executeUpdate();
//            em.createNamedQuery(ItemType.DELETE_ALL).executeUpdate();
        }
    }

    @Transactional
    public void deleteItemTypes() {
        if (em != null) {
            em.createNamedQuery(ItemType.DELETE_ALL).executeUpdate();
        }
    }

    @Transactional
    public String createOrUpdateLink(String key, JSONObject json) throws IOException {
        String fromItemKey = (String)json.get("start_item_key");
        String toItemKey = (String)json.get("end_item_key");
        String action = "UPDATED";
        TypedQuery<Link> query = em.createNamedQuery(Link.FIND_BY_KEY, Link.class);
        query.setParameter(Link.KEY_LINK, key);
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

            link.setKey(key);
        }
        link.setDescription((String)json.get("description"));
        link.setMeta(mapper.valueToTree(json.get("meta")));
        link.setTag((String)json.get("tag"));
        link.setRole((String)json.get("role"));

        link.setUpdated(time);

        em.persist(link);

        return action;
    }

    @Transactional
    public void deleteLink(String key) {
        if (em != null) {
            TypedQuery<Link> query = em.createNamedQuery(Link.FIND_BY_KEY, Link.class);
            query.setParameter(Link.KEY_LINK, key);
            try {
                List<Link> links = query.getResultList();
                for (Link link : links) {
                    em.remove(link);
                }
            }
            catch (NoResultException nre){
            }
        }
    }

    @Transactional
    public ItemData getItem(String key) {
        Item item = getItemModel(key);
        if (item == null) return null;
        return mapItemDatum(item);
    }

    private Item getItemModel(String key) {
        Item item = null;
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
        query.setParameter(Item.PARAM_KEY, key);
        try {
            item = query.getSingleResult();
        }
        catch (NoResultException nre){
        }
        return item;
    }

    private LinkedItemData getLinkedItemData(Item item, boolean isParent) {
        LinkedItemData liData = new LinkedItemData();
        liData.setDescription(item.getDescription());
        liData.setKey(item.getKey());
        liData.setName(item.getName());
        liData.setParent(isParent);
        return liData;
    }

    /***
     * Get a list of links departing from or arriving at the passed-in item.
     * @param item the configuration item connected to the links to find.
     * @param isParent true if the links arrive at the item and false if the links depart from the item.
     * @return a list of @see org.gatblau.onix.data.LinkData.Class
     */
    private List<LinkData> getLinksData(Item item, boolean isParent) {
        TypedQuery<Link> itemQuery = null;
        if (isParent) {
            itemQuery = em.createNamedQuery(Link.FIND_FROM_ITEM, Link.class);
        } else {
            itemQuery = em.createNamedQuery(Link.FIND_TO_ITEM, Link.class);
        }
        itemQuery.setParameter(Link.KEY_ITEM_ID, item.getId());
        List<LinkData> linksData = new ArrayList<>();
        List<Link> links = null;
        try {
            links = itemQuery.getResultList();
            links.forEach(new Consumer<Link>() {
                @Override
                public void accept(Link link) {
                    LinkData linkData = new LinkData();
                    linkData.setDescription(link.getDescription());
                    linkData.setKey(link.getKey());
                    linkData.setMeta(link.getMeta());
                    linkData.setTag(link.getTag());
                    linkData.setRole(link.getRole());

                    String itemKey = null;

                    if (isParent) {
                        itemKey = link.getEndItem().getKey();
                    } else {
                        itemKey = link.getStartItem().getKey();
                    }

                    linkData.setItem(
                        getLinkedItemData(
                            getItemModel(itemKey), !isParent)
                    );

                    linksData.add(linkData);
                }
            });
        }
        catch (NoResultException nre){
        }
        return linksData;
    }

    public List<ItemData> getItemsByType(Integer typeId, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_ITEM_TYPE_ID, typeId);
        List<ItemData> data = mapItemData(query.getResultList());
        return data;
    }

    public List<ItemData> getItemsByTag(String tag, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TAG, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_TAG, "%" + tag + "%");
        List<ItemData> data = mapItemData(query.getResultList());
        return data;
    }

    public List<ItemData> getItemsByDate(ZonedDateTime from, ZonedDateTime to, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_DATE, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_FROM_DATE, from);
        query.setParameter(Item.PARAM_TO_DATE, to);
        List<ItemData> data = mapItemData(query.getResultList());
        return data;
    }

    public List<ItemData> getItemsByTypeAndTag(Integer typeId, String tag, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE_AND_TAG, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_ITEM_TYPE_ID, typeId);
        query.setParameter(Item.PARAM_TAG, "%" + tag + "%");
        List<ItemData> data = mapItemData(query.getResultList());
        return data;
    }

    public List<ItemData> getItemsByTypeAndDate(Integer typeId, ZonedDateTime from, ZonedDateTime to, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE_AND_DATE, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_FROM_DATE, from);
        query.setParameter(Item.PARAM_TO_DATE, to);
        List<ItemData> data = mapItemData(query.getResultList());
        return data;
    }

    public List<ItemData> getItemsByTypeTagAndDate(Integer typeId, String tag, ZonedDateTime from, ZonedDateTime to, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE_TAG_AND_DATE, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_ITEM_TYPE_ID, typeId);
        query.setParameter(Item.PARAM_TAG, "%" + tag + "%");
        query.setParameter(Item.PARAM_FROM_DATE, from);
        query.setParameter(Item.PARAM_TO_DATE, to);
        List<ItemData> data = mapItemData(query.getResultList());
        return data;
    }

    public List<ItemData> getAllByDateDesc(int maxResultEntries) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_ALL_BY_DATE_DESC, Item.class);
        query.setMaxResults(maxResultEntries);
        List<ItemData> data = mapItemData(query.getResultList());
        return data;
    }

    private List<ItemData> mapItemData(List<Item> items) {
        List<ItemData> data = new ArrayList<>();
        items.forEach(new Consumer<Item> () {
            @Override
            public void accept(Item item) {
                data.add(mapItemDatum(item));
            }
        });
        return data;
    }

    private ItemData mapItemDatum(Item item) {
        if (item == null) return null;

        ItemData data = new ItemData();
        data.setKey(item.getKey());
        data.setName(item.getName());
        data.setDescription(item.getDescription());

        data.setCreated(item.getCreated().toString());
        data.setUpdated(item.getUpdated().toString());
        data.setVersion(item.getVersion());

        data.setStatus(item.getStatus());
        data.setItemType(item.getItemType().getName());
        data.setMeta(item.getMeta());
        data.setTag(item.getTag());

        item.getDimensions().forEach(new Consumer<Dimension>() {
            @Override
            public void accept(Dimension dimension) {
                data.getDimensions().put(dimension.getKey(), dimension.getValue());
            }
        });

        // populate linked items here
        List<LinkData> links = new ArrayList<>();
        links.addAll(getLinksData(item, true)); // to links
        links.addAll(getLinksData(item, false)); // from links
        data.setLinks(links);

        return data;
    }

    public List<ItemType> getItemTypes() {
        TypedQuery<ItemType> itemTypesQuery = em.createNamedQuery(ItemType.FIND_ALL, ItemType.class);
        return itemTypesQuery.getResultList();
    }

    @Transactional
    public void deleteItem(String key) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
        query.setParameter(Item.PARAM_KEY, key);
        Item item = query.getSingleResult();
        for (Dimension dim : item.getDimensions()){
            em.remove(dim);
        }
        em.remove(item);
    }

    @Transactional
    public void createItemType(JSONObject json) throws IOException {
        ZonedDateTime time = ZonedDateTime.now();
        ItemType itemType = new ItemType();
        itemType.setDescription((String)json.get("description"));
        itemType.setName((String)json.get("name"));
        itemType.setCustom((Boolean)json.get("custom"));
        em.persist(itemType);
    }
}
