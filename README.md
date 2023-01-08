# Status Checker API - REST API for monitoring system statuses

## Endpoints
### POST /systemstatus
* Add new systems to the database
### GET /systemstatus
* Get information about all systems 
### DELETE /systemstatus/:id
* Delete system from the monitor _(not yet implemented)_
### GET /systemstatus/:id
* Get information about a specific system. This also calls the endpoint of that system and fills out the status fields of the model, returning updated information.

## Alerting
Alerting can be performed through email or by sending a post request to some endpoint (like a slack webhook integration)

## General information
* https endpoints will get their certificate expiration time checked against the setting of the system on how many days before cert expiration a warning should be created.

* The current setup uses postgresql as the backend database. It can easily be stood up using docker:
```bash
    docker run -p 5432:5432 -e POSTGRES_USER=status -e POSTGRES_DB=status -e POSTGRES_PASSWORD=muchs3cretw0w postgres:latest
```
_Note:_ This command will not persist the data over container restarts. Please see the official documentation https://hub.docker.com/_/postgres on how to tell postgres where to save data. You also need to mount a volume on that filsystem location, https://docs.docker.com/storage/volumes/