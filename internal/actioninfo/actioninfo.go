// Package actioninfo выводит сводную информацию об активностях пользователя:
// тренировках и прогулках.
package actioninfo

import (
	"fmt"
	"log"
)

// DataParser описывает активность, которая умеет разбирать строку с данными
// и формировать по ним отчёт. Реализуется типами trainings.Training
// и daysteps.DaySteps.
type DataParser interface {
	Parse(datastring string) error
	ActionInfo() (string, error)
}

// Info разбирает каждую строку из dataset с помощью dp и выводит отчёт
// об активности в стандартный вывод. Строки, которые не удалось разобрать
// или обработать, пропускаются, а ошибки пишутся в лог.
func Info(dataset []string, dp DataParser) {
	for _, s := range dataset {
		err := dp.Parse(s)
		if err != nil {
			log.Printf("parse %q: %v", s, err)
			continue
		}

		output, err := dp.ActionInfo()
		if err != nil {
			log.Printf("action info for %q: %v", s, err)
			continue
		}
		fmt.Print(output)
	}
}
