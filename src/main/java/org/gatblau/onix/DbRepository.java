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
import org.json.simple.JSONObject;
import org.json.simple.parser.ParseException;

import java.io.IOException;
import java.sql.SQLException;
import java.time.ZonedDateTime;
import java.util.List;
import java.util.Map;

/*
    Provides an abstraction to the underlying onix cmdb database.
 */
public interface DbRepository {
    /*
        ITEMS
    */
    Result createOrUpdateItem(String key, JSONObject json);
    ItemData getItem(String key);
    Result deleteItem(String key);
    ItemList findItems(String itemTypeKey, List<String> tagList, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, Short status, Integer top);

    /*
        LINKS
    */
    LinkData getLink(String key);
    Result createOrUpdateLink(String key, JSONObject json);
    Result deleteLink(String key);
    LinkList findLinks();
    Result clear();

    /*
        ITEM TYPES
    */
    ItemTypeData getItemType(String key);
    Result deleteItemTypes();
    ItemTypeList getItemTypes(Map attribute, Boolean system, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo);

    Result createOrUpdateItemType(String key, JSONObject json);
    Result deleteItemType(String key);

    /*
        LINK TYPES
    */
    LinkTypeList getLinkTypes(Map attrMap, Boolean system, ZonedDateTime date, ZonedDateTime zonedDateTime, ZonedDateTime dateTime, ZonedDateTime time);
    Result createOrUpdateLinkType(String key, JSONObject json);
    Result deleteLinkType(String key);
    Result deleteLinkTypes();
    LinkTypeData getLinkType(String key);

    /*
        LINK RULES
    */
    LinkRuleList getLinkRules(String linkTypeKey, String startItemType, String endItemType, Boolean system, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo);
    Result createOrUpdateLinkRule(String key, JSONObject payload);
    Result deleteLinkRule(String key);
    Result deleteLinkRules();

    /*
        AUDIT
    */
    List<AuditItemData> findAuditItems();

    /*
        INVENTORY
     */
    Result createOrUpdateInventory(String key, String inventory);
    String getInventory(String key);

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
    String getGetItemTypeSQL();

    String getDeleteLinkTypeSQL();
    String getDeleteLinkTypes();
    String getSetLinkTypeSQL();
    String getFindLinkTypesSQL();

    String getDeleteLinkRuleSQL();
    String getDeleteLinkRulesSQL();
    String getSetLinkRuleSQL();
    String getFindLinkRulesSQL();

    String getFindChildItemsSQL();
}
