// package bot_log -- глобальный кольцевой буфер вывода бота.
// Перехватывает os.Stdout через pipe, параллельно пишет в оригинальный stdout.
// Авто-очистка каждые 20 минут. Макс 500 строк в памяти.
package bot_log

import (
	"bufio"
	"html"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	maxLines      = 500
	clearInterval = 20 * time.Minute
)

var (
	lines     []string
	mu        sync.Mutex
	lastClear time.Time
)

// Init -- инициализирует перехват stdout.
// Должен вызываться один раз в начале main().
func Init() {
	lastClear = time.Now()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return
	}
	os.Stdout = w

	go func() {
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 1024*64), 1024*64)
		for scanner.Scan() {
			line := scanner.Text()
			// Дублируем в оригинальный stdout (терминал / journald)
			_, _ = origStdout.WriteString(line + "\n")

			mu.Lock()
			// Авто-очистка каждые 20 минут
			if time.Since(lastClear) >= clearInterval {
				lines = nil
				lastClear = time.Now()
			}
			lines = append(lines, line)
			// Обрезаем если превысили лимит
			if len(lines) > maxLines {
				lines = lines[len(lines)-maxLines:]
			}
			mu.Unlock()
		}
	}()
}

// GetHTML -- возвращает последние строки лога в виде HTML (новые сверху).
func GetHTML() string {
	mu.Lock()
	defer mu.Unlock()

	if len(lines) == 0 {
		return "<em style='color:#888'>Лог пуст...</em>"
	}

	// Показываем последние 300 строк, новые сверху
	start := 0
	if len(lines) > 300 {
		start = len(lines) - 300
	}

	var sb strings.Builder
	for i := len(lines) - 1; i >= start; i-- {
		line := html.EscapeString(lines[i])
		// Подсветка ключевых слов
		color := "#ccc"
		switch {
		case strings.Contains(line, "ПАНИКА") || strings.Contains(line, "PANIC") || strings.Contains(line, "err="):
			color = "#ff6b6b"
		case strings.Contains(line, "выстрел") || strings.Contains(line, "атак") || strings.Contains(line, "урон"):
			color = "#ffd93d"
		case strings.Contains(line, "бой") || strings.Contains(line, "PVP") || strings.Contains(line, "CW") || strings.Contains(line, "DM"):
			color = "#6bcb77"
		case strings.Contains(line, "топливо") || strings.Contains(line, "серебро") || strings.Contains(line, "золото"):
			color = "#4d96ff"
		}
		sb.WriteString(`<div style="color:` + color + `">` + line + `</div>`)
	}
	return sb.String()
}

// NextClear -- возвращает сколько минут осталось до следующей очистки.
func NextClear() int {
	mu.Lock()
	defer mu.Unlock()
	remaining := clearInterval - time.Since(lastClear)
	if remaining < 0 {
		return 0
	}
	return int(remaining.Minutes())
}
