@echo off
setlocal

:: Define variables
set GOFILE=main.go
set OUTPUT=nftfetch.exe
set BIN_DIR=%USERPROFILE%\bin

:: Ensure the bin directory exists
if not exist "%BIN_DIR%" (
    mkdir "%BIN_DIR%"
)

:: Build the Go file
go build -o "%OUTPUT%" "%GOFILE%"

:: Check if the build was successful
if %ERRORLEVEL% neq 0 (
    echo Build failed
    exit /b %ERRORLEVEL%
)

:: Move the executable to the bin directory and overwrite if it exists
move /Y "%OUTPUT%" "%BIN_DIR%"

:: Check if the move was successful
if %ERRORLEVEL% neq 0 (
    echo Move failed
    exit /b %ERRORLEVEL%
)

echo Build and move completed successfully
endlocal
