@echo off
set APP_NAME=gophkeeper
set OUTPUT_DIR=build

mkdir %OUTPUT_DIR%

set TARGETS=windows/amd64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64

for %%T in (%TARGETS%) do (
    for /F "tokens=1,2 delims=/" %%G in ("%%T") do (
        set GOOS=%%G
        set GOARCH=%%H
        set OUTPUT_FILE=%OUTPUT_DIR%\%APP_NAME%_%%G_%%H
        if "%%G"=="windows" set OUTPUT_FILE=%OUTPUT_FILE%.exe

        echo Сборка для %%G/%%H...
        set GOOS=%%G
        set GOARCH=%%H
        go build -o %OUTPUT_FILE% cmd\server\main.go

        if %ERRORLEVEL%==0 (
            echo Собрано: %OUTPUT_FILE%
        ) else (
            echo Ошибка сборки для %%G/%%H
        )
    )
)

echo Сборка завершена!
