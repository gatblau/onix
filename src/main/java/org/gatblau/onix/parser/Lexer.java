package org.gatblau.onix.parser;

import java.io.BufferedReader;
import java.io.StringReader;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;

public class Lexer {
    private String source;
    private List<LexerRule> rules = new ArrayList<>();
    private List<String> skippedTokens = new ArrayList<>();

    public Lexer() {
    }

    public void init(String source) {
        this.source = source;
    }

    public void addRule(LexerRule rule) {
        rules.add(rule);
    }

    public void enableRule(String tokenType) {
        setRule(tokenType, true);
    }

    public void disableRule(String tokenType) {
        setRule(tokenType, false);
    }

    public void setRule(String tokenType, boolean state) {
        for (LexerRule rule : rules) {
            if (rule.getTokenType().equals(tokenType)) {
                rule.setEnabled(state);
                return;
            }
        }
    }

    public void enableRulesByPrefix(String tokenTypePrefix) {
        setRulesByPrefix(tokenTypePrefix, true);
    }

    public void disableRulesByPrefix(String tokenTypePrefix) {
        setRulesByPrefix(tokenTypePrefix, false);
    }

    public void setRulesByPrefix(String tokenTypePrefix, boolean state) {
        for (LexerRule rule : rules) {
            if (rule.getTokenType().startsWith(tokenTypePrefix)) {
                rule.setEnabled(state);
                return;
            }
        }
    }

    public List<LexerToken> getTokenStream() {
        List<LexerToken> tokens = new ArrayList<>();
        try {
            BufferedReader br = new BufferedReader(new StringReader(source));
            String line;
            while ((line = br.readLine()) != null) {
                Matcher match = null;
                for (LexerRule rule : rules) {
                    if (!rule.isEnabled()) continue;
                    match = rule.match(line);
                    if (match == null) continue;
                    LexerToken previousToken = (tokens.size() > 0) ? tokens.get(tokens.size() - 1) : null;
                    LexerToken nextToken = new LexerToken(rule.getLastMatched(), rule.getTokenType());
                    if (previousToken != null) {
                        previousToken.setNextToken(nextToken);
                    }
                    boolean skip = false;
                    for (String key : skippedTokens) {
                        if (nextToken.getType().equals(key)) {
                            skip = true;
                            break;
                        }
                    }
                    if (!skip) {
                        tokens.add(nextToken);
                    }
                    break;
                }
                if (match == null) {
                    tokens.add(new LexerToken(line));
                }
            }
            br.close();
        } catch (Exception e) {
        }
        return tokens;
    }

    public List<String> getSkippedTokens() {
        return skippedTokens;
    }
}
