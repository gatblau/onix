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
# Ansible Module: ox_import
# Description:
#   Import data into the CMDB.
#   Data can be models, item types, link types, link rules, items and links.
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

def putData(module):
    # parse the input variables
    wapi_uri = module.params['uri']
    access_token = module.params['token']
    src = module.params['src']

    # read the content of the data json file into a string
    file = open(src, "r")
    payload = file.read()
    file.close()

    # builds the URI required by the cmdb service
    data_uri = "{}/data".format(wapi_uri)

    # put the payload to the cmdb service
    stream = open_url(data_uri, method="PUT", data=payload, headers=getHeaders(access_token))

    # reads the returned stream
    result = json.loads(stream.read())

    return (result)

# module entry point
def main():
    params = {
        "uri": {"required": True, "type": "str"},
        "token": {"required": False, "type": "str", "default": "", "no_log": True},
        "src": {"required": True, "type": "str"}
    }

    # handle incoming parameters
    module = AnsibleModule(argument_spec=params, supports_check_mode=False)

    # push the data to the service
    result = putData(module)

    # handles errors
    if result['error']:
        module.fail_json(msg=result['message'], **result)
    else:
        # exit the module with a result (changed & meta json object)
        module.exit_json(changed=result['changed'], meta=result)

if __name__ == '__main__':
    main()
