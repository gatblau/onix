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

import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.gatblau.onix.data.*;
import org.json.simple.JSONObject;
import org.json.simple.parser.ParseException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.sql.SQLException;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@Api("ONIX CMDB Web API")
@RestController
public class WebAPI {

    @Autowired
    private DbRepository data;

    @Autowired
    private Info info;

    private DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyyMMddHHmm");

    @ApiOperation(
            value = "Returns a JSON payload if the service is alive.",
            notes = "Use as liveliness probe for the service.",
            response = JSONObject.class)
    @RequestMapping(value = "/", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<JSONObject> index() {
        return live();
    }

    @ApiOperation(
            value = "Returns a JSON payload if the service is alive.",
            notes = "Use as liveliness probe for the service.",
            response = JSONObject.class)
    @RequestMapping(value = "/live", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<JSONObject> live() {
        JSONObject response = new JSONObject();
        response.put("live", true);
        return ResponseEntity.ok(response);
    }

    @ApiOperation(
            value = "Returns the readyness status of the service.",
            notes = "Use as readyness probe for the service.",
            response = JSONObject.class)
    @RequestMapping(value = "/ready", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<JSONObject> ready() {
        return ResponseEntity.ok(data.getReadyStatus());
    }

    @ApiOperation(
            value = "Returns information about the service.",
            notes = "",
            response = String.class)
    @RequestMapping(value = "/info", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<Info> info() {
        return ResponseEntity.ok(info);
    }

    /*
        ITEMS
     */
    @ApiOperation(
        value = "Creates new item or updates an existing item based on the specified key.",
        notes = "Use this operation to create configuration item if it's not there or update it if it's there.")
    @RequestMapping(
            path = "/item/{key}", method = RequestMethod.PUT,
            consumes = {"application/json" },
            produces = {"application/json" })
    public ResponseEntity<Result> createOrUpdateItem(
            @PathVariable("key") String key,
            @RequestBody JSONObject payload) {
        return ResponseEntity.ok(data.createOrUpdateItem(key, payload));
    }

    @ApiOperation(
        value = "Deletes an existing configuration item.",
        notes = "Use this operation to remove a configuration item after it has been decommissioned.")
    @RequestMapping(
          path = "/item/{key}"
        , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteItem(
            @PathVariable("key") String key
    ) {
        return ResponseEntity.ok(data.deleteItem(key));
    }

    @ApiOperation(
            value = "Deletes all configuration items.",
            notes = "")
    @RequestMapping(
            path = "/item"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteAllItems(
    ) {
        return ResponseEntity.ok(data.deleteAllItems());
    }

    @ApiOperation(
            value = "Get a configuration item based on the specified key.",
            notes = "Use this search to retrieve a specific configuration item when its natural key is known.")
    @RequestMapping(
            path = "/item/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ItemData> getItem(
            @PathVariable("key") String key,
            @RequestParam(required = false, name = "links", defaultValue = "false" // true to retrieve link information
    ) boolean links) {
        return ResponseEntity.ok(data.getItem(key, links));
    }

    @ApiOperation(
            value = "Search for configuration items based on the specified filters (provided via a query string).",
            notes = "Use this function to retrieve configuration items based on type, tags and date range as required. " +
                    "Results are limited by the top parameter.")
    @RequestMapping(
            path = "/item"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<Wrapper> getItems(
            @RequestParam(value = "type", required = false) String itemTypeKey
            , @RequestParam(value = "tag", required = false) String tag
            , @RequestParam(value = "createdFrom", required = false) String createdFromDate
            , @RequestParam(value = "createdTo", required = false) String createdToDate
            , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
            , @RequestParam(value = "updatedTo", required = false) String updatedToDate
            , @RequestParam(value = "status", required = false) Short status
            , @RequestParam(value = "model", required = false) String modelKey
            , @RequestParam(value = "top", required = false, defaultValue = "100") Integer top
    ) {
        List<String> tagList = null;
        if (tag != null) {
            String[] tags = tag.split("[|]"); // separate tags using pipes in the query string
            tagList = Arrays.asList(tags);
        }
        ItemList list = data.findItems(
                itemTypeKey,
                tagList,
                getDate(createdFromDate),
                getDate(createdToDate),
                getDate(updatedFromDate),
                getDate(updatedToDate),
                status,
                modelKey,
                top
        );
        return ResponseEntity.ok(list);
    }

    @ApiOperation(
            value = "Gets the metadata associated with the specified item.",
            notes = "Use this function to retrieve the full metadata for an item.")
    @RequestMapping(
            path = "/item/{key}/meta"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<JSONObject> getItemMeta(
        @PathVariable("key") String key
    ){
        return ResponseEntity.ok(data.getItemMeta(key, null));
    }

    @ApiOperation(
            value = "Gets the portions of the metadata associated with the specified item and filter.",
            notes = "Use this function to retrieve portions of the metadata for an item.")
    @RequestMapping(
            path = "/item/{key}/meta/{filter}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<JSONObject> getFilteredItemMeta(
            @PathVariable("key") String key,
            @PathVariable("filter") String filter
    ){
        return ResponseEntity.ok(data.getItemMeta(key, filter));
    }

    /*
        ITEM TYPES
     */

    @ApiOperation(
        value = "Deletes all non-system specific configuration item types.",
        notes = "")
    @RequestMapping(
          path = "/itemtype"
        , method = RequestMethod.DELETE
    )
    public void deleteItemTypes() throws SQLException {
        data.deleteItemTypes();
    }

    @ApiOperation(
        value = "Deletes a configuration item type.",
        notes = "")
    @RequestMapping(
          path = "/itemtype/{key}"
        , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteItemType(@PathVariable("key") String key) throws SQLException {
        return ResponseEntity.ok(data.deleteItemType(key));
    }

    @ApiOperation(
        value = "Creates a new configuration item type.",
        notes = "")
    @RequestMapping(
          path = "/itemtype/{key}"
        , method = RequestMethod.PUT)
    public ResponseEntity<Result> createItemType(
            @PathVariable("key") String key,
            @RequestBody JSONObject payload
        ) throws IOException, SQLException {
        return ResponseEntity.ok(data.createOrUpdateItemType(key, payload));
    }

    @ApiOperation(
            value = "Get a configuration item type based on the specified key.",
            notes = "Use this search to retrieve a specific configuration item type when its natural key is known.")
    @RequestMapping(
            path = "/itemtype/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ItemTypeData> getItemType(@PathVariable("key") String key) throws SQLException, ParseException, IOException {
        return ResponseEntity.ok(data.getItemType(key));
    }

    @ApiOperation(
        value = "Get a list of available configuration item types.",
        notes = "Only item types marked as custom can be deleted.")
    @RequestMapping(
          path = "/itemtype"
        , method = RequestMethod.GET
        , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ItemTypeList> getItemTypes(
          @RequestParam(value = "attribute", required = false) String attribute
        , @RequestParam(value = "createdFrom", required = false) String createdFromDate
        , @RequestParam(value = "createdTo", required = false) String createdToDate
        , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
        , @RequestParam(value = "updatedTo", required = false) String updatedToDate
        , @RequestParam(value = "model", required = false) String modelKey
    ) {
        Map attrMap = null;
        if (attribute != null) {
            attrMap = new HashMap<String, String>();
            String[] items = attribute.split("[|]"); // separate tags using pipes in the query string
            for(String item : items) {
                String[] parts = item.split("->");
                attrMap.put(parts[0],parts[1]);
            }
        }
        ItemTypeList itemTypes = data.getItemTypes(
            attrMap,
            getDate(createdFromDate),
            getDate(createdToDate),
            getDate(updatedFromDate),
            getDate(updatedToDate),
            modelKey);
        return ResponseEntity.ok(itemTypes);
    }

    /*
        LINKS
     */
    @ApiOperation(
            value = "Creates new link or updates an existing link based on its natural key.",
            notes = "Use this operation to create a new link between two existing configuration items or to update such link if it already exists.")
    @RequestMapping(
            path = "/link/{key}", method = RequestMethod.PUT,
            consumes = {"application/json" },
            produces = {"application/json" })
    public ResponseEntity<Result> createOrUpdateLink(
            @PathVariable("key") String key,
            @RequestBody JSONObject payload) {
        return ResponseEntity.ok(data.createOrUpdateLink(key, payload));
    }

    @ApiOperation(
            value = "Deletes a link between two existing configuration items.",
            notes = "Use this operation to delete links between existing items.")
    @RequestMapping(path = "/link/{key}", method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteLink(
            @PathVariable("key") String key
    ) {
        return ResponseEntity.ok(data.deleteLink(key));
    }

    /*
        LINK TYPES
     */
    @ApiOperation(
            value = "Deletes all non-system specific item link types.",
            notes = "")
    @RequestMapping(
            path = "/linktype"
            , method = RequestMethod.DELETE
    )
    public void deleteLinkTypes() throws SQLException {
        data.deleteLinkTypes();
    }

    @ApiOperation(
            value = "Deletes an item link type.",
            notes = "")
    @RequestMapping(
            path = "/linktype/{key}"
            , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteLinkType(@PathVariable("key") String key) throws SQLException {
        return ResponseEntity.ok(data.deleteLinkType(key));
    }

    @ApiOperation(
            value = "Creates a new item link type.",
            notes = "")
    @RequestMapping(
            path = "/linktype/{key}"
            , method = RequestMethod.PUT)
    public ResponseEntity<Result> createOrUpdateLinkType(
            @PathVariable("key") String key,
            @RequestBody JSONObject payload
    ) {
        return ResponseEntity.ok(data.createOrUpdateLinkType(key, payload));
    }

    @ApiOperation(
            value = "Get an item link type based on the specified key.",
            notes = "Use this search to retrieve a specific item link type when its natural key is known.")
    @RequestMapping(
            path = "/linktype/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<LinkTypeData> getLinkType(@PathVariable("key") String key) throws SQLException, ParseException {
        return ResponseEntity.ok(data.getLinkType(key));
    }

    @ApiOperation(
            value = "Get a list of available link types.",
            notes = "")
    @RequestMapping(
            path = "/linktype"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<LinkTypeList> getLinkTypes(
            @RequestParam(value = "attribute", required = false) String attribute
            , @RequestParam(value = "createdFrom", required = false) String createdFromDate
            , @RequestParam(value = "createdTo", required = false) String createdToDate
            , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
            , @RequestParam(value = "updatedTo", required = false) String updatedToDate
            , @RequestParam(value = "model", required = false) String modelKey
    ) throws SQLException, ParseException, IOException {
        Map attrMap = null;
        if (attribute != null) {
            attrMap = new HashMap<String, String>();
            String[] items = attribute.split("[|]"); // separate tags using pipes in the query string
            for(String item : items) {
                String[] parts = item.split("->");
                attrMap.put(parts[0],parts[1]);
            }
        }
        LinkTypeList linkTypes = data.getLinkTypes(
                attrMap,
                getDate(createdFromDate),
                getDate(createdToDate),
                getDate(updatedFromDate),
                getDate(updatedToDate),
                modelKey);
        return ResponseEntity.ok(linkTypes);
    }

    /*
        LINK RULES
     */

    @ApiOperation(
            value = "Deletes all non-system specific item link types.",
            notes = "")
    @RequestMapping(
            path = "/linkrule"
            , method = RequestMethod.DELETE
    )
    public void deleteLinkRules() throws SQLException {
        data.deleteLinkRules();
    }

    @ApiOperation(
            value = "Creates a new or updates an existing link rule.",
            notes = "")
    @RequestMapping(
            path = "/linkrule/{key}"
            , method = RequestMethod.PUT)
    public ResponseEntity<Result> createOrUpdateLinkRule(
            @PathVariable("key") String key,
            @RequestBody JSONObject payload
    ) throws IOException, SQLException {
        return ResponseEntity.ok(data.createOrUpdateLinkRule(key, payload));
    }

    @ApiOperation(
            value = "Get a list of available link rules filtered by the specified query parameters.",
            notes = "")
    @RequestMapping(
            path = "/linkrule"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<LinkRuleList> getLinkRules(
              @RequestParam(value = "linkType", required = false) String linkType
            , @RequestParam(value = "startItemType", required = false) String startItemType
            , @RequestParam(value = "endItemType", required = false) String endItemType
            , @RequestParam(value = "createdFrom", required = false) String createdFromDate
            , @RequestParam(value = "createdTo", required = false) String createdToDate
            , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
            , @RequestParam(value = "updatedTo", required = false) String updatedToDate
    ) throws SQLException, ParseException {
        LinkRuleList linkRules = data.getLinkRules(
                linkType,
                startItemType,
                endItemType,
                getDate(createdFromDate),
                getDate(createdToDate),
                getDate(updatedFromDate),
                getDate(updatedToDate));
        return ResponseEntity.ok(linkRules);
    }

    /*
        MODEL
     */

    @ApiOperation(
            value = "Deletes a model for a specified key.",
            notes = "")
    @RequestMapping(
            path = "/model/{key}"
            , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteModel(@PathVariable("key") String key) {
        return ResponseEntity.ok(data.deleteModel(key));
    }

    @ApiOperation(
            value = "Creates a new model.",
            notes = "")
    @RequestMapping(
            path = "/model/{key}"
            , method = RequestMethod.PUT)
    public ResponseEntity<Result> createOrUpdateModel(
            @PathVariable("key") String key,
            @RequestBody JSONObject payload
    ) {
        return ResponseEntity.ok(data.createOrUpdateModel(key, payload));
    }

    @ApiOperation(
            value = "Get a model based on the specified key.",
            notes = "Use this search to retrieve a specific model when its natural key is known.")
    @RequestMapping(
            path = "/model/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ModelData> getModel(@PathVariable("key") String key) {
        return ResponseEntity.ok(data.getModel(key));
    }

    @ApiOperation(
            value = "Get an item link type based on the specified key.",
            notes = "Use this search to retrieve the list of models known to the system.")
    @RequestMapping(
            path = "/model"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ModelDataList> getModels() {
        return ResponseEntity.ok(data.getModels());
    }

    /*
        MISCELLANEOUS
     */

    @ApiOperation(
            value = "Removes ALL configuration items and links from the database.",
            notes = "Use at your own risk ONLY for testing of the CMDB!")
    @RequestMapping(path = "/clear", method = RequestMethod.DELETE)
    public void clear() throws SQLException {
        data.clear();
    }

    /*
        TAG
     */
    @ApiOperation(
            value = "Creates a new tag.",
            notes = "A tag is a set of items and their links at a specific point in time.")
    @RequestMapping(
            path = "/tag"
            , method = RequestMethod.POST)
    public ResponseEntity<Result> createTag(
            @RequestBody JSONObject payload
    ) {
        return ResponseEntity.ok(data.createTag(payload));
    }

    @ApiOperation(
            value = "Updates an existing tag.",
            notes = "A tag is a set of items and their links at a specific point in time.")
    @RequestMapping(
            path = "/tag/{root_item_key}/{label}"
            , method = RequestMethod.PUT)
    public ResponseEntity<Result> updateTag(
            @PathVariable("root_item_key") String rootItemKey,
            @PathVariable("label") String label,
            @RequestBody JSONObject payload
    ) {
        return ResponseEntity.ok(data.updateTag(rootItemKey, label, payload));
    }

    @ApiOperation(
            value = "Deletes an existing tag.",
            notes = "Takes the key of a root item and a tag label and deletes the matching tag.")
    @RequestMapping(
            path = "/tag/{root_item_key}/{label}"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteTag(
            @PathVariable("root_item_key") String rootItemKey,
            @PathVariable("label") String label
    ) {
        return ResponseEntity.ok(data.deleteTag(rootItemKey, label));
    }

    @ApiOperation(
            value = "Deletes all tags for an item.",
            notes = "Takes the key of a root item and deletes any associated tags.")
    @RequestMapping(
            path = "/tag/{root_item_key}"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteAllTags(
            @PathVariable("root_item_key") String rootItemKey
    ) {
        return ResponseEntity.ok(data.deleteTag(rootItemKey, null));
    }

    @ApiOperation(
            value = "Get a list of available tags for a specific item.",
            notes = "")
    @RequestMapping(
            path = "/tag/{root_item_key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<TagList> getItemTags(
            @PathVariable("root_item_key") String rootItemKey
    ) {
        TagList tags = data.getItemTags(rootItemKey);
        return ResponseEntity.ok(tags);
    }

    /*
        DATA
     */
    @ApiOperation(
            value = "Creates or updates a set of items and links.",
            notes = "")
    @RequestMapping(
            path = "/data"
            , method = RequestMethod.PUT
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ResultList> createOrUpdateData(
            @RequestBody JSONObject payload
    ) {
        ResultList results = data.createOrUpdateData(payload);
        return ResponseEntity.ok(results);
    }

    /*
       ITEM TREE
     */
    @ApiOperation(
            value = "Get a list of items and links in a specified item tag.",
            notes = "")
    @RequestMapping(
            path = "/tree/{root_item_key}/{label}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ItemTreeData> getItemTree(
            @PathVariable("root_item_key") String rootItemKey,
            @PathVariable("label") String label
    ) {
        ItemTreeData tree = data.getItemTree(rootItemKey, label);
        return ResponseEntity.ok(tree);
    }

    @ApiOperation(
            value = "Deletes an existing item tree.",
            notes = "")
    @RequestMapping(
            path = "/tree/{root_item_key}"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteTree(
        @PathVariable("root_item_key") String rootItemKey
    ) {
        return ResponseEntity.ok(data.deleteItemTree(rootItemKey));
    }

    /*
        helper methods
     */
    private ZonedDateTime getZonedDateTime(@RequestParam(value = "createdFrom", required = false) String createdFromDate) {
        ZonedDateTime createdFrom = null;
        if (createdFromDate != null) {
            createdFrom = ZonedDateTime.of(LocalDateTime.parse(createdFromDate, formatter), ZoneId.systemDefault());
        }
        return createdFrom;
    }

    private ZonedDateTime getDate(String dateString) {
        ZonedDateTime date = null;
        if (date != null) {
            date = ZonedDateTime.of(LocalDateTime.parse(dateString, formatter), ZoneId.systemDefault());
        }
        return date;
    }

}