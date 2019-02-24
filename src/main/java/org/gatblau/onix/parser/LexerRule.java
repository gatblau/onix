package org.gatblau.onix.parser;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class LexerRule {
    private Pattern pattern;
    private String regex;
    private boolean enabled;
    private String tokenType;
    private String lastMatched;

    public LexerRule(String tokenType, String regex) {
        this.tokenType = tokenType;
        this.regex = regex;
        this.pattern = Pattern.compile(regex);
        this.enabled = true;
    }

    public String getRegex() {
        return regex;
    }

    public void setRegex(String regex) {
        this.regex = regex;
        this.pattern = Pattern.compile(regex);
    }

    public Matcher match(String line) {
        Matcher matcher = this.pattern.matcher(line);
        if (matcher.find()) {
            if (matcher.groupCount() > 0) {
                lastMatched = matcher.group();
            } else {
                lastMatched = line.substring(matcher.start(), matcher.end());
            }
            return matcher;
        }
        return null;
    }

    public boolean isEnabled() {
        return enabled;
    }

    public void setEnabled(boolean enabled) {
        this.enabled = enabled;
    }

    public boolean toggle() {
        enabled = !enabled;
        return enabled;
    }

    public String getTokenType() {
        return tokenType;
    }

    public String getLastMatched() {
        return lastMatched;
    }
}
