# Database Releases

In order to provide intelligent database releases, a naming convention is required for installation and upgrades of the database schemas and database objects.

A git repository will hold a folder structure containing install and upgrade information for the database as follows:

```
|
+-- 1.0.0 (Current App Version)
|   +-- oxdb_ix.yml (index file for the current version)
|   |   +-- 1.0.0 (Target App Version)
|   |   |   +-- 1 (Database release number)
|   |   |   |   +-- install_schema.sql
|   |   |   |   +-- install_objects.sql
|   |   |   |   +-- upgrade_schema.sql
|   |   |   |   +-- upgrade_objects.sql
|   |   |   +-- 2 (Database release number)
|   |   |   |   +-- install_schema.sql
|   |   |   |   +-- ...
|   |   +-- 2.0.0 (Target App Version)
|   |   |   +-- 1 (Database release number)
|   |   |   |   +-- install_schema.sql
|   |   |   |   +-- ...
|   |   |   +-- 2 (Database release number)
|   |   |   |   +-- install_schema.sql
|   |   |   |   +-- ...
+-- 2.0.0 (Current App Version)
|   +-- oxdb_ix.yml (index file for the current version)
|   |   +-- 2.0.0 (Target App Version)
|   |   |   +-- 1 (Database release number)
|   |   |   |   +-- install_schema.sql
|   |   |   |   +-- ...
|   |   |   +-- 2 (Database release number)
|   |   |   |   +-- install_schema.sql
|   |   |   |   +-- ...
|   |   +-- 3.0.0 (Target App Version)
|   |   |   +-- 1 (Database release number)
|   |   |   |   +-- install_schema.sql
|   |   |   |   +-- ...
|   |   |   +-- 2 (Database release number)
|   |   |   |   +-- install_schema.sql
|   |   |   |   +-- ...
```

where:

| Item | Type | Description |
|---|---|---|
| __Current App Version__| Folder | The entry point folder for a specific application version. |
| __Index File__| File | The manifest for all available database scripts for the current application version. This file should be called __oxdb_ix.yml__ |
| __Target App Version__ | Folder | The version to upgrade to.  |
| __DB Release Number__ | Folder | A folder containing the database scripts for a particular release. |
| __Install Schema__ | File | The script that creates the tables and other schema objects for the database release. |
| __Install Objects__ | File | The script that creates the database objects such as functions, etc for the database release. |
| __Upgrade Schema__ | File | The script that creates the tables and other schema objects for the database release. |
| __Upgrade Objects__ | File | The script that creates the database objects such as functions, etc for the database release. |
