# Curriculum Tracker API Documentation

## Base URL

```md
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```md
Authorization: Bearer <your-jwt-token>
```

## Error Response Format

```json
{
  "success": false,
  "error": "Error message"
}
```

## Success Response Format

```json
{
  "success": true,
  "data": {
    // Response data
  }
}
```

---

## Authentication Endpoints

### Register User

**POST** `/auth/register`

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe"
}
```

**Response (201):**

```json
{
  "success": true,
  "data": {
    "token": "jwt-token-here",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2025-05-30T10:00:00Z",
      "updated_at": "2025-05-30T10:00:00Z"
    }
  }
}
```

### Login

**POST** `/auth/login`

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

**Response (200):**

```json
{
  "success": true,
  "data": {
    "token": "jwt-token-here",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2025-05-30T10:00:00Z",
      "updated_at": "2025-05-30T10:00:00Z"
    }
  }
}
```

### Get Current User

**GET** `/auth/me`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2025-05-30T10:00:00Z",
    "updated_at": "2025-05-30T10:00:00Z"
  }
}
```

---

## Curriculum Endpoints

### Create Curriculum

**POST** `/curricula`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
  "name": "C Programming Mastery",
  "description": "Complete C programming curriculum from basics to advanced"
}
```

**Response (201):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "C Programming Mastery",
    "description": "Complete C programming curriculum from basics to advanced",
    "created_at": "2025-05-30T10:00:00Z",
    "updated_at": "2025-05-30T10:00:00Z",
    "projects": []
  }
}
```

### Get All Curricula

**GET** `/curricula`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "name": "C Programming Mastery",
      "description": "Complete C programming curriculum from basics to advanced",
      "created_at": "2025-05-30T10:00:00Z",
      "updated_at": "2025-05-30T10:00:00Z",
      "total_projects": 15,
      "completed_projects": 5,
      "total_time_spent": 1200
    }
  ]
}
```

### Get Single Curriculum

**GET** `/curricula/{id}`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "C Programming Mastery",
    "description": "Complete C programming curriculum from basics to advanced",
    "created_at": "2025-05-30T10:00:00Z",
    "updated_at": "2025-05-30T10:00:00Z",
    "projects": [
      {
        "id": 1,
        "curriculum_id": 1,
        "identifier": "R1",
        "name": "Hello World Variations",
        "description": "Print patterns, ASCII art, formatted output",
        "learning_objectives": ["printf", "escape sequences", "basic I/O"],
        "estimated_time": "1 hour",
        "prerequisites": [],
        "project_type": "root",
        "position_order": 1,
        "created_at": "2025-05-30T10:00:00Z",
        "updated_at": "2025-05-30T10:00:00Z",
        "progress": {
          "id": 1,
          "user_id": 1,
          "project_id": 1,
          "status": "completed",
          "completion_percentage": 100,
          "started_at": "2025-05-30T09:00:00Z",
          "completed_at": "2025-05-30T10:00:00Z",
          "created_at": "2025-05-30T09:00:00Z",
          "updated_at": "2025-05-30T10:00:00Z"
        }
      }
    ]
  }
}
```

### Update Curriculum

**PUT** `/curricula/{id}`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
  "name": "Updated C Programming Mastery",
  "description": "Updated description"
}
```

**Response (200):** Same as create response with updated data.

### Delete Curriculum

**DELETE** `/curricula/{id}`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": {
    "message": "Curriculum deleted successfully"
  }
}
```

---

## Project Endpoints

### Create Project

**POST** `/curricula/{curriculumId}/projects`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
  "name": "Hello World Variations",
  "description": "Print patterns, ASCII art, formatted output",
  "learning_objectives": ["printf", "escape sequences", "basic I/O"],
  "estimated_time": "1 hour",
  "prerequisites": [],
  "project_type": "root",
  "position_order": 1
}
```

**Valid project_type values:**

- `root` - Foundation projects (generates identifiers: R1, R2, R3...)
- `rootTest` - Root test project (generates identifier: RT, only one allowed per curriculum)
- `base` - Core skill projects (generates identifiers: B1, B2, B3...)
- `baseTest` - Base test project (generates identifier: BT, only one allowed per curriculum)
- `lowerBranch` - Lower branch projects (generates identifiers: LB1, LB2, LB3...)
- `middleBranch` - Middle branch projects (generates identifiers: MB1, MB2, MB3...)
- `upperBranch` - Upper branch projects (generates identifiers: UB1, UB2, UB3...)
- `flowerMilestone` - Capstone projects (generates identifiers: F1, F2, F3...)

**Note:** The `identifier` field is automatically generated based on the project type and is not included in the request body.

**Response (201):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "curriculum_id": 1,
    "identifier": "R1",
    "name": "Hello World Variations",
    "description": "Print patterns, ASCII art, formatted output",
    "learning_objectives": ["printf", "escape sequences", "basic I/O"],
    "estimated_time": "1 hour",
    "prerequisites": [],
    "project_type": "root",
    "position_order": 1,
    "created_at": "2025-05-30T10:00:00Z",
    "updated_at": "2025-05-30T10:00:00Z"
  }
}
```

### Get Project

**GET** `/projects/{id}`

**Headers:** `Authorization: Bearer <token>`

**Response (200):** Same as create response.

### Update Project

**PUT** `/projects/{id}`

**Headers:** `Authorization: Bearer <token>`

**Request Body:** Same as create project (excluding identifier).

**Response (200):** Same as create response with updated data.

**Note:** The identifier cannot be changed through updates as it's automatically managed.

### Delete Project

**DELETE** `/projects/{id}`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": {
    "message": "Project deleted successfully"
  }
}
```

**Note:** Projects cannot be deleted if other projects depend on them as prerequisites.

### Get Project Notes

**GET** `/projects/{id}/notes`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "project_id": 1,
      "title": "Initial thoughts",
      "content": "This project was straightforward but taught me the basics of printf formatting.",
      "note_type": "reflection",
      "created_at": "2025-05-30T10:00:00Z",
      "updated_at": "2025-05-30T10:00:00Z"
    }
  ]
}
```

---

## Progress Endpoints

### Update Progress

**PUT** `/projects/{projectId}/progress`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
  "status": "in_progress",
  "completion_percentage": 75
}
```

**Valid status values:**

- `not_started` - Project not yet started (completion_percentage automatically set to 0)
- `in_progress` - Project currently being worked on
- `completed` - Project finished (completion_percentage automatically set to 100)
- `on_hold` - Project temporarily paused
- `abandoned` - Project discontinued

**Response (200):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "project_id": 1,
    "status": "in_progress",
    "completion_percentage": 75,
    "started_at": "2025-05-30T09:00:00Z",
    "completed_at": null,
    "created_at": "2025-05-30T09:00:00Z",
    "updated_at": "2025-05-30T10:00:00Z"
  }
}
```

### Get Project Progress

**GET** `/projects/{projectId}/progress`

**Headers:** `Authorization: Bearer <token>`

**Response (200):** Same as update progress response.

**Note:** If no progress exists, returns default progress with `not_started` status.

### Get Curriculum Progress

**GET** `/curricula/{curriculumId}/progress`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "project_id": 1,
      "status": "completed",
      "completion_percentage": 100,
      "started_at": "2025-05-30T09:00:00Z",
      "completed_at": "2025-05-30T10:00:00Z",
      "created_at": "2025-05-30T09:00:00Z",
      "updated_at": "2025-05-30T10:00:00Z"
    }
  ]
}
```

---

## Note Endpoints

### Create Note

**POST** `/projects/{projectId}/notes`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
  "title": "Initial thoughts",
  "content": "This project was straightforward but taught me the basics of printf formatting.",
  "note_type": "reflection"
}
```

**Valid note_type values:**

- `note` - General notes
- `reflection` - Learning reflections
- `learning` - Learning insights
- `question` - Questions and clarifications

**Response (201):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "project_id": 1,
    "title": "Initial thoughts",
    "content": "This project was straightforward but taught me the basics of printf formatting.",
    "note_type": "reflection",
    "created_at": "2025-05-30T10:00:00Z",
    "updated_at": "2025-05-30T10:00:00Z"
  }
}
```

### Get Note

**GET** `/notes/{id}`

**Headers:** `Authorization: Bearer <token>`

**Response (200):** Same as create note response.

### Update Note

**PUT** `/notes/{id}`

**Headers:** `Authorization: Bearer <token>`

**Request Body:** Same as create note.

**Response (200):** Same as create note response with updated data.

### Delete Note

**DELETE** `/notes/{id}`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": {
    "message": "Note deleted successfully"
  }
}
```

---

## Analytics Endpoints

### Create Time Entry

**POST** `/time-entries`

**Headers:** `Authorization: Bearer <token>`

**Request Body:**

```json
{
  "project_id": 1,
  "minutes": 60,
  "description": "Worked on hello world variations",
  "date": "2025-05-30"
}
```

**Response (201):**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "project_id": 1,
    "minutes": 60,
    "description": "Worked on hello world variations",
    "date": "2025-05-30T00:00:00Z",
    "created_at": "2025-05-30T10:00:00Z"
  }
}
```

### Get Project Time Entries

**GET** `/projects/{projectId}/time-entries`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "project_id": 1,
      "minutes": 60,
      "description": "Worked on hello world variations",
      "date": "2025-05-30T00:00:00Z",
      "created_at": "2025-05-30T10:00:00Z"
    }
  ]
}
```

### Get Curriculum Time Stats

**GET** `/curricula/{curriculumId}/time-stats`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": {
    "total_minutes": 300,
    "daily_breakdown": {
      "2025-05-30": 60,
      "2025-05-29": 120,
      "2025-05-28": 120
    },
    "project_breakdown": {
      "Hello World Variations": 60,
      "Number Guesser": 120,
      "Temperature Converter": 120
    },
    "weekly_average": 85.7
  }
}
```

### Get User Overall Stats

**GET** `/analytics/user-stats`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**

```json
{
  "success": true,
  "data": {
    "total_curricula": 2,
    "total_projects": 25,
    "completed_projects": 8,
    "in_progress_projects": 3,
    "total_time_minutes": 1200,
    "total_time_hours": 20,
    "total_notes": 15,
    "completion_rate": 32.0
  }
}
```

---

## Health Check

### Health Check Endpoint

**GET** `/health`

**Response (200):**

```md
OK
```

---

## Project Identifier System

The system automatically generates identifiers based on project type:

| Project Type | Identifier Pattern | Example | Notes |
|--------------|-------------------|---------|-------|
| `root` | R1, R2, R3... | R1, R2, R3 | Foundation projects |
| `rootTest` | RT | RT | Only one allowed per curriculum |
| `base` | B1, B2, B3... | B1, B2, B3 | Core skill projects |
| `baseTest` | BT | BT | Only one allowed per curriculum |
| `lowerBranch` | LB1, LB2, LB3... | LB1, LB2, LB3 | Lower specialization |
| `middleBranch` | MB1, MB2, MB3... | MB1, MB2, MB3 | Middle specialization |
| `upperBranch` | UB1, UB2, UB3... | UB1, UB2, UB3 | Upper specialization |
| `flowerMilestone` | F1, F2, F3... | F1, F2, F3 | Capstone projects |

**Important Notes:**

- Identifiers are automatically generated and cannot be manually set
- Test projects (`rootTest`, `baseTest`) can only have one instance per curriculum
- Each project type maintains its own counter within a curriculum
- Identifiers are used in the `prerequisites` array to reference other projects

---

## Setup Instructions

### Prerequisites

- Go 1.24 or later
- PostgreSQL 12 or later

### Environment Variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgres://username:password@localhost:5432/curriculum_tracker?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-in-production
PORT=8080
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
ENVIRONMENT=development
```

### Database Setup

1. Create a PostgreSQL database named `curriculum_tracker`
2. The application will automatically run migrations on startup

### Running the Application

```bash
go mod download
go run main.go
```

The server will start on `http://localhost:8080`

---

## Security Features

- **Argon2 Password Hashing**: Industry-standard password security
- **JWT Authentication**: Secure token-based authentication with 24-hour expiration
- **SQL Injection Prevention**: All queries use parameterized statements
- **User Isolation**: All data access is scoped to the authenticated user
- **CORS Configuration**: Configurable cross-origin resource sharing

---

## Frontend Integration Notes

### Authentication Flow

1. Register or login to get JWT token
2. Store token securely (localStorage/sessionStorage)
3. Include token in all API requests as `Authorization: Bearer <token>`
4. Handle token expiration (24 hours by default)

### Project Creation Workflow

1. Create a curriculum first
2. Add projects to the curriculum with appropriate types
3. Use generated identifiers in prerequisites arrays
4. Track progress and add notes as needed
5. Log time entries for analytics

### Key Integration Points

- **Automatic Identifiers**: Don't include `identifier` in project creation requests
- **Test Project Limits**: Handle errors when trying to create duplicate test projects
- **Prerequisites Validation**: Use existing project identifiers in prerequisites arrays
- **Progress States**: Implement UI states for all progress statuses
- **Real-time Updates**: Consider implementing optimistic updates for better UX

### State Management Recommendations

- Cache curriculum/project data for better performance
- Implement optimistic updates for progress tracking
- Handle prerequisite dependencies in project visualization
- Show identifier patterns to help users understand the system
