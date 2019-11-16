#!/usr/bin/python
#
# Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org
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
# Ansible Module: ox_setup
# Description:
#   sets up all variables required to connect to the Onix WAPI.
#
from ansible.module_utils.basic import *
from ansible.module_utils.urls import *

import base64

# creates a basic auth token using the passed in username and password
def get_basic_token(username, password):
    return "Basic %s" % (base64.b64encode("%s:%s" % (username, password)))

# following the OpenId Resource Owner Password Flow, gets a bearer token
def get_bearer_token(token_uri, clientId, secret, username, password):

    # creates a basic auth token using the authorisation server client id and secret
    basic_token = get_basic_token(clientId, secret)

    # prepares the headers for the post request to the token endpoint
    headers = {
        "accept":"application/json",
        "authorization":basic_token,
        "cache-control":"no-cache",
        "content-type":"application/x-www-form-urlencoded"
    }

    # with a payload indicating a client credentials flow and the onix scope
    payloadStr = 'grant_type=password&username={}&password={}&scope=openid%20onix'.format(username, password)

    # request the access token
    stream = open_url(token_uri, method="POST", data=payloadStr, headers=headers)

    # reads the returned token
    response = json.loads(stream.read())

    # returns a bearer token
    return "Bearer %s" % response["access_token"]

# returns an access token for the Onix WAPI
def get_access_token(data):
    client_id = data['client_id']
    secret = data['secret']
    username = data['username']
    password = data['password']
    auth_mode = data["auth_mode"]
    wapi_uri = data['uri']
    token_uri = data['token_uri']

    access_token = ""

    if auth_mode == "basic":
        # the access token is a basic access authentication token
        access_token = get_basic_token(username, password)

    elif auth_mode == "oidc":
        # the access token is a OAuth 2.0 bearer token
        access_token = get_bearer_token(token_uri, client_id, secret, username, password)

    return (access_token, wapi_uri)

# module entry point
def main():
    params = {
        "uri": {"required": True, "type": "str"},
        "client_id": {"required": False, "type": "str"},
        "secret": {"required": False, "type": "str", "no_log": True},
        "username": {"required": True, "type": "str"},
        "password": {"required": True, "type": "str", "no_log": True},
        "auth_mode": {"required": False, "type": "str", "default": "none"},
        "token_uri": {"required": False, "type": "str", "default": "none"}
    }

    # handle incoming parameters
    module = AnsibleModule(argument_spec=params, supports_check_mode=False)

    # obtains an access token
    access_token, wapi_uri = get_access_token(module.params)

    # exit the module with a result (changed & connection facts)
    module.exit_json(
        changed=False,
        ansible_facts={
            "ox_token": access_token,  # adds the Onix Access Token to the fact list
            "ox_uri": wapi_uri  # adds the Onix WAPI URL to the fact list
        },
        meta={
        }
    )

if __name__ == '__main__':
    main()
