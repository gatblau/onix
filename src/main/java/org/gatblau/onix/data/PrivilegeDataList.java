package org.gatblau.onix.data;

import java.util.List;

public class PrivilegeDataList extends Wrapper<PrivilegeData> {
    public PrivilegeDataList() {
    }

    public PrivilegeDataList(List<PrivilegeData> privilegeData){
        super(privilegeData);
    }
}
