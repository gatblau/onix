package org.gatblau.onix.inv;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import org.gatblau.onix.data.ItemData;
import org.gatblau.onix.data.ItemTreeData;
import org.gatblau.onix.data.LinkData;
import org.gatblau.onix.parser.Lexer;
import org.gatblau.onix.parser.LexerRule;
import org.gatblau.onix.parser.LexerToken;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;

import java.io.StringReader;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Inventory {
    private static final String EMPTY_LINE = "EMPTY_LINE";
    public static final String ANSIBLE_HOST = "ANSIBLE_HOST";
    public static final String ANSIBLE_HOST_GROUP = "ANSIBLE_HOST_GROUP";
    public static final String ANSIBLE_HOST_GROUP_SET = "ANSIBLE_HOST_GROUP_SET";
    private static final String ITEM = "ITEM";
    private static final String HOST_VARS = "HOST_VARS";
    private static final String COMMENT = "COMMENT";
    private String hostVars = "";
    private boolean readingHostVars;

    private NodeList nodes = new NodeList();

    /*
        builds an inventory from a flat inventory file
     */
    public Inventory(String inventory) {
        List<LexerToken> tokens = parse(inventory);
        populate(tokens);
    }

    /*
        builds an inventory from item tree data
     */
    public Inventory(ItemTreeData tree) {
        for (ItemData item : tree.getItems()) {
            processItem(tree, item, nodes);
        }
    }

    private void processItem(ItemTreeData tree, ItemData item, NodeList nodes) {
        String name = item.getName();
        Node node = this.nodes.find(name);
        if (node == null) {
            node = new Node(item.getName(), getType(item.getType()), (JSONObject)item.getMeta().get("hostvars"));
            nodes.add(node);
            List<ItemData> children = getChildren(tree, item);
            for (ItemData childData : children) {
                processItem(tree, childData, node.getChildren());
            }
        }
    }

    private List<ItemData> getChildren(ItemTreeData tree, ItemData item) {
        List<ItemData> items = new ArrayList<>();
        for (LinkData link : tree.getLinks()) {
            if (link.getStartItemKey().equals(item.getKey())) {
                String childItemKey = link.getEndItemKey();
                for (ItemData child : tree.getItems()){
                    if (child.getKey().equals(childItemKey)) {
                        items.add(child);
                    }
                }
            }
        }
        return items;
    }

    private Node.NodeType getType(String itemType) {
        switch (itemType){
            case ANSIBLE_HOST_GROUP:
                return Node.NodeType.ANSIBLE_HOST_GROUP;
            case ANSIBLE_HOST_GROUP_SET:
                return Node.NodeType.ANSIBLE_HOST_GROUP_SET;
            case ANSIBLE_HOST:
                return Node.NodeType.ANSIBLE_HOST;
        }
        throw new RuntimeException(String.format("Mapping not found for item type '%s'", itemType));
    }

    public List<Node> getNodes() {
        return nodes;
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder();
        WriteParentGroups(builder);
        WriteHostGroups(builder);
        WriteGroupsVars(builder);
        return builder.toString();
    }

    private void WriteGroupsVars(StringBuilder builder) {
        for (Node node : nodes) {
            switch (node.getType()) {
                case ANSIBLE_HOST_GROUP_SET: {
                    if (node.getVars().size() > 0) {
                        WriteGroupVars(builder, node);
                    }
                    for (Node hostGroup : node.getChildren()) {
                        WriteGroupVars(builder, hostGroup);
                    }
                    break;
                }
                case ANSIBLE_HOST_GROUP: {
                    WriteGroupVars(builder, node);
                    break;
                }
            }
        }
    }

    private void WriteGroupVars(StringBuilder builder, Node hostGroup) {
        if (hostGroup.getVars().size() > 0) {
            builder.append("[").append(hostGroup.getName()).append(":vars]").append(System.lineSeparator());
            builder.append(new JSONObject(hostGroup.getVars()).toJSONString()).append(System.lineSeparator());
        }
    }

    private void WriteParentGroups(StringBuilder builder) {
        for (Node node : nodes) {
            switch (node.getType()) {
                case ANSIBLE_HOST_GROUP_SET: {
                    builder.append("[").append(node.getName()).append(":children]").append(System.lineSeparator());
                    for (int i = 0; i < node.getChildren().size(); i++) {
                        Node hostGroup = node.getChildren().get(i);
                        builder.append(hostGroup.getName()).append(System.lineSeparator());
                        if (i == hostGroup.getChildren().size() - 1) {
                            builder.append(System.lineSeparator());
                        }
                    }
                    break;
                }
            }
        }
    }

    private void WriteHostGroups(StringBuilder builder) {
        for (Node node : nodes) {
            switch (node.getType()) {
                case ANSIBLE_HOST_GROUP_SET: {
                    for (Node hostGroup : node.getChildren()) {
                        WriteHostGroup(builder, hostGroup);
                    }
                    break;
                }
                case ANSIBLE_HOST_GROUP: {
                    WriteHostGroup(builder, node);
                    break;
                }
            }
        }
    }

    private void WriteHostGroup(StringBuilder builder, Node hostGroup) {
        builder.append("[").append(hostGroup.getName()).append("]").append(System.lineSeparator());
        for (int i = 0; i < hostGroup.getChildren().size(); i++) {
            Node host = hostGroup.getChildren().get(i);
            builder.append(host.getName());
            if (host.getVars().size() > 0) {
                builder.append(" ").append(new JSONObject(host.getVars()).toJSONString());
            }
            builder.append(System.lineSeparator());
            if (i == hostGroup.getChildren().size() - 1) {
                builder.append(System.lineSeparator());
            }
        }
    }

    private List<LexerToken> parse(String source) {
        Lexer lexer = new Lexer();
        // pass the text to tokenise to the lexer
        lexer.init(source);
        // skip tokens for empty lines and comments
        lexer.getSkippedTokens().add(EMPTY_LINE);
        lexer.getSkippedTokens().add(COMMENT);
        // add the tokenisation rules
        addRules(lexer);
        // return a list of tokens
        return lexer.getTokenStream();
    }

    private void populate(List<LexerToken> tokens) {
        Node currentParent = null;
        for (LexerToken token : tokens) {
            if (token.getType().equals(ANSIBLE_HOST_GROUP_SET)) {
                String name = token.getValue().substring(1, token.getValue().length() - ":children".length() - 1);
                Node newParent = new Node(name, Node.NodeType.ANSIBLE_HOST_GROUP_SET);
                nodes.add(newParent);
                currentParent = newParent;
            }
            else if (token.getType().equals(ANSIBLE_HOST_GROUP)) {
                String name = token.getValue().substring(1, token.getValue().length() - 1);
                // is this node already part of the tree
                Node hostGroup = nodes.find(name);
                // if it is then set it as the current parent
                if (hostGroup != null) {
                    currentParent = hostGroup;
                } else {
                    hostGroup = new Node(name, Node.NodeType.ANSIBLE_HOST_GROUP);
                    // if not, then it has to be added to the tree under the root
                    // and then made the current parent
                    nodes.add(hostGroup);
                    currentParent = hostGroup;
                }
            }
            else if (token.getType().equals(HOST_VARS)) {
                // get the host name for the vars
                String hostName = token.getValue().substring(1, token.getValue().length() - ":vars]".length());
                // find the node representing the host name
                Node host = nodes.find(hostName);
                if (host == null) {
                    throw new RuntimeException(
                        String.format("Failed to find host '%s' in inventory. " +
                            "Check it appears before host:vars statement in the inventory file.", hostName));
                }
                // makes the host the current parent
                currentParent = host;
                readingHostVars = true;
            }
            else if (token.getType().equals(ITEM)) {
                String item = token.getValue();
                switch (currentParent.getType()){
                    case ANSIBLE_HOST_GROUP_SET:
                        if (readingHostVars) {
                            // aggregates host vars
                            hostVars += item + System.lineSeparator();
                            // if it is the end of the vars section
                            if (token.getNextToken().getType().contains("GROUP")) {
                                // resets the accumulator flag
                                readingHostVars = false;

                                // gets vars in json format
                                JSONObject vars = new JSONObject();
                                try {
                                    vars = convertYamlToJson(hostVars);
                                } catch (Exception e) {
                                    e.printStackTrace();
                                }

                                // add vars to the current parent
                                currentParent.getVars().putAll(vars);
                                hostVars = "";
                            }
                        } else {
                            // add a new host node
                            currentParent.getChildren().add(new Node(item, Node.NodeType.ANSIBLE_HOST_GROUP));
                        }
                        break;
                    case ANSIBLE_HOST_GROUP:
                        String name = getItemName(item);
                        currentParent.getChildren().add(new Node(name, Node.NodeType.ANSIBLE_HOST, getItemVars(item)));
                        break;
                }
            }
        }
    }

    private void addRules(Lexer lexer) {
        lexer.addRule(new LexerRule(EMPTY_LINE, "^\\s*$"));
        lexer.addRule(new LexerRule(COMMENT, "#.*$"));
        lexer.addRule(new LexerRule(ANSIBLE_HOST_GROUP_SET, "^\\w*\\[\\w*(?<item>.*):children\\w*\\]\\w*$"));
        lexer.addRule(new LexerRule(HOST_VARS, "^\\w*\\[\\w*(?<item>.*):vars\\w*\\]\\w*$"));
        lexer.addRule(new LexerRule(ANSIBLE_HOST_GROUP, "^\\w*\\[\\w*(?<item>[^:]*)\\w*\\]\\w*$"));
        lexer.addRule(new LexerRule(ITEM, "^(?!\\[).+$"));
    }

    private String getItemName(String item) {
        int i = item.indexOf(" ");
        return (i > -1) ? item.substring(0, i) : item;
    }

    private Map<String, String> getItemVars(String item) {
        Map<String, String> vars = new HashMap<>();
        List<String> items = new ArrayList<>();
        String[] parts = item.split("=");
        for (String part : parts) {
            String[] p;
            if (part.startsWith("\"")) {
                p = new String[]{part};
            } else {
                int i = part.indexOf(" ");
                if (i > -1) {
                    p = new String[]{part.substring(0, i), part.substring(i, part.length())};
                } else {
                    p = new String[]{part};
                }
            }
            for (int j = 0; j < p.length; j++) {
                items.add(p[j]);
            }
        }
        for (int i = 1; i < items.size() - 1; i+=2) {
            vars.put(items.get(i), items.get(i+1));
        }
        return vars;
    }

    private JSONObject convertYamlToJson(String yaml) {
        String result = null;
        try {
            ObjectMapper yamlReader = new ObjectMapper(new YAMLFactory());
            Object obj = yamlReader.readValue(yaml, Object.class);
            ObjectMapper jsonWriter = new ObjectMapper();
            result = jsonWriter.writeValueAsString(obj);
        } catch (Exception e) {
            e.printStackTrace();
        }
        JSONParser parser = new JSONParser();
        JSONObject json = null;
        try {
            json = (JSONObject)parser.parse(new StringReader(result));
        } catch (Exception e) {
            e.printStackTrace();
        }
        return json;
    }
}
