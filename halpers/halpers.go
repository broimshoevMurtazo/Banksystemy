package halpers

import (
	"math/rand"
	"time"
)

func RandomSixDigit() int {
	// Устанавливаем seed для случайного генератора
	rand.Seed(time.Now().UnixNano())
	// Генерируем случайное число от 100000 до 999999
	return rand.Intn(900000) + 100000
}





  
  
  