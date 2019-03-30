# Importing Data Example

In order to import data into the CMDB, the [ox_import](../../../modules/ox_import.py) module is provided.

The module takes a json file containing data to be imported and send them to the Onix web 
API. 

This example shows how to import a custom [meta model](app_model.json) executing a single call to the Web API.


## Trying out the example
To run the example do the following:

1. Ensure Onix is up and running - Run docker compose with this [docker-compose.yml](../../docker-compose.yml) file.
2. Ensure Ansible is installed in the machine running the example.
3. Copy the [ox_import](../../../modules/ox_import.py) module to a location that can be found by Ansible (for instance, within a 
folder called "library" under the [site.yml](site.yml) file).
4. Open the terminal where the [site.yml](site.yml) file is and type the following:

```bash
$ ansible-playbook site.yml -i inventory
```

At the end of the execution, you should have the metal model imported into the CMDB.

## Checking the results

For example, to take a look at the created configuration items use the 
[Swagger UI /itemtype GET](http://localhost:8080/swagger-ui.html#/web-api/getItemTypesUsingGET) endpoint.

Note that the module is idempotent, if you run the playbook again, no changes should be made in the CMDB.

([back to index](../readme.md))