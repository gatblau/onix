# Query Module Example <img src="../../../../../docs/pics/ox.png" width="160" height="160" align="right">

This example shows how to query an item in the CMDB.

It first uses the **ox_import** module to import an item which can then be subsequently queried.

## Trying out the example
To run the example do the following:

1. Ensure Onix is up and running - Run docker compose with this [docker-compose.yml](../../docker-compose.yml) file.
2. Ensure Ansible is installed in the machine running the example.
3. Copy the Python scripts [here](../..) to a location that can be found by Ansible (for instance, within a 
folder called "library" under the [site.yml](site.yml) file - this is already provided by this project as a symlink).
4. Open the terminal where the [site.yml](site.yml) file is and type the following:

```bash
$ ansible-playbook site.yml -i inventory
```

At the end of the execution, you should see the results of the queries on the terminal window.

([back to index](../readme.md))
