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

package org.gatblau.onix;

import com.fasterxml.jackson.databind.ObjectMapper;
import io.swagger.annotations.*;
import org.gatblau.onix.conf.Info;
import org.gatblau.onix.data.*;
import org.gatblau.onix.db.DbRepository;
import org.gatblau.onix.security.Crypto;
import org.json.simple.JSONObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletRequest;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.*;

@Api("ONIX CMDB Web API")
@RestController
public class WebAPI {

    @Autowired
    private DbRepository data;

    @Autowired
    private org.gatblau.onix.conf.Info info;

    @Autowired
    Crypto crypto;

    private ObjectMapper mapper = new ObjectMapper();
    private DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyyMMddHHmm");

    @ApiOperation(
            value = "Returns the username of the logged on user.",
            notes = "",
            response = String.class)
    @RequestMapping(value = "/user", method = RequestMethod.GET, produces = "application/json")
    public synchronized final String user() {
        final String username = SecurityContextHolder.getContext().getAuthentication().getName();
        return String.format("You are logged on as: '%s'", username);
    }

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
        return ResponseEntity.ok(data.checkReady());
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
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requester does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
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
            @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
            HttpServletRequest request, // required to check on http headers
            Authentication authentication // required to authenticate user
        ) {
        // check the request integrity and de-serialise the payload
        Tuple<ResponseEntity<Result>, ItemData> req = prepareRequest(request, payloadStr, ItemData.class);
        // if the data integrity check or de-serialisation fails, returns
        if (req.response != null) { return req.response; }
        // now ready to process the request
        Result result = data.createOrUpdateItem(key, req.payload, getRole(authentication));
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
            @PathVariable("key") String key,
            Authentication authentication
    ) {
        Result result = data.deleteItem(key, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Deletes all configuration items.",
            notes = "")
    @RequestMapping(
            path = "/item"
            , method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteAllItems(
        Authentication authentication
    ) {
        Result result = data.deleteAllItems(getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
        @RequestParam(required = false, name = "links", defaultValue = "false") // true to retrieve link information
        boolean links,
        Authentication authentication
    ) {
        ItemData item = data.getItem(key, links, getRole(authentication));
        if (item != null) {
            return ResponseEntity.ok(item);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
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
                    name = "attrs",
                    value = "A list of atributes (key,value) pair separated by '|' in the query string. If no value is passed then all item types are retrieved.",
                    required = false,
                    example = "key1,value1|key2,value2|key3,value3"
            )
            @RequestParam(value = "attrs", required = false)
            String attrs,
            @ApiParam(
                name = "keyIX",
                value = "The index of the key used to encrypt data in the items. Values can be 0: no key, 1 or 2 for key 1 or 2 respectively.",
                required = false,
                example = "1"
            )
            @RequestParam(value = "keyIX", required = false)
            Short encKeyIx,
            @ApiParam(
                name = "top",
                value = "The maximum number of items to retrieve.",
                required = false,
                example = "12-03-18",
                defaultValue = "100"
            )
            @RequestParam(value = "top", required = false, defaultValue = "100")
            Integer top,
            Authentication authentication
    ) {
        List<String> tagList = null;
        if (tag != null) {
            String[] tags = tag.split("[|]"); // separate tags using pipes in the query string
            tagList = Arrays.asList(tags);
        }
        Map<String, String> attributes = new HashMap<>();
        if (attrs != null) {
            String[] attrsPairs = attrs.split("[|]"); // separate key value pairs using pipes in the query string
            for (int i=0; i<attrsPairs.length; i++){
                String[] slice = attrsPairs[i].split(",");
                attributes.put(slice[0], slice[1]);
            }
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
                attributes,
                encKeyIx,
                top,
                getRole(authentication)
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
        String key,
        Authentication authentication
    ) {
        JSONObject meta = data.getItemMeta(key, null, getRole(authentication));
        if (meta != null){
            return ResponseEntity.ok(meta);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
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
        String filter,
        Authentication authentication
    ) {
        JSONObject meta = data.getItemMeta(key, filter, getRole(authentication));
        if (meta != null){
            return ResponseEntity.ok(meta);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    @ApiOperation(
            value = "Gets the children associated with the specified item.",
            notes = "Use this function to retrieve a list of items linked to an item.")
    @RequestMapping(
            path = "/item/{key}/children"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<Wrapper> getItemChildren(
            @ApiParam(
                    name = "key",
                    value = "A string which uniquely identifies the item and never changes.",
                    required = true,
                    example = "item_01_abc"
            )
            @PathVariable("key")
                    String key,
            Authentication authentication
    ) {
        ItemList list = data.getItemChildren(key, getRole(authentication));
        return ResponseEntity.ok(list);
    }

    /*
        PARTITIONS
     */
    @ApiOperation(
            value = "Deletes a logical partition.",
            notes = "")
    @RequestMapping(
            path = "/partition/{key}"
            , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deletePartition(
            @ApiParam(
                name = "key",
                value = "A string which uniquely identifies the logical partition and never changes.",
                required = true,
                example = "partition_01"
            )
            @PathVariable("key")
            String key,
            Authentication authentication
    ) {
        Result result = data.deletePartition(key, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Creates a new logical partition for RBAC.",
            notes = "The role executing this action has to be an admin role.")
    @RequestMapping(
            path = "/partition/{key}"
            , method = RequestMethod.PUT)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "No changes where performed to the configuration item type."),
            @ApiResponse(code = 201, message = "The configuration item type was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdatePartition(
            @ApiParam(
                    name = "key",
                    value = "A string which uniquely identifies the partition and never changes.",
                    required = true,
                    example = "part_01"
            )
            @PathVariable("key") String key,
            @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
            HttpServletRequest request, // required to check on http headers
            Authentication authentication // required to authenticate user
    ) {
        // check the request integrity and de-serialise the payload
        Tuple<ResponseEntity<Result>, PartitionData> req = prepareRequest(request, payloadStr, PartitionData.class);
        // if the data integrity check or de-serialisation fails, returns
        if (req.response != null) { return req.response; }
        // now ready to process the request
        Result result = data.createOrUpdatePartition(key, req.payload, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Get all logical partitions.",
            notes = "")
    @RequestMapping(
            path = "/partition"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<PartitionDataList> getAllPartitions(Authentication authentication) {
        return ResponseEntity.ok(data.getAllPartitions(getRole(authentication)));
    }

    @ApiOperation(
            value = "Get a logical partition based on the specified key.",
            notes = "")
    @RequestMapping(
            path = "/partition/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<PartitionData> getParttition(
            @PathVariable("key")
            String key,
            Authentication authentication
    ) {
        PartitionData partition = data.getPartition(key, getRole(authentication));
        if (partition != null){
            return ResponseEntity.ok(partition);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    /*
        ROLES
     */
    @ApiOperation(
            value = "Deletes a logical partition.",
            notes = "")
    @RequestMapping(
            path = "/role/{key}"
            , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteRole(
            @ApiParam(
                name = "key",
                value = "A string which uniquely identifies the role and never changes.",
                required = true,
                example = "role_01"
            )
            @PathVariable("key")
            String key,
            Authentication authentication
    ) {
        Result result = data.deleteRole(key, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Creates a new role for RBAC.",
            notes = "The role executing this action has to be an admin role.")
    @RequestMapping(
            path = "/role/{key}"
            , method = RequestMethod.PUT)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "No changes where performed to the configuration item type."),
            @ApiResponse(code = 201, message = "The configuration item type was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateRole(
            @ApiParam(
                name = "key",
                value = "A string which uniquely identifies the role and never changes.",
                required = true,
                example = "role_01"
            )
            @PathVariable("key") String key,
            @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
            HttpServletRequest request, // required to check on http headers
            Authentication authentication // required to authenticate user
    ) {
        // check the request integrity and de-serialise the payload
        Tuple<ResponseEntity<Result>, RoleData> req = prepareRequest(request, payloadStr, RoleData.class);
        // if the data integrity check or de-serialisation fails, returns
        if (req.response != null) { return req.response; }
        // now ready to process the request
        Result result = data.createOrUpdateRole(key, req.payload, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Get all roles.",
            notes = "")
    @RequestMapping(
            path = "/role"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<RoleDataList> getAllRoles(Authentication authentication) {
        return ResponseEntity.ok(data.getAllRoles(getRole(authentication)));
    }

    @ApiOperation(
            value = "Get a logical partition based on the specified key.",
            notes = "")
    @RequestMapping(
            path = "/role/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<RoleData> getRole(
            @PathVariable("key")
            String key,
            Authentication authentication
    ) {
        RoleData role = data.getRole(key, getRole(authentication));
        if (role != null) {
            return ResponseEntity.ok(role);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    /*
        PRIVILEGES
     */
    @ApiOperation(
        value = "Creates a new or updates an existing privilege granting a role access to a partition.",
        notes = "")
    @RequestMapping(
        path = "/privilege/{key}"
        , method = RequestMethod.PUT)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "No changes where performed to the configuration privilege."),
            @ApiResponse(code = 201, message = "The configuration privilege was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the right to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> addPrivilege(
        @PathVariable("key") String key,
        @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
        HttpServletRequest request, // required to check on http headers
        Authentication authentication // required to authenticate user
    ) {
        // check the request integrity and de-serialise the payload
        Tuple<ResponseEntity<Result>, PrivilegeData> req = prepareRequest(request, payloadStr, PrivilegeData.class);
        // if the data integrity check or de-serialisation fails, returns
        if (req.response != null) { return req.response; }
        // now ready to process the request
        Result result = data.createOrUpdatePrivilege(key, req.payload, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Get a list of privileges for the specified role.",
            notes = "")
    @RequestMapping(
            path = "/role/{key}/privilege"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<PrivilegeDataList> getRolePrivileges(
            @PathVariable("key") String key,
            Authentication authentication
    ) {
        return ResponseEntity.ok(data.getPrivilegesByRole(key, getRole(authentication)));
    }

    @ApiOperation(
            value = "Get the privilege for the specified partition and role.",
            notes = "")
    @RequestMapping(
            path = "/privilege/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<PrivilegeData> getPrivilege(
            @PathVariable("key") String key,
            Authentication authentication
    ) {
        PrivilegeData privilege = data.getPrivilege(key, getRole(authentication));
        if (privilege != null) {
            return ResponseEntity.ok(privilege);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    @ApiOperation(
            value = "Deletes a privilege.",
            notes = "")
    @RequestMapping(
            path = "/privilege/{key}"
            , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deletePrivilege(
        @PathVariable("key")
        String key,
        Authentication authentication
    ) {
        Result result = data.removePrivilege(key, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
    public ResponseEntity<Result> deleteItemTypes(
        Authentication authentication
    ) {
        Result result = data.deleteItemTypes(getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
        Authentication authentication
    ) {
        Result result = data.deleteItemType(key, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requester does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateItemType(
        @ApiParam(
            name = "key",
            value = "A string which uniquely identifies the item type and never changes.",
            required = true,
            example = "item_type_01"
        )
        @PathVariable("key") String key,
        @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
        HttpServletRequest request, // required to check on http headers
        Authentication authentication // required to authenticate user
    ) {
        // check the request integrity and de-serialise the payload
        Tuple<ResponseEntity<Result>, ItemTypeData> req = prepareRequest(request, payloadStr, ItemTypeData.class);
        // if the data integrity check or de-serialisation fails, returns
        if (req.response != null) { return req.response; }
        // now ready to process the request
        Result result = data.createOrUpdateItemType(key, req.payload, getRole(authentication));
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
    public ResponseEntity<ItemTypeData> getItemType(
            @PathVariable("key")
            String key,
            Authentication authentication
    ) {
        ItemTypeData itemType = data.getItemType(key, getRole(authentication));
        if (itemType != null) {
            return ResponseEntity.ok(itemType);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
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
              @RequestParam(value = "createdFrom", required = false) String createdFromDate
            , @RequestParam(value = "createdTo", required = false) String createdToDate
            , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
            , @RequestParam(value = "updatedTo", required = false) String updatedToDate
            , @RequestParam(value = "model", required = false) String modelKey
            , Authentication authentication
    ) {
        ItemTypeList itemTypes = data.getItemTypes(
            getDate(createdFromDate),
            getDate(createdToDate),
            getDate(updatedFromDate),
            getDate(updatedToDate),
            modelKey,
            getRole(authentication)
        );
        return ResponseEntity.ok(itemTypes);
    }

    /*
        ITEM TYPE ATTRIBUTES
     */
    @ApiOperation(
            value = "Creates a new attribute for an item type.",
            notes = "")
    @RequestMapping(
            path = "/itemtype/{item_type_key}/attribute/{type_attr_key}"
            , method = RequestMethod.PUT)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "No changes where performed to the configuration item type attribute."),
            @ApiResponse(code = 201, message = "The configuration item type attribute was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
            @ApiResponse(code = 400, message = "Bad request. The request contained the wrong payload. Check the request body is well form, and if so, that the attributes passed in are or the correct type and name. "),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateItemTypeAttribute(
            @ApiParam(
                    name = "item type key",
                    value = "A string which uniquely identifies the item type for the attribute.",
                    required = true,
                    example = "item_type_01"
            )
            @PathVariable("item_type_key")
                    String itemTypeKey,
            @ApiParam(
                    name = "type attribute key",
                    value = "A string which uniquely identifies the attribute for item type.",
                    required = true,
                    example = "item_type_attr_01"
            )
            @PathVariable("type_attr_key")
                    String typeAttrKey,
            @RequestBody
                    ItemTypeAttrData typeAttr,
            Authentication authentication
    ) {
        Result result = data.createOrUpdateItemTypeAttr(itemTypeKey, typeAttrKey, typeAttr, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Deletes a configuration item type.",
            notes = "")
    @RequestMapping(
            path = "/itemtype/{item_type_key}/attribute/{type_attr_key}"
            , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteItemTypeAttribute(
            @ApiParam(
                    name = "item type key",
                    value = "A string which uniquely identifies the item type (never changes).",
                    required = true,
                    example = "item_type_01"
            )
            @PathVariable("item_type_key")
                    String itemTypeKey,
            @ApiParam(
                    name = "type attribute key",
                    value = "A string which uniquely identifies the type attribute associated with the item type (never changes).",
                    required = true,
                    example = "item_type_01"
            )
            @PathVariable("type_attr_key")
                    String typeAttrKey,
            Authentication authentication
    ) {
        Result result = data.deleteItemTypeAttr(itemTypeKey, typeAttrKey, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Get all the attributes for the specified item type.",
            notes = "Use this search to retrieve the specification of the attributes for an item based on its item type.")
    @RequestMapping(
            path = "/itemtype/{item_type_key}/attribute"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ItemTypeAttrList> getItemTypeAttrs(
            @PathVariable("item_type_key") String itemTypeKey,
            Authentication authentication
    ) {
        ItemTypeAttrList itemTypeAttrs = data.getItemTypeAttributes(itemTypeKey, getRole(authentication));
        if (itemTypeAttrs != null) {
            return ResponseEntity.ok(itemTypeAttrs);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    @ApiOperation(
            value = "Get the attribute for the specified item type and attribute key.",
            notes = "")
    @RequestMapping(
            path = "/itemtype/{item_type_key}/attribute/{type_attr_key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ItemTypeAttrData> getItemTypeAttr(
            @PathVariable("item_type_key") String itemTypeKey,
            @PathVariable("type_attr_key") String typeAttrKey,
            Authentication authentication
    ) {
        ItemTypeAttrData typeAttr = data.getItemTypeAttribute(itemTypeKey, typeAttrKey, getRole(authentication));
        if (typeAttr != null) {
            return ResponseEntity.ok(typeAttr);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
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
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateLink(
            @PathVariable("key") String key,
            @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
            HttpServletRequest request, // required to check on http headers
            Authentication authentication // required to authenticate user
    ) {
        // check the request integrity and de-serialise the payload
        Tuple<ResponseEntity<Result>, LinkData> req = prepareRequest(request, payloadStr, LinkData.class);
        // if the data integrity check or de-serialisation fails, returns
        if (req.response != null) { return req.response; }
        // now ready to process the request
        Result result = data.createOrUpdateLink(key, req.payload, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Deletes a link between two existing configuration items.",
            notes = "Use this operation to delete links between existing items.")
    @RequestMapping(path = "/link/{key}", method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteLink(
            @PathVariable("key") String key,
            Authentication authentication
    ) {
        Result result = data.deleteLink(key, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
        @PathVariable("key")
        String key,
        Authentication authentication
    ) {
        LinkData link = data.getLink(key, getRole(authentication));
        if (link != null) {
            return ResponseEntity.ok(link);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
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
            , @ApiParam(
                name = "keyIX",
                value = "The index of the key used to encrypt data in the links. Values can be 0: no key, 1 or 2 for key 1 or 2 respectively.",
                required = false,
                example = "1"
            )
              @RequestParam(value = "keyIX", required = false) Short encKeyIx
            , @RequestParam(value = "model", required = false) String modelKey
            , @RequestParam(value = "top", required = false, defaultValue = "100") Integer top
            , Authentication authentication
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
            encKeyIx,
            top,
            getRole(authentication)
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
    public ResponseEntity<Result> deleteLinkTypes(
            Authentication authentication
    ) {
        Result result = data.deleteLinkTypes(getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
            Authentication authentication
    ) {
        Result result = data.deleteLinkType(key, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateLinkType(
        @PathVariable("key") String key,
        @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
        HttpServletRequest request, // required to check on http headers
        Authentication authentication // required to authenticate user
    ) {
        // check the request integrity and de-serialise the payload
        Tuple<ResponseEntity<Result>, LinkTypeData> req = prepareRequest(request, payloadStr, LinkTypeData.class);
        // if the data integrity check or de-serialisation fails, returns
        if (req.response != null) { return req.response; }
        // now ready to process the request
        Result result = data.createOrUpdateLinkType(key, req.payload, getRole(authentication));
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
    public ResponseEntity<LinkTypeData> getLinkType(
            @PathVariable("key")
            String key,
            Authentication authentication
        ) {
        LinkTypeData linkType = data.getLinkType(key, getRole(authentication));
        if (linkType != null) {
            return ResponseEntity.ok(linkType);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
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
              @RequestParam(value = "createdFrom", required = false) String createdFromDate
            , @RequestParam(value = "createdTo", required = false) String createdToDate
            , @RequestParam(value = "updatedFrom", required = false) String updatedFromDate
            , @RequestParam(value = "updatedTo", required = false) String updatedToDate
            , @RequestParam(value = "model", required = false) String modelKey
            , Authentication authentication
    ) {
        LinkTypeList linkTypes = data.getLinkTypes(
                getDate(createdFromDate),
                getDate(createdToDate),
                getDate(updatedFromDate),
                getDate(updatedToDate),
                modelKey,
                getRole(authentication));
        return ResponseEntity.ok(linkTypes);
    }

    /*
        LINK TYPE ATTRIBUTES
     */
    @ApiOperation(
            value = "Creates a new attribute for a link type.",
            notes = "")
    @RequestMapping(
            path = "/linktype/{link_type_key}/attribute/{link_attr_key}"
            , method = RequestMethod.PUT)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "No changes where performed to the configuration link type attribute."),
            @ApiResponse(code = 201, message = "The configuration link type attribute was created or updated. The operation attribute in the response can be used to determined if an insert or an update was performed. ", response = Result.class),
            @ApiResponse(code = 400, message = "Bad request. The request contained the wrong payload. Check the request body is well form, and if so, that the attributes passed in are or the correct type and name. "),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateLinkTypeAttribute(
            @ApiParam(
                    name = "link type key",
                    value = "A string which uniquely identifies the link type for the attribute.",
                    required = true,
                    example = "link_type_01"
            )
            @PathVariable("link_type_key") String linkTypeKey,
            @ApiParam(
                    name = "type attribute key",
                    value = "A string which uniquely identifies the attribute for the link type.",
                    required = true,
                    example = "link_type_attr_01"
            )
            @PathVariable("link_attr_key") String typeAttrKey,
            @RequestBody LinkTypeAttrData typeAttr,
            Authentication authentication
    ) {
        Result result = data.createOrUpdateLinkTypeAttr(linkTypeKey, typeAttrKey, typeAttr, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    @ApiOperation(
            value = "Get the attribute for the specified link type and attribute key.",
            notes = "")
    @RequestMapping(
            path = "/linktype/{link_type_key}/attribute/{type_attr_key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<LinkTypeAttrData> getLinkTypeAttr(
            @PathVariable("link_type_key") String linkTypeKey,
            @PathVariable("type_attr_key") String typeAttrKey,
            Authentication authentication
    ) {
        LinkTypeAttrData typeAttr = data.getLinkTypeAttribute(linkTypeKey, typeAttrKey, getRole(authentication));
        if (typeAttr != null) {
            return ResponseEntity.ok(typeAttr);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    @ApiOperation(
            value = "Get all the attributes for the specified item type.",
            notes = "Use this search to retrieve the specification of the attributes for an item based on its item type.")
    @RequestMapping(
            path = "/linktype/{link_type_key}/attribute"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<LinkTypeAttrList> getLinkTypeAttrs(
            @PathVariable("link_type_key") String linkTypeKey,
            Authentication authentication
    ) {
        LinkTypeAttrList linkTypeAttrs = data.getLinkTypeAttributes(linkTypeKey, getRole(authentication));
        if (linkTypeAttrs != null) {
            return ResponseEntity.ok(linkTypeAttrs);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    @ApiOperation(
            value = "Deletes an item link type attribute.",
            notes = "")
    @RequestMapping(
            path = "/linktype/{link_type_key}/attribute/{type_attr_key}"
            , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteLinkTypeAttribute(
            @PathVariable("link_type_key") String linkTypeKey,
            @PathVariable("type_attr_key") String typeAttrKey,
            Authentication authentication
    ) {
        Result result = data.deleteLinkTypeAttr(linkTypeKey, typeAttrKey, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    /*
        LINK RULES
     */
    @ApiOperation(
            value = "Get a link rule based on the specified key.",
            notes = "")
    @RequestMapping(
            path = "/linkrule/{key}"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<LinkRuleData> getLinkRule(
            @PathVariable("key") String key,
            Authentication authentication
    ) {
        LinkRuleData linkRule = data.getLinkRule(key, getRole(authentication));
        if (linkRule != null) {
            return ResponseEntity.ok(linkRule);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    @ApiOperation(
            value = "Deletes all non-system specific item link types.",
            notes = "")
    @RequestMapping(
            path = "/linkrule"
            , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteLinkRules(
            Authentication authentication
    ) {
        Result result = data.deleteLinkRules(getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<Result> createOrUpdateLinkRule(
        @PathVariable("key") String key,
        @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
        HttpServletRequest request, // required to check on http headers
        Authentication authentication // required to authenticate user
    ) {
        // check the request integrity and de-serialise the payload
        Tuple<ResponseEntity<Result>, LinkRuleData> req = prepareRequest(request, payloadStr, LinkRuleData.class);
        // if the data integrity check or de-serialisation fails, returns
        if (req.response != null) { return req.response; }
        // now ready to process the request
        Result result = data.createOrUpdateLinkRule(key, req.payload, getRole(authentication));
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
            , Authentication authentication
    ) {
        LinkRuleList linkRules = data.getLinkRules(
                linkType,
                startItemType,
                endItemType,
                getDate(createdFromDate),
                getDate(createdToDate),
                getDate(updatedFromDate),
                getDate(updatedToDate),
                getRole(authentication));
        return ResponseEntity.ok(linkRules);
    }

    @ApiOperation(
            value = "Deletes a link rule for a specified key.",
            notes = "")
    @RequestMapping(
            path = "/linkrule/{key}"
            , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteLinkRule(
            @PathVariable("key") String key,
            Authentication authentication
    ) {
        Result result = data.deleteLinkRule(key, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
            Authentication authentication
    ) {
        Result result = data.deleteModel(key, getRole(authentication));
        return ResponseEntity.status(result.getStatus()).body(result);
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
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requestor does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
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
        @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
        HttpServletRequest request, // required to check on http headers
        Authentication authentication // required to authenticate user

    ) {
        // check the request integrity and de-serialise the payload
        Tuple<ResponseEntity<Result>, ModelData> req = prepareRequest(request, payloadStr, ModelData.class);
        // if the data integrity check or de-serialisation fails, returns
        if (req.response != null) { return req.response; }
        // now ready to process the request
        Result result = data.createOrUpdateModel(key, req.payload, getRole(authentication));
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
            @PathVariable("key") String key,
            Authentication authentication
    ) {
        ModelData model = data.getModel(key, getRole(authentication));
        // if the model is null then return 404 not found
        if (model == null) {
            return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
        }
        return ResponseEntity.ok(model);
    }

    @ApiOperation(
            value = "Get all models.",
            notes = "Use this search to retrieve the list of all models known to the system.")
    @RequestMapping(
            path = "/model"
            , method = RequestMethod.GET
            , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ModelDataList> getModels(
            Authentication authentication
    ) {
        return ResponseEntity.ok(data.getModels(getRole(authentication)));
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
            @PathVariable("key") String modelKey,
            Authentication authentication
    ) {
        TypeGraphData graph = data.getTypeDataByModel(modelKey, getRole(authentication));
        if (graph != null) {
            return ResponseEntity.ok(graph);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    /*
        MISCELLANEOUS
     */

    @ApiOperation(
            value = "Removes ALL configuration items and links from the database.",
            notes = "Use at your own risk ONLY for testing of the CMDB!")
    @RequestMapping(path = "/clear", method = RequestMethod.DELETE)
    public ResponseEntity<Result> clear(
            Authentication authentication
    ) {
        Result result = data.clear(getRole(authentication));
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
        Result result = data.createTag(payload);
        return ResponseEntity.status(result.getStatus()).body(result);
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
        Result result = data.updateTag(itemKey, tag, payload);
        return ResponseEntity.status(result.getStatus()).body(result);
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
        Result result = data.deleteTag(itemKey, tag);
        return ResponseEntity.status(result.getStatus()).body(result);
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
        Result result = data.deleteTag(itemKey, null);
        return ResponseEntity.status(result.getStatus()).body(result);
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
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "The request was received and processed."),
            @ApiResponse(code = 400, message = "The request payload is malformed."),
            @ApiResponse(code = 401, message = "The request was unauthorised. The requester does not have the privilege to execute the request. "),
            @ApiResponse(code = 404, message = "The request was made to an URI which does not exist on the server. "),
            @ApiResponse(code = 422, message = "The request failed MD5 checksum validation, only enabled if the Content-MD5 header is added to the request.. "),
            @ApiResponse(code = 500, message = "There was an internal side server error.", response = Result.class)}
    )
    public ResponseEntity<ResultList> createOrUpdateData(
            @RequestBody String payloadStr, // required to compute MD5 checksum and de-serialise data
            HttpServletRequest request, // required to check on http headers
            Authentication authentication // required to authenticate user
        ) {
        ResultList resultList = new ResultList();
        // if the payload checksum does not match the one provided in the header
        if (!payloadIntegrityOk(request, payloadStr)) {
            // does not attempt to process the request as the integrity checksum does not match
            // assume http body integrity has been compromised
            Result result = new Result();
            result.setError(true);
            result.setMessage("HTTP body has failed MD5 checksum validation: the checksum sent by the client does not match the one calculated by the server.");
            resultList.add(result);
            return ResponseEntity.status(HttpStatus.UNPROCESSABLE_ENTITY).body(resultList);
        }
        GraphData payload = null;
        try {
            // tries to de-serialise payload
            payload = mapper.readValue(payloadStr, GraphData.class);
        } catch (Exception ex) {
            // cannot de-serialise the payload so assume a bad http request
            Result result = new Result();
            result.setError(true);
            result.setMessage(String.format("Issue reading the request body, revise the content that is being sent to the server: %s ", ex.getMessage()));
            resultList.add(result);
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(resultList);
        }
        ResultList results = data.createOrUpdateData(payload, getRole(authentication));
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
        if (graph != null){
            return ResponseEntity.ok(graph);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
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
        Result result = data.deleteData(rootItemKey);
        return ResponseEntity.status(result.getStatus()).body(result);
    }

    /*
        ENCRYPTION KEYS
     */
    @ApiOperation(
            value = "Returns a new secret key for configuration data encryption.",
            notes = "Use this endpoint to generate encryption keys that can be used to encrypt/decrypt configuration data.",
            response = JSONObject.class)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "Successful request.", response = JSONObject.class)}
    )
    @RequestMapping(value = "/enckey/generate", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<JSONObject> key() {
        JSONObject response = new JSONObject();
        response.put("key", crypto.newKey());
        return ResponseEntity.ok(response);
    }

    @ApiOperation(
            value = "Gets the usage status of encryption keys for item meta and/or txt attributes.",
            notes = "Use this endpoint to understand the state of use of keys and progress on key rotation.",
            response = JSONObject.class)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "Successful request.", response = JSONObject.class)}
    )
    @RequestMapping(value = "/enckey/status", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<EncKeyStatusData> getKeyStatus(Authentication authentication) {
        EncKeyStatusData keyStatus = data.getKeyStatus(getRole(authentication));
        return ResponseEntity.ok(keyStatus);
    }

    @ApiOperation(
            value = "Invokes the key rotation routine specifying how many items to process at a time.",
            notes = "Use this endpoint to gradually rotate encryption keys for meta and txt fields. " +
                    "This function only works if the default key expiry date is in the past. " +
                    "The rotation is always from the default key to the secondary key. " +
                    "To understand the status of key usage invoke the \"/enckey/status\" endpoint.",
            response = JSONObject.class)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "Successful request.", response = JSONObject.class)}
    )
    @RequestMapping(value = "/enckey/rotate/item/{limit}", method = RequestMethod.POST, produces = "application/json")
    public ResponseEntity<ResultList> rotateItemKey(
            @PathVariable("limit") int limit,
            Authentication authentication) {
        return ResponseEntity.ok(data.rotateItemKeys(limit, getRole(authentication)));
    }

    @ApiOperation(
            value = "Invokes the key rotation routine specifying how many links to process at a time.",
            notes = "Use this endpoint to gradually rotate encryption keys for meta and txt fields. " +
                    "This function only works if the default key expiry date is in the past. " +
                    "The rotation is always from the default key to the secondary key. " +
                    "To understand the status of key usage invoke the \"/enckey/status\" endpoint.",
            response = JSONObject.class)
    @ApiResponses(value = {
            @ApiResponse(code = 200, message = "Successful request.", response = JSONObject.class)}
    )
    @RequestMapping(value = "/enckey/rotate/link/{limit}", method = RequestMethod.POST, produces = "application/json")
    public ResponseEntity<ResultList> rotateLinkKey(
            @PathVariable("limit") int limit,
            Authentication authentication) {
        return ResponseEntity.ok(data.rotateLinkKeys(limit, getRole(authentication)));
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

    private String[] getRole(Authentication authentication) {
        // if the service is configured not to use authentication
        if (authentication == null) {
            // then return the ADMIN role
            return new String[]{"ADMIN"};
        }
        String[] roles = new String[authentication.getAuthorities().size()];
        // otherwise uses the role in the first authority
        int ix = 0;
        for (GrantedAuthority authority : authentication.getAuthorities()){
            String r = authority.getAuthority();
            if (r.startsWith("ROLE_")) {
                roles[ix] = r.substring("ROLE_".length());
            }
            ix++;
        }
        return roles;
    }

    private <T> Tuple<ResponseEntity<Result>, T> prepareRequest(HttpServletRequest request, String payloadStr, Class<T> valueType) {
        // if the payload checksum does not match the one provided in the header
        if (!payloadIntegrityOk(request, payloadStr)) {
            // does not attempt to process the request as the integrity checksum does not match
            // assume http body integrity has been compromised
            Result result = new Result();
            result.setError(true);
            result.setMessage("HTTP body has failed MD5 checksum validation: the checksum sent by the client does not match the one calculated by the server.");
            return new Tuple(ResponseEntity.status(HttpStatus.UNPROCESSABLE_ENTITY).body(result), null);
        }
        T payload = null;
        try {
            // tries to de-serialise payload
            payload = mapper.readValue(payloadStr, valueType);
        } catch (Exception ex) {
            // cannot de-serialise the payload so assume a bad http request
            Result result = new Result();
            result.setError(true);
            result.setMessage(ex.getMessage());
            return new Tuple(ResponseEntity.status(HttpStatus.BAD_REQUEST).body(result), null);
        }
        return new Tuple(null, payload);
    }

    private boolean payloadIntegrityOk(HttpServletRequest request, String httpBody) {
        Boolean valid = false;
        String requestSum = request.getHeader("Content-MD5");
        if (requestSum != null) {
            // if there is a checksum in the request, checks it matches the payload's
            try {
                // gets the payload checksum
                String payloadSum = getMD5Hash(httpBody);
                // valid if the request and the actual sums match
                valid = payloadSum.equals(requestSum);
            } catch (Exception ex) {
                ex.printStackTrace();
            }
        } else {
            // if there is not a checksum in the request, assumes valid
            valid = true;
        }
        return valid;
    }

    private String getMD5Hash(String objStr) {
        String sum = null;
        try {
            // gets the MD5 hashing algorithm
            MessageDigest md = MessageDigest.getInstance("MD5");
            // creates an MD5 hash of the passed-in string
            byte[] hash = md.digest(objStr.getBytes(StandardCharsets.UTF_8));
            // Base64 encode the hash
            byte[] encoded = Base64.getEncoder().encode(hash);
            // converts the byte[] to UTF8 string
            sum = new String(encoded, StandardCharsets.UTF_8);
        } catch (Exception ex) {
            ex.printStackTrace();
        }
        return sum;
    }

    class Tuple<X, Y> {
        public final X response;
        public final Y payload;
        public Tuple(X response, Y payload) {
            this.response = response;
            this.payload = payload;
        }
    }
}