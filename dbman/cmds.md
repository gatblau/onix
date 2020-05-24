# DbMan command hierarchy

The command hierarchy is:

- config (shows DbMan configuration variables)
- release (release information)
    - plan (shows the release plan)
    - info (shows a specific release information)
- db (database maintenance)
    - version (shows the database version)
    - deploy (deploy the latest or a specific release)
    - upgrade (upgrades to a specific release)
    - backup (backups the database)
    - restore (restores the database)
- serve (starts dbman as an http service)
