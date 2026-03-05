@echo off
chcp 65001 >nul
title WarTank Bot Server

echo ============================
echo   WarTank Bot Server
echo ============================
echo.

:: Переменные окружения
set STAGE=local
set LOCAL_STORE_PATH=./store
set LOCAL_HTTP_URL=:18061
set SERVER_PORT=18050

:: Проверяем наличие exe
if not exist "bin\server_dev.exe" (
    echo [!] server_dev.exe не найден, собираю...
    go build -o bin\server_dev.exe .\cmd\server
    if errorlevel 1 (
        echo [ОШИБКА] Сборка не удалась!
        pause
        exit /b 1
    )
    echo [OK] Сборка завершена
)

echo.
echo [*] Запуск: STAGE=%STAGE%
echo [*] Store:  %LOCAL_STORE_PATH%
echo [*] HTTP:   http://localhost%LOCAL_HTTP_URL%
echo [*] Port:   %SERVER_PORT%
echo.
echo --- Логи бота ---
echo.

:loop
bin\server_dev.exe
echo.
echo [!] Сервер завершился (код: %errorlevel%). Перезапуск через 5 сек...
timeout /t 5 /nobreak >nul
goto loop
