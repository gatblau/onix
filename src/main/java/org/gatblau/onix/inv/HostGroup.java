package org.gatblau.onix.inv;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class HostGroup {
    private String name;
    private List<HostGroup> groups = new ArrayList<>();
    private List<Host> hosts = new ArrayList<>();
    private Map<String, String> vars = new HashMap<>();

    public HostGroup(String name) {
        this.name = name;
    }

    public Map<String, String> getVars() {
        return vars;
    }

    public List<HostGroup> getGroups() {
        return groups;
    }

    public List<Host> getHosts() {
        return hosts;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
