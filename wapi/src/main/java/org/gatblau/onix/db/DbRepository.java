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
    Result createOrUpdateItem(String key, ItemData json, String[] role);
    ItemData getItem(String key, boolean includeLinks, String[] role);
    Result deleteItem(String key, String[] role);
    ItemList findItems(String itemTypeKey, List<String> tagList, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, Short status, String modelKey, Map<String, String> attributes, Integer top, String[] role);
    JSONObject getItemMeta(String key, String filter, String[] role);
    Result deleteAllItems(String[] role);

    /*
        LINKS
    */
    LinkData getLink(String key, String[] role);
    Result createOrUpdateLink(String key, LinkData link, String[] role);
    Result deleteLink(String key, String[] role);
    LinkList findLinks(String linkTypeKey, String startItemKey, String endItemKey, List<String> tagList, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, String modelKey, Integer top, String[] role);

    /*
       MISC
     */
    Result clear(String[] role);

    /*
        ITEM TYPES
    */
    ItemTypeData getItemType(String key, String[] role);
    Result deleteItemTypes(String[] role);
    ItemTypeList getItemTypes(Map attribute, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, String modelType, String[] role);
    Result createOrUpdateItemType(String key, ItemTypeData json, String[] role);
    Result deleteItemType(String key, String[] role);

    /*
        ITEM TYPE ATTRIBUTES
    */

    /*
        LINK TYPES
     */
    LinkTypeList getLinkTypes(Map attrMap, ZonedDateTime date, ZonedDateTime zonedDateTime, ZonedDateTime dateTime, ZonedDateTime time, String modelKey, String[] role);
    Result createOrUpdateLinkType(String key, LinkTypeData json, String[] role);
    Result deleteLinkType(String key, String[] role);
    Result deleteLinkTypes(String[] role);
    LinkTypeData getLinkType(String key, String[] role);

    /*
        LINK RULES
    */
    LinkRuleList getLinkRules(String linkTypeKey, String startItemType, String endItemType, ZonedDateTime createdFrom, ZonedDateTime createdTo, ZonedDateTime updatedFrom, ZonedDateTime updatedTo, String[] role);
    Result createOrUpdateLinkRule(String key, LinkRuleData linkRule, String[] role);
    Result deleteLinkRule(String key, String[] role);
    Result deleteLinkRules(String[] role);

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
    ResultList createOrUpdateData(GraphData payload, String[] role);
    Result deleteData(String rootItemKey);

    TypeGraphData getTypeDataByModel(String modelKey, String[] role);

    String getGetModelItemTypesSQL();
    String getGetModelLinkTypesSQL();
    String getGetModelLinkRulesSQL();

    /* Readiness probe */
    JSONObject checkReady();

    /* Model */
    Result deleteModel(String key, String[] role);
    Result createOrUpdateModel(String key, ModelData json, String[] role);
    ModelData getModel(String key, String[] role);
    ModelDataList getModels(String[] role);

    String getDeleteModelSQL();
    String getSetModelSQL();
    String getGetModelsSQL();
    String getGetModelSQL();

    /* Partitions */
    String getDeletePartitionSQL();
    String getSetPartitionSQL();
    String getGetAllPartitionsSQL();
    String getGetPartitionSQL();

    Result deletePartition(String key, String[] role);
    Result createOrUpdatePartition(String key, PartitionData partition, String[] role);
    PartitionDataList getAllPartitions(String[] role);
    PartitionData getPartition(String key, String[] role);

    /* Roles */
    String getDeleteRoleSQL();
    String getSetRoleSQL();
    String getGetRoleSQL();
    String getGetAllRolesSQL();

    Result deleteRole(String key, String[] role);
    Result createOrUpdateRole(String key, RoleData role, String[] role1);
    RoleData getRole(String key, String[] role);
    RoleDataList getAllRoles(String[] role);

    /* Privileges */
    String getAddPrivilegeSQL();
    String getRemovePrivilegeSQL();
    String getGetAllPrivilegeByRoleSQL();

    Result addPrivilege(String partitionKey, String roleKey, NewPrivilegeData privilege, String[] role);
    Result removePrivilege(String partitionKey, String roleKey, String[] role);

    PrivilegeDataList getPrivilegesByRole(String roleKey, String[] loggedRoleKey);

    ItemList getItemChildren(String key, String[] role);

    String getGetItemChildrenSQL();

    /* Enc Keys */
    EncKeyStatusData getKeyStatus(String[] role);
    String getGetEncKeyUsageSQL();
}
