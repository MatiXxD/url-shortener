package tokengen

import "math/rand"

const symbols = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func GenerateToken(tokenSize int) string {
	l := len(symbols)
	token := []byte{}
	for range tokenSize {
		token = append(token, []byte(symbols)[rand.Intn(l)])
	}
	return string(token)
}
