package org.gatblau.onix.inv;

import java.util.ArrayList;

public class NodeList extends ArrayList<Node> {

    public NodeList() {
    }

    public Node find(String name) {
        for (Node n : this) {
            Node found = n.find(name);
            if (found != null) {
                return found;
            }
        }
        return null;
    }

    @Override
    public boolean add(Node node) {
        return super.add(node);
    }
}
