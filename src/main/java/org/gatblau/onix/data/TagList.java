package org.gatblau.onix.data;

import java.util.List;

public class TagList extends Wrapper<TagData> {
    public TagList() {
    }

    public TagList(List<TagData> tag){
        super(tag);
    }
}
