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
            value = "Returns information about the service.",
            notes = "",
            response = String.class)
    @RequestMapping(value = "/", method = RequestMethod.GET, produces = "text/html")
    public ResponseEntity<String> index() {
        return ResponseEntity.ok("OK");
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
            @RequestBody JSONObject payload) throws IOException, SQLException, ParseException {
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
    ) throws InterruptedException, SQLException {
        return ResponseEntity.ok(data.deleteItem(key));
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
                top
        );
        return ResponseEntity.ok(list);
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
        , @RequestParam(value = "system", required = false) Boolean system
        , @RequestParam(value = "createdFrom", required = false) String createdFromDate
        , @RequestParam(value = "createdTo", required = false) String createdToDate
        , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
        , @RequestParam(value = "updatedTo", required = false) String updatedToDate
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
        ItemTypeList itemTypes = data.getItemTypes(
            attrMap,
            system,
            getDate(createdFromDate),
            getDate(createdToDate),
            getDate(updatedFromDate),
            getDate(updatedToDate));
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
            @RequestBody JSONObject payload) throws SQLException, ParseException {
        return ResponseEntity.ok(data.createOrUpdateLink(key, payload));
    }

    @ApiOperation(
            value = "Deletes a link between two existing configuration items.",
            notes = "Use this operation to delete links between existing items.")
    @RequestMapping(path = "/link/{key}", method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteLink(
            @PathVariable("key") String key
    ) throws InterruptedException, SQLException {
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
    public ResponseEntity<Result> createLinkType(
            @PathVariable("key") String key,
            @RequestBody JSONObject payload
    ) throws IOException, SQLException {
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
            , @RequestParam(value = "system", required = false) Boolean system
            , @RequestParam(value = "createdFrom", required = false) String createdFromDate
            , @RequestParam(value = "createdTo", required = false) String createdToDate
            , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
            , @RequestParam(value = "updatedTo", required = false) String updatedToDate
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
                system,
                getDate(createdFromDate),
                getDate(createdToDate),
                getDate(updatedFromDate),
                getDate(updatedToDate));
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
            , @RequestParam(value = "system", required = false) Boolean system
            , @RequestParam(value = "createdFrom", required = false) String createdFromDate
            , @RequestParam(value = "createdTo", required = false) String createdToDate
            , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
            , @RequestParam(value = "updatedTo", required = false) String updatedToDate
    ) throws SQLException, ParseException {
        LinkRuleList linkRules = data.getLinkRules(
                linkType,
                startItemType,
                endItemType,
                system,
                getDate(createdFromDate),
                getDate(createdToDate),
                getDate(updatedFromDate),
                getDate(updatedToDate));
        return ResponseEntity.ok(linkRules);
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
        INVENTORY
     */
    @ApiOperation(
            value = "Creates a new or updates an existing inventory.",
            notes = "NOTE: this endpoint is EXPERIMENTAL. In future versions, it will be deprecated and its logic will be moved to the ox_cli tool.")
    @RequestMapping(
            path = "/inventory/{key}"
            , method = RequestMethod.PUT)
    public ResponseEntity<Result> createOrUpdateInventory(
            @PathVariable("key") String key,
            @RequestBody String inventory
    ) throws IOException, SQLException, ParseException {
        return ResponseEntity.ok(data.createOrUpdateInventory(key, inventory));
    }

    @ApiOperation(
            value = "Retrieves an existing inventory.",
            notes = "NOTE: this endpoint is EXPERIMENTAL. In future versions, it will be deprecated and its logic will be moved to the ox_cli tool.")
    @RequestMapping(
              path = "/inventory/{key}/{label}"
            , method = RequestMethod.GET)
    public ResponseEntity<String> getInventory(
            @PathVariable("key") String key,
            @PathVariable("label") String label
    ) throws SQLException, ParseException, IOException {
        return ResponseEntity.ok(data.getInventory(key, label));
    }

    /*
        SNAPSHOT
     */
    @ApiOperation(
            value = "Creates a new snapshot.",
            notes = "A snapshot is a set of items and their links at a specific point in time.")
    @RequestMapping(
            path = "/snapshot"
            , method = RequestMethod.POST)
    public ResponseEntity<Result> createSnapshot(
            @RequestBody JSONObject payload
    ) {
        return ResponseEntity.ok(data.createSnapshot(payload));
    }

    @ApiOperation(
            value = "Updates an existing snapshot.",
            notes = "A snapshot is a set of items and their links at a specific point in time.")
    @RequestMapping(
            path = "/snapshot/{root_item_key}/{label}"
            , method = RequestMethod.PUT)
    public ResponseEntity<Result> updateSnapshot(
            @PathVariable("root_item_key") String rootItemKey,
            @PathVariable("label") String label,
            @RequestBody JSONObject payload
    ) {
        return ResponseEntity.ok(data.updateSnapshot(rootItemKey, label, payload));
    }

    @ApiOperation(
            value = "Deletes an existing snapshot.",
            notes = "Takes the key of a root item and a snapshot label and deletes the matching snapshot.")
    @RequestMapping(
            path = "/snapshot/{root_item_key}/{label}"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteSnapshot(
            @PathVariable("root_item_key") String rootItemKey,
            @PathVariable("label") String label
    ) {
        return ResponseEntity.ok(data.deleteSnapshot(rootItemKey, label));
    }

    @ApiOperation(
            value = "Deletes all snapshots for an item.",
            notes = "Takes the key of a root item and deletes any associated snapshots.")
    @RequestMapping(
            path = "/snapshot/{root_item_key}"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteAllSnapshots(
            @PathVariable("root_item_key") String rootItemKey
    ) {
        return ResponseEntity.ok(data.deleteSnapshot(rootItemKey, null));
    }

    @ApiOperation(
            value = "Get a list of available snapshots for a specific item.",
            notes = "")
    @RequestMapping(
            path = "/snapshot/{root_item_key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<SnapshotList> getItemSnapshots(
            @PathVariable("root_item_key") String rootItemKey
    ) {
        SnapshotList snapshots = data.getItemSnapshots(rootItemKey);
        return ResponseEntity.ok(snapshots);
    }

    /*
       ITEM TREE
     */
    @ApiOperation(
            value = "Get a list of items and links in a specified item snapshot.",
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
            value = "Creates or updates a set of items and links.",
            notes = "")
    @RequestMapping(
            path = "/tree"
            , method = RequestMethod.PUT
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ResultList> createOrUpdateItemTree(
        @RequestBody JSONObject payload
    ) {
        ResultList results = data.createOrUpdateItemTree(payload);
        return ResponseEntity.ok(results);
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