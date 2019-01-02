package org.gatblau.onix;

import org.gatblau.onix.data.*;
import org.gatblau.onix.model.ItemType;
import org.json.simple.JSONObject;
import org.json.simple.parser.ParseException;

import java.io.IOException;
import java.sql.SQLException;
import java.time.ZonedDateTime;
import java.util.List;

/*
    Provides an abstraction to the underlying onix cmdb database.
 */
public interface DbRepository {
    /*
        ITEMS
    */
    Result createOrUpdateItem(String key, JSONObject json) throws IOException, SQLException, ParseException;
    ItemData getItem(String key) throws SQLException, ParseException;
    Result deleteItem(String key) throws SQLException;
    ItemList findItems(String itemTypeKey, List<String> tagList, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, Short status, Integer top) throws SQLException, ParseException;

    /*
        LINKS
    */
    LinkData getLink(String key) throws SQLException, ParseException;
    Result createOrUpdateLink(String key, JSONObject json) throws SQLException, ParseException;
    Result deleteLink(String key) throws SQLException;
    LinkList findLinks();
    Result clear() throws SQLException;

    /*
        ITEM TYPES
    */
    ItemTypeData getItemType(String key);
    Result deleteItemTypes() throws SQLException;
    List<ItemType> getItemTypes() throws SQLException;
    Result createOrUpdateItemType(String key, JSONObject json) throws SQLException;
    Result deleteItemType(String key) throws SQLException;

    /*
        LINK TYPES
    */
    List<ItemType> getLinkTypes();
    Result createOrUpdateLinkType(String key);
    Result deleteLinkType(String key) throws SQLException;
    Result deleteLinkTypes() throws SQLException;

    /*
        LINK RULES
    */
    List<ItemTypeData> getLinkRules();
    Result createOrUpdateLinkRule(String key);
    Result deleteLinkRule(String key) throws SQLException;
    Result deleteLinkRules() throws SQLException;

    /*
        AUDIT
    */
    List<AuditItemData> findAuditItems();

    /*
        Function Calls
     */
    String getGetItemSQL();
    String getSetItemSQL();
    String getFindItemsSQL();
    String getDeleteItemSQL();

    String getDeleteLinkSQL();
    String getGetLinkSQL();
    String getSetLinkSQL();
    String getFindLinksSQL();

    String getClearAllSQL();

    String getDeleteItemTypeSQL();
    String getDeleteItemTypes();
    String getFindItemTypesSQL();
    String getSetItemTypeSQL();

    String getDeleteLinkTypeSQL();
    String getDeleteLinkTypes();
    String getSetLinkTypeSQL();

    String getDeleteLinkRuleSQL();
    String getDeleteLinkRulesSQL();
}
