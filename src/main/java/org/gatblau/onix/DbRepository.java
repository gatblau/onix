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
    Result createOrUpdateItem(String key, ItemData json);
    ItemData getItem(String key, boolean includeLinks);
    Result deleteItem(String key);
    ItemList findItems(String itemTypeKey, List<String> tagList, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, Short status, String modelKey, Integer top);
    JSONObject getItemMeta(String key, String filter);
    Result deleteAllItems();

    /*
        LINKS
    */
    LinkData getLink(String key);
    Result createOrUpdateLink(String key, LinkData json);
    Result deleteLink(String key);
    LinkList findLinks(String linkTypeKey, String startItemKey, String endItemKey, List<String> tagList, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, String modelKey, Integer top);

    Result clear();

    /*
        ITEM TYPES
    */
    ItemTypeData getItemType(String key);
    Result deleteItemTypes();
    ItemTypeList getItemTypes(Map attribute, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, String modelType);

    Result createOrUpdateItemType(String key, ItemTypeData json);
    Result deleteItemType(String key);

    /*
        LINK TYPES
    */
    LinkTypeList getLinkTypes(Map attrMap, ZonedDateTime date, ZonedDateTime zonedDateTime, ZonedDateTime dateTime, ZonedDateTime time, String modelKey);
    Result createOrUpdateLinkType(String key, LinkTypeData json);
    Result deleteLinkType(String key);
    Result deleteLinkTypes();
    LinkTypeData getLinkType(String key);

    /*
        LINK RULES
    */
    LinkRuleList getLinkRules(String linkTypeKey, String startItemType, String endItemType, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo);
    Result createOrUpdateLinkRule(String key, LinkRuleData json);
    Result deleteLinkRule(String key);
    Result deleteLinkRules();

    /*
        CHANGE
    */
    List<ChangeItemData> findChangeItems();

    /*
        Function Calls
     */
    String getGetItemSQL();
    String getSetItemSQL();
    String getFindItemsSQL();
    String getDeleteItemSQL();
    String getDeleteAllItemsSQL();

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
    String getGetLinkTypeSQL();
    String getFindLinkTypesSQL();

    String getDeleteLinkRuleSQL();
    String getDeleteLinkRulesSQL();
    String getSetLinkRuleSQL();
    String getFindLinkRulesSQL();

    String getFindChildItemsSQL();

    String getCreateTagSQL();
    String getDeleteTagSQL();
    String getUpdateTagSQL();
    String getGetItemTagsSQL();

    String getGetTreeItemsForTagSQL();
    String getGetTreeLinksForTagSQL();

    String getDeleteItemTreeSQL();

    String getTableCountSQL();

    /* Tag */
    Result createTag(JSONObject payload);
    Result updateTag(String rootItemKey, String label, JSONObject payload);
    Result deleteTag(String rootItemKey, String label);
    TagList getItemTags(String rootItemKey);

    /* Graph Data */
    GraphData getData(String rootItemKey, String label);
    ResultList createOrUpdateData(GraphData payload);
    Result deleteData(String rootItemKey);

    TypeGraphData getTypeDataByModel(String modelKey);

    String getGetModelItemTypesSQL();
    String getGetModelLinkTypesSQL();
    String getGetModelLinkRulesSQL();

    /* Readiness probe */
    JSONObject getReadyStatus();

    /* Model */
    Result deleteModel(String key);
    Result createOrUpdateModel(String key, ModelData json);
    ModelData getModel(String key);
    ModelDataList getModels();

    String getDeleteModelSQL();
    String getSetModelSQL();
    String getGetModelsSQL();
    String getGetModelSQL();
}
