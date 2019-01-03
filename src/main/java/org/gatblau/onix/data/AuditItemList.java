package org.gatblau.onix.data;

import java.util.List;

public class AuditItemList extends Wrapper<AuditItemData> {
    public AuditItemList() {
    }

    public AuditItemList(List<AuditItemData> item){
        super(item);
    }
}

