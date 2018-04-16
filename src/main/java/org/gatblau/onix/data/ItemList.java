package org.gatblau.onix.data;

import java.util.List;

public class ItemList extends Wrapper<ItemData> {
    public ItemList() {
    }

    public ItemList(List<ItemData> item){
        super(item);
    }
}
