# Form2Mail

A simple Go service that handles HTML contact form submissions and sends email notifications to both the site owner and the customer.

## Features

- POST endpoint for contact form submissions
- Sends email to recipient (site owner) with form details
- Sends confirmation email to customer
- Supports both JSON and form-urlencoded data
- CORS enabled for cross-origin requests
- HTML formatted emails
- Clean, structured codebase following Go best practices

## Project Structure

```
form2mail/
├── cmd/
│   └── server/          # Application entry point
│       └── main.go
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── email/           # Email sending functionality
│   └── handler/         # HTTP request handlers
├── .github/
│   └── workflows/       # GitHub Actions workflows
│       └── docker-build.yml
├── .env.example         # Example environment variables
├── .dockerignore
├── .gitignore
├── AGENTS.md            # Guidelines for AI coding agents
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Multi-stage Docker build
├── go.mod
└── README.md
```

## Setup

1. Clone the repository

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Edit `.env` and configure your SMTP settings:
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=your-email@gmail.com
RECIPIENT_EMAIL=recipient@example.com
SERVER_PORT=8080
CORS_ORIGIN=*
```

**CORS Configuration:**
- Use `*` to allow all origins (default)
- Use a specific domain like `https://yourdomain.com` to restrict access
- Use comma-separated values for multiple domains (not currently supported, use `*` or single domain)

### Gmail Setup

If using Gmail, you'll need to create an App Password:
1. Enable 2-Factor Authentication on your Google account
2. Go to https://myaccount.google.com/apppasswords
3. Generate a new app password
4. Use this password in the `SMTP_PASSWORD` field

## Running

### Development:
```bash
export SMTP_USER="your-email@gmail.com"
export SMTP_PASSWORD="your-app-password"
export RECIPIENT_EMAIL="recipient@example.com"
export FROM_EMAIL="your-email@gmail.com"
go run cmd/server/main.go
```

### Using Docker:

**Pull from GitHub Container Registry:**
```bash
docker pull ghcr.io/aveseli/form2mail:latest
```

**Run the container:**
```bash
docker run -d \
  -p 8080:8080 \
  -e SMTP_HOST="smtp.gmail.com" \
  -e SMTP_PORT="587" \
  -e SMTP_USER="your-email@gmail.com" \
  -e SMTP_PASSWORD="your-app-password" \
  -e FROM_EMAIL="your-email@gmail.com" \
  -e RECIPIENT_EMAIL="recipient@example.com" \
  -e CORS_ORIGIN="*" \
  --name form2mail \
  ghcr.io/aveseli/form2mail:latest
```

**Or using docker-compose:**
```yaml
version: '3.8'
services:
  form2mail:
    image: ghcr.io/aveseli/form2mail:latest
    ports:
      - "8080:8080"
    environment:
      - SMTP_HOST=smtp.gmail.com
      - SMTP_PORT=587
      - SMTP_USER=your-email@gmail.com
      - SMTP_PASSWORD=your-app-password
      - FROM_EMAIL=your-email@gmail.com
      - RECIPIENT_EMAIL=recipient@example.com
      - CORS_ORIGIN=*
    restart: unless-stopped
```

**Build locally:**
```bash
docker build -t form2mail .
docker run -p 8080:8080 --env-file .env form2mail
```


## API Usage

### Endpoint
```
POST /contact
```

### Request Format

**JSON:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "subject": "Question about services",
  "message": "I would like to know more about your services."
}
```

**Form Data:**
```
name=John+Doe&email=john@example.com&subject=Question&message=Hello
```

### Response

**Success (200):**
```json
{
  "status": "success",
  "message": "Your message has been sent successfully"
}
```

**Error (4xx/5xx):**
```
Error message in plain text
```

## HTML Form Example

```html
<form action="http://localhost:8080/contact" method="POST">
  <input type="text" name="name" placeholder="Your Name" required>
  <input type="email" name="email" placeholder="Your Email" required>
  <input type="text" name="subject" placeholder="Subject">
  <textarea name="message" placeholder="Your Message" required></textarea>
  <button type="submit">Send</button>
</form>
```

## JavaScript Fetch Example

```javascript
fetch('http://localhost:8080/contact', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: 'John Doe',
    email: 'john@example.com',
    subject: 'Question',
    message: 'Hello, I have a question.'
  })
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));
```

## Building

### Local Build:
```bash
go build -o form2mail cmd/server/main.go
./form2mail
```

Or build and run in one step:
```bash
go build -o form2mail cmd/server/main.go && ./form2mail
```

### Docker Build:
```bash
docker build -t form2mail .
```

## GitHub Container Registry

This project automatically builds and pushes Docker images to GitHub Container Registry when you push to the main branch or create tags.

**Image naming:**
- Latest: `ghcr.io/aveseli/form2mail:latest`
- Tagged: `ghcr.io/aveseli/form2mail:v1.0.0`
- Branch: `ghcr.io/aveseli/form2mail:main`

**To use the image, you need to:**
1. Make the package public in GitHub (Settings > Packages > form2mail > Package settings > Change visibility)
2. Or authenticate: `echo $GITHUB_TOKEN | docker login ghcr.io -u aveseli --password-stdin`

## Development

### Prerequisites
- Go 1.23 or higher
- Docker (optional, for containerized deployment)

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

### Code Formatting
```bash
# Format all code (run before committing)
go fmt ./...

# Run linter
go vet ./...

# Tidy dependencies
go mod tidy
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and formatting
5. Submit a pull request

For detailed coding guidelines, see [AGENTS.md](AGENTS.md).

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SMTP_HOST` | No | `smtp.gmail.com` | SMTP server hostname |
| `SMTP_PORT` | No | `587` | SMTP server port |
| `SMTP_USER` | Yes | - | SMTP username/email |
| `SMTP_PASSWORD` | Yes | - | SMTP password or app password |
| `FROM_EMAIL` | Yes | - | Email address to send from |
| `RECIPIENT_EMAIL` | Yes | - | Email address to receive contact forms |
| `SERVER_PORT` | No | `8080` | HTTP server port |
| `CORS_ORIGIN` | No | `*` | CORS allowed origin (`*` for all, or specific domain) |

## License

MIT
