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
      "created_at": "2025-05-27T10:00:00Z",
      "updated_at": "2025-05-27T10:00:00Z"
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
      "created_at": "2025-05-27T10:00:00Z",
      "updated_at": "2025-05-27T10:00:00Z"
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
    "created_at": "2025-05-27T10:00:00Z",
    "updated_at": "2025-05-27T10:00:00Z"
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
    "created_at": "2025-05-27T10:00:00Z",
    "updated_at": "2025-05-27T10:00:00Z"
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
      "created_at": "2025-05-27T10:00:00Z",
      "updated_at": "2025-05-27T10:00:00Z",
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
    "created_at": "2025-05-27T10:00:00Z",
    "updated_at": "2025-05-27T10:00:00Z",
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
        "created_at": "2025-05-27T10:00:00Z",
        "updated_at": "2025-05-27T10:00:00Z",
        "progress": {
          "id": 1,
          "user_id": 1,
          "project_id": 1,
          "status": "completed",
          "completion_percentage": 100,
          "started_at": "2025-05-27T09:00:00Z",
          "completed_at": "2025-05-27T10:00:00Z",
          "created_at": "2025-05-27T09:00:00Z",
          "updated_at": "2025-05-27T10:00:00Z"
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
  "identifier": "R1",
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

- `root`
- `rootTest`
- `base`
- `baseTest`
- `lowerBranch`
- `middleBranch`
- `upperBranch`
- `flowerMilestone`

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
    "created_at": "2025-05-27T10:00:00Z",
    "updated_at": "2025-05-27T10:00:00Z"
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

**Request Body:** Same as create project.

**Response (200):** Same as create response with updated data.

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
      "created_at": "2025-05-27T10:00:00Z",
      "updated_at": "2025-05-27T10:00:00Z"
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

- `not_started`
- `in_progress`
- `completed`
- `on_hold`
- `abandoned`

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
    "started_at": "2025-05-27T09:00:00Z",
    "completed_at": null,
    "created_at": "2025-05-27T09:00:00Z",
    "updated_at": "2025-05-27T10:00:00Z"
  }
}
```

### Get Project Progress

**GET** `/projects/{projectId}/progress`

**Headers:** `Authorization: Bearer <token>`

**Response (200):** Same as update progress response.

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
      "started_at": "2025-05-27T09:00:00Z",
      "completed_at": "2025-05-27T10:00:00Z",
      "created_at": "2025-05-27T09:00:00Z",
      "updated_at": "2025-05-27T10:00:00Z"
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

- `note`
- `reflection`
- `learning`
- `question`

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
    "created_at": "2025-05-27T10:00:00Z",
    "updated_at": "2025-05-27T10:00:00Z"
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
  "date": "2025-05-27"
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
    "date": "2025-05-27T00:00:00Z",
    "created_at": "2025-05-27T10:00:00Z"
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
      "date": "2025-05-27T00:00:00Z",
      "created_at": "2025-05-27T10:00:00Z"
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
      "2025-05-27": 60,
      "2025-05-26": 120,
      "2025-05-25": 120
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

### Health Checker

**GET** `/health`

**Response (200):**

```md
OK
```

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

## Frontend Integration Notes

### Authentication Flow

1. Register or login to get JWT token
2. Store token securely (localStorage/sessionStorage)
3. Include token in all API requests
4. Handle token expiration (24 hours by default)

### Recommended Frontend Structure

```md
components/
├── auth/
│   ├── LoginForm.jsx
│   ├── RegisterForm.jsx
│   └── ProtectedRoute.jsx
├── curriculum/
│   ├── CurriculumList.jsx
│   ├── CurriculumDetail.jsx
│   └── CreateCurriculum.jsx
├── project/
│   ├── ProjectCard.jsx
│   ├── ProjectDetail.jsx
│   └── ProgressTracker.jsx
├── notes/
│   ├── NotesList.jsx
│   ├── NoteEditor.jsx
│   └── NoteViewer.jsx
└── analytics/
    ├── TimeTracker.jsx
    ├── StatsOverview.jsx
    └── ProgressCharts.jsx
```

### Key Features to Implement

- Curriculum tree visualization
- Progress tracking with visual indicators
- Time tracking with charts
- Note-taking with rich text editor
- Project dependencies visualization
- Analytics dashboard with charts
- Search and filtering
- Export functionality

### State Management Recommendations

- Use React Context or Redux for global state
- Cache curriculum/project data
- Implement optimistic updates for better UX
- Handle offline scenarios
