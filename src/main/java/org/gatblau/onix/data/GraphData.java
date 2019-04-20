package org.gatblau.onix.data;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

public class GraphData implements Serializable {
    private static final long serialVersionUID = 1L;
    private List<ItemTypeData> itemTypes = new ArrayList<>();
    private List<ItemData> items = new ArrayList<>();
    private List<LinkTypeData> linkTypes = new ArrayList<>();
    private List<LinkData> links = new ArrayList<>();
    private List<ModelData> models = new ArrayList<>();
    private List<LinkRuleData> linkRules = new ArrayList<>();

    public GraphData() {
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

    public List<ModelData> getModels() {
        return models;
    }

    public void setModels(List<ModelData> models) {
        this.models = models;
    }

    public List<ItemTypeData> getItemTypes() {
        return itemTypes;
    }

    public void setItemTypes(List<ItemTypeData> itemTypes) {
        this.itemTypes = itemTypes;
    }

    public List<LinkTypeData> getLinkTypes() {
        return linkTypes;
    }

    public void setLinkTypes(List<LinkTypeData> linkTypes) {
        this.linkTypes = linkTypes;
    }

    public List<LinkRuleData> getLinkRules() {
        return linkRules;
    }

    public void setLinkRules(List<LinkRuleData> linkRules) {
        this.linkRules = linkRules;
    }
}
