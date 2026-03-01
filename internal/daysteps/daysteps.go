package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	if data == "" {
		log.Println("передана пустая строка")
		return 0, 0, errors.New("передана пустая строка")
	}

	stepAndDuration := strings.Split(data, ",")

	if len(stepAndDuration) != 2 {
		log.Println("некорректный формат")
		return 0, 0, errors.New("длина слайса не равна 2")
	}

	steps, err := strconv.Atoi(stepAndDuration[0])
	if err != nil {
		return 0, 0, errors.New("не удалось преобразовать количество шагов в int")
	}

	if steps <= 0 {
		log.Println("шаги <= 0")
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(stepAndDuration[1])
	if err != nil {
		return 0, 0, errors.New("не удалось преобразовать строку во время")
	}

	if duration <= 0 {
		log.Println("продолжительность <= 0")
		return 0, 0, errors.New("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	if duration <= 0 {
		log.Println("продолжительность <= 0")
		return ""
	}

	distanceKm := float64(steps) * stepLength / mInKm

	kkal, err := spentcalories.Walking(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	str := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, kkal)
	return str
}
