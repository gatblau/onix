package org.gatblau.onix.data;

import java.util.List;

public class Wrapper {
    private List<ItemData> items;

    public Wrapper(){
    }

    public Wrapper(List<ItemData> items) {
       this.items = items;
    }

    public List<ItemData> getItems() {
        return items;
    }
}
