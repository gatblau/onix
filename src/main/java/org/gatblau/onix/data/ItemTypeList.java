package org.gatblau.onix.data;

import org.gatblau.onix.model.ItemType;

import java.util.List;

public class ItemTypeList extends Wrapper<ItemType> {
    public ItemTypeList() {
    }

    public ItemTypeList(List<ItemType> itemTypes){
        super(itemTypes);
    }
}
