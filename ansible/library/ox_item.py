#!/usr/bin/python
#
# Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Contributors to this project, hereby assign copyright in their code to the
# project, to be licensed under the same terms as the rest of the code.
#
# Ansible Module: ox_item
# Description:
#   creates a new, updates or deletes an existing configuration item.
#
from ansible.module_utils.basic import *
from ansible.module_utils.urls import *


def createOrUpdateItem(module):
    data = module.params

    # parse the input variables
    wapi_uri = data['uri']
    access_token = data['token']
    key = data['key']
    name = data['name']
    description = data['description']
    status = data['status']
    type = data['type']
    meta = data['meta']
    tag = data['tag']
    attribute = data['attribute']

    payload = {
        "name": name,
        "description": description,
        "type": type,
        "meta": meta,
        "tag": tag,
        "status": status,
        "attribute": attribute
    }

    if access_token == "":
        # if not access token is provided do not send it to the service
        headers = {"Content-Type": "application/json"}
    else:
        # if an access token exists then add it to the request headers
        headers = {"Content-Type": "application/json", "Authorization": access_token}

    payloadStr = json.dumps(payload).replace('"{', '{').replace('}"', '}').replace('\'', '\"')

    # use line below for testing posting payload
    # item_uri = "https://httpbin.org/put"

    # builds the URI required by the cmdb service
    item_uri = "{}/item/{}".format(wapi_uri, key)

    # put the payload to the cmdb service
    stream = open_url(item_uri, method="PUT", data=payloadStr, headers=headers)

    # reads the returned stream
    result = json.loads(stream.read())

    return (result)


def deleteItem(module):
    data = module.params

    # parse the input variables
    wapi_uri = data['uri']
    access_token = data['token']
    key = data['key']

    if access_token == "":
        # if not access token is provided do not send it to the service
        headers = {"Content-Type": "application/json"}
    else:
        # if an access token exists then add it to the request headers
        headers = {"Content-Type": "application/json", "Authorization": access_token}

    # use line below for testing posting payload
    # item_uri = "https://httpbin.org/delete"

    # builds the URI required by the cmdb service
    item_uri = "{}/item/{}".format(wapi_uri, key)

    # put the payload to the cmdb service
    stream = open_url(item_uri, method="DELETE", headers=headers)

    # reads the returned stream
    result = json.loads(stream.read())

    return (result)


# module entry point
def main():
    has_changed = False

    params = {
        "uri": {"required": True, "type": "str"},
        "token": {"required": False, "type": "str", "default": "", "no_log": True},
        "key": {"required": True, "type": "str"},
        "name": {"required": False, "type": "str"},
        "description": {"required": False, "type": "str", "default": ""},
        "status": {"required": False, "type": "int", "default": 0},
        "type": {"required": False, "type": "str"},
        "meta": {"required": False, "type": "str", "default": "{}"},
        "tag": {"required": False, "type": "list", "default": "None"},
        "attribute": {"required": False, "type": "str", "default": "{}"},
        "state": {"required": False, "type": "str", "default": "present"}
    }

    # handle incoming parameters
    module = AnsibleModule(
        argument_spec=params,
        supports_check_mode=False
    )

    state = module.params['state']

    if state == "absent":
        result = deleteItem(module)
    else:
        result = createOrUpdateItem(module)

    if result['error']:
        module.fail_json(msg=result['message'], **result)
    else:
        # exit the module with a result (changed & meta json object)
        module.exit_json(
            changed=result['changed'],
            meta=result
        )


if __name__ == '__main__':
    main()
