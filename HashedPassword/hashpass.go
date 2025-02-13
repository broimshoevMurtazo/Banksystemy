package hashedpasswod


import (
	"golang.org/x/crypto/bcrypt"
)
func HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// ? ============ Compare The Password
func CompareHashPasswords(HashedPasswordFromDB, PasswordToCampare string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPasswordFromDB), []byte(PasswordToCampare))
	return err == nil
}