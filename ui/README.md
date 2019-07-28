# Ui

This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 7.3.5.

## Development server

First make sure all the other services are running and the database/api are available.
`docker-compose up -d`

Then start then start angular building and watching the files from the ui folder
`npm run build -- --watch`
Make sure for this step to finish and be ready before launching the docker container

Then from the ui folder start an nginx container, mount the development angular files and attach it to the network with the api
`docker run --rm --network=kagstats_default -v $(pwd)/dist/ui:/var/www/html -p 4200:80 gcr.io/kagstats/ui`


## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory. Use the `--prod` flag for a production build.

## Running unit tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via [Protractor](http://www.protractortest.org/).

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI README](https://github.com/angular/angular-cli/blob/master/README.md).
