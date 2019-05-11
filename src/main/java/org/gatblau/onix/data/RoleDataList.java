package org.gatblau.onix.data;

import java.util.List;

public class RoleDataList extends Wrapper<RoleData> {
    public RoleDataList() {
    }

    public RoleDataList(List<RoleData> roleData){
        super(roleData);
    }
}
