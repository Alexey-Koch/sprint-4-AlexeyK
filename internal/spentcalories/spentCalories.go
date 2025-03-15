package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// 1. Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// 2. Проверить, чтобы длина слайса была равна 3
	if len(parts) != 3 {
		return 0, "", 0, errors.New("данные должны содержать три элемента")
	}

	// 3. Преобразовать первый элемент в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать количество шагов: %w", err)
	}

	// 4. Преобразовать третий элемент в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать продолжительность: %w", err)
	}

	// 5. Возврат значений без ошибок
	return steps, parts[1], duration, nil

}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	// ваш код ниже
	// 1. Умножаем количество шагов на длину шага
	// 2. Делим на количество метров в километре
	return float64(steps) * lenStep / mInKm // Дистанция в километрах
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	// ваш код ниже
	// 1. Проверить, что продолжительность больше 0
	if duration <= 0 {
		return 0
	}

	// 2. Вычислить дистанцию
	dist := distance(steps)

	// 3. Вычислить среднюю скорость (дистанция в км / продолжительность в часах)
	hours := duration.Hours()
	return dist / hours
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	// ваш код ниже
	// Разбираем входные данные
	//parts := strings.Split(data, ",")
	steps, activity, d, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Ошибка: %s", err)
	}

	// Вычисляем информацию о тренировке
	distanceKm := distance(steps)
	speed := meanSpeed(steps, d)
	var caloriesBurned float64

	// Рассчитываем калории в зависимости от типа активности
	switch strings.ToLower(activity) {
	case "ходьба":
		caloriesBurned = WalkingSpentCalories(steps, weight, height, d)
	case "бег":
		caloriesBurned = RunningSpentCalories(steps, weight, d)
	default:
		return "Неподдерживаемый тип активности"
	}
	// Формируем строку с информацией
	return fmt.Sprintf("Тип тренировки: %s\n Длительность: %.2f ч.\n Дистанция: %.2f км.\n Скорость: %.2f км/ч\n Сожгли калорий: %.2f\n", activity, d.Hours(), distanceKm, speed, caloriesBurned)
	//fmt.Sprintf("Длительность: %.2f ч.\n", d.Hours())
	//fmt.Sprintf("Дистанция: %.2f км.\n", distanceKm)
	//fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
	//fmt.Sprintf("Сожгли калорий: %.2f\n", caloriesBurned)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	// ваш код здесь
	// 1. Рассчитать среднюю скорость
	speed := meanSpeed(steps, duration)

	// 2. Рассчитать и вернуть количество калорий
	caloriesBurned := ((runningCaloriesMeanSpeedMultiplier * speed) - runningCaloriesMeanSpeedShift) * weight
	return caloriesBurned
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// ваш код здесь
	// 1. Рассчитать среднюю скорость
	speed := meanSpeed(steps, duration)

	// 2. Перевести продолжительность в часы
	hours := duration.Hours()

	// 3. Рассчитать и вернуть количество калорий
	caloriesBurned := ((walkingCaloriesWeightMultiplier * weight) + (speed*speed/height)*walkingSpeedHeightMultiplier) * hours
	return caloriesBurned
}
