package genertaepasswprd

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// generatePassword генерирует надёжный пароль заданной длины
func GeneratePassword(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("длина пароля должна быть больше 0")
	}

	// Определяем символы, которые будут использоваться для генерации пароля
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?/"
	charLength := big.NewInt(int64(len(characters)))
	var password []byte

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, charLength)
		if err != nil {
			return "", fmt.Errorf("ошибка генерации случайного числа: %v", err)
		}
		password = append(password, characters[index.Int64()])
	}

	return string(password), nil
}


