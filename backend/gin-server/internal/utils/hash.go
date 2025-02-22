package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// accepts byte slice and cost (controls how complex the hash is)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

// we can take the plain text password that we got on the log in route and see if this hash, which we stored in the database could have been generated from that password, which will tell us that the password is valid.
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	// if there is no error, the password is valid (true)
	// if there is an error, the password is invalid (false)
	return err == nil
}
