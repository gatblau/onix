/*
Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Contributors to this project, hereby assign copyright in their code to the
project, to be licensed under the same terms as the rest of the code.
*/

package org.gatblau.onix.data;

import io.swagger.annotations.ApiModelProperty;
import org.springframework.http.HttpStatus;

import java.io.Serializable;

public class Result implements Serializable {
    private static final long serialVersionUID = 1L;

    private boolean changed;
    private String message;
    private String operation;
    private boolean error;
    private String ref;

    public Result() {
        this(null);
    }

    public Result(String ref) {
        this.ref = ref;
        this.changed = false;
        this.error = false;
        this.operation = "N";
    }

    @ApiModelProperty(
            position = 1,
            value = "A reference which identifies the entity this result is for.",
            example = "entity_type:entity_instance_01"
    )
    public String getRef() {
        return ref;
    }

    public void setRef(String ref) {
        this.ref = ref;
    }

    @ApiModelProperty(
            position = 2,
            value = "A message describing an error associated with the response",
            notes = "This value is empty if no error occurred whilst processing the request.",
            example = "empty"
    )
    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }

    @ApiModelProperty(
            position = 3,
            value = "A flag indicating whether the resource was changed as a result of the request.",
            example = "false"
    )
    public boolean isChanged() {
        return changed;
    }

    public void setChanged(boolean changed) {
        this.changed = changed;
    }

    @ApiModelProperty(
            position = 4,
            value = "A character indicating the type of operation executed on the resource.",
            notes = "I indicates INSERT, U indicates UPDATE, D indicates DELETE and L indicates OPTIMISTIC LOCK",
            example = "N"
    )
    public String getOperation() {
        return operation;
    }

    public void setOperation(String operation) {
        this.operation = operation;
        changed = (operation.equals("I") || operation.equals("U") || operation.equals("D"));
    }

    @ApiModelProperty(
            position = 5,
            value = "A flag indicating if the request resulted in an error condition.",
            notes = "If the flag is true then the message property contains the detail of the error.",
            example = "false"
    )
    public boolean isError() {
        return error;
    }

    public void setError(boolean error) {
        this.error = error;
    }

    public int getStatus() {
        if (isError()) {
            // if there is a business logic error returns 400
            return HttpStatus.BAD_REQUEST.value();
        } else {
            if (isChanged()) {
                if (getOperation().equals("I")){
                    return HttpStatus.CREATED.value();
                }
                if (getOperation().equals("D") || getOperation().equals("U")){
                    return HttpStatus.OK.value();
                }
            } else {
                if (getOperation().equals("L")){
                    // optimistic lock case
                    return HttpStatus.CONFLICT.value();
                }
                if (getOperation().equals("D")) {
                    return HttpStatus.NOT_FOUND.value();
                }
                if (getOperation().equals("N")) {
                    return HttpStatus.OK.value();
                }
            }
        }
        throw new RuntimeException("Return Status not identified.");
    }

    public void setMessage(Exception e) {
        if(e != null){
            if(e.getMessage() != null && e.getMessage().contains("valid_user")){
                String msg = String.format("Name must follow below policies \n "+
                " 1) start with alphabet \n "+
                " 2) ends with alpha numeric \n "+
                " 3) can contain underscore _ \n "+
                " 4) must be of minimum 3 character and maximum of 200 character");
                this.message = msg;
            }else{
                this.message = e.getMessage();
            }
            e.printStackTrace();
            this.setError(true);
        }
    }
}