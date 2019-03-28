package org.gatblau.onix.data;

import java.util.List;

public class ResultList extends Wrapper<Result> {
    private boolean error = false;
    private boolean changed = false;
    private String message = "";

    public ResultList() {
    }

    public ResultList(List<Result> result) {
        super(result);
    }

    public boolean isError() {
        return error;
    }

    public boolean isChanged() {
        return changed;
    }

    public String getMessage() {
        return message;
    }

    public void add(Result result) {
        if (result != null){
            if (result.isError()) {
                error = true;
                message += String.format("| ERROR: %s ", result.getMessage());
            }
            if (result.isChanged()) {
                changed = true;
            }
            getValues().add(result);
        }
    }
}
