/*
Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org

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
package org.gatblau.onix.data;

import java.io.Serializable;
import java.sql.Types;
import java.util.ArrayList;
import java.util.List;

public class TabularData implements Serializable {
    private List<Column> columns = new ArrayList<>();
    private List<Row> rows = new ArrayList<>();

    public void addColumn(DataType type, String name) {
        getColumns().add(new Column(type, name));
    }

    public void addColumn(int type, String name) {
        getColumns().add(new Column(type, name));
    }

    public List<Column> getColumns() {
        return columns;
    }

    public void addRow(Row row) {
        getRows().add(row);
    }

    public List<Row> getRows() {
        return rows;
    }

    public enum DataType {
        String,
        Int,
        Date,
        Decimal,
        Boolean,
    }
    public static class Column {
        private DataType type;
        private String name;

        public Column(){
        }

        public Column(DataType type, String name){
            this.type = type;
            this.name = name;
        }

        public Column(int type, String name){
            this.type = inferType(type);
            this.name = name;
        }

        private DataType inferType(Integer type) {
            switch (type){
                case Types.BOOLEAN:
                    return DataType.Boolean;
                case Types.DATE:
                    return DataType.Date;
                case Types.DECIMAL:
                    return DataType.Decimal;
                case Types.INTEGER:
                    return DataType.Int;
                case Types.VARCHAR:
                    return DataType.String;
            }
            throw new RuntimeException(String.format("data type Sql.Types = '%s' not implemented", type.toString()));
        }

        public DataType getType() {
            return type;
        }

        public void setType(DataType type) {
            this.type = type;
        }

        public String getName() {
            return name;
        }

        public void setName(String name) {
            this.name = name;
        }
    }
    public static class Row extends ArrayList<Object> {
    }
}
