package org.gatblau.onix.data;

import java.io.Serializable;

public class Result implements Serializable {
    private static final long serialVersionUID = 1L;

    private boolean changed;
    private String message;
    private String operation;
    private boolean error;

    public Result() {
        this.changed = false;
        this.error = false;
        this.operation = "N";
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    public boolean isChanged() {
        return changed;
    }

    public void setChanged(boolean changed) {
        this.changed = changed;
    }

    public String getOperation() {
        return operation;
    }

    public void setOperation(String operation) {
        this.operation = operation;
        changed = (operation.equals("I") || operation.equals("U"));
    }

    public boolean isError() {
        return error;
    }

    public void setError(boolean error) {
        this.error = error;
    }
}
