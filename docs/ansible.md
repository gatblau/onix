# Onix Ansible Modules

In order to maintain accurate information in the CMDB it is important that:
- Changes are recorded in the CMDB as soon as configuration changes are made (If changes are driven by Ansible, then it is important to update the CMDB automatically whilst Ansible performs any configuration changes). Or
- Changes are driven by information in the CMDB. Or
- Both of the above cases.

In order to facilitate this, a set of Ansible modules are provided in the [library](../ansible/library) folder, as follows:

|Module| Description |
|---|---|
| [**onix_login**](../ansible/library/onix_login.py)| Connects to an OpenId enabled authentication server and requests an access token to authenticate requests to the Onix RESTful API. This module does not currently support automatic token refreshes after expiration. |
| [**onix_item**](../ansible/library/onix_item.py)| Creates a new or updates an existing configuration item. |
| [**onix_link**](../ansible/library/onix_item.py)| Creates a new or updates an existing link between two existing configuration items. |

 More modules will be added in the future.
 
## How to use the Ansible modules
 
For an example of how to use the above modules, take a look at the [playbook here](../ansible/site.yml).
 
To execute the playbook run the following command from the [ansible](../ansible) folder:

```bash
$ ansible-playbook -i inventory site.yml -vvv
```

**NOTE**: it is assumed Onix Service, Onix Database and Keycloak are running in the localhost under the default ports.
If this is not the case, update the variables in the [inventory](../ansible/inventory) file accordingly.

## How to install the Ansible modules

### In the same location of the playbooks / project

Copy the [library](../ansible/library) folder in the same folder where the playbook using the modules is (e.g. see [site.yml here](../ansible/site.yml)). 

### In a shared location

To share the modules across multiple projects, add an entry to the **/etc/ansible/ansible.cfg** file pointing to a shared library location as follows:

```bash
library = /usr/share/ansible/library
```
