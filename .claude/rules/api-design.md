---
name: api-design
paths:
  - "handler/**/*.go"
  - "main.go"
---

# API Design Rules

- All endpoints return JSON with consistent structure: `{"data": ..., "error": ...}`
- Use proper HTTP status codes (200, 201, 400, 404, 500)
- Validate all input at the handler layer before passing to service
- Handlers must not contain business logic — delegate to service layer
- Route naming: lowercase, hyphen-separated, RESTful (`/league/table`, `/matches/:id`)
