# Status Checker API - REST API for monitoring system statuses

## Endpoints
### GET /healthz
* Returns 200 healthy if the app is up.
### POST /systemstatus
* Add new systems to the database
### GET /systemstatus
* Get information about all systems 
### DELETE /systemstatus/:id
* Delete system from the monitor 
### GET /systemstatus/:id
* Get information about a specific system. This also calls the endpoint of that system and fills out the status fields of the model, returning updated information.
### Form POST /clientcert
* Post a form with FormFile, Form `name` and Form `password`. 
### GET /clientcerts
* Get a list of Cert id and name to fill up dropdowns when selecting clientcert for a system
### DELETE /clientcert/:id
* Delete a clientcertificate from the database. No control is made if it is used.
## Alerting
Alerting can be performed through email or by sending a post request to some endpoint (like a slack webhook integration)

## General information
* https endpoints will get their certificate expiration time checked against the setting of the system on how many days before cert expiration a warning should be created.

* The current setup uses postgresql as the backend database. It can easily be stood up using docker:
```bash
    docker run -p 5432:5432 -e POSTGRES_USER=status -e POSTGRES_DB=status -e POSTGRES_PASSWORD=muchs3cretw0w postgres:latest
```
_Note:_ This command will not persist the data over container restarts. Please see the official documentation https://hub.docker.com/_/postgres on how to tell postgres where to save data. You also need to mount a volume on that filesystem location, https://docs.docker.com/storage/volumes/

* An env variable with the encryption key is expected. Please generate a key with OpenSSL or similar, `openssl rand -hex 16`. The password for the p12 is encrypted at rest and can only be decrypted in the business layer in memory when calling an mTLS endpoint. There is no api endpoint to retrieve the actual cert or password once it has been saved. 