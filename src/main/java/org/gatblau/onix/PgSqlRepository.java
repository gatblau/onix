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

import org.gatblau.onix.data.*;
import org.gatblau.onix.inv.Inventory;
import org.gatblau.onix.inv.Node;
import org.json.simple.JSONObject;
import org.json.simple.parser.ParseException;
import org.postgresql.util.HStoreConverter;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.time.ZonedDateTime;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

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
        Result result = new Result();
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
            result.setMessage(ex.getMessage());
        }
        finally {
            db.close();
        }
        return result;
    }

    @Override
    public ItemData getItem(String key) {
        try {
            db.prepare(getGetItemSQL());
            db.setString(1, key);
            ItemData item = util.toItemData(db.executeQuerySingleRow());

            ResultSet set;

            db.prepare(getFindLinksSQL());
            db.setString(1, item.getKey()); // start_item
            db.setObjectRange(2, 9, null);
            set = db.executeQuery();
            while (set.next()) {
                item.getFromLinks().add(util.toLinkData(set));
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
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return new ItemData();
    }

    @Override
    public Result deleteItem(String key) {
        return delete(getDeleteItemSQL(), key);
    }

    @Override
    public ItemList findItems(String itemTypeKey, List<String> tagList, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, Short status, Integer top) {
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
            db.setObject(9, (top == null) ? 20 : top);
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
        Result result = new Result();
        try {
            String description = (String)json.get("description");
            String linkTypeKey = (String)json.get("linkType");
            String startItemKey = (String)json.get("startItem");
            String endItemKey = (String)json.get("endItem");
            String meta = util.toJSONString(json.get("meta"));
            String tag = util.toArrayString(json.get("tag"));
            Object attribute = json.get("attribute");
            Object version = json.get("version");

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
            result.setMessage(ex.getMessage());
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
            Result result = new Result();
            result.setError(true);
            result.setMessage(ex.getMessage());
            return result;
        }
    }

    private Result delete(String sql, String key) {
        Result result = new Result();
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
        try {
            db.prepare(getGetItemTypeSQL());
            db.setString(1, key);
            ResultSet set = db.executeQuerySingleRow();
            return util.toItemTypeData(set);
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return new ItemTypeData();
    }

    @Override
    public Result deleteItemTypes() {
        return delete(getDeleteItemTypes(), null);
    }

    @Override
    public ItemTypeList getItemTypes(Map attribute, Boolean system, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo) {
        try {
            ItemTypeList itemTypes = new ItemTypeList();
            db.prepare(getFindItemTypesSQL());
            db.setString(1, util.toHStoreString(attribute)); // attribute_param
            db.setObject(2, system);
            db.setObject(3, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(4, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(5, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(6, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
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
        Result result = new Result();
        Object name = json.get("name");
        Object description = json.get("description");
        Object attribute = json.get("attribute_validation");
        Object version = json.get("version");
        try {
            db.prepare(getSetItemTypeSQL());
            db.setString(1, key); // key_param
            db.setString(2, (name != null) ? (String) name : null); // name_param
            db.setString(3, (description != null) ? (String) description : null); // description_param
            db.setString(4, (attribute != null) ? HStoreConverter.toString((LinkedHashMap<String, String>) attribute) : null); // attribute_param
            db.setObject(5, version); // version_param
            db.setString(6, getUser()); // changed_by_param
            result.setOperation(db.executeQueryAndRetrieveStatus("set_item_type"));
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
    public Result deleteItemType(String key) {
        return delete(getDeleteItemTypeSQL(), key);
    }

    /*
        LINK TYPES
     */
    @Override
    public LinkTypeList getLinkTypes(Map attribute, Boolean system, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo) {
        LinkTypeList linkTypes = new LinkTypeList();
        try {
            db.prepare(getFindLinkTypesSQL());
            db.setString(1, util.toHStoreString(attribute)); // attribute_param
            db.setObject(2, system);
            db.setObject(3, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(4, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(5, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(6, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
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
        Result result = new Result();
        Object name = json.get("name");
        Object description = json.get("description");
        Object attribute = json.get("attribute_validation");
        Object version = json.get("version");
        try {
            db.prepare(getSetItemTypeSQL());
            db.setString(1, key); // key_param
            db.setString(2, (name != null) ? (String) name : null); // name_param
            db.setString(3, (description != null) ? (String) description : null); // description_param
            db.setString(4, (attribute != null) ? HStoreConverter.toString((LinkedHashMap<String, String>) attribute) : null); // attribute_param
            db.setObject(5, version); // version_param
            db.setString(6, getUser()); // changed_by_param
            result.setOperation(db.executeQueryAndRetrieveStatus("set_item_type"));
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
    public Result deleteLinkType(String key) {
        return delete(getDeleteLinkTypeSQL(), key);
    }

    @Override
    public Result deleteLinkTypes() {
        return delete(getDeleteLinkTypes(), null);
    }

    @Override
    public LinkTypeData getLinkType(String key) {
        // TODO: implement getLinkType()
        throw new UnsupportedOperationException("getLinkType");
    }

    /*
        LINK RULES
     */
    @Override
    public LinkRuleList getLinkRules(String linkType, String startItemType, String endItemType, Boolean system, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo) {
        LinkRuleList linkRules = new LinkRuleList();
        try {
            db.prepare(getFindLinkRulesSQL());
            db.setString(1, linkType); // link_type key
            db.setString(2, startItemType); // start item_type key
            db.setString(3, endItemType); // end item_type key
            db.setObject(4, system); // system
            db.setObject(5, (createdFrom != null) ? java.sql.Date.valueOf(createdFrom.toLocalDate()) : null);
            db.setObject(6, (createdTo != null) ? java.sql.Date.valueOf(createdTo.toLocalDate()) : null);
            db.setObject(7, (updatedFrom != null) ? java.sql.Date.valueOf(updatedFrom.toLocalDate()) : null);
            db.setObject(8, (updatedTo != null) ? java.sql.Date.valueOf(updatedTo.toLocalDate()) : null);
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
        Result result = new Result();
        Object name = json.get("name");
        Object description = json.get("description");
        Object linkType = json.get("linkType");
        Object startItemType = json.get("startItemType");
        Object endItemType = json.get("endItemType");
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

    @Override
    public Result createOrUpdateInventory(String key, String inventory) {
        try {
            Inventory inv = new Inventory(inventory);
            Result result = new Result();
            result = createOrUpdateItem(key, getItemData(key, "Inventory imported from Ansible inventory file.", "INVENTORY", new JSONObject()));
            if (result.isError()) return result;
            for (Node node : inv.getNodes()) {
                processNode(node, null, key);
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        return null;
    }

    private Result processNode(Node node, Node parent, String key) throws ParseException, SQLException, IOException {
        String nodeType = getNodeType(node);
        JSONObject vars = node.getVarsJSON();
        Result result = createOrUpdateItem(prefix(key, node.getName()), getItemData(node.getName(), String.format("%s imported from Ansible inventory.", node.getType()), nodeType, vars));
        if (result.isError()) return result;
        switch (node.getType()) {
            case PARENT_GROUP:
            case GROUP: {
                if (parent == null) { // parent is the inventory node - link to inventory
                    result = createOrUpdateLink(
                            prefix(key, String.format("%s->%s", key, node.getName())),
                            getLinkData("Link imported from Ansible inventory.", "INVENTORY", key, prefix(key, node.getName()))
                    );
                    if (result.isError()) return result;
                } else { // link to a group
                    result = createOrUpdateLink(
                            prefix(key, String.format("%s->%s", parent.getName(), node.getName())),
                            getLinkData("Link imported from Ansible inventory.", "INVENTORY", prefix(key, parent.getName()), prefix(key, node.getName()))
                    );
                    if (result.isError()) return result;
                }
                for (Node child : node.getChildren()) {
                    processNode(child, node, key);
                }
                break;
            }
            case HOST:{
                result = createOrUpdateLink(
                    prefix(key, String.format("%s->%s", parent.getName(), node.getName())),
                    getLinkData("Link imported from Ansible inventory.", "INVENTORY", prefix(key, parent.getName()), prefix(key, node.getName())));
                if (result.isError()) return result;
                break;
            }
        }
        return result;
    }

    private String getNodeType(Node node) {
        String groupType = "";
        switch (node.getType()) {
            case PARENT_GROUP:
                groupType = "HOST-GROUP-GROUP";
                break;
            case GROUP:
                groupType = "HOST-GROUP";
                break;
            case HOST:
                groupType = "HOST";
                break;
        }
        return groupType;
    }

    @Override
    public String getInventory(String key) {
        StringBuilder builder = new StringBuilder();
        ItemList items = getChildItems(key);
        for (ItemData item : items.getItems()) {
            if (item.getItemType().equalsIgnoreCase("HOST-GROUP-GROUP")) {
                builder
                        .append("[")
                        .append(item.getName())
                        .append(":children]")
                        .append(System.getProperty("line.separator"));
            } else {

            }
            builder.append(item.getName()).append(System.getProperty("line.separator"));
        }
        return null;
    }

    private ItemList getChildItems(String parentKey) {
        ItemList items = new ItemList();
        try {
            db.prepare(getFindChildItemsSQL());
            db.setString(1, parentKey); // parent_key_param
            db.setString(2, "INVENTORY"); // item_type_key_param
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
        json.put("linkType", linkType);
        json.put("startItem", startItem);
        json.put("endItem", endItem);
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
                "?::integer" + // max_items
                ")";
    }

    @Override
    public String getDeleteItemSQL() {
        return "SELECT delete_item(?::character varying)";
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
            "?::boolean," + // system
            "?::timestamp(6) with time zone," + // date created from
            "?::timestamp(6) with time zone," + // date created to
            "?::timestamp(6) with time zone," + // date updates from
            "?::timestamp(6) with time zone" + // date updated to
        ")";
    }

    @Override
    public String getSetItemTypeSQL() {
        return "SELECT set_item_type(" +
                "?::character varying," + // key
                "?::character varying," + // name
                "?::text," + // description
                "?::hstore," + // attr_valid
                "?::bigint," + // version
                "?::character varying" + // changed_by
                ")";
    }

    @Override
    public String getGetItemTypeSQL() {
        return "SELECT item_type(" +
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
                "?::boolean," + // system
                "?::timestamp(6) with time zone," + // date created from
                "?::timestamp(6) with time zone," + // date created to
                "?::timestamp(6) with time zone," + // date updates from
                "?::timestamp(6) with time zone" + // date updated to
                ")";
    }

    @Override
    public String getSetLinkTypeSQL() {
        return null;
    }

    @Override
    public String getDeleteLinkRuleSQL() {
        return "SELECT delete_link_rule(?::character varying)";
    }

    @Override
    public String getDeleteLinkRulesSQL() {
        return "SELECT delete_link_rules()";
    }

    /* snapshots */
    @Override
    public Result createSnapshot(JSONObject json) {
        Result result = new Result();
        Object name = json.get("name");
        Object description = json.get("description");
        Object label = json.get("label");
        Object rootItemKey = json.get("rootItemKey");
        try {
            db.prepare(getCreateSnapshotSQL());
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
    public Result updateSnapshot(String rootItemKey, String currentLabel, JSONObject json) {
        Result result = new Result();
        Object name = json.get("name");
        Object description = json.get("description");
        Object newLabel = json.get("label");
        Object version = json.get("version");
        try {
            db.prepare(getUpdateSnapshotSQL());
            db.setString(1, (rootItemKey != null) ? (String) rootItemKey : null); // root item key
            db.setString(2, (currentLabel != null) ? (String) currentLabel : null); // current_label
            db.setString(3, (newLabel != null) ? (String) newLabel : null); // new_label
            db.setString(4, (name != null) ? (String) name : null); // name_param
            db.setString(5, (description != null) ? (String) description : null); // description_param
            db.setString(6, getUser()); // changed_by_param
            db.setObject(7, version); // version_param
            result.setOperation(db.executeQueryAndRetrieveStatus("update_snapshot"));
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
    public Result deleteSnapshot(String rootItemKey, String label) {
        Result result = new Result();
        try {
            db.prepare(getDeleteSnapshotSQL());
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
    public SnapshotList getItemSnapshots(String rootItemKey) {
        SnapshotList snapshots = new SnapshotList();
        try {
            db.prepare(getGetItemSnapshotsSQL());
            db.setString(1, rootItemKey); // root_item_key_param
            ResultSet set = db.executeQuery();
            while (set.next()) {
                snapshots.getItems().add(util.toSnapshotData(set));
            }
        }
        catch (Exception ex) {
            ex.printStackTrace();
        }
        finally {
            db.close();
        }
        return snapshots;
    }

    @Override
    public ItemTreeData getItemTree(String rootItemKey, String label) {
        ItemTreeData tree = new ItemTreeData();
        try {
            db.prepare(getGetTreeItemsForSnapshotSQL());
            db.setString(1, rootItemKey); // root_item_key_param
            db.setString(2, label); // label_param
            ResultSet set = db.executeQuery();
            while (set.next()) {
                tree.getItems().add(util.toItemData(set));
            }
            db.prepare(getGetTreeLinksForSnapshotSQL());
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
                    "?::boolean," +
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
    public String getCreateSnapshotSQL() {
        return "SELECT create_snapshot(" +
                "?::character varying," + // root_item_key_param
                "?::character varying," + // snapshot_label_param
                "?::character varying," + // snapshot_name_param
                "?::text," + // snapshot_description_param
                "?::character varying" + // changed_by
                ")";
    }

    @Override
    public String getDeleteSnapshotSQL() {
        return "SELECT delete_snapshot(" +
                "?::character varying," + // root_item_key_param
                "?::character varying" + // snapshot_label_param
                ")";
    }

    @Override
    public String getUpdateSnapshotSQL() {
        return "SELECT update_snapshot(" +
                "?::character varying," + // root_item_key_param
                "?::character varying," + // current_label_param
                "?::character varying," + // new_label_param
                "?::character varying," + // snapshot_name_param
                "?::text," + // snapshot_description_param
                "?::character varying," + // changed_by_param
                "?::bigint" + // version_param
                ")";
    }

    @Override
    public String getGetItemSnapshotsSQL() {
        return "SELECT * FROM get_item_snapshots(" +
                "?::character varying" + // root_item_key_param
                ")";
    }

    @Override
    public String getGetTreeItemsForSnapshotSQL() {
        return "SELECT * FROM get_tree_items(" +
                "?::character varying," + // root_item_key_param
                "?::character varying" + // snapshot_label_param
                ")";
    }

    @Override
    public String getGetTreeLinksForSnapshotSQL() {
        return "SELECT * FROM get_tree_links(" +
                "?::character varying," + // root_item_key_param
                "?::character varying" + // snapshot_label_param
                ")";
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
