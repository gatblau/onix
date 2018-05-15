#!/usr/bin/python
#
# Onix - Copyright (c) 2018 gatblau.org
# Apache License Version 2 - https://www.apache.org/licenses/LICENSE-2.0
#
# Module: onix_item
# Description: creates a new or updates an existing configuration item
#
from ansible.module_utils.basic import *
from ansible.module_utils.urls import *

def createOrUpdateItem(data):
    # parse the input variables
    cmdb_host = data['cmdb_host']
    access_token = data['access_token']
    key = data['key']
    name = data['name']
    description = data['description']
    status = data['status']
    type = data['type']
    meta = data['meta']
    tag = data['tag']
    dimensions = data['dimensions']

    payload = {
        "name": name,
        "description": description,
        "type": type,
        "meta": meta,
        "tag": tag,
        "status": status,
        "dimensions": dimensions
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
    item_uri = "{}/item/{}".format(cmdb_host, key)

    # put the payload to the cmdb service
    stream = open_url(item_uri, method="PUT", data=payloadStr, headers=headers)

    # reads the returned stream
    result = json.loads(stream.read())

    has_changed = True

    return (has_changed, result)

# module entry point
def main():
    has_changed = False

    params = {
        "cmdb_host": {"required": True, "type": "str"},
        "access_token": {"required": False, "type": "str", "default": "", "no_log": True},
        "key": {"required": True, "type": "str"},
        "name": {"required": True, "type": "str"},
        "description": {"required": False, "type": "str", "default": ""},
        "status": {"required": False, "type": "int","default": 0},
        "type": {"required": True, "type": "str"},
        "meta": {"required": False, "type": "str", "default": "{}"},
        "tag": {"required": False, "type": "str", "default": ""},
        "dimensions": {"required": False, "type": "str", "default": "{}"}
    }

    # handle incoming parameters
    module = AnsibleModule(
        argument_spec = params,
        supports_check_mode = False
    )

    has_changed, result = createOrUpdateItem(module.params)

    # exit the module with a result (changed & meta json object)
    module.exit_json(
        changed = has_changed,
        meta = result
    )

if __name__ == '__main__':
    main()
