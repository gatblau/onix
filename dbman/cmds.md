# DbMan command hierarchy

The command hierarchy is:

- config (manages dbman's configuration)
    - show (shows the current configuration)
    - delete (deletes a configuration set)
    - list (list configuration sets)
    - set (set a configuration value in the current configuration)
    - use (use a specified configuration)
- release (release information)
    - plan (shows the release plan)
    - info (shows a specific release information)
- db (database maintenance)
    - version (shows the database version)
    - deploy (deploy the latest or a specific release)
    - upgrade (upgrades to a specific release)
    - backup (backups the database)
    - restore (restores the database)
- check (check that tools and connections are working for the current config set)
- serve (starts dbman as an http service)
