@echo off
echo 🚀 Starting Chak AI...
echo.

echo Starting backend server...
start "Chak Backend" cmd /k "cd server && go run main.go"

timeout /t 2 /nobreak >nul

echo Starting frontend server...
start "Chak Frontend" cmd /k "cd web && python -m http.server 1010"

timeout /t 2 /nobreak >nul

echo Opening browser...
start http://localhost:1010

echo.
echo ✅ Chak AI is running!
echo Backend:  http://localhost:5000
echo Frontend: http://localhost:1010
echo.
echo Close the terminal windows to stop the servers.
pause
