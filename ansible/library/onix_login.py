#!/usr/bin/python
#
# Onix - Copyright (c) 2018 gatblau.org
# Apache License Version 2 - https://www.apache.org/licenses/LICENSE-2.0
#
# Module: onix_login
# Description: request an openid connect access token by login on the authentication server
#
from ansible.module_utils.basic import *
from ansible.module_utils.urls import *

def login(data):
    has_changed = False

    # parse the input variables
    realm = data['realm']
    auth_host = data['auth_host']
    client_id = data['client_id']
    username = data['username']
    password = data['password']

    # builds payload with url encoded form parameters
    payload = 'client_id={}&grant_type=password&username={}&password={}'.format(client_id, username, password)

    # use line below for testing posting payload
    # auth_uri = "https://httpbin.org/post"

    # builds the URI required by the token service
    auth_uri = "{}/auth/realms/{}/protocol/openid-connect/token".format(auth_host, realm)

    # post form urlencoded data fields to the token service
    stream = open_url(auth_uri, method="POST", data=payload)

    # reads the returned stream
    result = json.loads(stream.read())

    return (has_changed, result, auth_uri, payload)

# module entry point
def main():
    params = {
        "auth_host": {"required": True, "type": "str"},
        "realm": {"required": True, "type": "str"},
        "client_id": {"required": True, "type": "str"},
        "username": {"required": True, "type": "str"},
        "password": {"required": True, "type": "str", "no_log": True}
    }

    # handle incoming parameters
    module = AnsibleModule(
        argument_spec = params,
        supports_check_mode = False
    )

    has_changed, result, auth_uri, payload = login(module.params)

    # exit the module with a result (changed & meta json object)
    module.exit_json(
        changed = has_changed,
        ansible_facts = {
            "onix_access_token": result['access_token'],
            "onix_auth_uri": auth_uri,
            "onix_auth_fields": payload
        }
    )

if __name__ == '__main__':
    main()
