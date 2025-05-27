# Test script for CurriculumTracker API endpoints
# This script tests the endpoints that were experiencing timeout issues

$baseUrl = "http://localhost:8080/api"

Write-Host "Testing CurriculumTracker API endpoints..." -ForegroundColor Green

# Test 1: Try to login and get a token
Write-Host "`n1. Testing login endpoint..." -ForegroundColor Yellow
try {
    $loginData = @{
        email = "admin@curricula.com"
        password = "password"  # Adjust this to the correct password
    } | ConvertTo-Json

    $response = Invoke-RestMethod -Uri "$baseUrl/login" -Method POST -Body $loginData -ContentType "application/json" -ErrorAction Stop
    $token = $response.token
    Write-Host "✓ Login successful, got token" -ForegroundColor Green
} catch {
    Write-Host "✗ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Note: You may need to check the password or create a test user" -ForegroundColor Yellow
    exit 1
}

# Test 2: Test curricula endpoint (was timing out)
Write-Host "`n2. Testing curricula endpoint..." -ForegroundColor Yellow
try {
    $headers = @{
        "Authorization" = "Bearer $token"
        "Content-Type" = "application/json"
    }
    
    $startTime = Get-Date
    $response = Invoke-RestMethod -Uri "$baseUrl/curricula" -Method GET -Headers $headers -ErrorAction Stop
    $endTime = Get-Date
    $duration = ($endTime - $startTime).TotalMilliseconds
    
    Write-Host "✓ Curricula endpoint successful in $($duration)ms" -ForegroundColor Green
    Write-Host "  Found $($response.Count) curricula" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Curricula endpoint failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3: Test analytics endpoint (was timing out)
Write-Host "`n3. Testing analytics endpoint..." -ForegroundColor Yellow
try {
    $startTime = Get-Date
    $response = Invoke-RestMethod -Uri "$baseUrl/analytics" -Method GET -Headers $headers -ErrorAction Stop
    $endTime = Get-Date
    $duration = ($endTime - $startTime).TotalMilliseconds
    
    Write-Host "✓ Analytics endpoint successful in $($duration)ms" -ForegroundColor Green
    Write-Host "  Total projects: $($response.total_projects)" -ForegroundColor Cyan
    Write-Host "  Completed projects: $($response.completed_projects)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Analytics endpoint failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: Test profile endpoint (should be fast)
Write-Host "`n4. Testing profile endpoint..." -ForegroundColor Yellow
try {
    $startTime = Get-Date
    $response = Invoke-RestMethod -Uri "$baseUrl/profile" -Method GET -Headers $headers -ErrorAction Stop
    $endTime = Get-Date
    $duration = ($endTime - $startTime).TotalMilliseconds
    
    Write-Host "✓ Profile endpoint successful in $($duration)ms" -ForegroundColor Green
    Write-Host "  User: $($response.name) ($($response.email))" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Profile endpoint failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nTest complete!" -ForegroundColor Green
Write-Host "If the curricula and analytics endpoints are now working without timeouts," -ForegroundColor Yellow
Write-Host "the timeout issues have been resolved!" -ForegroundColor Yellow
