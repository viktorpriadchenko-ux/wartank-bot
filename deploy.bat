@echo off
chcp 65001 >nul
title WarTank Bot - Deploy to Server

echo ============================
echo   Deploy WarTank Bot
echo ============================
echo.

:: 1. Сборка Linux-бинарника
echo [1/4] Сборка для Linux...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
"C:\Program Files\Go\bin\go.exe" build -o bin\server_linux .\cmd\server
set GOOS=
set GOARCH=
set CGO_ENABLED=
if errorlevel 1 (
    echo [ОШИБКА] Сборка не удалась!
    pause
    exit /b 1
)
echo [OK] Сборка завершена

:: 2. Загрузка на сервер
echo [2/4] Загрузка на сервер...
scp -o StrictHostKeyChecking=no bin\server_linux pavel@195.14.48.124:/home/pavel/text_rpg/wartank-bot/server
if errorlevel 1 (
    echo [ОШИБКА] Загрузка не удалась!
    pause
    exit /b 1
)
echo [OK] Загружено

:: 3. Перезапуск сервиса
echo [3/4] Перезапуск сервиса...
ssh -o StrictHostKeyChecking=no pavel@195.14.48.124 "chmod +x /home/pavel/text_rpg/wartank-bot/server && systemctl --user restart wartank-bot"
if errorlevel 1 (
    echo [ОШИБКА] Перезапуск не удался!
    pause
    exit /b 1
)
echo [OK] Перезапущен

:: 4. Проверка статуса
echo [4/4] Проверка статуса...
timeout /t 3 /nobreak >nul
ssh -o StrictHostKeyChecking=no pavel@195.14.48.124 "systemctl --user is-active wartank-bot"
echo.
echo ============================
echo   Деплой завершён!
echo ============================
pause
