package pkg

import (
	"math/big"
	"strings"
)

func generateYouNumber(number int64, g int64, p int64) int64 {
	a := big.NewInt(number)
	gBig := big.NewInt(g)
	pBig := big.NewInt(p)
	itog := new(big.Int).Exp(gBig, a, pBig)
	return itog.Int64()
}

func getSecretNumber(myNumber int64, otherNumber int64, p int64) int64 {
	a := big.NewInt(myNumber)
	b := big.NewInt(otherNumber)
	pBig := big.NewInt(p)
	itog := new(big.Int).Exp(b, a, pBig)
	return itog.Int64()
}

func encrypt(text string, shift int64) string {
	shift = shift % 26
	if shift == 0 {
		shift += 2
	}
	var encrypted strings.Builder
	for _, char := range text {
		if char >= 'A' && char <= 'Z' {
			char = 'A' + (char-'A'+rune(shift))%26
		} else if char >= 'a' && char <= 'z' {
			char = 'a' + (char-'a'+rune(shift))%26
		}
		encrypted.WriteRune(char)
	}
	return encrypted.String()
}

func decrypt(text string, shift int64) string {
	shift = shift % 26
	if shift == 0 {
		shift += 2
	}
	return encrypt(text, 26-shift)
}
