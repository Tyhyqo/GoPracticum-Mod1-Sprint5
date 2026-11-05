// Package personaldata описывает персональные данные пользователя,
// которые нужны для расчёта показателей активности.
package personaldata

import "fmt"

// Personal хранит персональные данные пользователя.
// Предназначена для встраивания в структуры с данными о тренировках и прогулках.
type Personal struct {
	Name   string  // имя пользователя
	Weight float64 // вес в килограммах
	Height float64 // рост в метрах
}

// Print выводит имя, вес и рост пользователя в стандартный вывод.
func (p Personal) Print() {
	fmt.Printf("Имя: %s\nВес: %.2f\nРост: %.2f\n", p.Name, p.Weight, p.Height)
}
