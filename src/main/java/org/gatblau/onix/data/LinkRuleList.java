package org.gatblau.onix.data;

import java.io.Serializable;
import java.util.List;

public class LinkRuleList extends Wrapper<LinkRuleData> {
    public LinkRuleList() {
    }

    public LinkRuleList(List<LinkRuleData> linkData){
        super(linkData);
    }
}

