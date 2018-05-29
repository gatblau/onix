package org.gatblau.onix;

import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.gatblau.onix.data.*;
import org.gatblau.onix.model.ItemType;
import org.json.simple.JSONObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;

@Api("ONIX CMDB Web API")
@RestController
public class WebAPI {

    @Autowired
    private Repository data;

    @ApiOperation(
        value = "Returns OK if the service is up and running.",
        notes = "Use it as a readiness probe for the service.",
        response = String.class)
    @RequestMapping(value = "/", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<Info> index() {
        return ResponseEntity.ok(new Info("Onix Configuration Management Database Service.", "1.0"));
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
            @RequestBody JSONObject payload) throws InterruptedException, IOException {
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
            @RequestBody JSONObject payload) throws InterruptedException, IOException {
        return ResponseEntity.ok(data.createOrUpdateLink(key, payload));
    }

    @ApiOperation(
        value = "Gets all links for the specified item.",
        notes = "Use this resource to find all of the links associated with a particular configuration item.")
    @RequestMapping(
          path = "/link/item/{key}"
        , method = RequestMethod.GET
        , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<LinkList> getLinksByItem(
            @PathVariable("key") String itemKey) throws InterruptedException, IOException {
        return ResponseEntity.ok(data.getLinksByItem(itemKey));
    }

    @ApiOperation(
        value = "Removes ALL configuration items and links from the database.",
        notes = "Use at your own risk ONLY for testing of the CMDB!")
    @RequestMapping(path = "/clear", method = RequestMethod.DELETE)
    public void clear() throws InterruptedException {
        data.clear();
    }

    @ApiOperation(
        value = "Deletes a link between two existing configuration items.",
        notes = "Use this operation to delete links between existing items.")
    @RequestMapping(path = "/link/{key}", method = RequestMethod.DELETE)
    public ResponseEntity<Result> deleteLink(
            @PathVariable("key") String key
    ) throws InterruptedException {
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
    ) throws InterruptedException {
        return ResponseEntity.ok(data.deleteItem(key));
    }

    @ApiOperation(
        value = "Deletes a configuration item type.",
        notes = "")
    @RequestMapping(
          path = "/itemtype"
        , method = RequestMethod.DELETE
    )
    public void deleteItemTypes() throws InterruptedException {
        data.deleteItemTypes();
    }

    @ApiOperation(
        value = "Deletes a configuration item type.",
        notes = "")
    @RequestMapping(
          path = "/itemtype/{key}"
        , method = RequestMethod.DELETE
    )
    public ResponseEntity<Result> deleteItemType(@PathVariable("key") String key) {
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
        ) throws InterruptedException, IOException {
        return ResponseEntity.ok(data.createOrUpdateItemType(key, payload));
    }

    @ApiOperation(
        value = "Get a list of available configuration item types.",
        notes = "Only item types marked as custom can be deleted.")
    @RequestMapping(
          path = "/itemtype"
        , method = RequestMethod.GET
        , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ItemTypeList> getItemTypes() throws InterruptedException {
        List<ItemType> itemTypes = data.getItemTypes();
        return ResponseEntity.ok(new ItemTypeList(itemTypes));
    }

    @ApiOperation(
        value = "Get a configuration item based on the specified key.",
        notes = "Use this search to retrieve a specific configuration item when its natural key is known.")
    @RequestMapping(
          path = "/item/{key}"
        , method = RequestMethod.GET
        , produces = {"application/json", "application/x-yaml"}
    )
    public ResponseEntity<ItemData> getItem(@PathVariable("key") String key) {
        ItemData item = data.getItem(key);
        return ResponseEntity.ok(item);
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
            , @RequestParam(value = "from", required = false) String fromDate
            , @RequestParam(value = "to", required = false) String toDate
            , @RequestParam(value = "top", required = false, defaultValue = "100") Integer top
    ) {
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyyMMddHHmm");
        ZonedDateTime from = null;
        if (fromDate != null) {
            from = ZonedDateTime.of(LocalDateTime.parse(fromDate, formatter), ZoneId.systemDefault());
        }
        ZonedDateTime to = null;
        if (toDate != null) {
            to = ZonedDateTime.of(LocalDateTime.parse(toDate, formatter), ZoneId.systemDefault());
        }

        List<ItemData> items = null;
        if (itemTypeKey != null && tag == null && fromDate == null && toDate == null) {
            items = data.getItemsByType(itemTypeKey, top);
        } else if (itemTypeKey == null && tag != null && fromDate == null && toDate == null) {
            items = data.getItemsByTag(tag, top);
        } else if (itemTypeKey == null && tag == null && fromDate != null && toDate == null) {
            to = ZonedDateTime.now();
            items = data.getItemsByDate(from, to, top);
        } else if (itemTypeKey == null && tag == null && fromDate != null && toDate != null) {
            items = data.getItemsByDate(from, to, top);
        } else if (itemTypeKey != null && tag != null && fromDate == null && toDate == null) {
            items = data.getItemsByTypeAndTag(itemTypeKey, tag, top);
        } else if (itemTypeKey != null && tag == null && fromDate != null && toDate != null) {
            items = data.getItemsByTypeAndDate(itemTypeKey, from, to, top);
        } else if (itemTypeKey != null && tag != null && fromDate != null && toDate != null) {
            items = data.getItemsByTypeTagAndDate(itemTypeKey, tag, from, to, top);
        }  else {
            items = data.getAllByDateDesc(top);
        }
        return ResponseEntity.ok(new ItemList(items));
    }
}