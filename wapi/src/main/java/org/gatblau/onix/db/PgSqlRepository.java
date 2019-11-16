/*
Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

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
package org.gatblau.onix.db;

import com.jayway.jsonpath.JsonPath;
import com.jayway.jsonpath.ReadContext;
import org.gatblau.onix.Lib;
import org.gatblau.onix.data.*;
import org.json.simple.JSONObject;
import org.postgresql.util.HStoreConverter;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.oauth2.jwt.JwtClaimAccessor;
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

    private JSONObject ready;

    public PgSqlRepository() {
    }

    /*
       ITEMS
     */

    @Override
    public synchronized Result createOrUpdateItem(String key, ItemData item, String[] role) {
        Result result = new Result(String.format("Item:%s", key));
        ResultSet set = null;
        try {
            db.prepare(getSetItemSQL());
            db.setString(1, key); // key_param
            db.setString(2, item.getName()); // name_param
            db.setString(3, item.getDescription()); // description_param
            db.setString(4, util.toJSONString(item.getMeta())); // meta_param
            db.setString(5, util.toArrayString(item.getTag())); // tag_param
            db.setString(6, getAttributeString(item.getAttribute())); // attribute_param
            db.setInt(7, item.getStatus()); // status_param
            db.setString(8, item.getType()); // item_type_key_param
            db.setObject(9, item.getVersion()); // version_param
            db.setString(10, getUser()); // changed_by_param
            db.setString(11, item.getPartition()); // partition_key_param
            db.setArray(12, role); // role_key_param
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_set_item"));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(String.format("Failed to create or update item with key '%s': %s", key, ex.getMessage()));
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public synchronized ItemData getItem(String key, boolean includeLinks, String[] role) {
        ItemData item = null;
        try {
            db.prepare(getGetItemSQL());
            db.setString(1, key);
            db.setArray(2, role);
            ResultSet set = db.executeQuerySingleRow();
            if (set != null) {
                item = util.toItemData(set);
            }

            if (includeLinks) {
                db.prepare(getFindLinksSQL());
                db.setString(1, item.getKey()); // start_item
                db.setObjectRange(2, 11, null);
                db.setArray(12, role);
                set = db.executeQuery();
                if (set != null) {
                    while (set.next()) {
                        item.getToLinks().add(util.toLinkData(set));
                    }
                }
                db.prepare(getFindLinksSQL());
                db.setString(1, null); // start_item
                db.setString(2, item.getKey()); // end_item
                db.setObjectRange(3, 11, null);
                db.setArray(12, role);
                set = db.executeQuery();
                if (set != null) {
                    while (set.next()) {
                        item.getFromLinks().add(util.toLinkData(set));
                    }
                }
            }
        } catch (Exception ex) {
            ex.printStackTrace();
        } finally {
            db.close();
            return item;
        }
    }

    @Override
    public Result deleteItem(String key, String[] role) {
        return delete(getDeleteItemSQL(), "ox_delete_item", key, role);
    }

    @Override
    public synchronized ItemList findItems(
            String itemTypeKey,
            List<String> tagList,
            ZonedDateTime createdFrom,
            ZonedDateTime createdTo,
            ZonedDateTime updatedFrom,
            ZonedDateTime updatedTo,
            Short status,
            String modelKey,
            Map<String, String> attributes,
            Integer top,
            String[] role
    ) {
        ItemList items = new ItemList();
        try {
            db.prepare(getFindItemsSQL());
            db.setString(1, (tagList != null && !tagList.isEmpty()) ? util.toArrayString(tagList) : null);
            db.setString(2, (attributes != null && !attributes.isEmpty()) ? getAttributeString(new JSONObject(attributes)) : null);
            db.setObject(3, status);
            db.setString(4, itemTypeKey);
            db.setObject(5, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(6, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(7, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(8, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
            db.setString(9, modelKey);
            db.setObject(10, (top == null) ? 20 : top);
            db.setArray(11, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                items.getValues().add(util.toItemData(set));
            }
        } catch (Exception ex) {
            ex.printStackTrace();
            throw new RuntimeException(String.format("Can't retrieve items: %s", ex.getMessage()));
        }
        return items;
    }

    @Override
    public synchronized JSONObject getItemMeta(String key, String filter, String[] role) {
        HashMap<String, Object> results = new HashMap<>();
        // gets the item in question
        ItemData item = getItem(key, false, role);
        if (filter == null) {
            // if the query does not specify a filter key then returns the plain metadata
            return item.getMeta();
        }
        // as a filter key has been passed in then tries and retrieves the filter expression for
        // the key from the itemType definition
        ItemTypeData itemType = getItemType(item.getType(), role);
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
                if (jsonPaths.size() > 1) {
                    // if there are more than one json paths defined, runs an extraction for each path
                    // and builds a map result object
                    for (JSONObject jsonPath : jsonPaths) {
                        HashMap.Entry<String, String> path = (HashMap.Entry<String, String>) jsonPath.entrySet().toArray()[0];
                        Object result = ctx.read(path.getValue());
                        results.put(path.getKey(), result);
                    }
                } else {
                    // if there is only on json path then return the single result as an object
                    HashMap.Entry<String, String> path = (HashMap.Entry<String, String>) jsonPaths.get(0).entrySet().toArray()[0];
                    return new JSONObject(ctx.read(path.getValue()));
                }
                break;
            }
        }
        return new JSONObject(results);
    }

    @Override
    public synchronized Result deleteAllItems(String[] role) {
        Result result = new Result();
        try {
            db.prepare(getDeleteAllItemsSQL());
            db.setArray(1, role);
            db.execute();
            result.setOperation("D");
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setMessage(ex.getMessage());
            result.setError(true);
        } finally {
            db.close();
        }
        return result;
    }

    /*
       LINKS
     */
    @Override
    public synchronized LinkData getLink(String key, String[] role) {
        LinkData link = null;
        try {
            db.prepare(getGetLinkSQL());
            db.setString(1, key);
            db.setArray(2, role);
            ResultSet set = db.executeQuerySingleRow();
            if (set != null) {
                link = util.toLinkData(set);
            }
        } catch (Exception ex) {
            ex.printStackTrace();
        } finally {
            db.close();
        }
        return link;
    }

    @Override
    public synchronized Result createOrUpdateLink(String key, LinkData link, String[] role) {
        Result result = new Result(String.format("Link:%s", key));
        try {
            db.prepare(getSetLinkSQL());
            db.setString(1, key);
            db.setString(2, link.getType());
            db.setString(3, link.getStartItemKey());
            db.setString(4, link.getEndItemKey());
            db.setString(5, link.getDescription());
            db.setString(6, util.toJSONString(link.getMeta()));
            db.setString(7, util.toArrayString(link.getTag()));
            db.setString(8, getAttributeString(link.getAttribute()));
            db.setObject(9, link.getVersion());
            db.setString(10, getUser());
            db.setArray(11, role);
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_set_link"));
        } catch (Exception ex) {
            result.setError(true);
            result.setMessage(String.format("Failed to create or update link with key '%s': %s", key, ex.getMessage()));
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public synchronized Result deleteLink(String key, String[] role) {
        return delete(getDeleteLinkSQL(), "ox_delete_link", key, role);
    }

    @Override
    public synchronized LinkList findLinks(
            String linkTypeKey,
            String startItemKey,
            String endItemKey,
            List<String> tagList,
            ZonedDateTime createdFrom,
            ZonedDateTime createdTo,
            ZonedDateTime updatedFrom,
            ZonedDateTime updatedTo,
            String modelKey,
            Integer top,
            String[] role
    ) {
        LinkList links = new LinkList();
        try {
            db.prepare(getFindLinksSQL());
            db.setString(1, startItemKey);
            db.setString(2, endItemKey);
            db.setString(3, util.toArrayString(tagList));
            db.setString(4, null); // attribute
            db.setString(5, linkTypeKey);
            db.setObject(6, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(7, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(8, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(9, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
            db.setString(10, modelKey);
            db.setObject(11, (top == null) ? 20 : top);
            db.setArray(12, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                links.getValues().add(util.toLinkData(set));
            }
        } catch (Exception ex) {
            ex.printStackTrace();
            throw new RuntimeException(String.format("Cant retrieve links: %s", ex.getMessage()));
        }
        return links;
    }

    @Override
    public synchronized Result clear(String[] role) {
        try {
            return delete(getClearAllSQL(), "ox_clear_all", null, role);
        } catch (Exception ex) {
            ex.printStackTrace();
            Result result = new Result("CLEAR_ALL");
            result.setError(true);
            result.setMessage(ex.getMessage());
            return result;
        }
    }

    private synchronized Result delete(String sql, String resultColName, String key, String[] role){
        return delete(sql, resultColName, key, false, role);
    }

    private synchronized Result delete(String sql, String resultColName, String key, boolean isType) {
        return delete(sql, resultColName, key, isType, null);
    }

    private synchronized Result delete(String sql, String resultColName, String key, boolean isType, String[] role) {
        Result result = new Result(String.format("Delete(%s)", key));
        try {
            db.prepare(sql);
            if (key != null) {
                int paramIx = 1;
                db.setString(paramIx, key);
                paramIx++;
                if (role != null) db.setArray(paramIx, role);
            } else {
                db.setArray(1, role);
            }
            result.setOperation(db.executeQueryAndRetrieveStatus(resultColName));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        } finally {
            db.close();
        }
        return result;
    }

    /*
        ITEM TYPES
     */
    @Override
    public synchronized ItemTypeData getItemType(String key, String[] role) {
        ItemTypeData itemType = null;
        try {
            db.prepare(getGetItemTypeSQL());
            db.setString(1, key);
            db.setArray(2, role);
            ResultSet set = db.executeQuerySingleRow();
            if (set != null) {
                itemType = util.toItemTypeData(set);
            }
        } catch (Exception ex) {
            throw new RuntimeException(String.format("Failed to get item type with key '%s': %s", key, ex.getMessage()), ex);
        } finally {
            db.close();
        }
        return itemType;
    }

    @Override
    public Result deleteItemTypes(String[] role) {
        return delete(getDeleteItemTypes(), "ox_delete_item_types", null, role);
    }

    @Override
    public synchronized ItemTypeList getItemTypes(
            Map attribute,
            ZonedDateTime createdFrom,
            ZonedDateTime createdTo,
            ZonedDateTime updatedFrom,
            ZonedDateTime updatedTo,
            String modelKey,
            String[] role
    ) {
        ItemTypeList itemTypes = new ItemTypeList();
        try {
            db.prepare(getFindItemTypesSQL());
            db.setString(1, util.toHStoreString(attribute)); // attribute_param
            db.setObject(2, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(3, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(4, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(5, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
            db.setString(6, modelKey);
            db.setArray(7, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                itemTypes.getValues().add(util.toItemTypeData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException(ex);
        } finally {
            db.close();
        }
        return itemTypes;
    }

    @Override
    public synchronized Result createOrUpdateItemType(String key, ItemTypeData itemType, String[] role) {
        Result result = new Result(String.format("ItemType:%s", key));
        try {
            db.prepare(getSetItemTypeSQL());
            db.setString(1, key); // key_param
            db.setString(2, itemType.getName()); // name_param
            db.setString(3, itemType.getDescription()); // description_param
            db.setString(4, getAttributeString(itemType.getAttrValid())); // attribute_param
            db.setString(5, util.toJSONString(itemType.getFilter()));
            db.setString(6, util.toJSONString(itemType.getMetaSchema()));
            db.setObject(7, itemType.getVersion()); // version_param
            db.setObject(8, itemType.getModelKey()); // meta model key
            db.setString(9, getUser()); // changed_by_param
            db.setArray(10, role);
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_set_item_type"));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setMessage(ex.getMessage());
            result.setError(true);
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result deleteItemType(String key, String[] role) {
        return delete(getDeleteItemTypeSQL(), "ox_delete_item_type", key, true, role);
    }

    /*
        LINK TYPES
     */
    @Override
    public synchronized LinkTypeList getLinkTypes(Map attribute, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, String modelKey, String[] role) {
        LinkTypeList linkTypes = new LinkTypeList();
        try {
            db.prepare(getFindLinkTypesSQL());
            db.setString(1, util.toHStoreString(attribute)); // attribute_param
            db.setObject(2, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(3, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(4, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(5, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
            db.setObject(6, modelKey);
            db.setArray(7, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                linkTypes.getValues().add(util.toLinkTypeData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException(ex);
        } finally {
            db.close();
        }
        return linkTypes;
    }

    @Override
    public synchronized Result createOrUpdateLinkType(String key, LinkTypeData linkType, String[] role) {
        Result result = new Result(String.format("LinkType:%s", key));
        try {
            db.prepare(getSetLinkTypeSQL());
            db.setString(1, key); // key_param
            db.setString(2, linkType.getName()); // name_param
            db.setString(3, linkType.getDescription()); // description_param
            db.setString(4, getAttributeString(linkType.getAttrValid())); // attribute_param
            db.setString(5, util.toJSONString(linkType.getMetaSchema()));
            db.setObject(6, linkType.getVersion()); // version_param
            db.setString(7, linkType.getModelKey()); // model_key_param
            db.setString(8, getUser()); // changed_by_param
            db.setArray(9, role);
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_set_link_type"));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setMessage(ex.getMessage());
            result.setError(true);
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result deleteLinkType(String key, String[] role) {
        return delete(getDeleteLinkTypeSQL(), "ox_delete_link_type", key, true, role);
    }

    @Override
    public Result deleteLinkTypes(String[] role) {
        return delete(getDeleteLinkTypes(), "ox_delete_link_types", null, role);
    }

    @Override
    public synchronized LinkTypeData getLinkType(String key, String[] role) {
        LinkTypeData linkType = null;
        try {
            db.prepare(getGetLinkTypeSQL());
            db.setString(1, key);
            db.setArray(2, role);
            ResultSet set = db.executeQuerySingleRow();
            if (set != null) {
                linkType = util.toLinkTypeData(set);
            }
        } catch (Exception ex) {
            throw new RuntimeException(String.format("Failed to get link type with key '%s': %s", key, ex.getMessage()), ex);
        } finally {
            db.close();
        }
        return linkType;
    }

    /*
        LINK RULES
     */
    @Override
    public synchronized LinkRuleList getLinkRules(
            String linkType,
            String startItemType,
            String endItemType,
            ZonedDateTime createdFrom,
            ZonedDateTime createdTo,
            ZonedDateTime updatedFrom,
            ZonedDateTime updatedTo,
            String[] role
        ) {
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
            db.setArray(8, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                linkRules.getValues().add(util.toLinkRuleData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException("Failed to get link rules", ex);
        } finally {
            db.close();
        }
        return linkRules;
    }

    @Override
    public synchronized Result createOrUpdateLinkRule(String key, LinkRuleData linkRule, String[] role) {
        Result result = new Result(String.format("LinkRule:%s", key));
        try {
            db.prepare(getSetLinkRuleSQL());
            db.setString(1, key); // key_param
            db.setString(2, linkRule.getName()); // name_param
            db.setString(3, linkRule.getDescription()); // description_param
            db.setString(4, linkRule.getLinkTypeKey()); // linkType_param
            db.setString(5, linkRule.getStartItemTypeKey()); // startItemType_param
            db.setString(6, linkRule.getEndItemTypeKey()); // endItemType_param
            db.setObject(7, linkRule.getVersion()); // version_param
            db.setString(8, getUser()); // changed_by_param
            db.setArray(9, role); // roel_key_param
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_set_link_rule"));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result deleteLinkRule(String key, String[] role) {
        return delete(getDeleteLinkRuleSQL(), "ox_delete_link_rule", key, role);
    }

    @Override
    public Result deleteLinkRules(String[] role) {
        return delete(getDeleteLinkRulesSQL(), "ox_delete_link_rules",null, role);
    }

    @Override
    public String getGetItemSQL() {
        return "SELECT * FROM ox_item(" +
                "?::character varying," + // key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getSetItemSQL() {
        return "SELECT ox_set_item(" +
                "?::character varying," +
                "?::character varying," +
                "?::text," +
                "?::jsonb," +
                "?::text[]," +
                "?::hstore," +
                "?::smallint," +
                "?::character varying," +
                "?::bigint," +
                "?::character varying," +
                "?::character varying," + // partition_key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getFindItemsSQL() {
        return "SELECT * FROM ox_find_items(" +
                "?::text[]," + // tag
                "?::hstore," + // attribute
                "?::smallint," + // status
                "?::character varying," + // item_type_key
                "?::timestamp with time zone," + // created_from
                "?::timestamp with time zone," + // created_to
                "?::timestamp with time zone," + // updated_from
                "?::timestamp with time zone," + // updated_to
                "?::character varying," + // model_key
                "?::integer," + // max_items
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getDeleteItemSQL() {
        return "SELECT ox_delete_item(" +
                "?::character varying," +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getDeleteAllItemsSQL() {
        return "SELECT ox_delete_all_items(" +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getDeleteLinkSQL() {
        return "SELECT ox_delete_link(" +
                "?::character varying," +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetLinkSQL() {
        return "SELECT * FROM ox_link(" +
                "?::character varying," +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getSetLinkSQL() {
        return "SELECT ox_set_link(" +
                "?::character varying," + // key
                "?::character varying," + // link_type_key
                "?::character varying," + // start_item_key
                "?::character varying," + // end_item_key
                "?::text," + // description
                "?::jsonb," + // meta
                "?::text[]," + // tag
                "?::hstore," + // attribute
                "?::bigint," + // version
                "?::character varying," + // changed_by
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getFindLinksSQL() {
        return "SELECT * FROM ox_find_links(" +
                "?::character varying," + // start_item_key_param
                "?::character varying," + // end_item_key_param
                "?::text[]," + // tag_param
                "?::hstore," + // attribute_param
                "?::character varying," + // link_type_key_param
                "?::timestamp with time zone," + // date_created_from_param
                "?::timestamp with time zone," + // date_created_to_param
                "?::timestamp with time zone," + // date_updated_from_param
                "?::timestamp with time zone," + // date_updated_to_param
                "?::character varying," + // model_key_param
                "?::integer," + // max_items
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getClearAllSQL() {
        return "SELECT ox_clear_all(" +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getDeleteItemTypeSQL() {
        return "SELECT ox_delete_item_type(" +
                "?::character varying," +
                "?::character varying[]" +
                ")";
    }

    @Override
    public String getDeleteItemTypes() {
        return "SELECT ox_delete_item_types(" +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getFindItemTypesSQL() {
        return "SELECT * FROM ox_find_item_types(" +
                "?::hstore," + // attr_valid
                "?::timestamp(6) with time zone," + // date created from
                "?::timestamp(6) with time zone," + // date created to
                "?::timestamp(6) with time zone," + // date updates from
                "?::timestamp(6) with time zone," + // date updated to
                "?::character varying," + // model key
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getSetItemTypeSQL() {
        return "SELECT ox_set_item_type(" +
                "?::character varying," + // key
                "?::character varying," + // name
                "?::text," + // description
                "?::hstore," + // attr_valid
                "?::jsonb," + // filter
                "?::jsonb," + // meta_schema
                "?::bigint," + // version
                "?::character varying," + // meta model key
                "?::character varying," + // changed_by
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetItemTypeSQL() {
        return "SELECT * FROM ox_item_type(" +
                "?::character varying," + // key
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getDeleteLinkTypeSQL() {
        return "SELECT ox_delete_link_type(" +
                "?::character varying," +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getDeleteLinkTypes() {
        return "SELECT ox_delete_link_types(" +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getFindLinkTypesSQL() {
        return "SELECT * FROM ox_find_link_types(" +
                "?::hstore," + // attr_valid
                "?::timestamp(6) with time zone," + // date created from
                "?::timestamp(6) with time zone," + // date created to
                "?::timestamp(6) with time zone," + // date updates from
                "?::timestamp(6) with time zone," + // date updated to
                "?::character varying," + // model key
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getSetLinkTypeSQL() {
        return "SELECT ox_set_link_type(" +
                "?::character varying," + // key
                "?::character varying," + // name
                "?::text," + // description
                "?::hstore," + // attr_valid
                "?::jsonb," + // meta_schema
                "?::bigint," + // version
                "?::character varying," + // model_key
                "?::character varying," + // changed_by
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetLinkTypeSQL() {
        return "SELECT * FROM ox_link_type(" +
                "?::character varying," + // key
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getDeleteLinkRuleSQL() {
        return "SELECT ox_delete_link_rule(?::character varying[])";
    }

    @Override
    public String getDeleteLinkRulesSQL() {
        return "SELECT ox_delete_link_rules(" +
                "?::character varying[]" + // role_key_param
                ")";
    }

    /* tags */
    @Override
    public synchronized Result createTag(JSONObject json) {
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
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_create_tag"));
            if (result.getOperation().equals("L")){
                result.setMessage(String.format("Tag data for label '%s' already exists and cannot be overridden.", label));
            }
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public synchronized Result updateTag(String rootItemKey, String currentLabel, JSONObject json) {
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
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_update_tag"));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public synchronized Result deleteTag(String rootItemKey, String label) {
        Result result = new Result(String.format("TAG:%s", rootItemKey));
        try {
            db.prepare(getDeleteTagSQL());
            db.setString(1, (rootItemKey != null) ? (String) rootItemKey : null); // root item key
            db.setString(2, (label != null) ? (String) label : null); // current_label
            result.setError(!db.execute());
            result.setOperation("D");
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public synchronized TagList getItemTags(String rootItemKey) {
        TagList tags = new TagList();
        try {
            db.prepare(getGetItemTagsSQL());
            db.setString(1, rootItemKey); // root_item_key_param
            ResultSet set = db.executeQuery();
            while (set.next()) {
                tags.getValues().add(util.toTagData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException(ex);
        } finally {
            db.close();
        }
        return tags;
    }

    @Override
    public synchronized GraphData getData(String rootItemKey, String label) {
        GraphData graph = new GraphData();
        try {
            db.prepare(getGetTreeItemsForTagSQL());
            db.setString(1, rootItemKey); // root_item_key_param
            db.setString(2, label); // label_param
            ResultSet set = db.executeQuery();
            while (set.next()) {
                graph.getItems().add(util.toItemData(set));
            }
            db.prepare(getGetTreeLinksForTagSQL());
            db.setString(1, rootItemKey); // root_item_key_param
            db.setString(2, label); // label_param
            set = db.executeQuery();
            while (set.next()) {
                graph.getLinks().add(util.toLinkData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException(ex);
        } finally {
            db.close();
        }
        return graph;
    }

    @Override
    public synchronized ResultList createOrUpdateData(GraphData payload, String[] role) {
        ResultList results = new ResultList();
        List<ModelData> models = payload.getModels();
        for (ModelData model : models) {
            Result result = createOrUpdateModel(model.getKey(), model, role);
            results.add(result);
        }
        List<ItemTypeData> itemTypes = payload.getItemTypes();
        for (ItemTypeData itemType : itemTypes) {
            Result result = createOrUpdateItemType(itemType.getKey(), itemType, role);
            results.add(result);
        }
        List<LinkTypeData> linkTypes = payload.getLinkTypes();
        for (LinkTypeData linkType : linkTypes) {
            Result result = createOrUpdateLinkType(linkType.getKey(), linkType, role);
            results.add(result);
        }
        List<LinkRuleData> linkRules = payload.getLinkRules();
        for (LinkRuleData linkRule : linkRules) {
            Result result = createOrUpdateLinkRule(linkRule.getKey(), linkRule, role);
            results.add(result);
        }
        List<ItemData> items = payload.getItems();
        for (ItemData item : items) {
            Result result = createOrUpdateItem(item.getKey(), item, role);
            results.add(result);
        }
        List<LinkData> links = payload.getLinks();
        for (LinkData link : links) {
            Result result = createOrUpdateLink(link.getKey(), link, role);
            results.add(result);
        }
        return results;
    }

    @Override
    public synchronized Result deleteData(String rootItemKey) {
        Result result = new Result(String.format("ItemTree:%s", rootItemKey));
        try {
            db.prepare(getDeleteItemTreeSQL());
            db.setString(1, (rootItemKey != null) ? (String) rootItemKey : null); // root item key
            result.setError(!db.execute());
            result.setOperation("D");
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(String.format("Failed to delete item tree for root item with key '%s': %s", rootItemKey, ex.getMessage()));
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public synchronized JSONObject checkReady() {
        if (ready == null) {
            ready = new JSONObject();
            Database.Version v;
            boolean freshInstall = false;
            int currentVersion, targetVersion = 0;
            try {
                // if db not created, then if auto-deploy=true try and create db and deploy schemas
                if (!db.exists()) {
                    freshInstall = true;
                    db.createDb();
                }
                // tries and gets the version information from the database
                v = db.getVersion();
                // if the schemas have not been deployed (no info found), prepares for a
                // fresh install by making the current version 0, otherwise the current version is
                // the one in the database
                currentVersion = (v.app == null) ? 0 : Integer.parseInt(v.db);
                // the target version comes from the manifest, and is the one required by the application
                targetVersion = db.getTargetDbVersion();
                if (currentVersion == targetVersion) {
                    // nothing to do
                } else if (currentVersion < targetVersion) {
                    // deploys the schemas/functions
                    db.deployDb(currentVersion, targetVersion);
                    // gets the deployed version
                    v = db.getVersion(true);
                } else {
                    // the db is newer than the app, the app must stop
                    throw new RuntimeException("The application is too old for this database: upgrade the application to a newer version.");
                }
            } catch (Exception ex) {
                // only if this is a failed fresh installation can delete db to go back to a clean state
                if (freshInstall) {
                    // if the process of deploying a brand new db failed, then remove the database
                    db.deleteDb();
                }
                throw new RuntimeException(ex);
            }
            ready.put("status", "ready");
            ready.put("appVersion", v.app);
            ready.put("dbVersion", v.db);
        }
        return ready;
    }

    @Override
    public synchronized Result deleteModel(String key, String[] role) {
        return delete(getDeleteModelSQL(), "ox_delete_model", key, true, role);
    }

    @Override
    public synchronized Result createOrUpdateModel(String key, ModelData model, String[] role) {
        Result result = new Result(String.format("Model:%s", key));
        try {
            db.prepare(getSetModelSQL());
            db.setString(1, key); // model key
            db.setString(2, model.getName()); // name_param
            db.setString(3, model.getDescription()); // description_param
            db.setObject(4, model.getVersion()); // version_param
            db.setString(5, getUser()); // changed_by_param
            db.setString(6, model.getPartition()); // partition_key_param
            db.setArray(7, role);
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_set_model"));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public synchronized ModelData getModel(String key, String[] role) {
        ModelData model = null;
        try {
            db.prepare(getGetModelSQL());
            db.setString(1, key);
            db.setArray(2, role);
            ResultSet set = db.executeQuerySingleRow();
            if (set != null) {
                model = util.toModelData(set);
            }
            db.close();
        } catch (Exception ex) {
            throw new RuntimeException(String.format("Failed to get model with key '%s': %s", key, ex.getMessage()), ex);
        } finally {
            db.close();
        }
        return model;
    }

    @Override
    public synchronized ModelDataList getModels(String[] role) {
        ModelDataList models = new ModelDataList();
        try {
            db.prepare(getGetModelsSQL());
            db.setArray(1, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                models.getValues().add(util.toModelData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException("Failed to retrieve models.", ex);
        }
        return models;
    }

    @Override
    public String getSetModelSQL() {
        return "SELECT ox_set_model(" +
                "?::character varying," + // key_param
                "?::character varying," + // name_param
                "?::text," + // description_param
                "?::bigint," + // version_param
                "?::character varying," + // changed_by
                "?::character varying," + // partition_key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetModelsSQL() {
        return "SELECT * FROM ox_get_models(" +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetModelSQL() {
        return "SELECT * FROM ox_model(" +
                "?::character varying," + // key
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public synchronized TypeGraphData getTypeDataByModel(String modelKey, String[] role) {
        TypeGraphData graph = new TypeGraphData();
        try {
            db.prepare(getGetModelItemTypesSQL());
            db.setString(1, modelKey); // model key param
            db.setArray(2, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                graph.getItemTypes().add(util.toItemTypeData(set));
            }
            db.prepare(getGetModelLinkTypesSQL());
            db.setString(1, modelKey); // model key param
            db.setArray(2, role);
            set = db.executeQuery();
            while (set.next()) {
                graph.getLinkTypes().add(util.toLinkTypeData(set));
            }
            db.prepare(getGetModelLinkRulesSQL());
            db.setString(1, modelKey); // model key param
            db.setArray(2, role);
            set = db.executeQuery();
            while (set.next()) {
                graph.getLinkRules().add(util.toLinkRuleData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException(ex);
        } finally {
            db.close();
        }
        return graph;
    }

    @Override
    public String getGetModelItemTypesSQL() {
        return "SELECT * FROM ox_get_model_item_types(" +
                "?::character varying," + // model_key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetModelLinkTypesSQL() {
        return "SELECT * FROM ox_get_model_link_types(" +
                "?::character varying," + // model_key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetModelLinkRulesSQL() {
        return "SELECT * FROM ox_get_model_link_rules(" +
                "?::character varying," + // model_key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getDeleteModelSQL() {
        return "SELECT ox_delete_model(" +
                "?::character varying, " + // model_key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getSetLinkRuleSQL() {
        return "SELECT ox_set_link_rule(" +
                "?::character varying," + // key
                "?::character varying," + // name
                "?::text," + // description
                "?::character varying," + // link_type
                "?::character varying," + // start_item_type
                "?::character varying," + // end_item_type
                "?::bigint," + // version
                "?::character varying," + // changed_by
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getFindLinkRulesSQL() {
        return "SELECT * FROM ox_find_link_rules(" +
                "?::character varying," +
                "?::character varying," +
                "?::character varying," +
                "?::timestamp(6) with time zone," +
                "?::timestamp(6) with time zone," +
                "?::timestamp(6) with time zone," +
                "?::timestamp(6) with time zone," +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getFindChildItemsSQL() {
        return "SELECT * FROM ox_find_child_items(" +
                "?::character varying," + // parent_item_key_param
                "?::character varying" + // link_type_key_param
                ")";
    }

    @Override
    public String getCreateTagSQL() {
        return "SELECT ox_create_tag(" +
                "?::character varying," + // root_item_key_param
                "?::character varying," + // tag_label_param
                "?::character varying," + // tag_name_param
                "?::text," + // tag_description_param
                "?::character varying" + // changed_by
                ")";
    }

    @Override
    public String getDeleteTagSQL() {
        return "SELECT ox_delete_tag(" +
                "?::character varying," + // root_item_key_param
                "?::character varying" + // tag_label_param
                ")";
    }

    @Override
    public String getUpdateTagSQL() {
        return "SELECT ox_update_tag(" +
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
        return "SELECT * FROM ox_get_item_tags(" +
                "?::character varying" + // root_item_key_param
                ")";
    }

    @Override
    public String getGetTreeItemsForTagSQL() {
        return "SELECT * FROM ox_get_tree_items(" +
                "?::character varying," + // root_item_key_param
                "?::character varying" + // tag_label_param
                ")";
    }

    @Override
    public String getGetTreeLinksForTagSQL() {
        return "SELECT * FROM ox_get_tree_links(" +
                "?::character varying," + // root_item_key_param
                "?::character varying" + // tag_label_param
                ")";
    }

    @Override
    public String getDeleteItemTreeSQL() {
        return "SELECT ox_delete_tree(" +
                "?::character varying" + // root_item_key_param
                ")";
    }

    @Override
    public String getTableCountSQL() {
        return "SELECT ox_get_table_count();";
    }

    private String getUser() {
        String username = null;
        Object principal = SecurityContextHolder.getContext().getAuthentication().getPrincipal();
        if (principal instanceof UserDetails) {
            UserDetails details = (UserDetails) principal;
            username = details.getUsername();
            for (GrantedAuthority a : details.getAuthorities()) {
                String r = a.getAuthority();
                if (r.startsWith("ROLE_")) {
                    r = r.substring("ROLE_".length());
                }
                username += "," + r;
            }
        } else if (principal instanceof JwtClaimAccessor){
            JwtClaimAccessor jwt = (JwtClaimAccessor)principal;
            username = jwt.getSubject();
            String[] roles = jwt.getClaimAsString("roles").split(",");
            for (String role : roles) {
                username += "," + role.trim();
            }
        } else if (principal instanceof String) {
            username = String.format("%s, ADMIN", principal);
        }
        return username;
    }

    private String getAttributeString(JSONObject json) {
        if (json != null) {
            return HStoreConverter.toString(json);
        }
        return null;
    }

    @Override
    public Result deletePartition(String key, String[] role) {
        return delete(getDeletePartitionSQL(), "ox_delete_partition", key, role);
    }

    @Override
    public Result createOrUpdatePartition(String key, PartitionData part, String[] role) {
        Result result = new Result(String.format("Partition:%s", key));
        ResultSet set = null;
        try {
            db.prepare(getSetPartitionSQL());
            db.setString(1, key); // key_param
            db.setString(2, part.getName()); // name_param
            db.setString(3, part.getDescription()); // description_param
            db.setObject(4, part.getVersion()); // version_param
            db.setString(5, getUser()); // changed_by_param
            db.setArray(6, role); // role_key_param
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_set_partition"));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(String.format("Failed to create or update partition with key '%s': %s", key, ex.getMessage()));
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public PartitionDataList getAllPartitions(String[] role) {
        PartitionDataList parts = new PartitionDataList();
        try {
            db.prepare(getGetAllPartitionsSQL());
            db.setArray(1, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                parts.getValues().add(util.toPartitionData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException("Failed to get partitions.", ex);
        } finally {
            db.close();
        }
        return parts;
    }

    @Override
    public PartitionData getPartition(String key, String[] role) {
        PartitionData part = null;
        try {
            db.prepare(getGetPartitionSQL());
            db.setString(1, key);
            db.setArray(2, role);
            ResultSet set = db.executeQuerySingleRow();
            part = util.toPartitionData(set);
            db.close();
        } catch (Exception ex) {
            throw new RuntimeException(String.format("Failed to get partition with key '%s': %s", key, ex.getMessage()), ex);
        } finally {
            db.close();
        }
        return part;
    }

    @Override
    public String getDeleteRoleSQL() {
        return "SELECT ox_delete_role(" +
                "?::character varying," + // key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getSetRoleSQL() {
        return "SELECT ox_set_role(" +
                "?::character varying," + // key_param
                "?::character varying," + // name_param
                "?::text," + // description_param
                "?::bigint," + // version_param
                "?::character varying," + // changed_by
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetRoleSQL() {
        return "SELECT * FROM ox_role(" +
                "?::character varying," + // key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetAllRolesSQL() {
        return "SELECT * FROM ox_get_roles(" +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public Result deleteRole(String key, String[] role) {
        return delete(getDeleteRoleSQL(), "ox_delete_role", key, role);
    }

    @Override
    public Result createOrUpdateRole(String key, RoleData roleData, String[] role) {
        Result result = new Result(String.format("Role:%s", key));
        ResultSet set = null;
        try {
            db.prepare(getSetRoleSQL());
            db.setString(1, key); // key_param
            db.setString(2, roleData.getName()); // name_param
            db.setString(3, roleData.getDescription()); // description_param
            db.setObject(4, roleData.getVersion()); // version_param
            db.setString(5, getUser()); // changed_by_param
            db.setArray(6, role); // role_key_param
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_set_role"));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(String.format("Failed to create or update role with key '%s': %s", key, ex.getMessage()));
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public RoleData getRole(String key, String[] role) {
        RoleData roleData = null;
        try {
            db.prepare(getGetRoleSQL());
            db.setString(1, key);
            db.setArray(2, role);
            ResultSet set = db.executeQuerySingleRow();
            roleData = util.toRoleData(set);
            db.close();
        } catch (Exception ex) {
            throw new RuntimeException(String.format("Failed to get role with key '%s': %s", key, ex.getMessage()), ex);
        } finally {
            db.close();
        }
        return roleData;
    }

    @Override
    public RoleDataList getAllRoles(String[] role) {
        RoleDataList parts = new RoleDataList();
        try {
            db.prepare(getGetAllRolesSQL());
            db.setArray(1, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                parts.getValues().add(util.toRoleData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException("Failed to get roles.", ex);
        } finally {
            db.close();
        }
        return parts;
    }

    @Override
    public Result addPrivilege(String partitionKey, String roleKey, NewPrivilegeData privilege, String[] role) {
        Result result = new Result(String.format("Privilege:%s:%s", roleKey, partitionKey));
        ResultSet set = null;
        try {
            db.prepare(getAddPrivilegeSQL());
            db.setString(1, partitionKey); // partition_key_param
            db.setString(2, roleKey); // role_key_param
            db.setObject(3, privilege.isCanCreate()); // can_create_param
            db.setObject(4, privilege.isCanRead()); // can_read_param
            db.setObject(5, privilege.isCanDelete()); // can_delete_param
            db.setString(6, getUser()); // changed_by_param
            db.setArray(7, role); // role_key_param
            result.setOperation(db.executeQueryAndRetrieveStatus("ox_add_privilege"));
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(String.format("Failed to add privilege with role '%s' and partition '%s': %s.", roleKey, partitionKey, ex.getMessage()));
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public Result removePrivilege(String partitionKey, String roleKey, String[] role) {
        Result result = new Result(String.format("Remove_Privilege_%s_%s", roleKey, partitionKey));
        try {
            db.prepare(getRemovePrivilegeSQL());
            db.setString(1, partitionKey);
            db.setString(2, roleKey);
            db.setArray(3, role);
            result.setOperation((db.execute()) ? "D" : "N");
            if (result.getOperation().equals("D")) {
                result.setChanged(true);
            }
        } catch (Exception ex) {
            ex.printStackTrace();
            result.setError(true);
            result.setMessage(ex.getMessage());
        } finally {
            db.close();
        }
        return result;
    }

    @Override
    public PrivilegeDataList getPrivilegesByRole(String roleKey, String[] loggedRoleKey) {
        PrivilegeDataList priv = new PrivilegeDataList();
        try {
            db.prepare(getGetAllPrivilegeByRoleSQL());
            db.setString(1, roleKey);
            db.setArray(2, loggedRoleKey);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                priv.getValues().add(util.toPrivilegeData(set));
            }
        } catch (Exception ex) {
            throw new RuntimeException("Failed to get privileges.", ex);
        } finally {
            db.close();
        }
        return priv;
    }

    @Override
    public ItemList getItemChildren(String key, String[] role) {
        ItemList items = new ItemList();
        try {
            db.prepare(getGetItemChildrenSQL());
            db.setString(1, key);
            db.setArray(2, role);
            ResultSet set = db.executeQuery();
            while (set.next()) {
                ItemData item = util.toItemData(set);
                item.setTypeName(set.getString("item_type_name"));
                items.getValues().add(item);
            }
        } catch (Exception ex) {
            throw new RuntimeException("Failed to get child items.", ex);
        } finally {
            db.close();
        }
        return items;
    }

    @Override
    public String getAddPrivilegeSQL() {
        return "SELECT ox_add_privilege(" +
                "?::character varying," + // role_key_param
                "?::character varying," + // privilege_key_param
                "?::boolean," + // can_create_param
                "?::boolean," + // can_read_param
                "?::boolean," + // can_delete_param
                "?::character varying," + // changed_by_param
                "?::character varying[]" + // logged_role_key_param
                ")";
    }

    @Override
    public String getGetItemChildrenSQL() {
        return "SELECT * FROM ox_get_item_children(" +
                "?::character varying," + // item_key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getRemovePrivilegeSQL() {
        return "SELECT ox_remove_privilege(" +
                "?::character varying," + // role_key_param
                "?::character varying," + // privilege_key_param
                "?::character varying[]" + // logged_role_key_param
                ")";
    }

    @Override
    public String getGetAllPrivilegeByRoleSQL() {
        return "SELECT * FROM ox_get_privileges_by_role(" +
                "?::character varying," + // privileges_role_key_param
                "?::character varying[]" + // logged_role_key_param
                ")";
    }

    @Override
    public String getDeletePartitionSQL() {
        return "SELECT ox_delete_partition(" +
                "?::character varying," + // key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getSetPartitionSQL() {
        return "SELECT ox_set_partition(" +
                "?::character varying," + // key_param
                "?::character varying," + // name_param
                "?::text," + // description_param
                "?::bigint," + // version_param
                "?::character varying," + // changed_by
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetAllPartitionsSQL() {
        return "SELECT * FROM ox_get_partitions(" +
                "?::character varying[]" + // role_key_param
                ")";
    }

    @Override
    public String getGetPartitionSQL() {
        return "SELECT * FROM ox_partition(" +
                "?::character varying," + // key_param
                "?::character varying[]" + // role_key_param
                ")";
    }

}
