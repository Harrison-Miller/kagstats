# King Arthur's Gold Stats Collector

KAG (King Arthur's Gold) is a 2d, action, platforming, medieval themed multiplayer game.
You can find out more about the game on the website: https://kag2d.com/en/

This is a stats collector for the game, it connects to the game server using the open
TCP socket and listens for specific events. The events are then logged into a database.

## Stats Mod

The game server runs a mod that captures events and writes formatted information to the TCP connection.

## Collector

The collector connects to the game server, captures events and puts new entries in the database.
A single collector can connect to multiple game servers.

## Monitoring

The collector provides the option to launch a monitoring web server.
That shows basic information about the Collector and all the status of connections
to game servers.

## Ingestion

Because we're collecting so much data, fetching stats directly from the database such as k/d per class
can be costly and slow. As such instead an ingestion service will take the raw data and incrementally
write out stats such as k/d, nemesis, most played server, etc.

## Webserver

The webserver is responsible for displaying and browsing stats. It will also have an API layer that talks to the database
to fetch stats compiled by the ingestion.
