## Installation

Fetch the project and install the dependencies:

```bash
git clone https://github.com/mystardustcaptain/mattodo.git
cd mattodo
go mod download
```




## Configuration

Environment variables are used for the service to successfuly operate.

Refer to the `.env_sample` file, replace and fill in valid values, and renaming it to `.env` file.

Make sure it is available before running the application.

`.env_sample` content is as below

```
GOOGLE_CLIENT_ID=yourGoogleClientId
GOOGLE_CLIENT_SECRET=yourGoogleClientSecret

FACEBOOK_CLIENT_ID=123456789
FACEBOOK_CLIENT_SECRET=123456789

GITHUB_CLIENT_ID=123456789
GITHUB_CLIENT_SECRET=123456789

SIGNING_KEY=1234

SERVICE_PORT=:9003
DB_TYPE=sqlite3
DB_PATH=./mainDB.sqlite3
```



## Running the App

### Docker Image

The application can be run by using the Docker image published to Docker Hub.

```bash
PS> docker run -p 9003:9003 -v ${PWD}/.env:/.env eiggub/mattodo
CMD> docker run -p 9003:9003 -v $(pwd)/.env:/.env eiggub/mattodo
```
`-v ${PWD}/.env:/.env` mounts the .env file from the current directory `(${PWD}/.env)` to the root directory in the container `(/.env)`. 

### Visual Studio (Not Recommended)

You can run the service directly from Visual Studio project locally.

Make sure .env is at the root level next to `main.go`.

```bash
go run .
```
However, since the project uses SQLite3, GCC libraries will be required.

### Docker Compose

It is highly recommended to use Docker Compose from Visual Studio to continue development.

Make sure `Dockerfile`, `docker-compose.yml` and `.env` are all available in the same level next to `main.go`.

```bash
docker compose up
```
`docker-compose.yml` has all the parameters configured, including loading the `.env` file.

Hence the much simpler command.




## Testing the App

Explain how to run the automated tests for this system.

```bash
go test ./...
```

## Interface Documentation

Details about the interfaces used in the project and how to use them.

```
```