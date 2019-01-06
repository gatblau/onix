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

import org.gatblau.onix.data.ItemData;
import org.gatblau.onix.data.ItemTypeData;
import org.gatblau.onix.data.LinkData;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.postgresql.jdbc.PgArray;
import org.postgresql.util.PGobject;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.stereotype.Component;

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

    public JSONObject toJSON(Object value) throws ParseException {
        JSONObject json = null;
        if (value instanceof PGobject) {
            PGobject pgObj = (PGobject)value;
            String strValue = pgObj.getValue();
            json = (JSONObject)jsonParser.parse(strValue);
        }
        else if (value instanceof LinkedHashMap || value instanceof HashMap) {
            json = new JSONObject((HashMap)value);
        }
        else {
            // the object is not a list, then create an empty JSON object
            json = new JSONObject();
            System.out.println(String.format("WARNING: incorrect map format found on item '%s', discarding item content.", value));
        }
        return json;
    }

    public String toJSONString(Object value) throws ParseException {
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

    public ItemData toItemData(ResultSet set) throws SQLException, ParseException {
        Date updated = set.getDate("updated");
        ItemData item = new ItemData();
        item.setKey(set.getString("key"));
        item.setName(set.getString("name"));
        item.setDescription(set.getString("description"));
        item.setStatus(set.getShort("status"));
        item.setItemType(set.getString("item_type_key"));
        item.setCreated(dateFormat.format(set.getDate("created")));
        item.setUpdated((updated != null) ? dateFormat.format(updated) : null);
        item.setMeta(toJSON(set.getObject("meta")));
        item.setTag(toList(set.getObject("tag")));
        item.setVersion(set.getInt("version"));
        item.setAttribute(toJSON(set.getObject("attribute")));
        return item;
    }

    public LinkData toLinkData(ResultSet set) throws SQLException, ParseException {
        LinkData link = new LinkData();
        link.setKey(set.getString("key"));
        link.setDescription(set.getString("description"));
        link.setEndItemKey(set.getString("end_item_key"));
        link.setStartItemKey(set.getString("start_item_key"));
        link.setMeta(toJSON(set.getObject("meta")));
        link.setTag(toList(set.getObject("tag")));
        link.setAttribute(toJSON(set.getObject("attribute")));
        return link;
    }

    public ItemTypeData toItemTypeData(ResultSet set) throws SQLException, ParseException {
        Date updated = set.getDate("updated");
        ItemTypeData itemType = new ItemTypeData();
        itemType.setKey(set.getString("key"));
        itemType.setName(set.getString("name"));
        itemType.setDescription(set.getString("description"));
        itemType.setCreated(dateFormat.format(set.getDate("created")));
        itemType.setUpdated((updated != null) ? dateFormat.format(updated) : null);
        itemType.setVersion(set.getInt("version"));
        itemType.setAttribute(toJSON(set.getObject("attr_valid")));
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
}
