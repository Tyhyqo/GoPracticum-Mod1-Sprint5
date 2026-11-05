// Package daysteps разбирает строку с данными о дневной прогулке
// и формирует по ней отчёт: количество шагов, дистанцию
// и количество потраченных калорий.
package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Ошибки разбора данных о прогулке.
var (
	ErrInvalidFormat   = errors.New("invalid data format")
	ErrInvalidSteps    = errors.New("invalid steps count")
	ErrInvalidDuration = errors.New("invalid duration")
)

// DaySteps содержит данные о дневной прогулке и персональные данные пользователя.
type DaySteps struct {
	Steps    int           // количество шагов за прогулку
	Duration time.Duration // продолжительность прогулки
	personaldata.Personal
}

// Parse разбирает строку формата "678,0h50m"
// (количество шагов, продолжительность)
// и записывает полученные значения в поля DaySteps.
// Возвращает ошибку, если формат строки неверный или если количество шагов
// или продолжительность не положительные. При ошибке поля DaySteps не меняются.
func (ds *DaySteps) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return fmt.Errorf("%w: expected 2 comma-separated values, got %d",
			ErrInvalidFormat, len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidSteps, err)
	}
	if steps <= 0 {
		return fmt.Errorf("%w: %d, must be positive", ErrInvalidSteps, steps)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDuration, err)
	}
	if duration <= 0 {
		return fmt.Errorf("%w: %v, must be positive", ErrInvalidDuration, duration)
	}

	ds.Steps = steps
	ds.Duration = duration

	return nil
}

// ActionInfo возвращает отчёт о прогулке: количество шагов,
// дистанцию в километрах и потраченные калории.
// Возвращает ошибку, если данные прогулки или пользователя некорректны.
func (ds DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)

	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", fmt.Errorf("calculate calories: %w", err)
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps,
		distance,
		calories,
	), nil
}
