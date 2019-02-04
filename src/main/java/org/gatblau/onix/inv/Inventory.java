package org.gatblau.onix.inv;

import java.util.ArrayList;
import java.util.List;

public class Inventory {
    private ParentType parent;
    private HostGroup lastParentGroup;
    private List<HostGroup> groups = new ArrayList<>();

    public Inventory(String inventory) {
        parse(inventory);
    }

    public List<HostGroup> getGroups() {
        return groups;
    }

    private void parse(String inventory) {
       String[] lines = inventory.split("\n");
       for (String line : lines) {
           // omits commented out lines
           if (line.trim().length() == 0 || line.trim().startsWith("#")) {
               continue;
           }
           parseLine(line);
       }
    }

    private void parseLine(String line) {
        line = line.trim();
        if (line.length() == 0) return;
        // is a group
        if (line.startsWith("[")) {
            // is a parent group
            if (line.endsWith(":children]")) {
                parseParentGroup(line);
            }
            // is a child group
            else if (line.endsWith("]")) {
                parseGroup(line);
            }
        }
        // is a host or group
        else {
            if (parent == ParentType.ParentGroup) {
                parseChildGroup(line);
            }
            else if (parent == ParentType.Group) {
                parseHost(line);
            }
        }
    }

    private void parseGroup(String line) {
        String name = line.substring(1, line.length() - "]".length());
        getGroups().add(new HostGroup(name));
        parent = ParentType.Group;
    }

    private void parseParentGroup(String line) {
        String name = line.substring(1, line.length() - ":children]".length());
        lastParentGroup = new HostGroup(name);
        getGroups().add(lastParentGroup);
        parent = ParentType.ParentGroup;
    }

    private void parseChildGroup(String line) {
        String name = line;
        lastParentGroup.getGroups().add(new HostGroup(name));
    }

    private void parseHost(String line) {
        String[] parts = line.split(" ");
        HostGroup lastGroup = getGroups().get(getGroups().size()-1);
        lastGroup.getHosts().add(new Host(parts[0]));
    }
}
