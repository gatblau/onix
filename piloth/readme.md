<img src="https://github.com/gatblau/onix/piloth/raw/master/pilot.png" width="150" align="right"/>

# Onix Host Pilot

## Build it

```bash
art run build
```

the binary can be found under the created bin folder.

## Run it

1. Launch the [Remote Service backend](https://github.com/gatblau/onix/tree/develop/rem/docker)
2. Set up the Pilot [configuration file](.pilot) in the folder where pilot is located
3. Run the pilot

```bash
./bin/pilot
```

## Example service

It may be that you wish to run Onix Host Pilot as a service - the following shows an example of how to do this utilising systemd on a Debian based OS where you have a copy of the Pilot binary in your working directory ready to use.

*NB. To perform the steps below, elevate privileges to root or run using SUDO*

### Create user

```
useradd -m piloth
```

### Copy pilot and set permissions
```
cp pilot /home/piloth/.
chown piloth:piloth /home/piloth/pilot
chmod 700 /home/piloth/pilot
```

### Create basic service

Create a new file `/lib/systemd/system/piloth.service` containing the following:

```
[Unit]
Description=Host Pilot service
ConditionPathExists=/home/piloth/pilot
After=network.target
 
[Service]
Type=simple
User=piloth
Group=piloth
LimitNOFILE=1024

Restart=on-failure
RestartSec=10
startLimitIntervalSec=60

WorkingDirectory=/home/piloth
ExecStart=/home/piloth/pilot

[Install]
WantedBy=multi-user.target
```