package org.gatblau.onix;

import com.fasterxml.jackson.databind.JsonNode;
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
    public Result createOrUpdateItem(String key, JSONObject json) throws IOException {
        Result result = new Result();
        result.setChanged(false);
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
        query.setParameter(Item.PARAM_KEY, key);
        Item item = null;
        ZonedDateTime time = ZonedDateTime.now();

        ItemType itemType;
        String type = "";
        try {
            type = (String)json.get("type");
            itemType = getItemType(type);
        }
        catch (NoResultException nre) {
            result.setChanged(false);
            result.setError(true);
            result.setMessage(String.format("Cannot create or update item %s: Item Type %s is not defined.", key, type));
            return result;
        }
        try {
            item = query.getSingleResult();

            if (!itemType.getKey().equals(item.getItemType().getKey())) {
                item.setItemType(itemType);
                result.setChanged(true);
            }

            String value = (String)json.get("description");

            if (!item.getDescription().equals(value)) {
                item.setDescription(value);
                result.setChanged(true);
            }

            value = (String)json.get("tag");
            if (!item.getTag().equals(value)) {
                item.setTag(value);
                result.setChanged(true);
            }

            value = (String)json.get("name");
            if (!item.getName().equals(value)) {
                item.setName(value);
                result.setChanged(true);
            }

            Short valueShort = Short.parseShort(json.get("status").toString());
            if (!item.getStatus().equals(valueShort)) {
                item.setStatus(valueShort);
                result.setChanged(true);
            }

            JsonNode node = mapper.valueToTree(json.get("meta"));
            if (!item.getMeta().equals(node)) {
                item.setMeta(node);
                result.setChanged(true);
            }
        }
        catch (NoResultException e) {
            item = new Item();
            item.setKey(key);
            item.setCreated(time);
            item.setItemType(itemType);
            item.setName((String)json.get("name"));
            item.setDescription(ifNullThenEmpty((String)json.get("description")));
            item.setTag(ifNullThenEmpty((String)json.get("tag")));
            item.setMeta(mapper.valueToTree(json.get("meta")));
            item.setStatus(Short.parseShort(json.get("status").toString()));

            result.setChanged(true);
            result.setMessage(String.format("Item %s has been CREATED.", key));
            result.setOperation("C");
        }
        catch (Exception ex) {
            result.setChanged(false);
            result.setError(true);
            result.setMessage(String.format("Failed to create or update Item %s: %s.", key, ex.getMessage()));
            return result;
        }

        if (result.isChanged()) {
            if (result.getOperation() == null) {
                result.setMessage(String.format("Item %s has been UPDATED.", key));
                result.setOperation("U");
            }
            item.setUpdated(time);
            try {
                em.persist(item);
            }
            catch (Exception ex) {
                result.setChanged(false);
                result.setError(true);
                result.setMessage(String.format("Failed to create or update Item %s: %s.", key, ex.getMessage()));
                return result;
            }
        }
        else {
            result.setMessage(String.format("Item %s has not changed and does not need updating.", key));
        }

        LinkedHashMap<String, String> dims =  (LinkedHashMap<String, String>)json.get("dimensions");
        if (dims != null) {
            if (result.getOperation() != null && result.getOperation().equals("C")) {
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
            if (result.getOperation() == null || result.getOperation().equals("U")) {
                Iterator<Map.Entry<String, String>> iterator = dims.entrySet().iterator();
                while (iterator.hasNext()) {
                    boolean changed = false;
                    Map.Entry<String, String> entry = iterator.next();
                    TypedQuery<Dimension> dimQuery = em.createNamedQuery(Dimension.FIND_BY_KEY, Dimension.class);
                    dimQuery.setParameter(Dimension.PARAM_KEY, entry.getKey());
                    dimQuery.setParameter(Dimension.PARAM_ITEM_KEY, item.getKey());
                    Dimension d;
                    try {
                          d = dimQuery.getSingleResult();
                          if (!d.getValue().equals(entry.getValue())) {
                              d.setValue(entry.getValue());
                              result.setChanged(true);
                              result.setOperation("U");
                          }
                    }
                    catch (NoResultException nre) {
                        // if the dimension does not exist, then creates one
                        d = new Dimension();
                        d.setKey(entry.getKey());
                        d.setValue(entry.getValue());
                        d.setItem(item);
                        result.setChanged(true);
                        result.setOperation("C");
                    }

                    if (result.isChanged()) {
                        em.persist(d);
                        String m;
                        if (result.getOperation().equals("C")) {
                            m = String.format("Dimension %s was created.", entry.getKey());
                        }
                        else {
                            m = String.format("Dimension %s was updated.", entry.getKey());
                        }
                        result.setMessage(String.format("%s %s", result.getMessage(), m));
                    }
                }
            }
        }
        return result;
    }

    private ItemType getItemType(String type) {
        TypedQuery<ItemType> itquery = em.createNamedQuery(ItemType.FIND_BY_KEY, ItemType.class);
        itquery.setParameter(ItemType.PARAM_KEY, type);
        return itquery.getSingleResult();
    }

    /***
     * Find all nodes of a particular type and with a specific tags.
     * @param itemTypeKey
     * @param tag the tags used to filter the result of the search.
     * @return
     */
    public List<Item> findItemsByTypeAndTag(String itemTypeKey, String tag) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE_AND_TAG, Item.class);
        query.setParameter(Item.PARAM_ITEM_TYPE_KEY, itemTypeKey);
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
    public Result createOrUpdateLink(String key, JSONObject json) throws IOException {
        String fromItemKey = (String)json.get("start_item_key");
        String toItemKey = (String)json.get("end_item_key");
        Result result = new Result();
        TypedQuery<Link> query = em.createNamedQuery(Link.FIND_BY_KEY, Link.class);
        query.setParameter(Link.KEY_LINK, key);
        Link link = null;
        ZonedDateTime time = ZonedDateTime.now();
        try {
            link = query.getSingleResult();
            String value = (String)json.get("description");

            if (!link.getDescription().equals(value)) {
                link.setDescription(value);
                result.setChanged(true);
            }

            JsonNode jsonValue = mapper.valueToTree(json.get("meta"));

            if (!link.getMeta().equals(jsonValue)) {
                link.setMeta(jsonValue);
                result.setChanged(true);
            }

            value = (String)json.get("tag");

            if (!link.getTag().equals(value)) {
                link.setTag(value);
                result.setChanged(true);
            }

            value = (String)json.get("role");
            if (!link.getRole().equals(value)) {
                link.setRole(value);
                result.setChanged(true);
            }

            if (result.isChanged()) {
                result.setMessage(String.format("Link %s has been UPDATED.", key));
                result.setOperation("U");
            }
            else {
                result.setMessage(String.format("Nothing to update. Link %s has not changed.", key));
                result.setOperation("U");
            }
        }
        catch (NoResultException e) {
            link = new Link();
            link.setCreated(time);

            TypedQuery<Item> startItemQuery = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
            startItemQuery.setParameter(Item.PARAM_KEY, fromItemKey);
            Item startItem = null;
            try {
                startItem = startItemQuery.getSingleResult();
            }
            catch (NoResultException nre) {
                result.setChanged(false);
                result.setError(true);
                result.setMessage("Could not create link to start configuration item with key '" + fromItemKey + "' as it does not exist.");
                return result;
            }

            TypedQuery<Item> endItemQuery = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
            endItemQuery.setParameter(Item.PARAM_KEY, toItemKey);
            Item endItem = null;
            try {
                endItem = endItemQuery.getSingleResult();
            }
            catch (NoResultException nre) {
                result.setChanged(false);
                result.setError(true);
                result.setMessage("Could not create link to end configuration item with key '" + toItemKey + "' as it does not exist.");
                return result;
            }

            link.setStartItem(startItem);
            link.setEndItem(endItem);
            link.setKey(key);
            link.setDescription(ifNullThenEmpty((String)json.get("description")));
            link.setMeta(mapper.valueToTree(json.get("meta")));
            link.setTag(ifNullThenEmpty((String)json.get("tag")));
            link.setRole(ifNullThenEmpty((String)json.get("role")));

            result.setChanged(true);
            result.setMessage(String.format("Link %s has been CREATED.", key));
            result.setOperation("C");
        }

        if (result.isChanged()) {
            try {
                em.persist(link);
                link.setUpdated(time);
            }
            catch (Exception ex) {
                result.setChanged(false);
                result.setError(true);
                result.setMessage(String.format("Failed to create or update link %s: %s.", key, ex.getMessage()));
                return result;
            }
        }
        return result;
    }

    @Transactional
    public Result deleteLink(String key) {
        Result result = new Result();
        result.setChanged(false);
        result.setOperation("D");
        TypedQuery<Link> query = em.createNamedQuery(Link.FIND_BY_KEY, Link.class);
        query.setParameter(Link.KEY_LINK, key);
        try {
            Link link = query.getSingleResult();
            em.remove(link);
            result.setChanged(true);
            result.setMessage(String.format("Link %s has been deleted.", key));
        }
        catch (NoResultException nre){
            result.setChanged(false);
            result.setError(false);
            result.setMessage(String.format("Nothing to delete. Link %s has not been found.", key));
        }
        catch (Exception ex) {
            result.setChanged(false);
            result.setError(true);
            result.setMessage(String.format("Failed to delete Link %s: %s.", key, ex.getMessage()));
        }
        return result;
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

    public List<ItemData> getItemsByType(String itemTypeKey, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_ITEM_TYPE_KEY, itemTypeKey);
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

    public List<ItemData> getItemsByTypeAndTag(String itemTypeKey, String tag, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE_AND_TAG, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_ITEM_TYPE_KEY, itemTypeKey);
        query.setParameter(Item.PARAM_TAG, "%" + tag + "%");
        List<ItemData> data = mapItemData(query.getResultList());
        return data;
    }

    public List<ItemData> getItemsByTypeAndDate(String itemTypeKey, ZonedDateTime from, ZonedDateTime to, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE_AND_DATE, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_ITEM_TYPE_KEY, itemTypeKey);
        query.setParameter(Item.PARAM_FROM_DATE, from);
        query.setParameter(Item.PARAM_TO_DATE, to);
        List<ItemData> data = mapItemData(query.getResultList());
        return data;
    }

    public List<ItemData> getItemsByTypeTagAndDate(String itemTypeKey, String tag, ZonedDateTime from, ZonedDateTime to, Integer top) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_TYPE_TAG_AND_DATE, Item.class);
        if (top != null) query.setMaxResults(top);
        query.setParameter(Item.PARAM_ITEM_TYPE_KEY, itemTypeKey);
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
    public Result deleteItem(String key) {
        Result result = new Result();
        result.setChanged(false);
        result.setOperation("D");

        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
        query.setParameter(Item.PARAM_KEY, key);
        Item item;

        try {
            item = query.getSingleResult();
        }
        catch (NoResultException nre) {
            result.setError(false);
            result.setChanged(false);
            result.setMessage(String.format("Nothing to delete. Cannot find Item %s.", key));
            return result;
        }

        // precondition: cant delete item if links exist
        if (getLinksData(item, false).size() > 0) {
            result.setError(true);
            result.setChanged(false);
            result.setMessage(String.format("Cannot delete Item %s because it is linked to other items.", key));
            return result;
        }
        if (getLinksData(item, true).size() > 0) {
            result.setError(true);
            result.setChanged(false);
            result.setMessage(String.format("Cannot delete Item %s because it is linked to other items.", key));
            return result;
        }

        try {
            for (Dimension dim : item.getDimensions()) {
                em.remove(dim);
            }
            em.remove(item);
            result.setError(false);
            result.setChanged(true);
            result.setMessage(String.format("Item %s has been deleted.", key));
        }
        catch (Exception ex) {
            result.setError(true);
            result.setChanged(false);
            result.setMessage(String.format("Failed to delete Item %s: %s", key, ex.getMessage()));
        }
        return result;
    }

    @Transactional
    public Result createOrUpdateItemType(String key, JSONObject json) throws IOException {
        Result result = new Result();
        result.setChanged(false);

        ZonedDateTime time = ZonedDateTime.now();
        TypedQuery<ItemType> itquery = em.createNamedQuery(ItemType.FIND_BY_KEY, ItemType.class);
        itquery.setParameter(ItemType.PARAM_KEY, key);

        ItemType itemType;
        try {
            itemType = itquery.getSingleResult();

            String value = (String)json.get("description");

            if (!itemType.getDescription().equals(value)) {
                itemType.setDescription((String)json.get("description"));
                result.setChanged(true);
            }

            value = (String)json.get("name");
            if (!itemType.getName().equals(value)) {
                itemType.setName(value);
                result.setChanged(true);
            }

            itemType.setCustom(true);
            itemType.setUpdated(time);

            if (result.isChanged()) {
                result.setMessage(String.format("Item Type %s has been UPDATED.", key));
            }
            else {
                result.setMessage(String.format("Item Type %s already exists and does not require any updates.", key));
            }

            result.setOperation("U");
        }
        catch (NoResultException nre) {
            itemType = new ItemType();
            itemType.setKey(key);
            itemType.setName((String)json.get("name"));
            itemType.setDescription(ifNullThenEmpty((String)json.get("description")));
            itemType.setCreated(time);
            result.setChanged(true);
            result.setMessage(String.format("Item Type %s has been CREATED.", key));
            result.setOperation("C");
        }

        if (result.isChanged()) {
            try {
                em.persist(itemType);
            }
            catch (Exception ex) {
                result.setError(true);
                result.setChanged(false);
                result.setMessage(String.format("Failed to create or update Item Type %s: %s.", key, ex.getMessage()));
            }
        }

        return result;
    }

    @Transactional
    public Result deleteItemType(String key) {
        Result result = new Result();
        result.setOperation("D");
        // precondition: cant delete type if items exist of the type
        if (getItemsByType(key, 1).size() > 0) {
            result.setChanged(false);
            result.setMessage(String.format("Cannot delete Item Type %s because there are items still using it.", key));
            return result;
        }
        TypedQuery<ItemType> itQuery = em.createNamedQuery(ItemType.FIND_BY_KEY, ItemType.class);
        itQuery.setParameter(ItemType.PARAM_KEY, key);
        ItemType itemType;
        try {
            itemType = itQuery.getSingleResult();
            em.remove(itemType);
            result.setChanged(true);
            result.setMessage(String.format("Item Type %s has been deleted.", key));
        }
        catch (NoResultException nre) {
            result.setError(false);
            result.setChanged(false);;
            result.setMessage(String.format("Item Type %s not found.", key));
        }
        catch (Exception ex) {
            result.setError(true);
            result.setChanged(false);
            result.setMessage(String.format("Failed to delete Item Type %s: %s.", key, ex.getMessage()));
        }
        return result;
    }

    private String ifNullThenEmpty(String value) {
        return (value == null) ? "" : value;
    }
}
