package org.gatblau.onix;

import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.swagger.models.auth.In;
import org.gatblau.onix.data.ItemData;
import org.gatblau.onix.data.Wrapper;
import org.gatblau.onix.model.Item;
import org.json.simple.JSONObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.time.LocalDateTime;
import java.time.ZonedDateTime;
import java.util.Date;
import java.util.List;

@Api("ONIX CMDB Web API")
@RestController
public class WebAPI {

    @Autowired
    private Repository data;

    @ApiOperation(value = "Returns OK if the service is up and running.", notes = "Use it as a readiness probe for the service.", response = String.class)
    @RequestMapping(value = "/", method = RequestMethod.GET, produces = "application/json")
    public ResponseEntity<Info> index() {
        return ResponseEntity.ok(new Info("Onix Configuration Management Database Service.", "1.0"));
    }

    @ApiOperation(value = "Creates new item or updates an existing item based on the specified key.", notes = "")
    @RequestMapping(
            path = "/item/{key}/", method = RequestMethod.PUT,
            consumes = {"application/json" },
            produces = {"application/json" })
    public ResponseEntity<Result> createOrUpdateItem(
            @PathVariable("key") String key,
            @RequestBody JSONObject payload) throws InterruptedException, IOException {
        String action = data.createOrUpdateItem(key, payload);
        return ResponseEntity.ok(new Result(action));
    }

    @ApiOperation(value = "Creates new link or updates an existing link based on two existing item keys.", notes = "")
    @RequestMapping(
            path = "/link/{fromItemKey}/{toItemKey}/", method = RequestMethod.PUT,
            consumes = {"application/json" },
            produces = {"application/json" })
    public ResponseEntity<Result> createOrUpdateLink(
            @PathVariable("fromItemKey") String fromItemKey,
            @PathVariable("toItemKey") String toItemKey,
            @RequestBody JSONObject payload) throws InterruptedException, IOException {
        String action = data.createOrUpdateLink(fromItemKey, toItemKey, payload);
        return ResponseEntity.ok(new Result(action));
    }

    @RequestMapping(path = "/clear/", method = RequestMethod.DELETE)
    public void clear() throws InterruptedException {
        data.clear();
    }

    @RequestMapping(path = "/link/{fromItemKey}/{toItemKey}/", method = RequestMethod.DELETE)
    public void deleteLink(
            @PathVariable("fromItemKey") String fromItemKey,
            @PathVariable("toItemKey") String toItemKey
    ) throws InterruptedException {
        data.deleteLink(fromItemKey, toItemKey);
    }

    @RequestMapping(path = "/itemtype/", method = RequestMethod.DELETE)
    public void deleteItemTypes() throws InterruptedException {
        data.deleteItemTypes();
    }

    @ApiOperation(value = "Get a configuration item based on the specified key.", notes = "")
    @RequestMapping(
        path = "/item/{key}/",
        method = RequestMethod.GET,
        consumes = {"application/json" },
        produces = {"application/json" }
    )
    public ResponseEntity<ItemData> getItem(@PathVariable("key") String key) {
        ItemData item = data.getItem(key);
        return ResponseEntity.ok(item);
    }

    @ApiOperation(value = "Get a configuration item based on the specified key.", notes = "")
    @RequestMapping(
        path = "/item/search"
        , method = RequestMethod.GET
        , consumes = {"application/json" }
        , produces = {"application/json", "application/x-yaml" }
        , headers = { "Accept: application/json" }
    )
    public ResponseEntity<Wrapper> getItems(
              @RequestParam(value = "typeId", required = false) Integer typeId
            , @RequestParam(value = "tag", required = false) String tag
            , @RequestParam(value = "from", required = false) ZonedDateTime from
            , @RequestParam(value = "to", required = false) ZonedDateTime to
            , @RequestParam(value = "top", required = false, defaultValue = "100") Integer top) {
        List<ItemData> items = null;
        if (typeId != null && tag == null && from == null && to == null) {
          items = data.getItemsByType(typeId, top);
        } else if (typeId == null && tag != null && from == null && to == null) {
            items = data.getItemsByTag(tag, top);
        } else if (typeId == null && tag == null && from != null && to == null) {
            to = ZonedDateTime.now();
            items = data.getItemsByDate(from, to, top);
        } else if (typeId == null && tag == null && from != null && to != null) {
            items = data.getItemsByDate(from, to, top);
        } else if (typeId != null && tag != null && from == null && to == null) {
            items = data.getItemsByTypeAndTag(typeId, tag, top);
        } else if (typeId != null && tag == null && from != null && to != null) {
            items = data.getItemsByTypeAndDate(typeId, from, to, top);
        } else if (typeId == null && tag != null && from == null && to == null) {
            items = data.getItemsByTypeTagAndDate(typeId, tag, from, to, top);
        }  else {
            items = data.getAllByDateDesc(top);
        }
        return ResponseEntity.ok(new Wrapper(items));
    }
}