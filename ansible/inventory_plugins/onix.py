from __future__ import (absolute_import, division, print_function)

__metaclass__ = type

DOCUMENTATION = '''
    name: onix
    plugin_type: inventory
    author: gatblau.org
    short_description: Ansible dynamic inventory plugin for Onix CMDB.
    version_added: "2.7"
    description:
        - Reads inventories from Onix CMDB.
        - Supports reading configuration from both YAML config file and environment variables.
        - If reading from the YAML file, the file name must end with onix.(yml|yaml) or onix_inventory.(yml|yaml),
          the path in the command would be /path/to/onix_inventory.(yml|yaml). If some arguments in the config file
          are missing, this plugin will try to fill in missing arguments by reading from environment variables.
        - If reading configurations from environment variables, the path in the command must be @onix_inventory.
    options:
        plugin:
            description: the name of this plugin, it should always be set to 'onix'
                for this plugin to recognize it as it's own.
            env:
                - name: ANSIBLE_INVENTORY_ENABLED
            required: True
            choices: ['onix']
        host:
            description: The network address of your Onix Web API host.
            type: string
            env:
                - name: OX_HOST
            required: True
        username:
            description: The user that you plan to use to access inventories on the Onix WAPI.
            type: string
            env:
                - name: OX_USERNAME
            required: True
        password:
            description: The password for your Onix WAPI user.
            type: string
            env:
                - name: OX_PASSWORD
            required: True
        inventory_id:
            description: The ID of the Onix inventory that you wish to import.
            type: string
            env:
                - name: OX_INVENTORY_ID
            required: True
        inventory_version:
            description: The label of the Onix inventory that you wish to import.
            type: string
            env:
                - name: OX_INVENTORY_VERSION
            required: True
        verify_ssl:
            description: Specify whether Ansible should verify the SSL certificate of the Onix WAPI host.
            type: bool
            default: True
            env:
                - name: OX_VERIFY_SSL
            required: False
'''

import re
import os
import json
from ansible.module_utils import six
from ansible.module_utils.urls import Request, urllib_error, ConnectionError, socket, httplib
from ansible.module_utils._text import to_native
from ansible.errors import AnsibleParserError
from ansible.plugins.inventory import BaseInventoryPlugin

# Python 2/3 Compatibility
try:
    from urlparse import urljoin
except ImportError:
    from urllib.parse import urljoin

class InventoryModule(BaseInventoryPlugin):

    NAME = 'onix'  # used internally by Ansible, it should match the file name but not required

    # If the user supplies '@onix_inventory' as path, the plugin will read from environment variables.
    no_config_file_supplied = False

    def make_request(self, request_handler, onix_url):
        """Makes the request to given URL, handles errors, returns JSON"""
        try:
            response = request_handler.get(onix_url)
        except (ConnectionError, urllib_error.URLError, socket.error, httplib.HTTPException) as e:
            error_msg = 'Connection to remote host failed: {err}'.format(err=e)
            # If onix gives a readable error message, display that message to the user.
            if callable(getattr(e, 'read', None)):
                error_msg += ' with message: {err_msg}'.format(err_msg=e.read())
            raise AnsibleParserError(to_native(error_msg))

        # Attempt to parse JSON.
        try:
            return json.loads(response.read())
        except (ValueError, TypeError) as e:
            # If the JSON parse fails, print the ValueError
            raise AnsibleParserError(to_native('Failed to parse json from host: {err}'.format(err=e)))

    # determines if the inventory source provided is usable by the plugin
    def verify_file(self, path):
        if path.endswith('@onix_inventory'):
            self.no_config_file_supplied = True
            return True
        elif super(InventoryModule, self).verify_file(path):
            return path.endswith(('onix_inventory.yml', 'onix_inventory.yaml', 'onix.yml', 'onix.yaml'))
        else:
            return False

    def parse(self, inventory, loader, path, cache=True):
        super(InventoryModule, self).parse(inventory, loader, path)
        if not self.no_config_file_supplied and os.path.isfile(path):
            self._read_config_data(path)

        # Read inventory from onix service
        # Note the environment variables will be handled automatically by InventoryManager.
        onix_host = self.get_option('host')
        if not re.match('(?:http|https)://', onix_host):
            onix_host = 'https://{onix_host}'.format(onix_host=onix_host)

        request_handler = Request(url_username=self.get_option('username'),
                                  url_password=self.get_option('password'),
                                  force_basic_auth=True,
                                  validate_certs=self.get_option('verify_ssl'))

        inventory_id = self.get_option('inventory_id').replace('/', '')
        inventory_version = self.get_option('inventory_version').replace('/', '')
        inventory_url = '/inventory/{inv_id}/{inv_ver}'.format(inv_id=inventory_id, inv_ver=inventory_version)
        inventory_url = urljoin(onix_host, inventory_url)

        inventory = self.make_request(request_handler, inventory_url)

        # TODO: complete plugin development here...
        raise Exception(inventory)
