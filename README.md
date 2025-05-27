# Curriculum Tracker

A modern curriculum management and progress tracking application built with Go backend and WebAssembly frontend.

## Features

- 📚 **Curriculum Management**: Create and organize learning curricula
- 🎯 **Project Tracking**: Track individual projects with detailed progress
- 📝 **Notes & Reflections**: Take notes and reflect on your learning journey
- ⏱️ **Time Tracking**: Log time spent on projects with detailed analytics
- 📊 **Analytics Dashboard**: Visualize your learning progress with comprehensive statistics
- 🔐 **Secure Authentication**: JWT-based authentication with Argon2 password hashing
- 🌊 **Interactive UI**: Modern, responsive interface with fun colors and smooth animations

## Tech Stack

### Backend

- **Go 1.24** - Latest LTS version
- **PostgreSQL** - Robust relational database
- **Chi Router** - Lightweight, fast HTTP router
- **JWT** - Secure token-based authentication
- **Argon2** - State-of-the-art password hashing

### Frontend

- **WebAssembly (WASM)** - Go compiled to run in the browser
- **Pure CSS** - Custom responsive design with modern animations
- **No JavaScript frameworks** - Lightweight and fast

## Quick Start

### Prerequisites

- Go 1.24 or later
- PostgreSQL 13 or later
- Make (for build automation)

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd curriculum-tracker
   ```

2. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your database credentials and JWT secret
   ```

3. **Set up the database**

   ```bash
   # Create database
   createdb curriculum_tracker
   
   # Run migrations
   psql curriculum_tracker < database/schema.sql
   ```

4. **Install dependencies**

   ```bash
   go mod tidy
   ```

5. **Build and run**

   ```bash
   make run
   ```

The application will be available at `http://localhost:8080`

## Development

### Build Commands

```bash
# Build everything (server + WebAssembly)
make build

# Build only the server
make server

# Build only WebAssembly frontend
make wasm

# Run the application
make run

# Clean build artifacts
make clean

# Set up database
make setup
```

### Project Structure

```
curriculum-tracker/
├── cmd/
│   ├── server/          # HTTP server entry point
│   └── wasm/            # WebAssembly frontend entry point
├── internal/
│   ├── auth/            # Authentication & JWT handling
│   ├── database/        # Database operations
│   ├── handlers/        # HTTP handlers
│   └── models/          # Data models
├── database/
│   └── schema.sql       # Database schema
├── web/                 # Static web assets
│   ├── index.html       # Main HTML file
│   ├── styles.css       # CSS styles
│   ├── main.wasm        # Generated WebAssembly (build artifact)
│   └── wasm_exec.js     # Go WebAssembly runtime (build artifact)
└── Makefile             # Build automation
```

## API Endpoints

### Authentication

- `POST /api/register` - User registration
- `POST /api/login` - User login
- `GET /api/profile` - Get user profile (authenticated)

### Curricula

- `GET /api/curricula` - List user's curricula
- `POST /api/curricula` - Create new curriculum
- `GET /api/curricula/{id}` - Get curriculum details
- `PUT /api/curricula/{id}` - Update curriculum
- `DELETE /api/curricula/{id}` - Delete curriculum
- `GET /api/curricula/{id}/projects` - List curriculum projects
- `POST /api/curricula/{id}/projects` - Create new project

### Projects

- `GET /api/projects/{id}` - Get project details
- `PUT /api/projects/{id}` - Update project
- `DELETE /api/projects/{id}` - Delete project
- `GET /api/projects/{id}/progress` - Get project progress
- `PUT /api/projects/{id}/progress` - Update project progress

### Notes & Reflections

- `GET /api/projects/{id}/notes` - List project notes
- `POST /api/projects/{id}/notes` - Create note
- `PUT /api/notes/{id}` - Update note
- `DELETE /api/notes/{id}` - Delete note
- `GET /api/projects/{id}/reflections` - List project reflections
- `POST /api/projects/{id}/reflections` - Create reflection
- `PUT /api/reflections/{id}` - Update reflection
- `DELETE /api/reflections/{id}` - Delete reflection

### Time Tracking

- `GET /api/projects/{id}/time-entries` - List time entries
- `POST /api/projects/{id}/time-entries` - Create time entry
- `GET /api/analytics` - Get user analytics

## Database Schema

The application uses a comprehensive PostgreSQL schema with the following main tables:

- **users** - User accounts with secure authentication
- **curricula** - Learning curricula owned by users
- **projects** - Individual projects within curricula
- **project_progress** - Track progress status and time spent
- **notes** - Project notes for documentation
- **reflections** - Learning reflections and insights
- **time_entries** - Detailed time tracking logs

## Curriculum Structure

The application supports a hierarchical curriculum structure inspired by the "Maplewood Tree" methodology:

- 🌱 **Root Projects** - Foundation and basic concepts
- 🌲 **Base Projects** - Core skills and fundamentals
- 🌿 **Lower Branch** - Short practical projects
- 🍃 **Middle Branch** - Medium-length projects
- 🌸 **Upper Branch** - Advanced, longer projects
- 🌺 **Flower Milestones** - Major capstone projects
- 📝 **Tests** - Knowledge validation projects

## Docker Deployment

### Build Docker image

```bash
make docker-build
```

### Run with Docker

```bash
# Copy environment file
cp .env.example .env
# Edit .env with your configuration

# Run container
make docker-run
```

### Docker Compose (recommended)

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://postgres:password@db:5432/curriculum_tracker?sslmode=disable
      - JWT_SECRET=your-super-secret-jwt-key
    depends_on:
      - db

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=curriculum_tracker
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database/schema.sql:/docker-entrypoint-initdb.d/schema.sql

volumes:
  postgres_data:
```

## Configuration

Environment variables:

- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT token signing (change in production!)
- `PORT` - Server port (default: 8080)

## Security Features

- **Argon2 Password Hashing** - Industry-standard password security
- **JWT Authentication** - Stateless, secure token-based auth
- **CORS Configuration** - Proper cross-origin resource sharing
- **SQL Injection Prevention** - Parameterized queries throughout
- **Input Validation** - Comprehensive request validation
- **Secure Headers** - Security-focused HTTP headers

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

MIT License - see LICENSE file for details

## Architecture Notes

This application demonstrates several modern patterns:

- **WebAssembly Frontend** - Leveraging Go's WebAssembly support for type-safe frontend code
- **Clean Architecture** - Separation of concerns with clear module boundaries
- **RESTful API Design** - Following REST principles for predictable endpoints
- **Progressive Enhancement** - Works without JavaScript, enhanced with WebAssembly
- **Responsive Design** - Mobile-first CSS with modern features
- **Security by Design** - Built-in security considerations throughout

The choice to use WebAssembly with Go provides several advantages:

- **Type Safety** - Share types between frontend and backend
- **Performance** - Near-native execution speed
- **Developer Experience** - Single language for full-stack development
- **Maintainability** - Easier to maintain with shared code patterns

## Troubleshooting

### WebAssembly Issues

If the WebAssembly isn't loading:

1. Ensure you've built with `make wasm`
2. Check that `web/main.wasm` and `web/wasm_exec.js` exist
3. Verify your browser supports WebAssembly
4. Check browser console for errors

### Database Connection Issues

If you can't connect to PostgreSQL:

1. Verify PostgreSQL is running: `pg_isready`
2. Check your DATABASE_URL in `.env`
3. Ensure the database exists: `createdb curriculum_tracker`
4. Run the schema: `psql curriculum_tracker < database/schema.sql`

### Build Issues

If the build fails:

1. Verify Go version: `go version` (should be 1.24+)
2. Clean and rebuild: `make clean && make build`
3. Check for missing dependencies: `go mod tidy`
