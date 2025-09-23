# Notifiction Service

### My main repo in [GitLab](https://gitlab.com/fmkv/fmkv-backend/-/tree/main/notificationService?ref_type=heads)

Used for email notifications to the user

```bash
go version go1.24.6
```

## Routes

POST `/emails/send`

```bash
{
  to: "badabum@badabum.com",
  "subject": "Ebat!",
  "body": "its revolution Vlados..."
}
```

Response

```bash
CODE: 200
{
  "status": "ok",
  "message": "email send",
}
```

## How to start?

```bash
go mod tidy

touch .env

echo "
PORT=..
SMTP_USER=...
SMTP_PASS="..."
SMTP_HOST=...
SMTP_PORT=...
" > .env

make run

or

go run ./cmd/notification/main.go
