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
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.postgresql.jdbc.PgArray;
import org.postgresql.util.PGobject;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.io.StringReader;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.*;

@Component
public class Lib implements InitializingBean {
    private DateFormat dateFormat = new SimpleDateFormat("dd-MM-yyyy HH:mm:ss Z");
    private JSONParser jsonParser = new JSONParser();

    @Override
    public void afterPropertiesSet() throws Exception {
        dateFormat.setTimeZone(TimeZone.getTimeZone("GMT"));
    }

    public JSONObject toJSON(Object value) throws ParseException, IOException {
        JSONObject json = null;
        if (value instanceof PGobject) {
            PGobject pgObj = (PGobject)value;
            String strValue = pgObj.getValue();
            json = (JSONObject)jsonParser.parse(strValue);
        }
        else if (value instanceof LinkedHashMap || value instanceof HashMap) {
            json = new JSONObject((HashMap)value);
        }
        else if (value instanceof String) {
            JSONParser parser = new JSONParser();
            StringReader reader = new StringReader((String)value);
            return (JSONObject)parser.parse(reader);
        }
        else {
            // the object is not a list, then create an empty JSON object
            json = new JSONObject();
            System.out.println(String.format("WARNING: incorrect map format found on item '%s', discarding item content.", value));
        }
        return json;
    }

    public String toJSONString(Object value) throws ParseException, IOException {
        return toJSON(value).toJSONString();
    }

    public String toArrayString(Object value) {
        String arrayString = null;
        if (value == null) {
            arrayString = toArrayString(new ArrayList<>());
        }
        try {
            arrayString = toArrayString((List<String>)value);
        }
        catch (Exception ex) {
            System.out.println(String.format("WARNING: incorrect array format found on item '%'.", value));
            ex.printStackTrace();
            List<String> list = new ArrayList<>();
            String tagStr = (String) value;
            String[] strs = tagStr.split("[ ]|[|]|[:]|[,]"); // valid tag separators are blank space, pipe, colon or comma
            for (String s : strs) {
                list.add(s);
            }
            arrayString = toArrayString(list);
        }
        return arrayString;
    }

    public String toArrayString(List<String> list) {
        if (list == null) {
            return "{}";
        }
        StringBuilder sb = new StringBuilder();
        sb.append("{");
        for (int i = 0; i < list.size(); i++){
            sb.append("\"").append(list.get(i)).append("\"");
            if (i < list.size() - 1) {
                sb.append(",");
            }
        }
        sb.append("}");
        return sb.toString();
    }

    public List<String> toList(Object value) throws SQLException {
        if (value instanceof PgArray) {
            PgArray pgArray = (PgArray) value;
            String[] array = (String[])pgArray.getArray();
            return Arrays.asList(array);
        }
        throw new RuntimeException("Conversion not implemented.");
    }

    public ItemData toItemData(ResultSet set) throws SQLException, ParseException, IOException {
        Date updated = set.getDate("updated");
        ItemData item = new ItemData();
        item.setKey(set.getString("key"));
        item.setName(set.getString("name"));
        item.setDescription(set.getString("description"));
        item.setStatus(set.getShort("status"));
        item.setType(set.getString("item_type_key"));
        item.setCreated(dateFormat.format(set.getDate("created")));
        item.setUpdated((updated != null) ? dateFormat.format(updated) : null);
        item.setMeta(toJSON(set.getObject("meta")));
        item.setTag(toList(set.getObject("tag")));
        item.setVersion(set.getInt("version"));
        item.setAttribute(toJSON(set.getObject("attribute")));
        return item;
    }

    public LinkData toLinkData(ResultSet set) throws SQLException, ParseException, IOException {
        Date updated = set.getDate("updated");
        LinkData link = new LinkData();
        link.setKey(set.getString("key"));
        link.setType(set.getString("link_type_key"));
        link.setDescription(set.getString("description"));
        link.setEndItemKey(set.getString("end_item_key"));
        link.setStartItemKey(set.getString("start_item_key"));
        link.setMeta(toJSON(set.getObject("meta")));
        link.setTag(toList(set.getObject("tag")));
        link.setAttribute(toJSON(set.getObject("attribute")));
        link.setUpdated((updated != null) ? dateFormat.format(updated) : null);
        link.setVersion(set.getInt("version"));
        link.setChangedBy(set.getString("changed_by"));
        return link;
    }

    public ItemTypeData toItemTypeData(ResultSet set) throws SQLException, ParseException, IOException {
        Date updated = set.getDate("updated");

        ItemTypeData itemType = new ItemTypeData();
        itemType.setKey(set.getString("key"));
        itemType.setName(set.getString("name"));
        itemType.setDescription(set.getString("description"));
        itemType.setCreated(dateFormat.format(set.getDate("created")));
        itemType.setUpdated((updated != null) ? dateFormat.format(updated) : null);
        itemType.setVersion(set.getInt("version"));
        itemType.setAttribute(toJSON(set.getObject("attr_valid")));
        itemType.setFilter(toJSON(set.getObject("filter")));
        itemType.setMetaSchema(toJSON(set.getObject("meta_schema")));
        itemType.setModelKey(set.getString("model_key"));
        return itemType;
    }

    public String toHStoreString(Map<String, String> map) {
        String result = null;
        if (map == null) {
            result = null;
        }
        else {
            StringBuilder sb = new StringBuilder();
            int count = 0;
            for (Map.Entry<String, String> entry : map.entrySet()) {
                sb.append(entry.getKey()).append("=>").append(entry.getValue());
                if (count < map.entrySet().size() - 1) {
                    sb.append(",");
                }
                count++;
            }
            result = sb.toString();
        }
        return result;
    }

    public LinkTypeData toLinkTypeData(ResultSet set) throws SQLException, ParseException, IOException {
        Date updated = set.getDate("updated");

        LinkTypeData linkType = new LinkTypeData();
        linkType.setKey(set.getString("key"));
        linkType.setName(set.getString("name"));
        linkType.setDescription(set.getString("description"));
        linkType.setCreated(dateFormat.format(set.getDate("created")));
        linkType.setUpdated((updated != null) ? dateFormat.format(updated) : null);
        linkType.setVersion(set.getInt("version"));
        linkType.setAttribute(toJSON(set.getObject("attr_valid")));
        linkType.setMetaSchema(toJSON(set.getObject("meta_schema")));
        linkType.setModelKey(set.getString("model_key"));
        linkType.setChangedBy(set.getString("changed_by"));
        return linkType;
    }

    public LinkRuleData toLinkRuleData(ResultSet set) throws SQLException {
        Date updated = set.getDate("updated");

        LinkRuleData linkRule = new LinkRuleData();
        linkRule.setKey(set.getString("key"));
        linkRule.setName(set.getString("name"));
        linkRule.setDescription(set.getString("description"));
        linkRule.setLinkTypeKey(set.getString("link_type_key"));
        linkRule.setStartItemTypeKey(set.getString("start_item_type_key"));
        linkRule.setEndItemTypeKey(set.getString("end_item_type_key"));
        linkRule.setCreated(dateFormat.format(set.getDate("created")));
        linkRule.setUpdated((updated != null) ? dateFormat.format(updated) : null);
        linkRule.setVersion(set.getInt("version"));
        linkRule.setChangedBy(set.getString("changed_by"));
        return linkRule;
    }

    public TagData toTagData(ResultSet set) throws SQLException {
        Date updated = set.getDate("updated");

        TagData tag = new TagData();
        tag.setLabel(set.getString("label"));
        tag.setName(set.getString("name"));
        tag.setDescription(set.getString("description"));
        tag.setRootItemKey(set.getString("root_item_key"));
        tag.setCreated(dateFormat.format(set.getDate("created")));
        tag.setUpdated((updated != null) ? dateFormat.format(updated) : null);
        tag.setVersion(set.getInt("version"));
        tag.setChangedBy(set.getString("changed_by"));
        return tag;
    }

    public ModelData toModelData(ResultSet set) throws SQLException {
        Date updated = set.getDate("updated");

        ModelData model = new ModelData();
        model.setKey(set.getString("key"));
        model.setName(set.getString("name"));
        model.setDescription(set.getString("description"));
        model.setCreated(dateFormat.format(set.getDate("created")));
        model.setUpdated((updated != null) ? dateFormat.format(updated) : null);
        model.setVersion(set.getInt("version"));
        model.setChangedBy(set.getString("changed_by"));
        return model;
    }
}
