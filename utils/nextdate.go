package utils

import (
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date time.Time, repeatStr string) (time.Time, error) {
	var nextDate time.Time
	// Обработка повторения по дням
	if strings.HasPrefix(repeatStr, "d ") {
		daysStr := strings.TrimPrefix(repeatStr, "d ")
		days, err := strconv.Atoi(daysStr)
		if err != nil || days < 0 || days > 400 {
			return time.Time{}, nil // Вернуть пустое значение, если days недопустимо
		}
		nextDate = date.AddDate(0, 0, days) // Добавляем дни	// Проверяем, чтобы следующая дата была больше текущей
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
	} else if strings.HasPrefix(repeatStr, "y") {
		// Обработка повторения по годам
		nextDate = date.AddDate(1, 0, 0) // Добавляем 1 год	// Проверяем, чтобы следующая дата была больше текущей
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	} else {
		return time.Time{}, nil // Вернуть пустое значение, если repeat недопустимо
	}

	return nextDate, nil

}
