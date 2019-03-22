package org.gatblau.onix.data;

import java.util.List;

public class ModelDataList extends Wrapper<ModelData> {
    public ModelDataList() {
    }

    public ModelDataList(List<ModelData> modelData){
        super(modelData);
    }
}
