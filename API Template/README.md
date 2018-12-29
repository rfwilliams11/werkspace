# RTR API Template #

This package contains the tools necessary to create an API off of MongoDB for Centene's Real Time Repositories. This
repository can be both run as a stand alone application, or imported into another application as a package.

### Versions ###

### Dependency Management  ###
* The project uses `dep` to explicitly call out all dependencies
* If using this project as a package it is **STRONGLY** encouraged that the importing project also utilizes `dep`.

### Configuration ###
* Default configuration is kept in `./config/base.yml`
* Environment variables will be used to pass **all environment specific configurations**. They should be upper and snake case. See the `config` struct in `./env/env.go` for specifics
* Configuration priority is chosen in order of environment variable then base file
* During local testing environment variables will be loaded from `.env` files in the base project directory.
    * These files ARE NOT included in source control as they typically contain secrets.
    * The `.env` files are in Docker environment file syntax where `VAR=VALUE`.
* Certain environment variables are required 

### Endpoint Routes ###
#### As a Package
The following endpoints come standard from this library when imported as a package:  

**GET /version** Returns the latest git tag, branch name, sha1, and build time of the deployed binary  

**GET /health** Returns the status of the api. Currently only checks the database connection  

#### As a Stand Alone
The following endpoints can be executed if this application is run as it's own application:  

**GET /mgo/{collection}** Executes a basic Mongo DB query using the provided query parameters, as well as `select` for projection

**GET /mgo/{collection}/count** Returns the count of all documents in the collection provided  

#### Adding Routes
New routes should be added to the stand alone router in `main.go` unless there is a need to share this functionality 
across all Centene RTR APIs.

### Logging ###  
* Log levels can be `debug`, `info`, `warn`, `error`, `fatal` and `panic`.
    * The default is `info`.
* Log level can be set in the environment using `LOG_LEVEL`.
* Logrus can be hooked into `ElasticSearch`, `Logstash` and `Graylog`.
* The log facility is set to the Git project name. 

### Monitoring ###  
* The Kibana APM agent is added to the router middleware for endpoint monitoring.
* A sanity check is performed that required the environment variable `ELASTIC_APM_SERVER_URL` to be set to add the agent
* It is recommended but not currently required to add the following environment variables for the APM Agent
    * `ELASTIC_APM_SERVICE_VERSION`: The version of the application.
    * `ELASTIC_APM_ENVIRONMENT`: The environment the application is running in. 
    * `ELASTIC_APM_SERVICE_NAME`: The name of the application.
        * **NOTE:** It is best practice to include the application version and environment in the service. The recommended format is `service-name-service-environment:service-version` **ex.** rtr-api-dev:1.0
    
### SplitPea ###  
* KEY1, KEY2, SPLITPEA_USERNAME and SPLITPEA_PASSWORD are passed as environment variables (MONGO_KEY1, MONGO_KEY2, SPLITPEA_USERNAME and SPLITPEA_PASSWORD)
* Currently SplitPea is only being utilized to connect to Mongo RTR
 
### Makefile Commands ###
**default** Runs `help` which will print out a detailed list of all possible commands.
  
**deps** Create the build container and pull down necessary dependencies.

**test** Test the Go code.

**test-coverage** Test the Go code and generate coverage results.

**clean** Remove all build artifacts and generated files.

**swagger** Generate the swagger documentation.

**server-docker** Build the application container.

**start-docker** Starts the container image.

**stop-docker** Stops the container, and removes the image from Docker.

**deploy-docker** Builds and pushes the API image up to the rancher repository.
