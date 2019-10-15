# Hosting

You will need root access to a server where you
want to host the site. It is recommened that this server is running a linux distribution like ubuntu or Centos7 for ease of setup.

You can use any VPS such as Vultr, AWS or GCP.
It is not recommened to use a home computer to host a site, but technically possible.

The rest of this document will assume the host is running Centos7 however you can easily find equivalent steps.

## Hardware Requirements

Recommened:
* 1 CPU
* 2GB Memory
* 10GB Disk

## Docker and docker-compose

KAG Stats runs using a docker-compose file. You will need to install and configure both.

### Install Docker CE

https://docs.docker.com/v17.09/engine/installation/linux/docker-ce/centos/#uninstall-old-versions

### Install docker-compose

https://docs.docker.com/compose/install/

```
sudo curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

sudo chmod +x /usr/local/bin/docker-compose
```

## KAG Stats

### Install

Make a directory where you can store your configuration and download the docker-compose file needed to run the services.
```
mkdir -p /opt/kagstats
cd /opt/kagstats

curl https://raw.githubusercontent.com/Harrison-Miller/kagstats/master/docker-compose.yaml --output docker-compose.yaml
```

### Configuration

Create settings.json file modify the servers section to point to the KAG servers you wish to monitor.

**NOTE:** Your KAG server must be configured with an rcon password and sv_tcpr set to 1

```
cat >> settings.json <<-EOF
{
    "name": "<site-name> collector",
    "motd": "visit <site-name> to view leaderboards",
    "server": [
        {
            "address": "<kag-server-ip>",
            "port": 50301,
            "password": <rcon-password>,
            "tags": ["REGION", "GAMEMODE", "SOMETHING"]
        }
    ]
}
EOF
```

### Run KAG Stats
```
docker-compose up -d
```

Now put the ip of your host in the browser and it'll start working.

## Troubleshooting

If you aren't seeing kills from server appear you can view the logs by going to the docker-compose directory and running:

```
docker-compose logs collector
```

If something is wrong with the collector like the configuration of the KAG server. You can modify the configuration and restart the collector with:

```
docker-compose restart collector
```

## Shutting down

If you wish to no longer host the site simply run:

```
docker-compose stop
```

## Additional host configuration

If you are serious about hosting KAG Stats for your own servers you will want to take the time to do these things.

* Buy a domain and point your domain at your host
* Change the docker https://linuxconfig.org/how-to-move-docker-s-default-var-lib-docker-to-another-directory-on-ubuntu-debian-linux
* Turn off auto-updates for docker (This will reset the above changes and take down your site)
* Setup portainer (a management UI for docker)
* Setup ssl using nginx and letsencrypt

