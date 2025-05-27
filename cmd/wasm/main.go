package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"syscall/js"
	"time"
)

var (
	document             js.Value
	window               js.Value
	authToken            string
	currentUser          *User
	currentView          string
	selectedCurriculumID string
	selectedProjectID    string
)

func main() {
	window = js.Global()
	document = window.Get("document")

	// Set up router
	setupRouter()

	// Check for existing auth token
	if token := getLocalStorage("authToken"); token != "" {
		authToken = token
		loadUserProfile()
	} else {
		showLoginView()
	}

	// Keep the Go program running
	select {}
}

func setupRouter() {
	// Handle browser back/forward
	window.Call("addEventListener", "popstate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handleRoute()
		return nil
	}))

	// Initial route
	handleRoute()
}

func handleRoute() {
	path := window.Get("location").Get("pathname").String()

	switch {
	case path == "/" || path == "/login":
		if authToken == "" {
			showLoginView()
		} else {
			showDashboardView()
		}
	case path == "/register":
		showRegisterView()
	case path == "/dashboard":
		showDashboardView()
	case path == "/curricula":
		showCurriculaView()
	case strings.HasPrefix(path, "/curriculum/"):
		parts := strings.Split(path, "/")
		if len(parts) >= 3 {
			showCurriculumDetailView(parts[2])
		}
	case strings.HasPrefix(path, "/project/"):
		parts := strings.Split(path, "/")
		if len(parts) >= 3 {
			showProjectDetailView(parts[2])
		}
	case path == "/analytics":
		showAnalyticsView()
	default:
		showDashboardView()
	}
}

func navigate(path string) {
	window.Get("history").Call("pushState", nil, "", path)
	handleRoute()
}

func setLocalStorage(key, value string) {
	window.Get("localStorage").Call("setItem", key, value)
}

func getLocalStorage(key string) string {
	val := window.Get("localStorage").Call("getItem", key)
	if val.IsNull() || val.IsUndefined() {
		return ""
	}
	return val.String()
}

func removeLocalStorage(key string) {
	window.Get("localStorage").Call("removeItem", key)
}

func apiRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, "/api"+endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(respBody.Bytes(), &errResp); err == nil {
			return nil, fmt.Errorf("%s", errResp.Error)
		}
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return respBody.Bytes(), nil
}

func loadUserProfile() {
	data, err := apiRequest("GET", "/profile", nil)
	if err != nil {
		authToken = ""
		removeLocalStorage("authToken")
		showLoginView()
		return
	}

	if err := json.Unmarshal(data, &currentUser); err != nil {
		showError("Failed to load user profile")
		return
	}

	updateNavbar()
}

func updateNavbar() {
	navbar := document.Call("getElementById", "navbar")
	if navbar.IsNull() {
		return
	}

	if currentUser != nil {
		navbar.Set("innerHTML", fmt.Sprintf(`
			<div class="nav-container">
				<div class="nav-brand">
					<span class="nav-logo">🌲</span>
					<span class="nav-title">Curriculum Tracker</span>
				</div>
				<nav class="nav-links">
					<a href="/dashboard" onclick="event.preventDefault(); navigate('/dashboard')">Dashboard</a>
					<a href="/curricula" onclick="event.preventDefault(); navigate('/curricula')">Curricula</a>
					<a href="/analytics" onclick="event.preventDefault(); navigate('/analytics')">Analytics</a>
				</nav>
				<div class="nav-user">
					<span class="user-name">%s</span>
					<button class="btn-logout" onclick="logout()">Logout</button>
				</div>
			</div>
		`, currentUser.Name))
	} else {
		navbar.Set("innerHTML", `
			<div class="nav-container">
				<div class="nav-brand">
					<span class="nav-logo">🌲</span>
					<span class="nav-title">Curriculum Tracker</span>
				</div>
			</div>
		`)
	}
}

func showLoginView() {
	currentView = "login"
	updateNavbar()

	content := document.Call("getElementById", "content")
	content.Set("innerHTML", `
		<div class="auth-container">
			<div class="auth-card">
				<h1>Welcome Back!</h1>
				<p class="auth-subtitle">Log in to continue your learning journey</p>
				<form id="loginForm" class="auth-form">
					<div class="form-group">
						<label for="email">Email</label>
						<input type="email" id="email" name="email" required>
					</div>
					<div class="form-group">
						<label for="password">Password</label>
						<input type="password" id="password" name="password" required>
					</div>
					<button type="submit" class="btn btn-primary">Log In</button>
				</form>
				<p class="auth-link">
					Don't have an account? 
					<a href="/register" onclick="event.preventDefault(); navigate('/register')">Sign up</a>
				</p>
			</div>
		</div>
	`)

	// Attach event listener
	form := document.Call("getElementById", "loginForm")
	form.Call("addEventListener", "submit", js.FuncOf(handleLogin))
}

func showRegisterView() {
	currentView = "register"
	updateNavbar()

	content := document.Call("getElementById", "content")
	content.Set("innerHTML", `
		<div class="auth-container">
			<div class="auth-card">
				<h1>Start Your Journey</h1>
				<p class="auth-subtitle">Create an account to track your learning progress</p>
				<form id="registerForm" class="auth-form">
					<div class="form-group">
						<label for="name">Name</label>
						<input type="text" id="name" name="name" required>
					</div>
					<div class="form-group">
						<label for="email">Email</label>
						<input type="email" id="email" name="email" required>
					</div>
					<div class="form-group">
						<label for="password">Password</label>
						<input type="password" id="password" name="password" required>
					</div>
					<button type="submit" class="btn btn-primary">Sign Up</button>
				</form>
				<p class="auth-link">
					Already have an account? 
					<a href="/login" onclick="event.preventDefault(); navigate('/login')">Log in</a>
				</p>
			</div>
		</div>
	`)

	// Attach event listener
	form := document.Call("getElementById", "registerForm")
	form.Call("addEventListener", "submit", js.FuncOf(handleRegister))
}

func showDashboardView() {
	if authToken == "" {
		navigate("/login")
		return
	}

	currentView = "dashboard"
	updateNavbar()

	content := document.Call("getElementById", "content")
	content.Set("innerHTML", `
		<div class="dashboard-container">
			<h1>Welcome back, <span id="userName"></span>!</h1>
			<div class="dashboard-grid">
				<div class="dashboard-card" id="statsCard">
					<h2>📊 Quick Stats</h2>
					<div class="stats-loading">Loading...</div>
				</div>
				<div class="dashboard-card" id="recentCard">
					<h2>🕐 Recent Activity</h2>
					<div class="recent-loading">Loading...</div>
				</div>
				<div class="dashboard-card" id="progressCard">
					<h2>🎯 Current Progress</h2>
					<div class="progress-loading">Loading...</div>
				</div>
			</div>
			<div class="dashboard-actions">
				<button class="btn btn-primary" onclick="navigate('/curricula')">
					<span class="btn-icon">📚</span> View Curricula
				</button>
				<button class="btn btn-secondary" onclick="navigate('/analytics')">
					<span class="btn-icon">📈</span> Detailed Analytics
				</button>
			</div>
		</div>
	`)

	if currentUser != nil {
		document.Call("getElementById", "userName").Set("textContent", currentUser.Name)
	}

	loadDashboardData()
}

func loadDashboardData() {
	// Load analytics data
	go func() {
		data, err := apiRequest("GET", "/analytics", nil)
		if err != nil {
			showError("Failed to load analytics: " + err.Error())
			return
		}

		var analytics AnalyticsResponse
		if err := json.Unmarshal(data, &analytics); err != nil {
			showError("Failed to parse analytics")
			return
		}

		// Update stats card
		statsCard := document.Call("getElementById", "statsCard")
		if !statsCard.IsNull() {
			statsCard.Set("innerHTML", fmt.Sprintf(`
				<h2>📊 Quick Stats</h2>
				<div class="stats-grid">
					<div class="stat-item">
						<div class="stat-value">%d</div>
						<div class="stat-label">Total Projects</div>
					</div>
					<div class="stat-item">
						<div class="stat-value">%d</div>
						<div class="stat-label">Completed</div>
					</div>
					<div class="stat-item">
						<div class="stat-value">%d</div>
						<div class="stat-label">In Progress</div>
					</div>
					<div class="stat-item">
						<div class="stat-value">%dh %dm</div>
						<div class="stat-label">Time Invested</div>
					</div>
				</div>
			`, analytics.TotalProjects, analytics.CompletedProjects,
				analytics.InProgressProjects, analytics.TotalTimeSpent/60, analytics.TotalTimeSpent%60))
		}

		// Update recent activity
		recentCard := document.Call("getElementById", "recentCard")
		if !recentCard.IsNull() && len(analytics.RecentActivity) > 0 {
			var activityHTML strings.Builder
			activityHTML.WriteString(`<h2>🕐 Recent Activity</h2><div class="activity-list">`)

			for i, entry := range analytics.RecentActivity {
				if i >= 5 { // Show only 5 most recent
					break
				}
				activityHTML.WriteString(fmt.Sprintf(`
					<div class="activity-item">
						<div class="activity-time">%d min</div>
						<div class="activity-desc">%s</div>
						<div class="activity-date">%s</div>
					</div>
				`, entry.Minutes, entry.Description, entry.LoggedAt.Format("Jan 2, 3:04 PM")))
			}

			activityHTML.WriteString("</div>")
			recentCard.Set("innerHTML", activityHTML.String())
		} else if !recentCard.IsNull() {
			recentCard.Set("innerHTML", `
				<h2>🕐 Recent Activity</h2>
				<p class="empty-state">No recent activity. Start tracking your time!</p>
			`)
		}
	}()

	// Load curricula for progress
	go func() {
		data, err := apiRequest("GET", "/curricula", nil)
		if err != nil {
			return
		}

		var curricula []Curriculum
		if err := json.Unmarshal(data, &curricula); err != nil {
			return
		}

		progressCard := document.Call("getElementById", "progressCard")
		if !progressCard.IsNull() && len(curricula) > 0 {
			progressCard.Set("innerHTML", `
				<h2>🎯 Current Progress</h2>
				<div class="progress-list" id="progressList">Loading projects...</div>
			`)

			// Load projects for first curriculum
			if len(curricula) > 0 {
				loadCurriculumProgress(curricula[0].ID.String())
			}
		} else if !progressCard.IsNull() {
			progressCard.Set("innerHTML", `
				<h2>🎯 Current Progress</h2>
				<p class="empty-state">No curricula yet. Create one to start tracking!</p>
			`)
		}
	}()
}

func loadCurriculumProgress(curriculumID string) {
	data, err := apiRequest("GET", "/curricula/"+curriculumID+"/projects", nil)
	if err != nil {
		return
	}

	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return
	}

	progressList := document.Call("getElementById", "progressList")
	if progressList.IsNull() || len(projects) == 0 {
		return
	}

	var progressHTML strings.Builder
	inProgressCount := 0

	for _, project := range projects {
		// We'd need to load progress for each project, but for now show first few
		if inProgressCount >= 3 {
			break
		}

		progressHTML.WriteString(fmt.Sprintf(`
			<div class="progress-item" onclick="navigate('/project/%s')">
				<div class="progress-icon">%s</div>
				<div class="progress-info">
					<div class="progress-name">%s</div>
					<div class="progress-type">%s</div>
				</div>
			</div>
		`, project.ID.String(), getProjectIcon(project.ProjectType),
			project.Name, getProjectTypeName(project.ProjectType)))

		inProgressCount++
	}

	if inProgressCount == 0 {
		progressHTML.WriteString(`<p class="empty-state">No projects in progress</p>`)
	}

	progressList.Set("innerHTML", progressHTML.String())
}

func showCurriculaView() {
	if authToken == "" {
		navigate("/login")
		return
	}

	currentView = "curricula"
	updateNavbar()

	content := document.Call("getElementById", "content")
	content.Set("innerHTML", `
		<div class="curricula-container">
			<div class="page-header">
				<h1>My Curricula</h1>
				<button class="btn btn-primary" onclick="showCreateCurriculumModal()">
					<span class="btn-icon">➕</span> New Curriculum
				</button>
			</div>
			<div class="curricula-grid" id="curriculaGrid">
				<div class="loading">Loading curricula...</div>
			</div>
		</div>
	`)

	loadCurricula()
}

func loadCurricula() {
	data, err := apiRequest("GET", "/curricula", nil)
	if err != nil {
		showError("Failed to load curricula: " + err.Error())
		return
	}

	var curricula []Curriculum
	if err := json.Unmarshal(data, &curricula); err != nil {
		showError("Failed to parse curricula")
		return
	}

	grid := document.Call("getElementById", "curriculaGrid")
	if grid.IsNull() {
		return
	}

	if len(curricula) == 0 {
		grid.Set("innerHTML", `
			<div class="empty-state-card">
				<div class="empty-icon">📚</div>
				<h3>No curricula yet</h3>
				<p>Create your first curriculum to start tracking your learning journey!</p>
				<button class="btn btn-primary" onclick="showCreateCurriculumModal()">
					Create Curriculum
				</button>
			</div>
		`)
		return
	}

	var html strings.Builder
	for _, curr := range curricula {
		html.WriteString(fmt.Sprintf(`
			<div class="curriculum-card" onclick="navigate('/curriculum/%s')">
				<div class="curriculum-icon">🌲</div>
				<h3>%s</h3>
				<p>%s</p>
				<div class="curriculum-meta">
					Created %s
				</div>
				<div class="curriculum-actions">
					<button class="btn btn-sm" onclick="event.stopPropagation(); editCurriculum('%s')">Edit</button>
					<button class="btn btn-sm btn-danger" onclick="event.stopPropagation(); deleteCurriculum('%s')">Delete</button>
				</div>
			</div>
		`, curr.ID.String(), curr.Name, curr.Description,
			curr.CreatedAt.Format("Jan 2, 2006"), curr.ID.String(), curr.ID.String()))
	}

	grid.Set("innerHTML", html.String())
}

func showCurriculumDetailView(curriculumID string) {
	if authToken == "" {
		navigate("/login")
		return
	}

	selectedCurriculumID = curriculumID
	currentView = "curriculum-detail"
	updateNavbar()

	content := document.Call("getElementById", "content")
	content.Set("innerHTML", `
		<div class="curriculum-detail-container">
			<div class="loading">Loading curriculum...</div>
		</div>
	`)

	// Load curriculum details
	go func() {
		data, err := apiRequest("GET", "/curricula/"+curriculumID, nil)
		if err != nil {
			showError("Failed to load curriculum: " + err.Error())
			navigate("/curricula")
			return
		}

		var curriculum Curriculum
		if err := json.Unmarshal(data, &curriculum); err != nil {
			showError("Failed to parse curriculum")
			return
		}

		// Load projects
		projectsData, err := apiRequest("GET", "/curricula/"+curriculumID+"/projects", nil)
		if err != nil {
			showError("Failed to load projects: " + err.Error())
			return
		}

		var projects []Project
		if err := json.Unmarshal(projectsData, &projects); err != nil {
			showError("Failed to parse projects")
			return
		}

		// Group projects by type
		projectsByType := make(map[ProjectType][]Project)
		for _, p := range projects {
			projectsByType[p.ProjectType] = append(projectsByType[p.ProjectType], p)
		}

		content.Set("innerHTML", fmt.Sprintf(`
			<div class="curriculum-detail-container">
				<div class="page-header">
					<div>
						<h1>%s</h1>
						<p class="curriculum-desc">%s</p>
					</div>
					<button class="btn btn-primary" onclick="showCreateProjectModal()">
						<span class="btn-icon">➕</span> Add Project
					</button>
				</div>
				<div class="project-tree" id="projectTree"></div>
			</div>
		`, curriculum.Name, curriculum.Description))

		renderProjectTree(projectsByType)
	}()
}

func renderProjectTree(projectsByType map[ProjectType][]Project) {
	tree := document.Call("getElementById", "projectTree")
	if tree.IsNull() {
		return
	}

	var html strings.Builder

	// Define the order of project types
	typeOrder := []ProjectType{
		ProjectTypeRoot,
		ProjectTypeRootTest,
		ProjectTypeBase,
		ProjectTypeBaseTest,
		ProjectTypeLowerBranch,
		ProjectTypeMiddleBranch,
		ProjectTypeUpperBranch,
		ProjectTypeFlowerMilestone,
	}

	for _, projectType := range typeOrder {
		projects, exists := projectsByType[projectType]
		if !exists || len(projects) == 0 {
			continue
		}

		html.WriteString(fmt.Sprintf(`
			<div class="project-section">
				<h3 class="project-section-title">
					<span class="section-icon">%s</span>
					%s
				</h3>
				<div class="project-list">
		`, getProjectIcon(projectType), getProjectTypeName(projectType)))

		for _, project := range projects {
			html.WriteString(fmt.Sprintf(`
				<div class="project-card" onclick="navigate('/project/%s')">
					<div class="project-header">
						<span class="project-id">%s</span>
						<h4>%s</h4>
					</div>
					<p class="project-desc">%s</p>
					<div class="project-meta">
						<span class="project-time">⏱️ %s</span>
					</div>
				</div>
			`, project.ID.String(), project.Identifier, project.Name,
				project.Description, project.EstimatedTime))
		}

		html.WriteString("</div></div>")
	}

	if html.Len() == 0 {
		html.WriteString(`
			<div class="empty-state-card">
				<div class="empty-icon">🌱</div>
				<h3>No projects yet</h3>
				<p>Add your first project to start building your curriculum!</p>
				<button class="btn btn-primary" onclick="showCreateProjectModal()">
					Add Project
				</button>
			</div>
		`)
	}

	tree.Set("innerHTML", html.String())
}

func showProjectDetailView(projectID string) {
	if authToken == "" {
		navigate("/login")
		return
	}

	selectedProjectID = projectID
	currentView = "project-detail"
	updateNavbar()

	content := document.Call("getElementById", "content")
	content.Set("innerHTML", `
		<div class="project-detail-container">
			<div class="loading">Loading project...</div>
		</div>
	`)

	// Load project details
	go func() {
		// Load project
		projectData, err := apiRequest("GET", "/projects/"+projectID, nil)
		if err != nil {
			showError("Failed to load project: " + err.Error())
			window.Get("history").Call("back")
			return
		}

		var project Project
		if err := json.Unmarshal(projectData, &project); err != nil {
			showError("Failed to parse project")
			return
		}

		// Load progress
		progressData, err := apiRequest("GET", "/projects/"+projectID+"/progress", nil)
		if err != nil {
			showError("Failed to load progress: " + err.Error())
		}

		var progress ProjectProgress
		if progressData != nil {
			json.Unmarshal(progressData, &progress)
		}

		// Render project details
		content.Set("innerHTML", fmt.Sprintf(`
			<div class="project-detail-container">
				<div class="project-detail-header">
					<div>
						<span class="project-type-badge">%s %s</span>
						<h1>%s %s</h1>
						<p class="project-desc">%s</p>
					</div>
					<div class="project-status">
						<select id="statusSelect" onchange="updateProjectStatus()" class="status-select status-%s">
							<option value="not_started" %s>Not Started</option>
							<option value="in_progress" %s>In Progress</option>
							<option value="completed" %s>Completed</option>
							<option value="paused" %s>Paused</option>
						</select>
					</div>
				</div>
				
				<div class="project-info-grid">
					<div class="info-card">
						<h3>⏱️ Time Tracking</h3>
						<div class="time-info">
							<div class="time-stat">
								<span class="time-label">Estimated:</span>
								<span class="time-value">%s</span>
							</div>
							<div class="time-stat">
								<span class="time-label">Logged:</span>
								<span class="time-value">%d hours %d minutes</span>
							</div>
						</div>
						<button class="btn btn-sm" onclick="showLogTimeModal()">Log Time</button>
					</div>
					
					<div class="info-card">
						<h3>🎯 Learning Objectives</h3>
						<ul class="objectives-list">%s</ul>
					</div>
					
					<div class="info-card">
						<h3>📋 Prerequisites</h3>
						<ul class="prerequisites-list">%s</ul>
					</div>
				</div>
				
				<div class="project-tabs">
					<button class="tab-btn active" onclick="showProjectTab('notes')">📝 Notes</button>
					<button class="tab-btn" onclick="showProjectTab('reflections')">💭 Reflections</button>
					<button class="tab-btn" onclick="showProjectTab('time')">⏱️ Time Entries</button>
				</div>
				
				<div class="tab-content" id="tabContent">
					<div class="loading">Loading...</div>
				</div>
			</div>
		`, getProjectIcon(project.ProjectType), getProjectTypeName(project.ProjectType),
			project.Identifier, project.Name, project.Description, progress.Status,
			getSelectedAttr(string(progress.Status), "not_started"),
			getSelectedAttr(string(progress.Status), "in_progress"),
			getSelectedAttr(string(progress.Status), "completed"),
			getSelectedAttr(string(progress.Status), "paused"),
			project.EstimatedTime, progress.TimeSpentMinutes/60, progress.TimeSpentMinutes%60,
			formatList(project.LearningObjectives), formatList(project.Prerequisites)))

		// Load default tab (notes)
		showProjectTab("notes")
	}()
}

func showProjectTab(tab string) {
	// Update active tab
	tabs := document.Call("querySelectorAll", ".tab-btn")
	for i := 0; i < tabs.Length(); i++ {
		t := tabs.Index(i)
		if strings.Contains(t.Get("textContent").String(), getTabName(tab)) {
			t.Get("classList").Call("add", "active")
		} else {
			t.Get("classList").Call("remove", "active")
		}
	}

	tabContent := document.Call("getElementById", "tabContent")
	if tabContent.IsNull() {
		return
	}

	switch tab {
	case "notes":
		loadNotesTab()
	case "reflections":
		loadReflectionsTab()
	case "time":
		loadTimeEntriesTab()
	}
}

func loadNotesTab() {
	tabContent := document.Call("getElementById", "tabContent")

	// Load notes
	data, err := apiRequest("GET", "/projects/"+selectedProjectID+"/notes", nil)
	if err != nil {
		tabContent.Set("innerHTML", `<div class="error">Failed to load notes</div>`)
		return
	}

	var notes []Note
	json.Unmarshal(data, &notes)

	var html strings.Builder
	html.WriteString(`
		<div class="tab-header">
			<h3>Notes</h3>
			<button class="btn btn-sm" onclick="showCreateNoteModal()">Add Note</button>
		</div>
		<div class="notes-list">
	`)

	if len(notes) == 0 {
		html.WriteString(`<p class="empty-state">No notes yet. Add your first note!</p>`)
	} else {
		for _, note := range notes {
			html.WriteString(fmt.Sprintf(`
				<div class="note-card">
					<div class="note-header">
						<h4>%s</h4>
						<div class="note-actions">
							<button class="btn-icon" onclick="editNote('%s')">✏️</button>
							<button class="btn-icon" onclick="deleteNote('%s')">🗑️</button>
						</div>
					</div>
					<div class="note-content">%s</div>
					<div class="note-date">%s</div>
				</div>
			`, note.Title, note.ID.String(), note.ID.String(),
				note.Content, note.CreatedAt.Format("Jan 2, 2006 at 3:04 PM")))
		}
	}

	html.WriteString("</div>")
	tabContent.Set("innerHTML", html.String())
}

func loadReflectionsTab() {
	tabContent := document.Call("getElementById", "tabContent")

	// Load reflections
	data, err := apiRequest("GET", "/projects/"+selectedProjectID+"/reflections", nil)
	if err != nil {
		tabContent.Set("innerHTML", `<div class="error">Failed to load reflections</div>`)
		return
	}

	var reflections []Reflection
	json.Unmarshal(data, &reflections)

	var html strings.Builder
	html.WriteString(`
		<div class="tab-header">
			<h3>Reflections</h3>
			<button class="btn btn-sm" onclick="showCreateReflectionModal()">Add Reflection</button>
		</div>
		<div class="reflections-list">
	`)

	if len(reflections) == 0 {
		html.WriteString(`<p class="empty-state">No reflections yet. Share your thoughts!</p>`)
	} else {
		for _, reflection := range reflections {
			html.WriteString(fmt.Sprintf(`
				<div class="reflection-card">
					<div class="reflection-content">%s</div>
					<div class="reflection-footer">
						<span class="reflection-date">%s</span>
						<div class="reflection-actions">
							<button class="btn-icon" onclick="editReflection('%s')">✏️</button>
							<button class="btn-icon" onclick="deleteReflection('%s')">🗑️</button>
						</div>
					</div>
				</div>
			`, reflection.Content, reflection.CreatedAt.Format("Jan 2, 2006 at 3:04 PM"),
				reflection.ID.String(), reflection.ID.String()))
		}
	}

	html.WriteString("</div>")
	tabContent.Set("innerHTML", html.String())
}

func loadTimeEntriesTab() {
	tabContent := document.Call("getElementById", "tabContent")

	// Load time entries
	data, err := apiRequest("GET", "/projects/"+selectedProjectID+"/time-entries", nil)
	if err != nil {
		tabContent.Set("innerHTML", `<div class="error">Failed to load time entries</div>`)
		return
	}

	var entries []TimeEntry
	json.Unmarshal(data, &entries)

	var html strings.Builder
	html.WriteString(`
		<div class="tab-header">
			<h3>Time Entries</h3>
			<button class="btn btn-sm" onclick="showLogTimeModal()">Log Time</button>
		</div>
		<div class="time-entries-list">
	`)

	if len(entries) == 0 {
		html.WriteString(`<p class="empty-state">No time logged yet. Start tracking!</p>`)
	} else {
		totalMinutes := 0
		for _, entry := range entries {
			totalMinutes += entry.Minutes
			html.WriteString(fmt.Sprintf(`
				<div class="time-entry">
					<div class="time-entry-main">
						<span class="time-duration">%d min</span>
						<span class="time-desc">%s</span>
					</div>
					<span class="time-date">%s</span>
				</div>
			`, entry.Minutes, entry.Description,
				entry.LoggedAt.Format("Jan 2, 3:04 PM")))
		}

		html.WriteString(fmt.Sprintf(`
			<div class="time-total">
				<strong>Total:</strong> %d hours %d minutes
			</div>
		`, totalMinutes/60, totalMinutes%60))
	}

	html.WriteString("</div>")
	tabContent.Set("innerHTML", html.String())
}

func showAnalyticsView() {
	if authToken == "" {
		navigate("/login")
		return
	}

	currentView = "analytics"
	updateNavbar()

	content := document.Call("getElementById", "content")
	content.Set("innerHTML", `
		<div class="analytics-container">
			<h1>Analytics Dashboard</h1>
			<div class="analytics-grid" id="analyticsGrid">
				<div class="loading">Loading analytics...</div>
			</div>
		</div>
	`)

	loadAnalytics()
}

func loadAnalytics() {
	data, err := apiRequest("GET", "/analytics", nil)
	if err != nil {
		showError("Failed to load analytics: " + err.Error())
		return
	}

	var analytics AnalyticsResponse
	if err := json.Unmarshal(data, &analytics); err != nil {
		showError("Failed to parse analytics")
		return
	}

	grid := document.Call("getElementById", "analyticsGrid")
	if grid.IsNull() {
		return
	}

	// Calculate completion percentage
	completionRate := 0.0
	if analytics.TotalProjects > 0 {
		completionRate = float64(analytics.CompletedProjects) / float64(analytics.TotalProjects) * 100
	}

	// Build project type distribution
	var typeDistHTML strings.Builder
	for _, pType := range []ProjectType{
		ProjectTypeRoot, ProjectTypeBase,
		ProjectTypeLowerBranch, ProjectTypeMiddleBranch,
		ProjectTypeUpperBranch, ProjectTypeFlowerMilestone,
	} {
		count := analytics.ProjectsByType[pType]
		completion := analytics.CompletionByType[pType]
		if count > 0 {
			typeDistHTML.WriteString(fmt.Sprintf(`
				<div class="type-stat">
					<div class="type-header">
						<span class="type-icon">%s</span>
						<span class="type-name">%s</span>
						<span class="type-count">%d</span>
					</div>
					<div class="progress-bar">
						<div class="progress-fill" style="width: %.1f%%"></div>
					</div>
					<span class="type-completion">%.1f%% completed</span>
				</div>
			`, getProjectIcon(pType), getProjectTypeName(pType), count, completion, completion))
		}
	}

	// Build weekly chart (simplified text representation)
	var weeklyChartHTML strings.Builder
	weeklyChartHTML.WriteString(`<div class="weekly-chart">`)
	maxMinutes := 0
	for _, minutes := range analytics.WeeklyTimeSpent {
		if minutes > maxMinutes {
			maxMinutes = minutes
		}
	}

	for i, minutes := range analytics.WeeklyTimeSpent {
		height := 0
		if maxMinutes > 0 {
			height = int(float64(minutes) / float64(maxMinutes) * 100)
		}
		weeklyChartHTML.WriteString(fmt.Sprintf(`
			<div class="week-bar">
				<div class="bar" style="height: %d%%"></div>
				<div class="week-label">W%d</div>
				<div class="week-value">%dh</div>
			</div>
		`, height, i+1, minutes/60))
	}
	weeklyChartHTML.WriteString(`</div>`)

	grid.Set("innerHTML", fmt.Sprintf(`
		<div class="analytics-card main-stats">
			<h2>Overview</h2>
			<div class="main-stats-grid">
				<div class="main-stat">
					<div class="stat-icon">📚</div>
					<div class="stat-value">%d</div>
					<div class="stat-label">Total Projects</div>
				</div>
				<div class="main-stat">
					<div class="stat-icon">✅</div>
					<div class="stat-value">%d</div>
					<div class="stat-label">Completed</div>
				</div>
				<div class="main-stat">
					<div class="stat-icon">🚀</div>
					<div class="stat-value">%d</div>
					<div class="stat-label">In Progress</div>
				</div>
				<div class="main-stat">
					<div class="stat-icon">⏱️</div>
					<div class="stat-value">%dh %dm</div>
					<div class="stat-label">Total Time</div>
				</div>
			</div>
			<div class="completion-rate">
				<h3>Overall Completion Rate</h3>
				<div class="big-progress-bar">
					<div class="big-progress-fill" style="width: %.1f%%"></div>
				</div>
				<span class="completion-text">%.1f%% Complete</span>
			</div>
		</div>
		
		<div class="analytics-card">
			<h2>Projects by Type</h2>
			%s
		</div>
		
		<div class="analytics-card">
			<h2>Weekly Time Investment</h2>
			%s
		</div>
	`, analytics.TotalProjects, analytics.CompletedProjects,
		analytics.InProgressProjects, analytics.TotalTimeSpent/60,
		analytics.TotalTimeSpent%60, completionRate, completionRate,
		typeDistHTML.String(), weeklyChartHTML.String()))
}

// Event handlers
func handleLogin(this js.Value, args []js.Value) interface{} {
	args[0].Call("preventDefault")

	email := document.Call("getElementById", "email").Get("value").String()
	password := document.Call("getElementById", "password").Get("value").String()

	go func() {
		req := LoginRequest{
			Email:    email,
			Password: password,
		}

		data, err := apiRequest("POST", "/login", req)
		if err != nil {
			showError("Login failed: " + err.Error())
			return
		}

		var resp AuthResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			showError("Login failed")
			return
		}

		authToken = resp.Token
		currentUser = &resp.User
		setLocalStorage("authToken", authToken)

		navigate("/dashboard")
	}()

	return nil
}

func handleRegister(this js.Value, args []js.Value) interface{} {
	args[0].Call("preventDefault")

	name := document.Call("getElementById", "name").Get("value").String()
	email := document.Call("getElementById", "email").Get("value").String()
	password := document.Call("getElementById", "password").Get("value").String()

	go func() {
		req := RegisterRequest{
			Name:     name,
			Email:    email,
			Password: password,
		}

		data, err := apiRequest("POST", "/register", req)
		if err != nil {
			showError("Registration failed: " + err.Error())
			return
		}

		var resp AuthResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			showError("Registration failed")
			return
		}

		authToken = resp.Token
		currentUser = &resp.User
		setLocalStorage("authToken", authToken)

		navigate("/dashboard")
	}()

	return nil
}

// Helper functions
func showError(message string) {
	// Create a toast notification
	toast := document.Call("createElement", "div")
	toast.Set("className", "toast toast-error")
	toast.Set("textContent", message)
	document.Get("body").Call("appendChild", toast)

	// Remove after 3 seconds
	window.Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		toast.Call("remove")
		return nil
	}), 3000)
}

func showSuccess(message string) {
	toast := document.Call("createElement", "div")
	toast.Set("className", "toast toast-success")
	toast.Set("textContent", message)
	document.Get("body").Call("appendChild", toast)

	window.Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		toast.Call("remove")
		return nil
	}), 3000)
}

func getProjectIcon(projectType ProjectType) string {
	icons := map[ProjectType]string{
		ProjectTypeRoot:            "🌱",
		ProjectTypeBase:            "🌲",
		ProjectTypeLowerBranch:     "🌿",
		ProjectTypeMiddleBranch:    "🍃",
		ProjectTypeUpperBranch:     "🌸",
		ProjectTypeFlowerMilestone: "🌺",
		ProjectTypeRootTest:        "🧪",
		ProjectTypeBaseTest:        "📝",
	}

	if icon, ok := icons[projectType]; ok {
		return icon
	}
	return "📋"
}

func getProjectTypeName(projectType ProjectType) string {
	names := map[ProjectType]string{
		ProjectTypeRoot:            "Root Projects",
		ProjectTypeBase:            "Base Projects",
		ProjectTypeLowerBranch:     "Lower Branch",
		ProjectTypeMiddleBranch:    "Middle Branch",
		ProjectTypeUpperBranch:     "Upper Branch",
		ProjectTypeFlowerMilestone: "Flower Milestones",
		ProjectTypeRootTest:        "Root Tests",
		ProjectTypeBaseTest:        "Base Tests",
	}

	if name, ok := names[projectType]; ok {
		return name
	}
	return "Projects"
}

func getSelectedAttr(current, value string) string {
	if current == value {
		return "selected"
	}
	return ""
}

func formatList(items []string) string {
	if len(items) == 0 {
		return "<li>None specified</li>"
	}

	var html strings.Builder
	for _, item := range items {
		html.WriteString(fmt.Sprintf("<li>%s</li>", item))
	}
	return html.String()
}

func getTabName(tab string) string {
	switch tab {
	case "notes":
		return "Notes"
	case "reflections":
		return "Reflections"
	case "time":
		return "Time"
	default:
		return ""
	}
}

// Global functions exported to JS
func init() {
	// Export functions that need to be called from HTML
	js.Global().Set("navigate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			navigate(args[0].String())
		}
		return nil
	}))

	js.Global().Set("logout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		authToken = ""
		currentUser = nil
		removeLocalStorage("authToken")
		navigate("/login")
		return nil
	}))

	js.Global().Set("showCreateCurriculumModal", js.FuncOf(showCreateCurriculumModal))
	js.Global().Set("showCreateProjectModal", js.FuncOf(showCreateProjectModal))
	js.Global().Set("showCreateNoteModal", js.FuncOf(showCreateNoteModal))
	js.Global().Set("showCreateReflectionModal", js.FuncOf(showCreateReflectionModal))
	js.Global().Set("showLogTimeModal", js.FuncOf(showLogTimeModal))
	js.Global().Set("updateProjectStatus", js.FuncOf(updateProjectStatus))
	js.Global().Set("showProjectTab", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			showProjectTab(args[0].String())
		}
		return nil
	}))

	// Export closeModal function
	js.Global().Set("closeModal", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		closeModal()
		return nil
	}))
	// Export edit and delete functions - placeholder implementations
	js.Global().Set("editCurriculum", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			// TODO: Implement edit curriculum functionality
			showError("Edit curriculum functionality not yet implemented")
		}
		return nil
	}))

	js.Global().Set("deleteCurriculum", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			// TODO: Implement delete curriculum functionality
			showError("Delete curriculum functionality not yet implemented")
		}
		return nil
	}))

	js.Global().Set("editNote", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			// TODO: Implement edit note functionality
			showError("Edit note functionality not yet implemented")
		}
		return nil
	}))

	js.Global().Set("deleteNote", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			// TODO: Implement delete note functionality
			showError("Delete note functionality not yet implemented")
		}
		return nil
	}))

	js.Global().Set("editReflection", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			// TODO: Implement edit reflection functionality
			showError("Edit reflection functionality not yet implemented")
		}
		return nil
	}))

	js.Global().Set("deleteReflection", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			// TODO: Implement delete reflection functionality
			showError("Delete reflection functionality not yet implemented")
		}
		return nil
	}))
}

// Modal functions
func showCreateCurriculumModal(this js.Value, args []js.Value) interface{} {
	modal := createModal("Create New Curriculum", `
		<form id="createCurriculumForm">
			<div class="form-group">
				<label for="curriculumName">Name</label>
				<input type="text" id="curriculumName" required>
			</div>
			<div class="form-group">
				<label for="curriculumDesc">Description</label>
				<textarea id="curriculumDesc" rows="3"></textarea>
			</div>
			<div class="modal-actions">
				<button type="button" class="btn btn-secondary" onclick="closeModal()">Cancel</button>
				<button type="submit" class="btn btn-primary">Create</button>
			</div>
		</form>
	`)

	form := modal.Call("querySelector", "#createCurriculumForm")
	form.Call("addEventListener", "submit", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")

		name := modal.Call("querySelector", "#curriculumName").Get("value").String()
		desc := modal.Call("querySelector", "#curriculumDesc").Get("value").String()

		go func() {
			req := CurriculumRequest{
				Name:        name,
				Description: desc,
			}

			_, err := apiRequest("POST", "/curricula", req)
			if err != nil {
				showError("Failed to create curriculum: " + err.Error())
				return
			}

			showSuccess("Curriculum created successfully!")
			closeModal()
			loadCurricula()
		}()

		return nil
	}))

	return nil
}

func showCreateProjectModal(this js.Value, args []js.Value) interface{} {
	modal := createModal("Add New Project", fmt.Sprintf(`
		<form id="createProjectForm">
			<div class="form-group">
				<label for="projectType">Project Type</label>
				<select id="projectType" required>
					<option value="">Select a type</option>
					<option value="root">🌱 Root Project</option>
					<option value="base">🌲 Base Project</option>
					<option value="lower_branch">🌿 Lower Branch</option>
					<option value="middle_branch">🍃 Middle Branch</option>
					<option value="upper_branch">🌸 Upper Branch</option>
					<option value="flower_milestone">🌺 Flower Milestone</option>
					<option value="root_test">🧪 Root Test</option>
					<option value="base_test">📝 Base Test</option>
				</select>
			</div>
			<div class="form-group">
				<label for="projectIdentifier">Identifier</label>
				<input type="text" id="projectIdentifier" placeholder="e.g., R1, B2">
			</div>
			<div class="form-group">
				<label for="projectName">Name</label>
				<input type="text" id="projectName" required>
			</div>
			<div class="form-group">
				<label for="projectDesc">Description</label>
				<textarea id="projectDesc" rows="3"></textarea>
			</div>
			<div class="form-group">
				<label for="projectTime">Estimated Time</label>
				<input type="text" id="projectTime" placeholder="e.g., 2-3 hours">
			</div>
			<div class="form-group">
				<label for="projectObjectives">Learning Objectives (one per line)</label>
				<textarea id="projectObjectives" rows="3"></textarea>
			</div>
			<div class="form-group">
				<label for="projectPrereqs">Prerequisites (one per line)</label>
				<textarea id="projectPrereqs" rows="2"></textarea>
			</div>
			<div class="modal-actions">
				<button type="button" class="btn btn-secondary" onclick="closeModal()">Cancel</button>
				<button type="submit" class="btn btn-primary">Create</button>
			</div>
		</form>
	`))

	form := modal.Call("querySelector", "#createProjectForm")
	form.Call("addEventListener", "submit", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")

		// Get form values
		projectType := modal.Call("querySelector", "#projectType").Get("value").String()
		identifier := modal.Call("querySelector", "#projectIdentifier").Get("value").String()
		name := modal.Call("querySelector", "#projectName").Get("value").String()
		desc := modal.Call("querySelector", "#projectDesc").Get("value").String()
		estimatedTime := modal.Call("querySelector", "#projectTime").Get("value").String()

		// Parse objectives and prerequisites
		objectivesText := modal.Call("querySelector", "#projectObjectives").Get("value").String()
		prereqsText := modal.Call("querySelector", "#projectPrereqs").Get("value").String()

		objectives := parseLines(objectivesText)
		prerequisites := parseLines(prereqsText)

		go func() {
			req := ProjectRequest{
				Identifier:         identifier,
				Name:               name,
				Description:        desc,
				LearningObjectives: objectives,
				EstimatedTime:      estimatedTime,
				Prerequisites:      prerequisites,
				ProjectType:        ProjectType(projectType),
				OrderIndex:         0, // Could be made configurable
			}

			_, err := apiRequest("POST", "/curricula/"+selectedCurriculumID+"/projects", req)
			if err != nil {
				showError("Failed to create project: " + err.Error())
				return
			}

			showSuccess("Project created successfully!")
			closeModal()
			showCurriculumDetailView(selectedCurriculumID)
		}()

		return nil
	}))

	return nil
}

func showCreateNoteModal(this js.Value, args []js.Value) interface{} {
	modal := createModal("Add Note", `
		<form id="createNoteForm">
			<div class="form-group">
				<label for="noteTitle">Title</label>
				<input type="text" id="noteTitle">
			</div>
			<div class="form-group">
				<label for="noteContent">Content</label>
				<textarea id="noteContent" rows="5" required></textarea>
			</div>
			<div class="modal-actions">
				<button type="button" class="btn btn-secondary" onclick="closeModal()">Cancel</button>
				<button type="submit" class="btn btn-primary">Save</button>
			</div>
		</form>
	`)

	form := modal.Call("querySelector", "#createNoteForm")
	form.Call("addEventListener", "submit", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")

		title := modal.Call("querySelector", "#noteTitle").Get("value").String()
		content := modal.Call("querySelector", "#noteContent").Get("value").String()

		go func() {
			req := NoteRequest{
				Title:   title,
				Content: content,
			}

			_, err := apiRequest("POST", "/projects/"+selectedProjectID+"/notes", req)
			if err != nil {
				showError("Failed to create note: " + err.Error())
				return
			}

			showSuccess("Note saved!")
			closeModal()
			loadNotesTab()
		}()

		return nil
	}))

	return nil
}

func showCreateReflectionModal(this js.Value, args []js.Value) interface{} {
	modal := createModal("Add Reflection", `
		<form id="createReflectionForm">
			<div class="form-group">
				<label for="reflectionContent">What did you learn? What insights did you gain?</label>
				<textarea id="reflectionContent" rows="6" required></textarea>
			</div>
			<div class="modal-actions">
				<button type="button" class="btn btn-secondary" onclick="closeModal()">Cancel</button>
				<button type="submit" class="btn btn-primary">Save</button>
			</div>
		</form>
	`)

	form := modal.Call("querySelector", "#createReflectionForm")
	form.Call("addEventListener", "submit", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")

		content := modal.Call("querySelector", "#reflectionContent").Get("value").String()

		go func() {
			req := ReflectionRequest{
				Content: content,
			}

			_, err := apiRequest("POST", "/projects/"+selectedProjectID+"/reflections", req)
			if err != nil {
				showError("Failed to create reflection: " + err.Error())
				return
			}

			showSuccess("Reflection saved!")
			closeModal()
			loadReflectionsTab()
		}()

		return nil
	}))

	return nil
}

func showLogTimeModal(this js.Value, args []js.Value) interface{} {
	modal := createModal("Log Time", `
		<form id="logTimeForm">
			<div class="form-group">
				<label for="timeMinutes">Minutes</label>
				<input type="number" id="timeMinutes" min="1" required>
			</div>
			<div class="form-group">
				<label for="timeDesc">Description</label>
				<input type="text" id="timeDesc" placeholder="What did you work on?">
			</div>
			<div class="modal-actions">
				<button type="button" class="btn btn-secondary" onclick="closeModal()">Cancel</button>
				<button type="submit" class="btn btn-primary">Log Time</button>
			</div>
		</form>
	`)
	form := modal.Call("querySelector", "#logTimeForm")
	form.Call("addEventListener", "submit", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		args[0].Call("preventDefault")

		minutes := modal.Call("querySelector", "#timeMinutes").Get("value").Int()
		desc := modal.Call("querySelector", "#timeDesc").Get("value").String()

		go func() {
			req := TimeEntryRequest{
				Minutes:     minutes,
				Description: desc,
			}

			_, err := apiRequest("POST", "/projects/"+selectedProjectID+"/time-entries", req)
			if err != nil {
				showError("Failed to log time: " + err.Error())
				return
			}

			showSuccess("Time logged successfully!")
			closeModal()

			// Reload current view
			if currentView == "project-detail" {
				showProjectDetailView(selectedProjectID)
			}
		}()

		return nil
	}))

	return nil
}

func updateProjectStatus(this js.Value, args []js.Value) interface{} {
	select_ := document.Call("getElementById", "statusSelect")
	if select_.IsNull() {
		return nil
	}

	newStatus := select_.Get("value").String()

	go func() {
		// First get current progress
		progressData, err := apiRequest("GET", "/projects/"+selectedProjectID+"/progress", nil)
		if err != nil {
			showError("Failed to load progress")
			return
		}

		var progress ProjectProgress
		json.Unmarshal(progressData, &progress)

		// Update status
		req := ProgressUpdateRequest{
			Status:           ProjectStatus(newStatus),
			TimeSpentMinutes: progress.TimeSpentMinutes,
		}

		_, err = apiRequest("PUT", "/projects/"+selectedProjectID+"/progress", req)
		if err != nil {
			showError("Failed to update status: " + err.Error())
			return
		}

		showSuccess("Status updated!")

		// Update select class
		select_.Set("className", "status-select status-"+newStatus)
	}()

	return nil
}

// Helper modal functions
func createModal(title, content string) js.Value {
	modal := document.Call("createElement", "div")
	modal.Set("className", "modal-overlay")
	modal.Set("innerHTML", fmt.Sprintf(`
		<div class="modal">
			<div class="modal-header">
				<h2>%s</h2>
				<button class="modal-close" onclick="closeModal()">×</button>
			</div>
			<div class="modal-content">
				%s
			</div>
		</div>
	`, title, content))

	document.Get("body").Call("appendChild", modal)

	// Close on overlay click
	modal.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("target").Equal(modal) {
			closeModal()
		}
		return nil
	}))

	return modal
}

func closeModal() {
	overlay := document.Call("querySelector", ".modal-overlay")
	if !overlay.IsNull() {
		overlay.Call("remove")
	}
}

func parseLines(text string) []string {
	lines := strings.Split(text, "\n")
	var result []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

// Export closeModal to global scope
func init() {
	js.Global().Set("closeModal", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		closeModal()
		return nil
	}))
}
