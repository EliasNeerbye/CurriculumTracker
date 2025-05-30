# Curriculum Tracker Backend

A comprehensive backend API for tracking and managing learning curricula, built with Go 1.24, PostgreSQL, and JWT authentication.

## Features

- **User Authentication**: Secure registration and login with JWT tokens and Argon2 password hashing
- **Curriculum Management**: Create, read, update, and delete curricula
- **Smart Project Organization**: Organize projects within curricula with type-based identifiers and dependencies
- **Progress Tracking**: Track completion status and percentage for each project with automatic state transitions
- **Note Taking**: Add notes and reflections to projects with different types
- **Time Tracking**: Log time spent on projects with comprehensive analytics
- **Analytics**: Get detailed stats on learning progress and time investment

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
   ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
   ENVIRONMENT=development
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

1. **Create a curriculum**

   ```bash
   curl -X POST http://localhost:8080/api/v1/curricula \
     -H "Authorization: Bearer <your-token>" \
     -H "Content-Type: application/json" \
     -d '{"name":"C Programming Mastery","description":"Complete C programming curriculum"}'
   ```

2. **Add projects with automatic identifier generation**

   ```bash
   # Create a root project (will get identifier "R1")
   curl -X POST http://localhost:8080/api/v1/curricula/1/projects \
     -H "Authorization: Bearer <your-token>" \
     -H "Content-Type: application/json" \
     -d '{"name":"Hello World","description":"Basic output","project_type":"root","position_order":1}'
   
   # Create a base project (will get identifier "B1")
   curl -X POST http://localhost:8080/api/v1/curricula/1/projects \
     -H "Authorization: Bearer <your-token>" \
     -H "Content-Type: application/json" \
     -d '{"name":"Variables","description":"Variable declaration","project_type":"base","prerequisites":["R1"],"position_order":2}'
   ```

3. **Track progress on projects**

   ```bash
   curl -X PUT http://localhost:8080/api/v1/projects/1/progress \
     -H "Authorization: Bearer <your-token>" \
     -H "Content-Type: application/json" \
     -d '{"status":"completed","completion_percentage":100}'
   ```

4. **Add notes and reflections**

   ```bash
   curl -X POST http://localhost:8080/api/v1/projects/1/notes \
     -H "Authorization: Bearer <your-token>" \
     -H "Content-Type: application/json" \
     -d '{"title":"Learning Notes","content":"This was straightforward","note_type":"reflection"}'
   ```

5. **Log time spent**

   ```bash
   curl -X POST http://localhost:8080/api/v1/time-entries \
     -H "Authorization: Bearer <your-token>" \
     -H "Content-Type: application/json" \
     -d '{"project_id":1,"minutes":60,"description":"Completed hello world","date":"2025-05-30"}'
   ```

6. **View analytics**

   ```bash
   curl -H "Authorization: Bearer <your-token>" \
     http://localhost:8080/api/v1/analytics/user-stats
   ```

See `API_DOCUMENTATION.md` for complete endpoint documentation.

## Architecture Highlights

### Modern Go Practices

- **Go 1.24**: Uses latest Go features and standard library
- **Minimal Dependencies**: Only essential third-party packages
- **Modular Design**: Clean separation of concerns
- **Error Handling**: Comprehensive error handling with proper HTTP status codes

### Security

- **Argon2**: Industry-standard password hashing (upgraded from bcrypt)
- **JWT**: Secure authentication with configurable expiration (24 hours default)
- **SQL Injection Prevention**: Parameterized queries throughout
- **CORS**: Configurable cross-origin resource sharing
- **Input Validation**: Comprehensive request validation

### Database Design

- **PostgreSQL**: Robust relational database with ACID compliance
- **Foreign Keys**: Proper relational integrity with CASCADE deletes
- **Indexes**: Optimized query performance on common operations
- **Migrations**: Automatic schema management on startup
- **JSON Support**: PostgreSQL arrays for flexible data storage

### Performance

- **Connection Pooling**: Optimized database connections (25 max open, 5 max idle)
- **Prepared Statements**: Efficient query execution
- **Middleware**: Efficient request processing pipeline
- **Structured Logging**: Performance monitoring ready

## Key Models

### Curriculum Structure

Based on the tree learning methodology with:

- **Roots (R1, R2, R3...)**: Foundation building blocks
- **Root Tests (RT)**: Assessment of foundation skills
- **Base (B1, B2, B3...)**: Core engineering skills  
- **Base Tests (BT)**: Assessment of core skills
- **Branches**: Domain specialization
  - **Lower Branch (LB1, LB2...)**: Basic specialization projects
  - **Middle Branch (MB1, MB2...)**: Intermediate specialization
  - **Upper Branch (UB1, UB2...)**: Advanced specialization
- **Flowers (F1, F2, F3...)**: Capstone achievements

### Smart Project Identifier System

The system automatically generates type-based identifiers:

| Project Type | Identifier Pattern | Example | Restrictions |
|--------------|-------------------|---------|--------------|
| `root` | R1, R2, R3... | R1, R2, R3 | Foundation projects |
| `rootTest` | RT | RT | Only one per curriculum |
| `base` | B1, B2, B3... | B1, B2, B3 | Core skill projects |
| `baseTest` | BT | BT | Only one per curriculum |
| `lowerBranch` | LB1, LB2, LB3... | LB1, LB2, LB3 | Lower specialization |
| `middleBranch` | MB1, MB2, MB3... | MB1, MB2, MB3 | Middle specialization |
| `upperBranch` | UB1, UB2, UB3... | UB1, UB2, UB3 | Upper specialization |
| `flowerMilestone` | F1, F2, F3... | F1, F2, F3 | Capstone projects |

**Key Features:**

- Identifiers automatically generated based on project type
- Each project type maintains its own counter within a curriculum
- Test projects (`rootTest`, `baseTest`) limited to one per curriculum
- Prerequisites reference other projects using their identifiers

### Progress Tracking

- **Automatic State Management**: Status transitions automatically update timestamps
- **Smart Timestamps**:
  - `started_at` set when moving from `not_started` to any active status
  - `completed_at` set when status becomes `completed`
- **Status Types**: `not_started`, `in_progress`, `completed`, `on_hold`, `abandoned`
- **Completion Percentage**: Automatically managed based on status
- **Prerequisite Validation**: Ensures learning path integrity

### Advanced Features

- **Dependency Management**: Projects can specify prerequisites using identifiers
- **Automatic Progress Calculation**: Curriculum completion rates calculated in real-time
- **Time Analytics**: Detailed breakdowns by project, daily activity, and trends
- **Note Categories**: Different note types for various learning activities
- **Data Integrity**: Comprehensive validation and error handling

## Development

### Code Quality

- **Self-Documenting Code**: No comments needed, clear naming conventions
- **Consistent Error Handling**: Structured error responses throughout
- **Comprehensive Input Validation**: All user inputs validated
- **RESTful API Design**: Clean, predictable endpoint structure
- **Proper HTTP Status Codes**: Meaningful response codes

### Testing

Run the health check to verify the server:

```bash
curl http://localhost:8080/health
```

### Deployment

For production deployment:

1. **Security Configuration**

   ```env
   JWT_SECRET=generate-a-strong-random-secret-key
   ENVIRONMENT=production
   ALLOWED_ORIGINS=https://yourdomain.com
   ```

2. **Database Optimization**
   - Use connection pooling
   - Configure proper indexes
   - Set up regular backups

3. **Infrastructure**
   - Enable HTTPS/TLS
   - Configure reverse proxy (nginx/Apache)
   - Set up monitoring and logging
   - Use environment-specific configurations

## Frontend Integration

This backend is designed to support rich frontend applications with:

### Key Integration Points

- **Automatic Identifier Generation**: Don't send `identifier` in project creation requests
- **Type-Based Organization**: Use project types to organize UI components
- **Real-Time Progress**: Complete CRUD operations for all entities
- **Comprehensive Analytics**: Rich data for dashboards and visualizations
- **Flexible Note System**: Support different note types and categories

### Recommended Frontend Architecture

```md
src/
├── components/
│   ├── auth/              # Authentication components
│   ├── curriculum/        # Curriculum management
│   ├── project/          # Project components with type-aware UI
│   ├── progress/         # Progress tracking and visualization
│   ├── notes/           # Note-taking interface
│   └── analytics/       # Time tracking and statistics
├── hooks/
│   ├── useAuth.js       # Authentication state management
│   ├── useCurriculum.js # Curriculum data management
│   └── useProgress.js   # Progress tracking
├── services/
│   ├── api.js          # API client with automatic token handling
│   └── auth.js         # Authentication service
└── utils/
    ├── identifiers.js  # Helper functions for identifier patterns
    └── validation.js   # Frontend validation helpers
```

### State Management Recommendations

- **Optimistic Updates**: Update UI immediately, sync with server
- **Real-Time Sync**: Consider WebSocket integration for live updates
- **Offline Support**: Cache data for offline project work
- **Type-Aware Components**: Create different UI components for different project types

## Technology Stack

- **Language**: Go 1.24 with latest features
- **Database**: PostgreSQL 15+ with advanced features
- **Authentication**: JWT with Argon2 password hashing
- **HTTP Router**: Gorilla Mux for flexible routing
- **Configuration**: godotenv for environment management
- **Database Driver**: pq (Pure Go PostgreSQL driver)

## Performance Characteristics

- **Response Times**: Sub-millisecond for cached queries
- **Concurrency**: Handles hundreds of concurrent connections
- **Memory Usage**: Efficient memory management with connection pooling
- **Scalability**: Horizontal scaling ready with stateless design

## License

This project is built for educational purposes and curriculum tracking. Feel free to adapt and extend for your learning needs.

---

## Contributing

When contributing to this project:

1. Follow Go best practices and formatting
2. Write self-documenting code with clear naming
3. Add appropriate error handling
4. Update documentation for API changes
5. Test all endpoints before submitting changes

For questions or suggestions, please open an issue or submit a pull request.
