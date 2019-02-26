package org.gatblau.onix.parser;

public class LexerToken {
    private String match;
    private String tokenType;
    private final boolean isNullMatch;
    private final String nullValueData;
    private LexerToken nextToken;

    public LexerToken(String match, String tokenType) {
        this.match = match;
        this.tokenType = tokenType;
        this.isNullMatch = false;
        this.nullValueData = null;
    }

    public LexerToken(String unknownData) {
        this.isNullMatch = true;
        this.nullValueData = unknownData;
        this.tokenType = "UNKNOWN";
    }

    public String getType() {
        return tokenType;
    }

    public String getValue() {
        return match;
    }

    @Override
    public String toString() {
        return (match != null) ? match : nullValueData;
    }

    public LexerToken getNextToken() {
        return nextToken;
    }

    void setNextToken(LexerToken nextToken) {
        this.nextToken = nextToken;
    }
}
