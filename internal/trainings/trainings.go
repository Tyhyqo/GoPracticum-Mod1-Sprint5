// Package trainings разбирает строку с данными о тренировке
// и формирует по ней отчёт: длительность, дистанцию, среднюю скорость
// и количество потраченных калорий.
//
// Поддерживаются два вида тренировок: «Бег» и «Ходьба».
package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Ошибки разбора и обработки данных о тренировке.
var (
	ErrInvalidFormat   = errors.New("invalid data format")
	ErrInvalidSteps    = errors.New("invalid steps count")
	ErrInvalidDuration = errors.New("invalid duration")
	ErrUnknownTraining = errors.New("неизвестный тип тренировки") // текст задан в ТЗ
)

// Виды тренировок.
const (
	trainingRunning = "Бег"
	trainingWalking = "Ходьба"
)

// Training содержит данные о тренировке и персональные данные пользователя.
type Training struct {
	Steps        int           // количество шагов за тренировку
	TrainingType string        // вид тренировки: «Бег» или «Ходьба»
	Duration     time.Duration // продолжительность тренировки
	personaldata.Personal
}

// Parse разбирает строку формата "3456,Ходьба,3h00m"
// (количество шагов, вид тренировки, продолжительность)
// и записывает полученные значения в поля Training.
// Возвращает ошибку, если формат строки неверный или если количество шагов
// или продолжительность не положительные.
func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return fmt.Errorf("%w: expected 3 comma-separated values, got %d", ErrInvalidFormat, len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidSteps, err)
	}
	if steps <= 0 {
		return fmt.Errorf("%w: %d, must be positive", ErrInvalidSteps, steps)
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidDuration, err)
	}
	if duration <= 0 {
		return fmt.Errorf("%w: %v, must be positive", ErrInvalidDuration, duration)
	}

	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration

	return nil
}

// ActionInfo возвращает отчёт о тренировке: вид, длительность в часах,
// дистанцию в километрах, среднюю скорость в км/ч и потраченные калории.
// Для неизвестного вида тренировки возвращает ErrUnknownTraining.
func (t Training) ActionInfo() (string, error) {
	var (
		calories float64
		err      error
	)
	switch t.TrainingType {
	case trainingRunning:
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case trainingWalking:
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", ErrUnknownTraining
	}
	if err != nil {
		return "", fmt.Errorf("calculate calories: %w", err)
	}

	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType,
		t.Duration.Hours(),
		distance,
		speed,
		calories,
	), nil
}
