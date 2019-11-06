# Granular Ansible Modules Example <img src="../../../../../docs/pics/ox.png" width="160" height="160" align="right">

This sample shows how to use the granular Ansible module set for Onix.
These modules perform CRUD operations on entities in the CMDB.

For example, they allow to create or update and delete items, links and their associated types.

A full example using all these modules is provided here.

## Trying out the example
To run the example do the following:

1. Ensure Onix is up and running - Run docker compose with this [docker-compose.yml](../../docker-compose.yml) file.
2. Ensure Ansible is installed in the machine running the example.
3. Copy the Python scripts [here](../../modules) to a location that can be found by Ansible (for instance, within a 
folder called "library" under the [site.yml](site.yml) file).
4. Open the terminal where the [site.yml](site.yml) file is and type the following:

```bash
$ ansible-playbook site.yml -i inventory
```

At the end of the execution, you should have configuration items in the CMDB.

## Checking the results

For example, to take a look at the created configuration items use the 
[Swagger UI /item GET](http://localhost:8080/swagger-ui.html#/web-api/getItemsUsingGET) endpoint.

Note that the modules are idempotent, if you run the playbook again, no changes should be made in the CMDB.

([back to index](../readme.md))