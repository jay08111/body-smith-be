# body-smith-be

## Structure

```text
body-smith-be/
├── cmd/server/main.go
├── internal/
│   ├── config/
│   ├── db/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   └── service/
├── migrations/
├── .env.example
└── README.md
```

## Run

```bash
cp .env.example .env
go mod tidy
go run ./cmd/server
```

## Migrations

The server runs embedded `golang-migrate` migrations automatically on startup.

Example CLI usage:

```bash
migrate -path migrations -database "mysql://root:password@tcp(127.0.0.1:3306)/body_smith" up
migrate -path migrations -database "mysql://root:password@tcp(127.0.0.1:3306)/body_smith" down 1
```

## Sample curl

Login:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@bodysmith.com",
    "password": "ChangeMe123!"
  }'
```

Get current admin:

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <TOKEN>"
```

Create post:

```bash
curl -X POST http://localhost:8080/api/v1/admin/posts \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "How to Build a Sustainable Fitness Routine",
    "content": "Consistency beats intensity over time.",
    "thumbnail": "https://example.com/post.jpg",
    "meta_title": "Sustainable Fitness Routine",
    "meta_description": "Practical advice for building a sustainable fitness routine."
  }'
```

List admin posts:

```bash
curl "http://localhost:8080/api/v1/admin/posts?page=1&per_page=10" \
  -H "Authorization: Bearer <TOKEN>"
```

Update post:

```bash
curl -X PUT http://localhost:8080/api/v1/admin/posts/1 \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "How to Build a Sustainable Fitness Routine",
    "content": "Updated content.",
    "thumbnail": "https://example.com/post-updated.jpg",
    "meta_title": "Updated SEO title",
    "meta_description": "Updated SEO description."
  }'
```

Delete post:

```bash
curl -X DELETE http://localhost:8080/api/v1/admin/posts/1 \
  -H "Authorization: Bearer <TOKEN>"
```

List public posts:

```bash
curl "http://localhost:8080/api/v1/posts?page=1&per_page=10"
```

Get public post:

```bash
curl http://localhost:8080/api/v1/posts/how-to-build-a-sustainable-fitness-routine
```
