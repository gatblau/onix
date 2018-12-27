/*
Onix CMDB - Copyright (c) 2018 by www.gatblau.org

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Contributors to this project, hereby assign copyright in their code to the
project, to be licensed under the same terms as the rest of the code.
*/

package org.gatblau.onix;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.gatblau.onix.data.ItemData;
import org.gatblau.onix.data.LinkData;
import org.gatblau.onix.data.LinkList;
import org.gatblau.onix.data.LinkedItemData;
import org.gatblau.onix.model.Item;
import org.gatblau.onix.model.ItemType;
import org.gatblau.onix.model.Link;
import org.json.simple.JSONObject;
import org.postgresql.util.HStoreConverter;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import javax.persistence.*;
import java.io.IOException;
import java.sql.*;
import java.time.ZonedDateTime;
import java.util.*;
import java.util.function.Consumer;

@Service
public class Repository {
    private ObjectMapper mapper = new ObjectMapper();

    @Autowired
    private EntityManager em;

    @Autowired
    private Database db;

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

    public Result createOrUpdateItem(String key, JSONObject json) throws IOException, SQLException {
        Result result = new Result();

        Object name = json.get("name");
        Object description = json.get("description");
        Object meta = json.get("meta");
        Object tag = json.get("tag");
        Object attribute = json.get("attribute");
        Object status = json.get("status");
        Object type = json.get("type");

        String sql = "SELECT set_item(" +
            "?::character varying,\n" +
            "?::character varying,\n" +
            "?::text,\n" +
            "?::jsonb,\n" +
            "?::text[],\n" +
            "?::hstore,\n" +
            "?::smallint,\n" +
            "?::character varying,\n" +
            "?::bigint,\n" +
            "?::character varying)";
        Connection conn = null;
        PreparedStatement stmt = null;
        ResultSet set = null;
        try {
            conn = db.createConnection();
            stmt = conn.prepareStatement(sql);
            stmt.setString(1, key); // key_param
            stmt.setString(2, (name != null) ? (String) name : null); // name_param
            stmt.setString(3, (description != null) ? (String) description : null); // description_param
            stmt.setString(4, (meta != null) ? new JSONObject((LinkedHashMap<String, String>) meta).toJSONString() : null); // meta_param
            stmt.setString(5, (tag != null) ? toArrayStr((ArrayList<String>) tag) : null); // tag_param
            stmt.setString(6, (attribute != null) ? HStoreConverter.toString((LinkedHashMap<String, String>) attribute) : null); // attribute_param
            stmt.setInt(7, (status != null) ? (int) status : null); // status_param
            stmt.setString(8, (type != null) ? (String) type : null); // item_type_key_param
            stmt.setObject(9, null); // version_param
            stmt.setString(10, getUser()); // changedby_param
            set = stmt.executeQuery();
            if (set.next()){
                String r = set.getString("set_item");
                result.setOperation(r);
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            set.close();
            stmt.close();
            conn.close();
        }
        return result;
    }

    private String toArrayStr(List tag) {
        return null;
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

    @Transactional
    public LinkList getLinksByItem(String itemKey) {
        TypedQuery<Item> query = em.createNamedQuery(Item.FIND_BY_KEY, Item.class);
        query.setParameter(Item.PARAM_KEY, itemKey);
        Item item;

        try {
            item = query.getSingleResult();
            List<LinkData> ll = getLinksData(item, false);
            List<LinkData> lr = getLinksData(item, true);
            ll.forEach(new Consumer<LinkData>() {
                @Override
                public void accept(LinkData linkData) {
                    lr.add(linkData);
                }
            });
            return new LinkList(lr);
        }
        catch (NoResultException nre) {
        }
        return null;
    }

    private String getUser() {
        String username = null;
        Object principal = SecurityContextHolder.getContext().getAuthentication().getPrincipal();
        if (principal instanceof UserDetails) {
            UserDetails details = (UserDetails)principal;
            username = details.getUsername();
            for (GrantedAuthority a : details.getAuthorities()){
                username += "|" + a.getAuthority();
            };
        }
        else {
            username = principal.toString();
        }
        return username;
    }
}
