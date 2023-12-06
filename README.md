## Installation

Clone the project and install dependencies:

```bash
git clone https://github.com/mystardustcaptain/mattodo.git
cd mattodo
go mod download
```




## Configuration

Set up environment variables for successful operation:

- Copy `.env_sample` to `.env` and fill in the appropriate values.
- Ensure `.env` is available before running the application.

Example `.env_sample`:

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
DB_TYPE=sqlite
DB_PATH=./mainDB.db
```




## Running the App

### Using Docker (Recommended)

Run the app using the published Docker image:


```bash
docker run -p 9003:9003 --env-file .env -e DOCKER_ENV_SET=true eiggub/mattodo
```


### Using Visual Studio (Not Recommended)

- Run the service from Visual Studio.
- Ensure `.env` is at the root level next to `main.go`.

```bash
go run .
```


### Using Docker Compose

For development, use Docker Compose:

```bash
docker compose up
```
`docker-compose.yml` includes all necessary configurations.




## Testing the App (Unit Test)

Run unit tests located in `./pkg/model/todo_test.go`:

```bash
go test ./pkg/model/...
```

Visual Studio's Testing panel offers a visual interface for running and reviewing tests. (Recommended)





## Building the App

### Using  Visual Studio
Build the executable:
```bash
go build
```
Ensure `.env` and database files are alongside the generated `mattodo.exe`.


### Using Docker
Build with Docker:
```bash
docker build -t mattodo . 
```

Or use Docker Compose:
```bash
docker compose up
```

## API Documentation
Access Todo features via the following APIs after starting the service.



### Authentication
Supported providers: Google, Github (Facebook coming soon).

Visit the following links in a browser for OAuth flow:



- Google: [http://localhost:9003/auth/login?provider=google](http://localhost:9003/auth/login?provider=google)
- Github: [http://localhost:9003/auth/login?provider=github](http://localhost:9003/auth/login?provider=github)
- Facebook (Soon): [http://localhost:9003/auth/login?provider=facebook](http://localhost:9003/auth/login?provider=facebook)

A JWT token is provided upon successful login.



### Get All Todo Items
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:9003/todo
```

### Create Todo Item
```bash
curl -X POST -H "Authorization: Bearer YOUR_JWT_TOKEN" -H "Content-Type: application/json" --data "{'title': 'New Task', 'completed': false}" http://localhost:9003/todo
```


### Delete Todo Item
```bash
curl -X DELETE -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:9003/todo/{id}
```


### Mark Todo Item as Completed
```bash
curl -X PUT -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:9003/todo/{id}/complete
```

Replace `YOUR_JWT_TOKEN` and `{id}` with actual values.

