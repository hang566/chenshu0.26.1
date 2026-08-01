@echo off
chcp 65001 >nul 2>&1
title YunShu

set FRONTEND_PORT=3000
set BACKEND_PORT=8080

echo ========================================
echo    YunShu - Service Launcher
echo ========================================
echo.
echo   Frontend port: %FRONTEND_PORT%
echo   Backend port:  %BACKEND_PORT%
echo.

rem === Check Python ===
set "PYTHON_CMD=python"
where python >nul 2>&1
if errorlevel 1 (
    set "PYTHON_CMD=py"
    where py >nul 2>&1
    if errorlevel 1 (
        echo [ERROR] Python not found. Please install Python and add to PATH.
        pause
        exit /b 1
    )
)

rem === Check Go ===
where go >nul 2>&1
if errorlevel 1 (
    echo [WARN] Go not found, trying compiled binary...
    if not exist "bg\yunshu-bg.exe" (
        echo [ERROR] Go not found and bg\yunshu-bg.exe does not exist.
        echo         Install Go: https://go.dev/dl/
        pause
        exit /b 1
    )
)

rem === Kill old processes ===
echo [1/3] Cleaning old processes...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%FRONTEND_PORT% " ^| findstr "LISTENING"') do (
    taskkill /F /PID %%a >nul 2>&1
)
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%BACKEND_PORT% " ^| findstr "LISTENING"') do (
    taskkill /F /PID %%a >nul 2>&1
)

rem === Start backend ===
echo [2/3] Starting backend (port %BACKEND_PORT%)...
start "YunShu-Backend (Go :%BACKEND_PORT%)" cmd /k "cd /d %~dp0bg && go run . || yunshu-bg.exe"

timeout /t 3 /nobreak >nul

rem === Start frontend ===
echo [3/3] Starting frontend (port %FRONTEND_PORT%)...
start "YunShu-Frontend (:%FRONTEND_PORT%)" cmd /k "cd /d %~dp0 && %PYTHON_CMD% -m http.server %FRONTEND_PORT%"

timeout /t 2 /nobreak >nul

echo.
echo ========================================
echo   Done!
echo   Frontend: http://localhost:%FRONTEND_PORT%
echo   Backend:  http://localhost:%BACKEND_PORT%
echo ========================================
echo.
echo   Close the windows to stop services.
echo.

start http://localhost:%FRONTEND_PORT%
