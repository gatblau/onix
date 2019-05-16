/*
Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org

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
