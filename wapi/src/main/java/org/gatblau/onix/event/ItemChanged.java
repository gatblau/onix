/*
    Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org

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
package org.gatblau.onix.event;

import org.gatblau.onix.data.ItemData;

class ItemChanged {
    private final ItemData item;
    private final char notifyType;

    ItemChanged(char notifyType, ItemData item) {
        this.item = item;
        this.notifyType = notifyType;
    }

    @Override
    public String toString() {
        switch (notifyType) {
            case 'T':
                return String.format("key=%s" , item.getKey());
            case 'I':
                return String.format("changed");
            default:
                throw new RuntimeException(String.format("notify type '%s' not recognised", notifyType));
        }

    }

    String getTopicName() {
        switch (notifyType) {
            case 'T':
                return String.format("ITEM_TYPE_%s" , item.getType().toUpperCase());
            case 'I':
                return String.format("ITEM_INSTANCE_%s" , item.getType().toUpperCase());
            default:
                throw new RuntimeException(String.format("notify type '%s' not recognised", notifyType));
        }
    }
}
