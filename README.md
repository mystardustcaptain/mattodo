## Installation

Fetch the project and install the dependencies:

```bash
git clone https://github.com/mystardustcaptain/mattodo.git
cd mattodo
go mod download
```

Please note that, if you plan to run the service locally, `gcc` dependencies will be required.
Highly recommended to work with Docker.


## Configuration

Environment variables are used for the service to successfuly operate.

Refer to the `.env_sample` file, replace and fill in valid values, and renaming it to `.env` file.

Make sure it is available before running the application.

`.env_sample` content is as below

```
GOOGLE_CLIENT_ID=yourGoogleClientId
GOOGLE_CLIENT_SECRET=yourGoogleClientSecret
GOOGLE_REDIRECT_URL=http://localhost:9003/auth/callback?provider=google

FACEBOOK_CLIENT_ID=123456789
FACEBOOK_CLIENT_SECRET=123456789
FACEBOOK_REDIRECT_URL=http://localhost:9003/auth/callback?provider=facebook

GITHUB_CLIENT_ID=123456789
GITHUB_CLIENT_SECRET=123456789
GITHUB_REDIRECT_URL=http://localhost:9003/auth/callback?provider=github

SIGNING_KEY=1234

SERVICE_PORT=:9003
DB_TYPE=sqlite3
DB_PATH=./mainDB.sqlite3
```




## Running the App

### Docker Image

The application can be run by using the Docker image published to Docker Hub.

```bash
PS> docker run -p 9003:9003 --env-file .env -e DOCKER_ENV_SET=true eiggub/mattodo
CMD> docker run -p 9003:9003 --env-file .env -e DOCKER_ENV_SET=true eiggub/mattodo
```


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




## Testing the App (Unit Test)

Currently only one test file is created under `./pkg/model` in `todo_test.go`.

### From Terminal
To run the test from command prompt, from the root level of the project, execute:

```bash
go test ./pkg/model/...
```

### From Visual Studio > Testing
From Visual Studio > Testing, the comprehensive list of all available tests in the projects are listed.

You can choose to run all or selectively a few tests, and see the results all visually from the panel.



## Building the App

### From Visual Studio
```bash
go build
```
`mattodo.exe` will be created. Together with the `.env` and `db` files, the service will ru successfully.

(GCC dependencies required)

### Docker
At the root level of the project next to Dockerfile:
```bash
docker build -t mattodo . 
```

Or to Build + Run using Docker Compose with `docker-compose.yml`:
```bash
docker compose up
```

## API Documentation
After making sure the service is up, start using the Todo features with APIs available below.

### Login
There are currently 2 Authentication Provider supported: Google, Github. (Facebook in progress)

You are recommended to visit the link via Browser, as it will direct you to the HTML login interface at the provider platform.

[http://localhost:9003/auth/login?provider=google](http://localhost:9003/auth/login?provider=google)

[http://localhost:9003/auth/login?provider=github](http://localhost:9003/auth/login?provider=github)

[http://localhost:9003/auth/login?provider=facebook  (coming soon)](http://localhost:9003/auth/login?provider=facebook)

A JWT Token will be provided when login is successful.

Save it down, replace `my_JWT_token` with the real token received in the subsequent calls below.
.

### Get All Todo Items
To retrieve all the Todo items associated to you (by Email address):
```bash
curl -H "Authorization: Bearer my_JWT_token" http://localhost:9003/todo
```


If successful, you will receive all the Todo items belong to you:
```json
[
    {
        "id": 1,
        "user_id": 1,
        "title": "Item 1",
        "completed": false,
        "created_at": "2023-12-03T13:08:04Z",
        "updated_at": "2023-12-03T13:10:56Z"
    },
    {
        "id": 14,
        "user_id": 1,
        "title": "Item 2",
        "completed": false,
        "created_at": "2023-12-05T07:49:50.535983103Z",
        "updated_at": "2023-12-05T07:49:50.535983162Z"
    }
]
```
### Create Todo Item
To create a Todo item associated to you:
```bash
curl -X POST -H "Authorization: Bearer my_JWT_token" http://localhost:9003/todo -H "Content-Type: application/json" --data "{'title': 'My Birthday', 'completed': false}"
```

If successful, you will received the Todo item created:
```json
{
    "id": 15,
    "user_id": 1,
    "title": "My Birthday",
    "completed": false,
    "created_at": "2023-12-05T08:48:03.59392094Z",
    "updated_at": "2023-12-05T08:48:03.59392094Z"
}
```

### Delete Todo Item
To delete a Todo item associated to you, of given id:
```bash
curl -X DELETE -H "Authorization: Bearer my_JWT_token" http://localhost:9003/todo/{id}
```
Example:
```bash
curl -X DELETE -H "Authorization: Bearer my_JWT_token" http://localhost:9003/todo/2
```

If successful, you will receve `StatusCode: 204` with no other content.

If the Todo item you're trying to delete does not belong to you, the delete will fail.

### Mark Todo Item as Completed
To mark a Todo item associated to you, of given id, to Completed:
```bash
curl -X PUT -H "Authorization: Bearer my_JWT_token" http://localhost:9003/todo/{id}/complete
```
Example:
```bash
curl -X PUT -H "Authorization: Bearer my_JWT_token" http://localhost:9003/todo/2/complete
```
If successful, you will received the complete Todo item updated:
```json
{
    "id": 15,
    "user_id": 1,
    "title": "My Birthday",
    "completed": true,
    "created_at": "2023-12-05T08:48:03.59392094Z",
    "updated_at": "2023-12-10T10:00:00.39291084Z"
}
```