package org.gatblau.onix.data;

import java.util.List;

public abstract class Wrapper<T> {
    private List<T> items;

    public Wrapper(){
    }

    public Wrapper(List<T> items) {
       this.items = items;
    }

    public List<T> getItems() {
        return items;
    }
}
