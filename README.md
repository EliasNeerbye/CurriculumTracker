# Curriculum Tracker Backend

A comprehensive backend API for tracking and managing learning curricula, built with Go 1.24, PostgreSQL, and JWT authentication.

## Features

- **User Authentication**: Secure registration and login with JWT tokens and Argon2 password hashing
- **Curriculum Management**: Create, read, update, and delete curricula
- **Project Organization**: Organize projects within curricula with dependencies and types
- **Progress Tracking**: Track completion status and percentage for each project
- **Note Taking**: Add notes and reflections to projects with different types
- **Time Tracking**: Log time spent on projects with analytics
- **Analytics**: Get comprehensive stats on learning progress and time investment

## Project Structure

```md
curriculum-tracker/
├── main.go                     # Application entry point
├── go.mod                      # Go module dependencies
├── .env.example                # Environment variables template
├── README.md                   # Project documentation
├── API_DOCUMENTATION.md        # Complete API documentation
├── config/
│   └── config.go              # Configuration management
├── database/
│   ├── connection.go          # Database connection
│   └── migrations.go          # Database schema migrations
├── models/
│   ├── user.go               # User data models
│   ├── curriculum.go         # Curriculum data models
│   ├── project.go            # Project data models
│   ├── progress.go           # Progress tracking models
│   ├── note.go               # Note data models
│   └── time_entry.go         # Time tracking models
├── utils/
│   ├── password.go           # Argon2 password hashing
│   ├── jwt.go                # JWT token utilities
│   └── response.go           # HTTP response helpers
├── middleware/
│   ├── auth.go               # JWT authentication middleware
│   ├── cors.go               # CORS middleware
│   └── logging.go            # Request logging middleware
├── services/
│   ├── auth.go               # Authentication business logic
│   ├── curriculum.go         # Curriculum business logic
│   ├── project.go            # Project business logic
│   ├── progress.go           # Progress tracking logic
│   ├── note.go               # Note management logic
│   └── analytics.go          # Analytics and time tracking
├── handlers/
│   ├── auth.go               # Authentication HTTP handlers
│   ├── curriculum.go         # Curriculum HTTP handlers
│   ├── project.go            # Project HTTP handlers
│   ├── progress.go           # Progress HTTP handlers
│   ├── note.go               # Note HTTP handlers
│   └── analytics.go          # Analytics HTTP handlers
└── routes/
    └── routes.go             # HTTP route configuration
```

## Quick Start

### Prerequisites

- Go 1.24 or later
- PostgreSQL 12 or later

### Installation

1. **Clone and setup**

   ```bash
   git clone <repository-url>
   cd curriculum-tracker
   go mod download
   ```

2. **Database setup**

   ```sql
   CREATE DATABASE curriculum_tracker;
   ```

3. **Environment configuration**

   ```bash
   cp .env.example .env
   ```

   Edit `.env` with your database credentials:

   ```env
   DATABASE_URL=postgres://username:password@localhost:5432/curriculum_tracker?sslmode=disable
   JWT_SECRET=your-super-secret-jwt-key-change-in-production
   PORT=8080
   ```

4. **Run the application**

   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080` and automatically run database migrations.

## API Usage

### Authentication

1. **Register a user**

   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"password123","name":"John Doe"}'
   ```

2. **Login**

   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"user@example.com","password":"password123"}'
   ```

3. **Use the returned JWT token in subsequent requests**

   ```bash
   curl -H "Authorization: Bearer <your-token>" \
     http://localhost:8080/api/v1/curricula
   ```

### Example Workflow

1. Create a curriculum
2. Add projects to the curriculum
3. Track progress on projects
4. Add notes and reflections
5. Log time spent
6. View analytics

See `API_DOCUMENTATION.md` for complete endpoint documentation.

## Architecture Highlights

### Modern Go Practices

- **Go 1.24**: Uses latest Go features and standard library
- **Minimal Dependencies**: Only essential third-party packages
- **Modular Design**: Clean separation of concerns
- **Error Handling**: Comprehensive error handling with proper HTTP status codes

### Security

- **Argon2**: Industry-standard password hashing (upgraded from bcrypt)
- **JWT**: Secure authentication with configurable expiration
- **SQL Injection Prevention**: Parameterized queries throughout
- **CORS**: Configurable cross-origin resource sharing
- **Input Validation**: Comprehensive request validation

### Database Design

- **PostgreSQL**: Robust relational database with ACID compliance
- **Foreign Keys**: Proper relational integrity
- **Indexes**: Optimized query performance
- **Migrations**: Automatic schema management
- **JSON Support**: PostgreSQL arrays for flexible data storage

### Performance

- **Connection Pooling**: Optimized database connections
- **Prepared Statements**: Efficient query execution
- **Middleware**: Efficient request processing pipeline
- **Structured Logging**: Performance monitoring ready

## Key Models

### Curriculum Structure

Based on the tree learning methodology with:

- **Roots**: Foundation projects
- **Base**: Core engineering skills  
- **Branches**: Domain specialization (Lower, Middle, Upper)
- **Flowers**: Capstone projects

### Project Types

- `root` - Foundation building blocks
- `base` - Core skill development
- `lowerBranch` - Basic specialization projects
- `middleBranch` - Intermediate specialization
- `upperBranch` - Advanced specialization
- `flowerMilestone` - Capstone achievements

### Progress Tracking

- Status: `not_started`, `in_progress`, `completed`, `on_hold`, `abandoned`
- Completion percentage (0-100)
- Start and completion timestamps
- Automatic progress calculation

## Development

### Code Quality

- No comments in code (self-documenting)
- Consistent error handling patterns
- Comprehensive input validation
- RESTful API design
- Proper HTTP status codes

### Testing

Run the health check to verify the server:

```bash
curl http://localhost:8080/health
```

### Deployment

For production deployment:

1. Set strong JWT secret
2. Use connection pooling for database
3. Enable HTTPS
4. Configure proper CORS origins
5. Set up monitoring and logging
6. Use environment-specific configurations

## Frontend Integration

This backend is designed to support rich frontend applications with:

- Complete CRUD operations for all entities
- Comprehensive analytics endpoints
- Real-time progress tracking
- Flexible note-taking system
- Time tracking with detailed breakdowns

See `API_DOCUMENTATION.md` for detailed frontend integration guidance.

## Technology Stack

- **Language**: Go 1.24
- **Database**: PostgreSQL 15+
- **Authentication**: JWT with Argon2 password hashing
- **HTTP Router**: Gorilla Mux
- **Environment**: godotenv
- **Database Driver**: pq (PostgreSQL driver)

## License

This project is built for educational purposes and curriculum tracking.
