@echo off

REM Generate Swagger documentation
echo Generating Swagger documentation...

REM Install swag if not installed
where swag >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Installing swag...
    go install github.com/swaggo/swag/cmd/swag@latest
)

REM Generate docs
swag init -g cmd/app/main.go -o docs --parseDependency --parseInternal

echo Swagger documentation generated successfully!
echo API Documentation available at: http://localhost:8080/swagger/index.html
echo API Docs available at: http://localhost:8080/docs
