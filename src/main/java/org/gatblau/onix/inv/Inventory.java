package org.gatblau.onix.inv;

import org.gatblau.onix.parser.Lexer;
import org.gatblau.onix.parser.LexerRule;
import org.gatblau.onix.parser.LexerToken;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Inventory {
    private static String HOST_GROUP = "HOST-GROUP";
    private static String GROUP_OF_HOST_GROUPS = "GROUP-OF-HOST-GROUPS";
    private static String ITEM = "ITEM";
    private static String HOST_VARS = "HOST-VARS";
    private static String COMMENT = "COMMENT";

    private NodeList nodes = new NodeList();

    public Inventory(String inventory) {
        List<LexerToken> tokens = parse(inventory);
        populate(tokens);
    }

    public List<Node> getNodes() {
        return nodes;
    }

    private List<LexerToken> parse(String source) {
        Lexer lexer = new Lexer();
        lexer.init(source);
        addRules(lexer);
        return lexer.getTokenStream();
    }

    private void populate(List<LexerToken> tokens) {
        Node currentParent = null;
        for (LexerToken token : tokens) {
            if (token.getType().equals(GROUP_OF_HOST_GROUPS)) {
                String name = token.getValue().substring(1, token.getValue().length() - ":children".length() - 1);
                Node newParent = new Node(name, Node.NodeType.PARENT_GROUP);
                nodes.add(newParent);
                currentParent = newParent;
            }
            else if (token.getType().equals(HOST_GROUP)) {
                String name = token.getValue().substring(1, token.getValue().length() - 1);
                // is this node already part of the tree
                Node hostGroup = nodes.find(name);
                // if it is then set it as the current parent
                if (hostGroup != null) {
                    currentParent = hostGroup;
                } else {
                    hostGroup = new Node(name, Node.NodeType.GROUP);
                    // if not, then it has to be added to the tree under the root
                    // and then made the current parent
                    nodes.add(hostGroup);
                    currentParent = hostGroup;
                }
            }
            else if (token.getType().equals(HOST_VARS)) {
            }
            else if (token.getType().equals(ITEM)) {
                String item = token.getValue();
                switch (currentParent.getType()){
                    case PARENT_GROUP:
                        currentParent.getChildren().add(new Node(item, Node.NodeType.GROUP));
                        break;
                    case GROUP:
                        String name = getItemName(item);
                        Map<String, String> vars = getItemVars(item);
                        currentParent.getChildren().add(new Node(name, Node.NodeType.HOST, vars));
                        break;
                }
            }
        }
    }

    private void addRules(Lexer lexer) {
        lexer.addRule(new LexerRule(COMMENT, "#.*$"));
        lexer.addRule(new LexerRule(GROUP_OF_HOST_GROUPS, "^\\w*\\[\\w*(?<item>.*):children\\w*\\]\\w*$"));
        lexer.addRule(new LexerRule(HOST_VARS, "^\\w*\\[\\w*(?<item>.*):vars\\w*\\]\\w*$"));
        lexer.addRule(new LexerRule(HOST_GROUP, "^\\w*\\[\\w*(?<item>[^:]*)\\w*\\]\\w*$"));
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
}
