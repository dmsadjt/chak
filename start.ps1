# Start Chak AI System

Write-Host "🚀 Starting Chak AI..." -ForegroundColor Cyan

# Start backend in new window
Write-Host "Starting backend server..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd server; go run main.go"

# Wait a bit for backend to start
Start-Sleep -Seconds 2

# Start frontend in new window
Write-Host "Starting frontend server..." -ForegroundColor Green
Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd web; python -m http.server 1010"

# Wait a bit for frontend to start
Start-Sleep -Seconds 2

# Open browser
Write-Host "Opening browser..." -ForegroundColor Green
Start-Process "http://localhost:1010"

Write-Host "`n✅ Chak AI is running!" -ForegroundColor Cyan
Write-Host "Backend:  http://localhost:5000" -ForegroundColor Yellow
Write-Host "Frontend: http://localhost:1010" -ForegroundColor Yellow
Write-Host "`nPress Ctrl+C in the terminal windows to stop." -ForegroundColor Gray
