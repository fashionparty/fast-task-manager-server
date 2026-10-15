# Fast Task Manager — Server

REST API for the [Fast Task Manager](https://github.com/fashionparty/fast-task-manager-app) Android app.

## Running

```bash
docker compose up --build
```

## API

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/auth/register` | Create an account |
| POST | `/api/auth/login` | Log in, returns a JWT token |
| GET | `/api/tasks` | List tasks |
| POST | `/api/tasks` | Create a task |
| GET | `/api/tasks/{id}` | Get a task |
| PUT | `/api/tasks/{id}` | Update a task (only the fields sent) |
| DELETE | `/api/tasks/{id}` | Delete a task |
| GET | `/health` | Health check |
