package org.gatblau.onix.data;

import java.util.List;

public class PartitionDataList extends Wrapper<PartitionData> {
    public PartitionDataList() {
    }

    public PartitionDataList(List<PartitionData> partitionData){
        super(partitionData);
    }
}
