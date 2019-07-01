# King Arthur's Gold Stats Collector

KAG (King Arthur's Gold) is a 2d, action, platforming, medieval themed multiplayer game.
You can find out more about the game on the website: https://kag2d.com/en/

This is a stats collector for the game, it connects to the game server using the open
TCP socket and listens for specific events. The events are then logged into a database.

You can view the trello board with road map of items to be done and backlog: https://trello.com/b/WR8dcqD7/kag-stats

## Stats Mod

The game server runs a mod that captures events and writes formatted information to the TCP connection.

## Collector

The collector connects to the game server, captures events and puts new entries in the database.
A single collector can connect to multiple game servers.

## Monitoring

The collector provides the option to launch a monitoring web server.
That shows basic information about the Collector and all the status of connections
to game servers.

## Indexer

Because we're collecting so much data, fetching stats directly from the database such as k/d per class
can be costly and slow. As such instead an indexer service will take the raw data and incrementally
write out stats such as k/d, nemesis, most played server, etc.

## Webserver

The webserver is responsible for displaying and browsing stats. It will also have an API layer that talks to the database
to fetch stats compiled by the ingestion.

## Building and Running
Install docker and docker-compose

`docker-compose build` will build all the images
`docker-compose build <service name>` will build the image for a specific service

`docker-compose up -d` Will start all the services
After that you can do docker logs -f kagstats_basic-indexer_1 to see the logs of one of the indexers.
Or you can use the mysql command line to access the database. The API is also setup by the docker-compose file
at localhost:8080


To tear down the environment run:
`docker-compose down` and `docker-compose stop`
