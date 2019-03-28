package org.gatblau.onix.data;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

public class TypeGraphData implements Serializable {
    private static final long serialVersionUID = 1L;
    private List<ItemTypeData> itemTypes = new ArrayList<>();
    private List<LinkTypeData> linkTypes = new ArrayList<>();
    private List<LinkRuleData> linkRules = new ArrayList<>();

    public TypeGraphData() {
    }

    public void setLinkTypes(List<LinkTypeData> linkTypes) {
        this.linkTypes = linkTypes;
    }

    public List<LinkTypeData> getLinkTypes() {
        return linkTypes;
    }

    public List<ItemTypeData> getItemTypes() {
        return itemTypes;
    }

    public void setItemTypes(List<ItemTypeData> itemTypes) {
        this.itemTypes = itemTypes;
    }

    public List<LinkRuleData> getLinkRules() {
        return linkRules;
    }

    public void setLinkRules(List<LinkRuleData> linkRules) {
        this.linkRules = linkRules;
    }
}
