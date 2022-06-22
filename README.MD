# Stuning Couscous

## Requirment & Considerations

- Having Docker and Docker Compose installed
- To test the APIs you need to seed the database with seeder CLI, because the database is empty.

## Start

To start the app in development mode run:

```sh
cd PROGECT_DIRECTORY
```

```sh
docker compose up
```

Info:

- RestAPI: 127.0.0.1:8080
- POSTGRESQL IP: 127.0.0.1:5432
- POSTGRESQL PASSWORD=gHteuivwdvkew4wt
- POSTGRESQL USERNAME=postgres
- POSTGRESQL DATABASE=stunning-couscous

---

## Run Test

To run tests:

```sh
# start the database
docker compose up db
```

```sh
cd PROGECT_DIRECTORY
```

```sh
./test.sh
```

---

## Build

Build script:

- Creates a ready to use docker image for production
- Creates binary seeder cli
- Creates partnerapi binary file to use

To run the script:

```sh
cd PROGECT_DIRECTORY

```

```sh
./build.sh
```

---

## Seeder

The seeder CLI generates random partners and writes them to the database.
Make sure the database is up and running then run the seeder:

```sh
cd PROGECT_DIRECTORY
```

```sh
./cmd/seeder
```

---

## Postman collectin

[Postman colletion](https://github.com/pysf/stunning-couscous/blob/9de451c63117db38f7199fcbbb82df765881b2f8/stunning-couscous.postman_collection.json) in the project main directory is provided to use and call the APIs.

## Call the API using Curl

### GetPartner

- GET /api/partner/:id

example:

```sh
curl --location --request GET 'http://127.0.0.1:8080/api/partner/1'

```

---

### FindBestMatch

- GET /api/search/partner/best-match

#### Qeury String Parameters:

- material e.g. wood
- latitude e.g. 52.49304178513987
- longitude e.g. 13.40497255
- phone e.g. 1123456789
- square e.g. 70

example:

```sh
curl --location --request GET 'http://127.0.0.1:8080/api/search/partner/best-match?material=carpet&latitude=52.49304178513987&longitude=13.40497255&phone=1123456789&square=70'
```
