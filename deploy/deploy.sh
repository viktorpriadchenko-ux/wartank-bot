#!/bin/bash
# deploy.sh — Скрипт деплоя wartank-бота на Ubuntu сервер
#
# Использование:
#   ./deploy/deploy.sh              — собрать + задеплоить
#   ./deploy/deploy.sh --build-only — только собрать бинарник
#   ./deploy/deploy.sh --no-build   — задеплоить уже собранный бинарник
#
# Требования:
#   - Go установлен локально (для кросс-компиляции)
#   - SSH-ключ настроен для pavel@<server>
#   - На сервере: /home/pavel/wartank/ существует

set -euo pipefail

# === Конфигурация ===
SERVER_USER="pavel"
SERVER_HOST="${WARTANK_SERVER_HOST:-your-server-ip}"
SERVER_DIR="/home/pavel/wartank"
SERVICE_NAME="wartank-bot"
BINARY_NAME="server"
LOCAL_BIN="bin/${BINARY_NAME}_linux"

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# === Парсинг аргументов ===
DO_BUILD=true
DO_DEPLOY=true
case "${1:-}" in
    --build-only) DO_DEPLOY=false ;;
    --no-build)   DO_BUILD=false ;;
esac

# === 1. Сборка под Linux ===
if [ "$DO_BUILD" = true ]; then
    log_info "Кросс-компиляция для Linux amd64..."
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "$LOCAL_BIN" ./cmd/server
    log_info "Бинарник собран: $LOCAL_BIN ($(du -h "$LOCAL_BIN" | cut -f1))"
fi

if [ "$DO_DEPLOY" = false ]; then
    log_info "Готово (--build-only)"
    exit 0
fi

# === 2. Проверка подключения ===
if [ "$SERVER_HOST" = "your-server-ip" ]; then
    log_error "Укажите IP сервера: export WARTANK_SERVER_HOST=1.2.3.4"
    exit 1
fi

log_info "Проверяю подключение к ${SERVER_USER}@${SERVER_HOST}..."
ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no "${SERVER_USER}@${SERVER_HOST}" "echo ok" > /dev/null
log_info "Подключение OK"

# === 3. Загрузка бинарника ===
log_info "Загружаю бинарник на сервер..."
ssh "${SERVER_USER}@${SERVER_HOST}" "mkdir -p ${SERVER_DIR}/bin ${SERVER_DIR}/store"
scp "$LOCAL_BIN" "${SERVER_USER}@${SERVER_HOST}:${SERVER_DIR}/bin/${BINARY_NAME}"
ssh "${SERVER_USER}@${SERVER_HOST}" "chmod +x ${SERVER_DIR}/bin/${BINARY_NAME}"
log_info "Бинарник загружен"

# === 4. Установка systemd сервиса (если первый раз) ===
log_info "Обновляю systemd сервис..."
scp deploy/wartank-bot.service "${SERVER_USER}@${SERVER_HOST}:/tmp/${SERVICE_NAME}.service"
ssh "${SERVER_USER}@${SERVER_HOST}" "sudo mv /tmp/${SERVICE_NAME}.service /etc/systemd/system/ && sudo systemctl daemon-reload"

# === 5. Перезапуск ===
log_info "Перезапускаю сервис..."
ssh "${SERVER_USER}@${SERVER_HOST}" "sudo systemctl restart ${SERVICE_NAME}"
sleep 3

# === 6. Проверка статуса ===
log_info "Статус сервиса:"
ssh "${SERVER_USER}@${SERVER_HOST}" "sudo systemctl status ${SERVICE_NAME} --no-pager -l" || true

log_info "Последние логи:"
ssh "${SERVER_USER}@${SERVER_HOST}" "sudo journalctl -u ${SERVICE_NAME} --no-pager -n 20" || true

echo ""
log_info "===== Деплой завершён ====="
echo ""
echo "Полезные команды:"
echo "  ssh ${SERVER_USER}@${SERVER_HOST} 'sudo systemctl status ${SERVICE_NAME}'"
echo "  ssh ${SERVER_USER}@${SERVER_HOST} 'sudo journalctl -u ${SERVICE_NAME} -f'"
echo "  ssh ${SERVER_USER}@${SERVER_HOST} 'sudo systemctl stop ${SERVICE_NAME}'"
echo "  ssh ${SERVER_USER}@${SERVER_HOST} 'sudo systemctl enable ${SERVICE_NAME}'  # автозапуск при ребуте"
