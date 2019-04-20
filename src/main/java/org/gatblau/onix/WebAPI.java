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

import io.swagger.annotations.*;
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
    @RequestMapping(value = "/", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<Info> index() {
        return ResponseEntity.ok(info);
    }

    @ApiOperation(
            value = "Returns a JSON payload if the service is alive.",
            notes = "In Kubernetes, it is used as a liveliness probe for the service. " +
                    "That is, to know when the web api container should be restarted, as the web service process " +
                    "is not receiving requests.",
            response = JSONObject.class)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "Successful connection to the web service endpoint.", response = JSONObject.class)}
    )
    @RequestMapping(value = "/live", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<JSONObject> live() {
        JSONObject response = new JSONObject();
        response.put("live", true);
        return ResponseEntity.ok(response);
    }

    @ApiOperation(
            value = "Returns 200 if the service is ready, i.e. can establish a successful connection to the database.",
            notes = "In Kubernetes, it is used as a readyness probe. " +
                    "That is, to know when the web api container is ready to start accepting traffic. " +
                    "The web api pod is considered ready when the database container is ready and the web api can establish " +
                    "a database connection. ",
            response = JSONObject.class)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "Successful connection to the database.", response = JSONObject.class),
            @ApiResponse(code = 500, message = "Internal server error")}
    )
    @RequestMapping(value = "/ready", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<JSONObject> ready() {
        return ResponseEntity.ok(data.getReadyStatus());
    }

    /*
        ITEMS
     */
    @ApiOperation(
            value = "Creates a new configuration item or updates an existing configuration item based on the passed-in key.",
            notes = "This operation is idempotent.")
    @RequestMapping(
            path = "/item/{key}", method = RequestMethod.PUT,
            consumes = {"application/json"},
            produces = {"application/json"})
    @ApiResponses(value = {
        @ApiResponse(code = 200, message = "No changes where performed to the configuration item."),
        @ApiResponse(code = 201, message = "The configuration item was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
        @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
        @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
        @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateItem(
            @ApiParam(
                    name = "key",
                    value = "A string which uniquely identifies the item and never changes.",
                    required = true,
                    example = "item_01_abc"
            )
            @PathVariable("key") String key,
            @RequestBody ItemData payload) {
        Result result = data.createOrUpdateItem(key, payload);
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Deletes an existing configuration item.",
            notes = "Use this operation to remove a configuration item after it has been decommissioned.")
    @RequestMapping(
            path = "/item/{key}"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteItem(
            @ApiParam(
                name = "key",
                value = "A string which uniquely identifies the item and never changes.",
                required = true,
                example = "item_01_abc"
            )
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
        @ApiParam(
            name = "key",
            value = "A string which uniquely identifies the item and never changes.",
            required = true,
            example = "item_01_abc"
        )
        @PathVariable("key") String key,
        @ApiParam(
            name = "links",
            value = "If present in the query string, returns the links that are related to the item.",
            required = false,
            example = "?links"
        )
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
            @ApiParam(
                name = "type",
                value = "The type of items to retrieve. If no value is passed then all item types are retrieved.",
                required = false,
                example = "test_item_type"
            )
            @RequestParam(value = "type", required = false)
            String itemTypeKey,
            @ApiParam(
                name = "tag",
                value = "A list of search item tags separated by '|' in the query string. If no value is passed then all item types are retrieved.",
                required = false,
                example = "Europe|VM|Large"
            )
            @RequestParam(value = "tag", required = false)
            String tag,
            @ApiParam(
                name = "createdFrom",
                value = "The minimum creation date for the items to find. If no value is passed then all item types are retrieved.",
                required = false,
                example = "12-03-18"
            )
            @RequestParam(value = "createdFrom", required = false)
            String createdFromDate,
            @ApiParam(
                name = "createdTo",
                value = "The maximum creation date for the items to find. If no value is passed then all item types are retrieved.",
                required = false,
                example = "12-03-18"
            )
            @RequestParam(value = "createdTo", required = false)
            String createdToDate,
            @ApiParam(
                name = "updatedFrom",
                value = "The minimum last update date of the items to find. If no value is passed then all item types are retrieved.",
                required = false,
                example = "12-03-18"
            )
            @RequestParam(value = "updatedFrom", required = false)
            String updatedFromDate,
            @ApiParam(
                name = "updatedTo",
                value = "The maximum last update date of the items to find. If no value is passed then all item types are retrieved.",
                required = false,
                example = "12-03-18"
            )
            @RequestParam(value = "updatedTo", required = false)
            String updatedToDate,
            @ApiParam(
                name = "status",
                value = "The status number of the items to find. If no value is passed then all item types are retrieved.",
                required = false,
                example = "5"
            )
            @RequestParam(value = "status", required = false)
            Short status,
            @ApiParam(
                name = "modelKey",
                value = "The key of the model containing the items to find. If no value is passed then all item types are retrieved.",
                required = false,
                example = "test_model_01"
            )
            @RequestParam(value = "model", required = false)
            String modelKey,
            @ApiParam(
                name = "top",
                value = "The maximum number of items to retrieve.",
                required = false,
                example = "12-03-18",
                defaultValue = "100"
            )
            @RequestParam(value = "top", required = false, defaultValue = "100")
            Integer top
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
        @ApiParam(
            name = "key",
            value = "A string which uniquely identifies the item and never changes.",
            required = true,
            example = "item_01_abc"
        )
        @PathVariable("key")
        String key
    ) {
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
        @ApiParam(
                name = "key",
                value = "A string which uniquely identifies the item and never changes.",
                required = true,
                example = "item_01_abc"
        )
        @PathVariable("key")
        String key,
        @ApiParam(
            name = "filter",
            value = "A string which uniquely identifies a filter applied to the meta field content.",
            required = false,
            example = "books"
        )
        @PathVariable("filter")
        String filter
    ) {
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
    public void deleteItemTypes() {
        data.deleteItemTypes();
    }

    @ApiOperation(
        value = "Deletes a configuration item type.",
        notes = "")
    @RequestMapping(
        path = "/itemtype/{key}"
        , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteItemType(
        @ApiParam(
            name = "key",
            value = "A string which uniquely identifies the item type and never changes.",
            required = true,
            example = "item_type_01"
        )
        @PathVariable("key")
        String key,
        @ApiParam(
            name = "force",
            value = "If true, it forces the deletion of existing items linked to the item type.",
            required = false,
            example = "?force"
        )
        @RequestParam(required = false, name = "force", defaultValue = "false") // true to force deletion of related items
        boolean force
    ) {
        return ResponseEntity.ok(data.deleteItemType(key, force));
    }

    @ApiOperation(
        value = "Creates a new configuration item type.",
        notes = "")
    @RequestMapping(
        path = "/itemtype/{key}"
        , method = RequestMethod.PUT)
    @ApiResponses(value = {
        @ApiResponse(code = 200, message = "No changes where performed to the configuration item type."),
        @ApiResponse(code = 201, message = "The configuration item type was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
        @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
        @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
        @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateItemType(
        @ApiParam(
            name = "key",
            value = "A string which uniquely identifies the item type and never changes.",
            required = true,
            example = "item_type_01"
        )
        @PathVariable("key")
        String key,
        @RequestBody
        ItemTypeData itemType
    ) {
        Result result = data.createOrUpdateItemType(key, itemType);
        return ResponseEntity.status(result.getStatus()).body(result);
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
            for (String item : items) {
                String[] parts = item.split("->");
                attrMap.put(parts[0], parts[1]);
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
        consumes = {"application/json"},
        produces = {"application/json"})
    @ApiResponses(value = {
        @ApiResponse(code = 200, message = "No changes where performed to the link."),
        @ApiResponse(code = 201, message = "The configuration link was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
        @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
        @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
        @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateLink(
            @PathVariable("key") String key,
            @RequestBody LinkData link) {
        Result result = data.createOrUpdateLink(key, link);
        return ResponseEntity.status(result.getStatus()).body(result);
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

    @ApiOperation(
            value = "Get an item link based on the specified key.",
            notes = "Use this search to retrieve a specific item link when its natural key is known.")
    @RequestMapping(
            path = "/link/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<LinkData> getLink(
            @PathVariable("key") String key) {
        return ResponseEntity.ok(data.getLink(key));
    }

    @ApiOperation(
            value = "Search for item linkss based on the specified filters (provided via a query string).",
            notes = "Use this function to retrieve item links based on type, tags and date range as required. " +
                    "Results are limited by the top parameter.")
    @RequestMapping(
            path = "/link"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<Wrapper> getLinks(
            @RequestParam(value = "type", required = false) String linkTypeKey
            , @RequestParam(value = "tag", required = false) String tag
            , @RequestParam(value = "startItemKey", required = false) String startItemKey
            , @RequestParam(value = "endItemKey", required = false) String endItemKey
            , @RequestParam(value = "createdFrom", required = false) String createdFromDate
            , @RequestParam(value = "createdTo", required = false) String createdToDate
            , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
            , @RequestParam(value = "updatedTo", required = false) String updatedToDate
            , @RequestParam(value = "model", required = false) String modelKey
            , @RequestParam(value = "top", required = false, defaultValue = "100") Integer top
    ) {
        List<String> tagList = null;
        if (tag != null) {
            String[] tags = tag.split("[|]"); // separate tags using pipes in the query string
            tagList = Arrays.asList(tags);
        }
        LinkList list = data.findLinks(
                linkTypeKey,
                startItemKey,
                endItemKey,
                tagList,
                getDate(createdFromDate),
                getDate(createdToDate),
                getDate(updatedFromDate),
                getDate(updatedToDate),
                modelKey,
                top
        );
        return ResponseEntity.ok(list);
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
    public ResponseEntity<Result> deleteLinkType(
            @PathVariable("key")
            String key,
            @ApiParam(
                name = "force",
                value = "If true, it forces the deletion of existing links of the link type.",
                required = false,
                example = "?force"
            )
            @RequestParam(value = "force", defaultValue = "false", required = false)
            boolean force
    ) {
        return ResponseEntity.ok(data.deleteLinkType(key, force));
    }

    @ApiOperation(
        value = "Creates a new item link type.",
        notes = "")
    @RequestMapping(
        path = "/linktype/{key}"
        , method = RequestMethod.PUT)
    @ApiResponses(value = {
        @ApiResponse(code = 200, message = "No changes where performed to the link type."),
        @ApiResponse(code = 201, message = "The configuration link type was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
        @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
        @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
        @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateLinkType(
        @PathVariable("key")
        String key,
        @RequestBody
        LinkTypeData linkType
    ) {
        Result result = data.createOrUpdateLinkType(key, linkType);
        return ResponseEntity.status(result.getStatus()).body(result);
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
            for (String item : items) {
                String[] parts = item.split("->");
                attrMap.put(parts[0], parts[1]);
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
    @ApiResponses(value = {
        @ApiResponse(code = 200, message = "No changes where performed to the link rule."),
        @ApiResponse(code = 201, message = "The configuration link rule was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
        @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
        @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
        @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateLinkRule(
        @PathVariable("key")
        String key,
        @RequestBody
        LinkRuleData linkRule
    ) {
        Result result = data.createOrUpdateLinkRule(key, linkRule);
        return ResponseEntity.status(result.getStatus()).body(result);
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
    public ResponseEntity<Result> deleteModel(
            @PathVariable("key")
            String key,
            @ApiParam(
                name = "force",
                value = "If true, it forces the deletion of existing item and link types associated with this model.",
                required = false,
                example = "?force"
            )
            @RequestParam(value="force", required = false, defaultValue = "false")
            boolean force
    ) {
        return ResponseEntity.ok(data.deleteModel(key, force));
    }

    @ApiOperation(
        value = "Creates a new model.",
        notes = "")
    @RequestMapping(
        path = "/model/{key}"
        , method = RequestMethod.PUT)
    @ApiResponses(value = {
        @ApiResponse(code = 200, message = "No changes where performed to the model."),
        @ApiResponse(code = 201, message = "The configuration model was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
        @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
        @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
        @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateModel(
        @ApiParam(
            name = "key",
            value = "A string which uniquely identifies the model and never changes.",
            required = true,
            example = "model_01_abc"
        )
        @PathVariable("key") String key,
        @RequestBody ModelData model
    ) {
        Result result = data.createOrUpdateModel(key, model);
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Get a model based on the specified key.",
            notes = "Use this search to retrieve a specific model when its natural key is known.")
    @RequestMapping(
            path = "/model/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ModelData> getModel(
            @ApiParam(
                name = "key",
                value = "A string which uniquely identifies the model and never changes.",
                required = true,
                example = "model_01_abc"
            )
            @PathVariable("key") String key) {
        return ResponseEntity.ok(data.getModel(key));
    }

    @ApiOperation(
            value = "Get all models.",
            notes = "Use this search to retrieve the list of all models known to the system.")
    @RequestMapping(
            path = "/model"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ModelDataList> getModels() {
        return ResponseEntity.ok(data.getModels());
    }

    @ApiOperation(
            value = "Get a list of item types, link types and link rules for a specified model.",
            notes = "")
    @RequestMapping(
            path = "/model/{key}/data"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<TypeGraphData> getTypeData(
            @ApiParam(
                name = "key",
                value = "A string which uniquely identifies the model and never changes.",
                required = true,
                example = "model_01_abc"
            )
            @PathVariable("key") String modelKey
    ) {
        TypeGraphData graph = data.getTypeDataByModel(modelKey);
        return ResponseEntity.ok(graph);
    }

    /*
        MISCELLANEOUS
     */

    @ApiOperation(
            value = "Removes ALL configuration items and links from the database.",
            notes = "Use at your own risk ONLY for testing of the CMDB!")
    @RequestMapping(path = "/clear", method = RequestMethod.DELETE)
    public ResponseEntity<Result> clear() {
        Result result = data.clear();
        return ResponseEntity.status(result.getStatus()).body(result);
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
            path = "/tag/{item_key}/{tag}"
            , method = RequestMethod.PUT)
    public ResponseEntity<Result> updateTag(
            @PathVariable("item_key") String itemKey,
            @PathVariable("tag") String tag,
            @RequestBody JSONObject payload
    ) {
        return ResponseEntity.ok(data.updateTag(itemKey, tag, payload));
    }

    @ApiOperation(
            value = "Deletes an existing tag.",
            notes = "Takes the key of a root item and a tag label and deletes the matching tag.")
    @RequestMapping(
            path = "/tag/{item_key}/{tag}"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteTag(
            @PathVariable("item_key") String itemKey,
            @PathVariable("tag") String tag
    ) {
        return ResponseEntity.ok(data.deleteTag(itemKey, tag));
    }

    @ApiOperation(
            value = "Deletes all tags for an item.",
            notes = "Takes the key of a root item and deletes any associated tags.")
    @RequestMapping(
            path = "/tag/{item_key}"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteAllTags(
            @PathVariable("item_key") String itemKey
    ) {
        return ResponseEntity.ok(data.deleteTag(itemKey, null));
    }

    @ApiOperation(
            value = "Get a list of available tags for a specific item.",
            notes = "")
    @RequestMapping(
            path = "/tag/{item_key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<TagList> getItemTags(
            @PathVariable("item_key") String itemKey
    ) {
        TagList tags = data.getItemTags(itemKey);
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
    public ResponseEntity<ResultList> createOrUpdateData(@RequestBody GraphData graphData) {
        ResultList results = data.createOrUpdateData(graphData);
        return ResponseEntity.ok(results);
    }

    @ApiOperation(
            value = "Get a list of items and links that are children of a specified item for a specified item tag.",
            notes = "")
    @RequestMapping(
            path = "/data/{item_key}/tag/{tag}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<GraphData> getDataWithTag(
            @PathVariable("item_key") String itemKey,
            @PathVariable("tag") String tag
    ) {
        GraphData graph = data.getData(itemKey, tag);
        return ResponseEntity.ok(graph);
    }

    @ApiOperation(
            value = "Deletes an existing item and all its children.",
            notes = "")
    @RequestMapping(
            path = "/data/{item_key}"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteData(
            @PathVariable("item_key") String rootItemKey
    ) {
        return ResponseEntity.ok(data.deleteData(rootItemKey));
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