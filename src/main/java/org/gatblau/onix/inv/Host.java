package org.gatblau.onix.inv;

public class Host {
    private String name;
    private HostGroup group;

    public Host(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public HostGroup getGroup() {
        return group;
    }

    public void setGroup(HostGroup group) {
        this.group = group;
    }
}
