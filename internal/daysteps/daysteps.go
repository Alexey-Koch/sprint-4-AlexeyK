package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разделяем строку на слайс.
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("недопустимый формат данных")
	}

	// Преобразуем количество шагов в int.
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования количества шагов: %v", err)
	}

	// Проверяем, что количество шагов больше 0.
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	// Преобразуем продолжительность в time.Duration.
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования продолжительности: %v", err)
	}

	return steps, duration, nil

}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	// Получаем количество шагов и продолжительность прогулки.
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Проверяем, чтобы количество шагов было больше 0.
	if steps <= 0 {
		return ""
	}

	// Вычисляем дистанцию в метрах и переводим в километры.
	distanceMeters := float64(steps) * StepLength
	distanceKilometers := distanceMeters / 1000

	// Вычисляем количество калорий, потраченных на прогулке.
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	// Формируем строку для возврата.
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distanceKilometers, calories)
}
