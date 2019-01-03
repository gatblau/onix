package org.gatblau.onix.data;

import java.util.List;

public class LinkTypeList extends Wrapper<LinkTypeData> {
    public LinkTypeList() {
    }

    public LinkTypeList(List<LinkTypeData> linkData){
        super(linkData);
    }
}
