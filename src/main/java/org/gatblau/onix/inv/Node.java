package org.gatblau.onix.inv;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Node {
    private NodeType type;
    private String name;
    private NodeList children;
    private Map<String, String> vars = new HashMap<>();

    public Node(String name, NodeType type) {
        this(name, type, new HashMap<>());
    }

    public Node(String name, NodeType type, Map<String, String> vars) {
        this.type = type;
        this.name = name;
        this.children = new NodeList();
        this.vars = vars;
    }

    public NodeType getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public List<Node> getChildren() {
        return children;
    }

    public Map<String, String> getVars() {
        return vars;
    }

    public Node find(String name) {
        if (this.name.equals(name)) {
            return this;
        } else {
            Node found = children.find(name);
            if (found != null) {
                return found;
            }
        }
        return null;
    }

    public enum NodeType {
        HOST,
        GROUP,
        PARENT_GROUP,
        ROOT,
    }

    @Override
    public String toString() {
        return name;
    }
}
