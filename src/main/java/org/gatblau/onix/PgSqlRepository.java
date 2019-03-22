/*
Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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

import com.jayway.jsonpath.JsonPath;
import com.jayway.jsonpath.ReadContext;
import org.gatblau.onix.data.*;
import org.json.simple.JSONObject;
import org.postgresql.util.HStoreConverter;
import org.postgresql.util.PSQLException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;

import java.sql.ResultSet;
import java.time.ZonedDateTime;
import java.util.*;

@Service
public class PgSqlRepository implements DbRepository {

    @Autowired
    private Lib util;

    @Autowired
    private Database db;

    public PgSqlRepository() {
    }

    /*
       ITEMS
     */

    @Override
    public Result createOrUpdateItem(String key, JSONObject json) {
        Result result = new Result(String.format("Item:%s", key));
        ResultSet set = null;
        try {
            Object name = json.get("name");
            Object description = json.get("description");
            String meta = util.toJSONString(json.get("meta"));
            String tag = util.toArrayString(json.get("tag"));
            Object attribute = json.get("attribute");
            Object status = json.get("status");
            Object type = json.get("type");
            Object version = json.get("version");

            db.prepare(getSetItemSQL());
            db.setString(1, key); // key_param
            db.setString(2, (name != null) ? (String) name : null); // name_param
            db.setString(3, (description != null) ? (String) description : null); // description_param
            db.setString(4, meta); // meta_param
            db.setString(5, tag); // tag_param
            db.setString(6, (attribute != null) ? HStoreConverter.toString((LinkedHashMap<String, String>) attribute) : null); // attribute_param
            db.setInt(7, (status != null) ? (int) status : 0); // status_param
            db.setString(8, (type != null) ? (String) type : null); // item_type_key_param
            db.setObject(9, version); // version_param
            db.setString(10, getUser()); // changed_by_param
            result.setOperation(db.executeQueryAndRetrieveStatus("set_item"));
        }
        catch(Exception ex) {
            result.setError(true);
            result.setMessage(
                    String.format(
                        "Failed to create or update item with key '%s': %s", key, ex.getMessage()));
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public ItemData getItem(String key, boolean includeLinks) {
        ItemData item = new ItemData();
        try {
            db.prepare(getGetItemSQL());
            db.setString(1, key);
            item = util.toItemData(db.executeQuerySingleRow());

            if (includeLinks) {
                ResultSet set;

                db.prepare(getFindLinksSQL());
                db.setString(1, item.getKey()); // start_item
                db.setObjectRange(2, 9, null);
                set = db.executeQuery();
                while (set.next()) {
                    item.getToLinks().add(util.toLinkData(set));
                }

                db.prepare(getFindLinksSQL());
                db.setString(1, null); // start_item
                db.setString(2, item.getKey()); // end_item
                db.setObjectRange(3, 9, null);
                set = db.executeQuery();
                while (set.next()) {
                    item.getFromLinks().add(util.toLinkData(set));
                }
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
            return item;
        }
    }

    @Override
    public Result deleteItem(String key) {
        return delete(getDeleteItemSQL(), key);
    }

    @Override
    public ItemList findItems(String itemTypeKey, List<String> tagList, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, Short status, String modelKey, Integer top) {
        ItemList items = new ItemList();
        try {
            db.prepare(getFindItemsSQL());
            db.setString(1, util.toArrayString(tagList));
            db.setString(2, null); // attribute
            db.setObject(3, status);
            db.setString(4, itemTypeKey);
            db.setObject(5, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(6, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(7, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(8, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
            db.setString(9, modelKey);
            db.setObject(10, (top == null) ? 20 : top);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                items.getItems().add(util.toItemData(set));
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        return items;
    }

    @Override
    public JSONObject getItemMeta(String key, String filter) {
        HashMap<String, Object> results = new HashMap<>();
        // gets the item in question
        ItemData item = getItem(key, false);
        if (filter == null) {
            // if the query does not specify a filter key then returns the plain metadata
            return item.getMeta();
        }
        // as a filter key has been passed in then tries and retrieves the filter expression for
        // the key from the itemType definition
        ItemTypeData itemType = getItemType(item.getType());
        JSONObject f = itemType.getFilter();
        if (f == null) {
            // if the itemType does not define a filter then returns the plain whole metadata
            return item.getMeta();
        }
        // parses the json metadata into a read context in order to apply the json paths later
        ReadContext ctx = JsonPath.parse(item.getMeta());
        // starts processing the filter expression
        ArrayList<JSONObject> filters = (ArrayList) f.get("filters");
        for (JSONObject json : filters) {
            // each filter can have a set of values (json path expressions)
            // matches the filter key with the key in the list of predefined filters
            ArrayList<JSONObject> jsonPaths = (ArrayList) json.get(filter);
            if (jsonPaths != null) {
                // if there are json paths defined, runs an extraction for each path
                for (JSONObject jsonPath: jsonPaths) {
                    HashMap.Entry<String, String> path = (HashMap.Entry<String, String>) jsonPath.entrySet().toArray()[0];
                    Object result = ctx.read(path.getValue());
                    results.put(path.getKey(), result);
                }
                break;
            }
        }
        return new JSONObject(results);
    }

    @Override
    public Result deleteAllItems() {
        Result result = new Result();
        try {
            db.prepare(getDeleteAllItemsSQL());
            db.execute();
            result.setOperation("D");
        }
        catch (Exception ex) {
            ex.printStackTrace();
            result.setMessage(ex.getMessage());
            result.setError(true);
        }
        return result;
    }

    /*
       LINKS
     */
    @Override
    public LinkData getLink(String key) {
        LinkData link = null;
        try {
            db.prepare(getGetLinkSQL());
            db.setString(1, key);
            ResultSet set = db.executeQuerySingleRow();
            link = util.toLinkData(set);
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return link;
    }

    @Override
    public Result createOrUpdateLink(String key, JSONObject json) {
        Result result = new Result(String.format("Link:%s", key));
        try {
            String description = (String)json.get("description");
            String linkTypeKey = (String)json.get("type");
            String startItemKey = (String)json.get("startItemKey");
            String endItemKey = (String)json.get("endItemKey");
            String meta = util.toJSONString(json.get("meta"));
            String tag = util.toArrayString(json.get("tag"));
            Object attribute = json.get("attribute");
            Object version = json.get("version");
            Object model = json.get("modelKey");

            db.prepare(getSetLinkSQL());
            db.setString(1, key);
            db.setString(2, linkTypeKey);
            db.setString(3, startItemKey);
            db.setString(4, endItemKey);
            db.setString(5, description);
            db.setString(6, meta);
            db.setString(7, tag);
            db.setString(8, (attribute != null) ? HStoreConverter.toString((LinkedHashMap<String, String>) attribute) : null);
            db.setObject(9, version);
            db.setString(10, getUser());
            result.setOperation(db.executeQueryAndRetrieveStatus("set_link"));
        }
        catch (Exception ex) {
            result.setError(true);
            result.setMessage(String.format("Failed to create or update link with key '%s': %s", key, ex.getMessage()));
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result deleteLink(String key) {
        return delete(getDeleteLinkSQL(), key);
    }

    @Override
    public LinkList findLinks() {
        // TODO: implement findLinks()
        throw new UnsupportedOperationException("findLinks");
    }

    @Override
    public Result clear() {
        try {
            return delete(getClearAllSQL(), null);
        }
        catch (Exception ex) {
            ex.printStackTrace();
            Result result = new Result("CLEAR_ALL");
            result.setError(true);
            result.setMessage(ex.getMessage());
            return result;
        }
    }

    private Result delete(String sql, String key) {
        Result result = new Result(String.format("Delete(%s)", key));
        try {
            db.prepare(sql);
            if (key != null) {
                db.setString(1, key);
            }
            boolean deleted = db.execute();
            result.setOperation((deleted) ? "D" : "N");
        }
        catch (Exception ex) {
            result.setError(true);
            result.setMessage(ex.getMessage());
        }
        finally {
            db.close();
        }
        return result;
    }
    /*
        ITEM TYPES
     */
    @Override
    public ItemTypeData getItemType(String key) {
        ItemTypeData itemType = null;
        try {
            db.prepare(getGetItemTypeSQL());
            db.setString(1, key);
            ResultSet set = db.executeQuerySingleRow();
            itemType = util.toItemTypeData(set);
            db.close();
        }
        catch (Exception ex) {
            db.close();
            ex.printStackTrace();
            throw new RuntimeException(String.format("Failed to get item type with key '%s': %s", key, ex.getMessage()), ex);
        }
        return itemType;
    }

    @Override
    public Result deleteItemTypes() {
        return delete(getDeleteItemTypes(), null);
    }

    @Override
    public ItemTypeList getItemTypes(Map attribute, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, String modelKey) {
        try {
            ItemTypeList itemTypes = new ItemTypeList();
            db.prepare(getFindItemTypesSQL());
            db.setString(1, util.toHStoreString(attribute)); // attribute_param
            db.setObject(2, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(3, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(4, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(5, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
            db.setString(6, modelKey);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                itemTypes.getItems().add(util.toItemTypeData(set));
            }
            return itemTypes;
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return new ItemTypeList();
    }

    @Override
    public Result createOrUpdateItemType(String key, JSONObject json) {
        Result result = new Result(String.format("ItemType:%s", key));
        try {
            Object name = json.get("name");
            Object description = json.get("description");
            Object attribute = json.get("attrValid");
            String filter = util.toJSONString(json.get("filter"));
            String metaSchema = util.toJSONString(json.get("metaSchema"));
            Object modelKey = json.get("modelKey");
            Object version = json.get("version");

            db.prepare(getSetItemTypeSQL());
            db.setString(1, key); // key_param
            db.setString(2, (name != null) ? (String) name : null); // name_param
            db.setString(3, (description != null) ? (String) description : null); // description_param
            db.setString(4, (attribute != null) ? HStoreConverter.toString((LinkedHashMap<String, String>) attribute) : null); // attribute_param
            db.setString(5, filter);
            db.setString(6, metaSchema);
            db.setObject(7, version); // version_param
            db.setObject(8, modelKey); // meta model key
            db.setString(9, getUser()); // changed_by_param
            result.setOperation(db.executeQueryAndRetrieveStatus("set_item_type"));
        }
        catch (Exception ex) {
            ex.printStackTrace();
            result.setMessage(ex.getMessage());
            result.setError(true);
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result deleteItemType(String key) {
        return delete(getDeleteItemTypeSQL(), key);
    }

    /*
        LINK TYPES
     */
    @Override
    public LinkTypeList getLinkTypes(Map attribute, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, String modelKey) {
        LinkTypeList linkTypes = new LinkTypeList();
        try {
            db.prepare(getFindLinkTypesSQL());
            db.setString(1, util.toHStoreString(attribute)); // attribute_param
            db.setObject(2, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(3, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(4, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(5, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
            db.setObject(6, modelKey);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                linkTypes.getItems().add(util.toLinkTypeData(set));
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return linkTypes;
    }

    @Override
    public Result createOrUpdateLinkType(String key, JSONObject json) {
        Result result = new Result(String.format("LinkType:%s", key));
        try {
            Object name = json.get("name");
            Object description = json.get("description");
            Object attribute = json.get("attrValid");
            String metaSchema = util.toJSONString(json.get("metaSchema"));
            Object version = json.get("version");
            Object modelKey = json.get("modelKey");
            db.prepare(getSetLinkTypeSQL());
            db.setString(1, key); // key_param
            db.setString(2, (name != null) ? (String) name : null); // name_param
            db.setString(3, (description != null) ? (String) description : null); // description_param
            db.setString(4, (attribute != null) ? HStoreConverter.toString((LinkedHashMap<String, String>) attribute) : null); // attribute_param
            db.setString(5, metaSchema);
            db.setObject(6, version); // version_param
            db.setString(7, (modelKey != null) ? (String)modelKey : null); // model_key_param
            db.setString(8, getUser()); // changed_by_param
            result.setOperation(db.executeQueryAndRetrieveStatus("set_link_type"));
        }
        catch (Exception ex) {
            ex.printStackTrace();
            result.setMessage(ex.getMessage());
            result.setError(true);
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result deleteLinkType(String key) {
        return delete(getDeleteLinkTypeSQL(), key);
    }

    @Override
    public Result deleteLinkTypes() {
        return delete(getDeleteLinkTypes(), null);
    }

    @Override
    public LinkTypeData getLinkType(String key) {
        LinkTypeData linkType = null;
        try {
            db.prepare(getGetLinkTypeSQL());
            db.setString(1, key);
            ResultSet set = db.executeQuerySingleRow();
            linkType = util.toLinkTypeData(set);
            db.close();
        }
        catch (Exception ex) {
            db.close();
            ex.printStackTrace();
            throw new RuntimeException(String.format("Failed to get link type with key '%s': %s", key, ex.getMessage()), ex);
        }
        return linkType;
    }

    /*
        LINK RULES
     */
    @Override
    public LinkRuleList getLinkRules(String linkType, String startItemType, String endItemType, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo) {
        LinkRuleList linkRules = new LinkRuleList();
        try {
            db.prepare(getFindLinkRulesSQL());
            db.setString(1, linkType); // link_type key
            db.setString(2, startItemType); // start item_type key
            db.setString(3, endItemType); // end item_type key
            db.setObject(4, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(5, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(6, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(7, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                linkRules.getItems().add(util.toLinkRuleData(set));
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return linkRules;
    }

    @Override
    public Result createOrUpdateLinkRule(String key, JSONObject json) {
        Result result = new Result(String.format("LinkRule:%s", key));
        Object name = json.get("name");
        Object description = json.get("description");
        Object linkType = json.get("linkTypeKey");
        Object startItemType = json.get("startItemTypeKey");
        Object endItemType = json.get("endItemTypeKey");
        Object version = json.get("version");
        try {
            db.prepare(getSetLinkRuleSQL());
            db.setString(1, key); // key_param
            db.setString(2, (name != null) ? (String) name : null); // name_param
            db.setString(3, (description != null) ? (String) description : null); // description_param
            db.setString(4, (linkType != null) ? (String) linkType : null); // linkType_param
            db.setString(5, (startItemType != null) ? (String) startItemType : null); // startItemType_param
            db.setString(6, (endItemType != null) ? (String) endItemType : null); // endItemType_param
            db.setObject(7, version); // version_param
            db.setString(8, getUser()); // changed_by_param
            result.setOperation(db.executeQueryAndRetrieveStatus("set_link_rule"));
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result deleteLinkRule(String key) {
        return delete(getDeleteLinkRuleSQL(), key);
    }

    @Override
    public Result deleteLinkRules() {
        return delete(getDeleteLinkRulesSQL(), null);
    }

    /*
        CHANGE
     */
    @Override
    public List<ChangeItemData> findChangeItems() {
        // TODO: implement findChangeItems()
        throw new UnsupportedOperationException("findChangeItems");
    }

    private ItemList getChildItems(String parentKey) {
        ItemList items = new ItemList();
        try {
            db.prepare(getFindChildItemsSQL());
            db.setString(1, parentKey); // parent_key_param
            db.setString(2, "ANSIBLE_INVENTORY"); // item_type_key_param
            ResultSet set = db.executeQuery();
            while (set.next()) {
                items.getItems().add(util.toItemData(set));
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return items;
    }

    private String prefix(String prefix, String str) {
        return String.format("%s::%s", prefix, str);
    }

    private JSONObject getLinkData(String description, String linkType, String startItem, String endItem) {
        JSONObject json = new JSONObject();
        json.put("description", description);
        json.put("type", linkType);
        json.put("startItemKey", startItem);
        json.put("endItemKey", endItem);
        return json;
    }

    private JSONObject getItemData(String name, String description, String type, JSONObject meta) {
        JSONObject json = new JSONObject();
        json.put("name", name);
        json.put("description", description);
        json.put("type", type);
        json.put("meta", meta);
        return json;
    }

    @Override
    public String getGetItemSQL() {
        return "SELECT * FROM item(?::character varying)";
    }

    @Override
    public String getSetItemSQL() {
        return "SELECT set_item(" +
                "?::character varying," +
                "?::character varying," +
                "?::text," +
                "?::jsonb," +
                "?::text[]," +
                "?::hstore," +
                "?::smallint," +
                "?::character varying," +
                "?::bigint," +
                "?::character varying)";
    }

    @Override
    public String getFindItemsSQL() {
        return "SELECT * FROM find_items(" +
                "?::text[]," + // tag
                "?::hstore," + // attribute
                "?::smallint," + // status
                "?::character varying," + // item_type_key
                "?::timestamp with time zone," + // created_from
                "?::timestamp with time zone," + // created_to
                "?::timestamp with time zone," + // updated_from
                "?::timestamp with time zone," + // updated_to
                "?::character varying," + // model_key
                "?::integer" + // max_items
                ")";
    }

    @Override
    public String getDeleteItemSQL() {
        return "SELECT delete_item(?::character varying)";
    }

    @Override
    public String getDeleteAllItemsSQL() {
        return "SELECT delete_all_items()";
    }

    @Override
    public String getDeleteLinkSQL() {
        return "SELECT delete_link(?::character varying)";
    }

    @Override
    public String getGetLinkSQL() {
        return "SELECT * FROM link(?::character varying)";
    }

    @Override
    public String getSetLinkSQL() {
        return "SELECT set_link(" +
                "?::character varying," + // key
                "?::character varying," + // link_type_key
                "?::character varying," + // start_item_key
                "?::character varying," + // end_item_key
                "?::text," + // description
                "?::jsonb," + // meta
                "?::text[]," + // tag
                "?::hstore," + // attribute
                "?::bigint," + // version
                "?::character varying" + // changed_by
                ")";
    }

    @Override
    public String getFindLinksSQL() {
        return "SELECT * FROM find_links(" +
                "?::character varying," +
                "?::character varying," +
                "?::text[]," +
                "?::hstore," +
                "?::character varying," +
                "?::timestamp with time zone," +
                "?::timestamp with time zone," +
                "?::timestamp with time zone," +
                "?::timestamp with time zone" +
                ")";
    }

    @Override
    public String getClearAllSQL() {
        return "SELECT clear_all()";
    }

    @Override
    public String getDeleteItemTypeSQL() {
        return "SELECT delete_item_type(?::character varying)";
    }

    @Override
    public String getDeleteItemTypes() {
        return "SELECT delete_item_types()";
    }

    @Override
    public String getFindItemTypesSQL() {
        return "SELECT * FROM find_item_types(" +
            "?::hstore," + // attr_valid
            "?::timestamp(6) with time zone," + // date created from
            "?::timestamp(6) with time zone," + // date created to
            "?::timestamp(6) with time zone," + // date updates from
            "?::timestamp(6) with time zone," + // date updated to
            "?::character varying" + // model key
        ")";
    }

    @Override
    public String getSetItemTypeSQL() {
        return "SELECT set_item_type(" +
                "?::character varying," + // key
                "?::character varying," + // name
                "?::text," + // description
                "?::hstore," + // attr_valid
                "?::jsonb," + // filter
                "?::jsonb," + // meta_schema
                "?::bigint," + // version
                "?::character varying," + // meta model key
                "?::character varying" + // changed_by
                ")";
    }

    @Override
    public String getGetItemTypeSQL() {
        return "SELECT * FROM item_type(" +
                "?::character varying" + // key
                ")";
    }

    @Override
    public String getDeleteLinkTypeSQL() {
        return "SELECT delete_link_type(?::character varying)";
    }

    @Override
    public String getDeleteLinkTypes() {
        return "SELECT delete_link_types()";
    }

    @Override
    public String getFindLinkTypesSQL() {
        return "SELECT * FROM find_link_types(" +
                "?::hstore," + // attr_valid
                "?::timestamp(6) with time zone," + // date created from
                "?::timestamp(6) with time zone," + // date created to
                "?::timestamp(6) with time zone," + // date updates from
                "?::timestamp(6) with time zone," + // date updated to
                "?::character varying" + // model key
                ")";
    }

    @Override
    public String getSetLinkTypeSQL() {
        return "SELECT set_link_type(" +
                "?::character varying," + // key
                "?::character varying," + // name
                "?::text," + // description
                "?::hstore," + // attr_valid
                "?::jsonb," + // meta_schema
                "?::bigint," + // version
                "?::character varying," + // model_key
                "?::character varying" + // changed_by
                ")";
    }

    @Override
    public String getGetLinkTypeSQL() {
        return "SELECT * FROM link_type(" +
                "?::character varying" + // key
                ")";
    }

    @Override
    public String getDeleteLinkRuleSQL() {
        return "SELECT delete_link_rule(?::character varying)";
    }

    @Override
    public String getDeleteLinkRulesSQL() {
        return "SELECT delete_link_rules()";
    }

    /* tags */
    @Override
    public Result createTag(JSONObject json) {
        Result result = new Result("CREATE_TAG");
        Object name = json.get("name");
        Object description = json.get("description");
        Object label = json.get("label");
        Object rootItemKey = json.get("rootItemKey");
        try {
            db.prepare(getCreateTagSQL());
            db.setString(1, (rootItemKey != null) ? (String) rootItemKey : null); // root item key
            db.setString(3, (name != null) ? (String) name : null); // name_param
            db.setString(4, (description != null) ? (String) description : null); // description_param
            db.setString(2, (label != null) ? (String) label : null); // label
            db.setString(5, getUser()); // changed_by_param
            result.setError(!db.execute());
            result.setOperation("I");
        }
        catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result updateTag(String rootItemKey, String currentLabel, JSONObject json) {
        Result result = new Result(String.format("TAG:%s", rootItemKey));
        Object name = json.get("name");
        Object description = json.get("description");
        Object newLabel = json.get("label");
        Object version = json.get("version");
        try {
            db.prepare(getUpdateTagSQL());
            db.setString(1, (rootItemKey != null) ? (String) rootItemKey : null); // root item key
            db.setString(2, (currentLabel != null) ? (String) currentLabel : null); // current_label
            db.setString(3, (newLabel != null) ? (String) newLabel : null); // new_label
            db.setString(4, (name != null) ? (String) name : null); // name_param
            db.setString(5, (description != null) ? (String) description : null); // description_param
            db.setString(6, getUser()); // changed_by_param
            db.setObject(7, version); // version_param
            result.setOperation(db.executeQueryAndRetrieveStatus("update_tag"));
        }
        catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result deleteTag(String rootItemKey, String label) {
        Result result = new Result(String.format("TAG:%s", rootItemKey));
        try {
            db.prepare(getDeleteTagSQL());
            db.setString(1, (rootItemKey != null) ? (String) rootItemKey : null); // root item key
            db.setString(2, (label != null) ? (String) label : null); // current_label
            result.setError(!db.execute());
            result.setOperation("D");
        }
        catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public TagList getItemTags(String rootItemKey) {
        TagList tags = new TagList();
        try {
            db.prepare(getGetItemTagsSQL());
            db.setString(1, rootItemKey); // root_item_key_param
            ResultSet set = db.executeQuery();
            while (set.next()) {
                tags.getItems().add(util.toTagData(set));
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return tags;
    }

    @Override
    public ItemTreeData getItemTree(String rootItemKey, String label) {
        ItemTreeData tree = new ItemTreeData();
        try {
            db.prepare(getGetTreeItemsForTagSQL());
            db.setString(1, rootItemKey); // root_item_key_param
            db.setString(2, label); // label_param
            ResultSet set = db.executeQuery();
            while (set.next()) {
                tree.getItems().add(util.toItemData(set));
            }
            db.prepare(getGetTreeLinksForTagSQL());
            db.setString(1, rootItemKey); // root_item_key_param
            db.setString(2, label); // label_param
            set = db.executeQuery();
            while (set.next()) {
                tree.getLinks().add(util.toLinkData(set));
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return tree;
    }

    @Override
    public ResultList createOrUpdateData(JSONObject payload) {
        ResultList results = new ResultList();
        Object modelsObject = payload.get("models");
        if (modelsObject != null) {
            ArrayList<LinkedHashMap> models = (ArrayList<LinkedHashMap>)modelsObject;
            for (Map model : models) {
                String key = (String)model.get("key");
                Result result = createOrUpdateModel(key, new JSONObject(model));
                results.getItems().add(result);
            }
        }
        Object itemTypesObject = payload.get("itemTypes");
        if (itemTypesObject != null) {
            ArrayList<LinkedHashMap> itemTypes = (ArrayList<LinkedHashMap>)itemTypesObject;
            for (Map itemType : itemTypes) {
                String key = (String)itemType.get("key");
                Result result = createOrUpdateItemType(key, new JSONObject(itemType));
                results.getItems().add(result);
            }
        }
        Object linkTypesObject = payload.get("linkTypes");
        if (linkTypesObject != null) {
            ArrayList<LinkedHashMap> linkTypes = (ArrayList<LinkedHashMap>)linkTypesObject;
            for (Map linkType : linkTypes) {
                String key = (String)linkType.get("key");
                Result result = createOrUpdateLinkType(key, new JSONObject(linkType));
                results.getItems().add(result);
            }
        }
        Object linkRulesObject = payload.get("linkRules");
        if (linkRulesObject != null) {
            ArrayList<LinkedHashMap> linkRules = (ArrayList<LinkedHashMap>)linkRulesObject;
            for (Map linkRule : linkRules) {
                String key = (String)linkRule.get("key");
                Result result = createOrUpdateLinkRule(key, new JSONObject(linkRule));
                results.getItems().add(result);
            }
        }
        Object itemsObject = payload.get("items");
        if (itemsObject != null) {
            ArrayList<LinkedHashMap> items = (ArrayList<LinkedHashMap>)itemsObject;
            for (Map item: items) {
                String key = (String)item.get("key");
                Result result = createOrUpdateItem(key, new JSONObject(item));
                results.getItems().add(result);
            }
        }
        Object linksObject = payload.get("links");
        if (linksObject != null) {
            ArrayList<LinkedHashMap> links = (ArrayList<LinkedHashMap>) linksObject;
            for (Map link : links) {
                String key = (String) link.get("key");
                Result result = createOrUpdateLink((String) link.get("key"), new JSONObject(link));
                results.getItems().add(result);
            }
        }
        return results;
    }

    @Override
    public Result deleteItemTree(String rootItemKey) {
        Result result = new Result(String.format("ItemTree:%s", rootItemKey));
        try {
            db.prepare(getDeleteItemTreeSQL());
            db.setString(1, (rootItemKey != null) ? (String) rootItemKey : null); // root item key
            result.setError(!db.execute());
            result.setOperation("D");
        }
        catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(
                String.format("Failed to delete item tree for root item with key '%s': %s", rootItemKey, ex.getMessage()));
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public JSONObject getReadyStatus() {
        JSONObject status = new JSONObject();
        try{
            db.prepare(getTableCountSQL());
            ResultSet set = db.executeQuerySingleRow();
            while (set.next()) {
                int count = set.getInt("get_table_count");
                if (count == 0) {
                    throw new RuntimeException("No tables found in the database.");
                }
            }
            status.put("ready", true);
        }
        catch (Exception ex) {
            status.put("ready", false);
            status.put("error", ex.getMessage());
        }
        return status;
    }

    @Override
    public Result deleteModel(String key) {
        Result result = new Result(String.format("Model:%s", key));
        try {
            db.prepare(getDeleteModelSQL());
            db.setString(1, key); // meta model key
            result.setError(!db.execute());
            result.setOperation("D");
        }
        catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(
                String.format("Failed to delete model for key '%s': %s", key, ex.getMessage()));
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result createOrUpdateModel(String key, JSONObject json) {
        Result result = new Result(String.format("Model:%s", key));
        Object name = json.get("name");
        Object description = json.get("description");
        Object version = json.get("version");
        try {
            db.prepare(getSetModelSQL());
            db.setString(1, key); // model key
            db.setString(2, (name != null) ? (String) name : null); // name_param
            db.setString(3, (description != null) ? (String) description : null); // description_param
            db.setObject(4, version); // version_param
            db.setString(5, getUser()); // changed_by_param
            result.setError(!db.execute());
            result.setOperation("I");
        }
        catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public ModelData getModel(String key) {
        ModelData model = null;
        try {
            db.prepare(getGetModelSQL());
            db.setString(1, key);
            ResultSet set = db.executeQuerySingleRow();
            model = util.toModelData(set);
            db.close();
        }
        catch (Exception ex) {
            db.close();
            ex.printStackTrace();
            throw new RuntimeException(String.format("Failed to get model with key '%s': %s", key, ex.getMessage()), ex);
        }
        return model;
    }

    @Override
    public ModelDataList getModels() {
        ModelDataList models = new ModelDataList();
        try {
            db.prepare(getGetModelsSQL());
            ResultSet set = db.executeQuery();
            while (set.next()) {
                models.getItems().add(util.toModelData(set));
            }
        } catch (Exception ex) {
            ex.printStackTrace();
            throw new RuntimeException("Failed to retrieve models.", ex);
        }
        return models;
    }

    @Override
    public String getSetModelSQL() {
        return "SELECT set_model(" +
                "?::character varying," + // key_param
                "?::character varying," + // name_param
                "?::text," + // description_param
                "?::bigint," + // version_param
                "?::character varying" + // changed_by
                ")";
    }

    @Override
    public String getGetModelsSQL() {
        return "SELECT * FROM get_models()";
    }

    @Override
    public String getGetModelSQL() {
        return "SELECT * FROM model(" +
                "?::character varying" + // key
                ")";
    }

    @Override
    public String getDeleteModelSQL() {
        return "SELECT delete_model(" +
                "?::character varying" + // model_key_param
                ")";
    }

    @Override
    public String getSetLinkRuleSQL() {
        return "SELECT set_link_rule(" +
                "?::character varying," + // key
                "?::character varying," + // name
                "?::text," + // description
                "?::character varying," + // link_type
                "?::character varying," + // start_item_type
                "?::character varying," + // end_item_type
                "?::bigint," + // version
                "?::character varying" + // changed_by
                ")";
    }

    @Override
    public String getFindLinkRulesSQL() {
        return "SELECT * FROM find_link_rules(" +
                    "?::character varying," +
                    "?::character varying," +
                    "?::character varying," +
                    "?::timestamp(6) with time zone," +
                    "?::timestamp(6) with time zone," +
                    "?::timestamp(6) with time zone," +
                    "?::timestamp(6) with time zone" +
                ")";
    }

    @Override
    public String getFindChildItemsSQL() {
        return "SELECT * FROM find_child_items(" +
                    "?::character varying," + // parent_item_key_param
                    "?::character varying" + // link_type_key_param
                ")";
    }

    @Override
    public String getCreateTagSQL() {
        return "SELECT create_tag(" +
                "?::character varying," + // root_item_key_param
                "?::character varying," + // tag_label_param
                "?::character varying," + // tag_name_param
                "?::text," + // tag_description_param
                "?::character varying" + // changed_by
                ")";
    }

    @Override
    public String getDeleteTagSQL() {
        return "SELECT delete_tag(" +
                "?::character varying," + // root_item_key_param
                "?::character varying" + // tag_label_param
                ")";
    }

    @Override
    public String getUpdateTagSQL() {
        return "SELECT update_tag(" +
                "?::character varying," + // root_item_key_param
                "?::character varying," + // current_label_param
                "?::character varying," + // new_label_param
                "?::character varying," + // tag_name_param
                "?::text," + // tag_description_param
                "?::character varying," + // changed_by_param
                "?::bigint" + // version_param
                ")";
    }

    @Override
    public String getGetItemTagsSQL() {
        return "SELECT * FROM get_item_tags(" +
                "?::character varying" + // root_item_key_param
                ")";
    }

    @Override
    public String getGetTreeItemsForTagSQL() {
        return "SELECT * FROM get_tree_items(" +
                "?::character varying," + // root_item_key_param
                "?::character varying" + // tag_label_param
                ")";
    }

    @Override
    public String getGetTreeLinksForTagSQL() {
        return "SELECT * FROM get_tree_links(" +
                "?::character varying," + // root_item_key_param
                "?::character varying" + // tag_label_param
                ")";
    }

    @Override
    public String getDeleteItemTreeSQL() {
        return "SELECT delete_tree(" +
                "?::character varying" + // root_item_key_param
                ")";
    }

    @Override
    public String getTableCountSQL() {
        return "SELECT get_table_count();";
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
