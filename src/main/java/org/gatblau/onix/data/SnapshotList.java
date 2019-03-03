package org.gatblau.onix.data;

import java.util.List;

public class SnapshotList extends Wrapper<SnapshotData> {
    public SnapshotList() {
    }

    public SnapshotList(List<SnapshotData> snapshot){
        super(snapshot);
    }
}
