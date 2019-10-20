# Angular Development

The Angular portion of KAG Stats serves the UI and interfaces with the API. There are two ways to work on the angular application
depending on what you need to accomplish. 

The first way, we'll call "Local Development", is for when you don't need to change the API or how the site interacts with it. You can run the angular application in a development mode from the CLI and point to the production api at https://kagstats.com/api. This will allow you to change the display, style and functionality of the site without modifications to the API.

If you need to modify the API and how the site interfaces with it you'll use the second method which we will call  "Dockerized Development".  This will allow you to quickly modify the API and point your local angular site to your own API.

Don't worry both modes work with angular hot reloading!

## Setup Angular

In order to install angular CLI you will first need to install Node.js and npm.

### Node.js and npm

Visit https://nodejs.org/en/download and follow the instruction to download.
npm will be installed with Node.js.

### Angular

Visit https://cli.angular.io and follow the instruction to download.
TL;DR `npm install -g @angular/cli`

### Setup

Run: `npm install` to download all the node modules.

## Local Development

API at: https://kagstats.com/api

Open a terminal to the ui directory of the kagstats repo.
Run: `ng serve --watch --poll 250 -c=dev`. Then open your browser to
http://localhost:4200.

**NOTE:** Defining poll and watch together because often watch doesn't work when an editor has a lock on the file first.

## Dockerized Development

API hosted in a docker container. This works by building and watching the dist folder and launching the ui docker container with the dist folder mounted. The ui container serves using nginx and has a proxy to the API running another container on the same network.

If you haven't already you should read the Getting Started Guide before continuing.

After starting the docker stack for KAG stats open a terminal to the ui folder of the kagstats repo.
Run: `npm run build -- --watch --poll 250`.

Then open a second terminal to the ui directory and run:
`docker run --rm --network=kagstats_default -v $(pwd)/dist/ui:/var/www/html -p 4200:80 gcr.io/kagstats/ui`

This command starts the ui container, attaches it to the existing docker network and mounts the dist directory.

Open up your browser to http://localhost:4200.

**NOTE:** Unlike the other method this will not automatically reload your browser page on rebuild.

## Angular Development Tips

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

Make sure to do this from the appropriate directory

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).