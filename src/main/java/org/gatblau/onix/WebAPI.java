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

    private DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyyMMddHHmm");

    @ApiOperation(
        value = "Returns OK if the service is up and running.",
        notes = "Use it as a readiness probe for the service.",
        response = String.class)
    @RequestMapping(value = "/", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<Info> index() {
        return ResponseEntity.ok(new Info("Onix Configuration Management Database Service.", "0.2"));
    }

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
        value = "Removes ALL configuration items and links from the database.",
        notes = "Use at your own risk ONLY for testing of the CMDB!")
    @RequestMapping(path = "/clear", method = RequestMethod.DELETE)
    public void clear() throws SQLException {
        data.clear();
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
        value = "Deletes a configuration item type.",
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
        value = "Get a list of available configuration item types.",
        notes = "Only item types marked as custom can be deleted.")
    @RequestMapping(
          path = "/itemtypes"
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
    ) throws SQLException, ParseException {
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

    @ApiOperation(
        value = "Get a configuration item based on the specified key.",
        notes = "Use this search to retrieve a specific configuration item when its natural key is known.")
    @RequestMapping(
          path = "/item/{key}"
        , method = RequestMethod.GET
        , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ItemData> getItem(@PathVariable("key") String key) throws SQLException, ParseException {
        return ResponseEntity.ok(data.getItem(key));
    }

    @ApiOperation(
        value = "Search for configuration items based on the specified filters (provided via a query string).",
        notes = "Use this function to retrieve configuration items based on type, tags and date range as required. " +
                "Results are limited by the top parameter.")
    @RequestMapping(
          path = "/items"
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
    ) throws SQLException, ParseException {
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