// Package spentenergy рассчитывает показатели физической активности:
// пройденную дистанцию, среднюю скорость и количество калорий,
// потраченных при ходьбе и беге.
//
// Расчёты основаны на количестве шагов, росте (в метрах), весе (в килограммах)
// и продолжительности активности. Функции расчёта калорий возвращают ошибку
// при некорректных входных данных; проверить её тип можно через errors.Is
// с одной из ошибок ErrInvalid*.
package spentenergy

import (
	"errors"
	"fmt"
	"time"
)

// Ошибки, возвращаемые при некорректных входных параметрах.
var (
	ErrInvalidSteps    = errors.New("invalid steps count")
	ErrInvalidWeight   = errors.New("invalid weight")
	ErrInvalidHeight   = errors.New("invalid height")
	ErrInvalidDuration = errors.New("invalid duration")
)

// Основные константы, необходимые для расчётов.
const (
	mInKm                      = 1000 // количество метров в километре
	minInH                     = 60   // количество минут в часе
	stepLengthCoefficient      = 0.45 // коэффициент для расчёта длины шага на основе роста
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчёта калорий при ходьбе
)

// validateCaloriesInput проверяет, что все параметры для расчёта калорий положительные.
func validateCaloriesInput(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return fmt.Errorf("%w: %d, must be positive", ErrInvalidSteps, steps)
	case weight <= 0:
		return fmt.Errorf("%w: %.2f, must be positive", ErrInvalidWeight, weight)
	case height <= 0:
		return fmt.Errorf("%w: %.2f, must be positive", ErrInvalidHeight, height)
	case duration <= 0:
		return fmt.Errorf("%w: %v, must be positive", ErrInvalidDuration, duration)
	}
	return nil
}

// WalkingSpentCalories возвращает количество калорий, потраченных при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	runningSpentCalories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return runningSpentCalories * walkingCaloriesCoefficient, nil
}

// RunningSpentCalories возвращает количество калорий, потраченных при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateCaloriesInput(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := MeanSpeed(steps, height, duration)
	return weight * speed * duration.Minutes() / minInH, nil
}

// MeanSpeed возвращает среднюю скорость в км/ч.
// Для неположительной продолжительности или количества шагов возвращает 0.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 || steps <= 0 {
		return 0
	}
	return Distance(steps, height) / duration.Hours()
}

// Distance возвращает дистанцию в километрах. Длина шага зависит от роста в метрах.
func Distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}
