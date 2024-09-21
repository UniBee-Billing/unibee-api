package utility

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func PasswordEncrypt(pwd string) string {
	if len(pwd) == 0 {
		return ""
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePasswords(encryptedPwd string, plainPwd string) bool {
	if len(encryptedPwd) == 0 && len(plainPwd) == 0 {
		return true
	} else if len(encryptedPwd) == 0 || len(plainPwd) == 0 {
		return false
	}
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(encryptedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		fmt.Printf("comparePasswords err:%s\n", err.Error())
		return false
	}
	return true
}
