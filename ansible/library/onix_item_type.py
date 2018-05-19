#!/usr/bin/python
#
# Onix - Copyright (c) 2018 gatblau.org
# Apache License Version 2 - https://www.apache.org/licenses/LICENSE-2.0
#
# Module: onix_item_type
# Description: creates a new or updates an existing configuration item type
#
from ansible.module_utils.basic import *
from ansible.module_utils.urls import *

def deleteItemType(data):
    # parse the input variables
    cmdb_host = data['cmdb_host']
    access_token = data['access_token']
    key = data['key']

    if access_token == "":
        # if not access token is provided do not send it to the service
        headers = {"Content-Type": "application/json"}
    else:
        # if an access token exists then add it to the request headers
        headers = {"Content-Type": "application/json", "Authorization": "bearer {}".format(access_token)}

    # use line below for testing posting payload
    # item_uri = "https://httpbin.org/delete"

    # builds the URI required by the cmdb service
    item_type_uri = "{}/itemtype/{}".format(cmdb_host, key)

    # put the payload to the cmdb service
    stream = open_url(item_type_uri, method="DELETE", headers=headers)

    # reads the returned stream
    result = json.loads(stream.read())

    return (result)

def createOrUpdateItemType(data):
    # parse the input variables
    cmdb_host = data['cmdb_host']
    access_token = data['access_token']
    key = data['key']
    name = data['name']
    description = data['description']

    payload = {
        "name": name,
        "description": description
    }

    if access_token == "":
        # if not access token is provided do not send it to the service
        headers = {"Content-Type": "application/json"}
    else:
        # if an access token exists then add it to the request headers
        headers = {"Content-Type": "application/json", "Authorization": "bearer {}".format(access_token)}

    payloadStr = json.dumps(payload).replace('"{','{').replace('}"', '}').replace('\'', '\"')

    # use line below for testing posting payload
    # item_uri = "https://httpbin.org/put"

    # builds the URI required by the cmdb service
    item_type_uri = "{}/itemtype/{}".format(cmdb_host, key)

    # put the payload to the cmdb service
    stream = open_url(item_type_uri, method="PUT", data=payloadStr, headers=headers)

    # reads the returned stream
    result = json.loads(stream.read())

    return (result)

# module entry point
def main():
    params = {
        "cmdb_host": {"required": True, "type": "str"},
        "access_token": {"required": False, "type": "str", "default": "", "no_log": True},
        "key": {"required": True, "type": "str"},
        "name": {"required": False, "type": "str"},
        "description": {"required": False, "type": "str", "default": ""},
        "state": {"required": False, "type": "str", "default": "present"}
    }

    # handle incoming parameters
    module = AnsibleModule(
        argument_spec = params,
        supports_check_mode = False
    )

    state = module.params['state']

    if state == "absent":
        result = deleteItemType(module.params)
    else:
        result = createOrUpdateItemType(module.params)

    # exit the module with a result (changed & meta json object)
    module.exit_json(
        changed = result['changed'],
        meta = result
    )

if __name__ == '__main__':
    main()
