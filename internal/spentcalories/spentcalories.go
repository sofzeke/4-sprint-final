package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	// lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	stepActivityDuration := strings.Split(data, ",")

	if len(stepActivityDuration) != 3 {
		return 0, "", 0, errors.New("длина слайса не равна 3")
	}

	steps, err := strconv.Atoi(stepActivityDuration[0])
	if err != nil {
		return 0, "", 0, errors.New("не удалось преобразовать количество шагов в int")
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(stepActivityDuration[2])
	if err != nil {
		return 0, "", 0, errors.New("не удалось преобразовать строку во время")
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность активности должна быть больше 0")
	}

	return steps, stepActivityDuration[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceKm := float64(steps) * stepLength / mInKm
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)

	return distance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if activity != "Ходьба" && activity != "Бег" {
		errActivity := errors.New("неизвестный тип тренировки")
		log.Println(errActivity)
		return "", errActivity
	}

	distance := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, duration)
	kkal := 0.0

	switch activity {
	case "Ходьба":
		kkal, err = Walking(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		kkal, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	}
	str := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distance, meanSpeed, kkal)
	return str, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("все входные параметры должны быть больше 0")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * meanSpeed * durationInMinutes) / minInH, nil
}

func Walking(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("все входные параметры должны быть больше 0")
	}

	meanSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * meanSpeed * durationInMinutes) / minInH * walkingCaloriesCoefficient, nil
}
