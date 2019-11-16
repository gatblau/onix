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
# Ansible Module: ox_model
# Description:
#   creates a new or updates an existing meta model.
#   deletes an existing meta model if absent is set to true.
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

def createOrUpdateModel(module):

    # parse the input variables
    wapi_uri = module.params['uri']
    access_token = module.params['token']
    key = module.params['key']
    name = module.params['name']
    description = module.params['description']

    payload = {
        "name": name,
        "description": description
    }

    payloadStr = json.dumps(payload).replace('"{', '{').replace('}"', '}').replace('\'', '\"')

    # builds the URI required by the cmdb service
    model_uri = "{}/model/{}".format(wapi_uri, key)

    # put the payload to the cmdb service
    stream = open_url(model_uri, method="PUT", data=payloadStr, headers=getHeaders(access_token))

    # reads the returned stream
    result = json.loads(stream.read())

    return (result)

def deleteModel(module):

    # parse the input variables
    wapi_uri = module.params['uri']
    access_token = module.params['token']
    key = module.params['key']

    # builds the URI required by the cmdb service
    model_uri = "{}/model/{}".format(wapi_uri, key)

    # put the payload to the cmdb service
    stream = open_url(model_uri, method="DELETE", headers=getHeaders(access_token))

    # reads the returned stream
    result = json.loads(stream.read())

    return (result)

# module entry point
def main():
    params = {
        "uri": {"required": True, "type": "str"},
        "token": {"required": False, "type": "str", "default": "", "no_log": True},
        "key": {"required": True, "type": "str"},
        "name": {"required": False, "type": "str"},
        "description": {"required": False, "type": "str", "default": ""},
        "state": {"required": False, "type": "str", "default": "present"}
    }

    # handle incoming parameters
    module = AnsibleModule(argument_spec=params, supports_check_mode=False)

    state = module.params['state']

    if state == "absent":
        result = deleteModel(module)
    else:
        result = createOrUpdateModel(module)

    if result['error']:
        module.fail_json(msg=result['message'], **result)
    else:
        # exit the module with a result (changed & meta json object)
        module.exit_json(changed=result['changed'], meta=result)


if __name__ == '__main__':
    main()
