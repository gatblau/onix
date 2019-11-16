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
# Ansible Module: ox_query
# Description:
#   query configuration data
#
from ansible.module_utils.basic import *
from ansible.module_utils.urls import *

def getHeaders(token):
    if token == "":
        # if not access token is provided do not send it to the service
        headers = {"Content-Type": "application/json"}
    else:
        # if an access token exists then add it to the request headers
        headers = {"Content-Type": "application/json", "Authorization": token}
    return headers

def query(module):
    # parse the input variables
    wapi_uri = module.params['uri']
    access_token = module.params['token']
    key = module.params['key']
    type = module.params['type']

    # use line below for testing posting payload
    # item_uri = "https://httpbin.org/get"

    # builds the URI required by the cmdb service
    resource_uri = "{}/{}/{}".format(wapi_uri, type, key)

    # put the payload to the cmdb service
    stream = open_url(resource_uri, method="GET", headers=getHeaders(access_token))

    # reads the returned stream
    result = json.loads(stream.read())

    return (result)

# module entry point
def main():
    params = {
        "uri": {"required": True, "type": "str"},
        "token": {"required": False, "type": "str", "default": "", "no_log": True},
        "key": {"required": True, "type": "str"},
        "type": {"required": False, "type": "str", "default": "item", "choices": ["item", "link", "item_type", "link_type", "link_rule", "model"]}
    }

    # handle incoming parameters
    module = AnsibleModule(argument_spec=params, supports_check_mode=False)

    result = query(module)

    # exit the module with a result (changed & meta json object)
    module.exit_json(changed=False, **result)

if __name__ == '__main__':
    main()
