package org.gatblau.onix.data;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

public class ItemTreeData implements Serializable {
    private static final long serialVersionUID = 1L;
    private List<ItemData> items = new ArrayList<>();
    private List<LinkData> links = new ArrayList<>();

    public ItemTreeData() {
    }

    public void setLinks(List<LinkData> links) {
        this.links = links;
    }

    public List<LinkData> getLinks() {
        return links;
    }

    public List<ItemData> getItems() {
        return items;
    }

    public void setItems(List<ItemData> items) {
        this.items = items;
    }
}
