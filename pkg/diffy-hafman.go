package pkg

import (
	"math/big"
	"strings"
)

func ChetMod(number int64, g int64, p int64) int64 { // функция для подсчета числа в степени по модулю
	a := big.NewInt(number)                 // передаем свой секретный номер в функцию для создания эземпляра для работы с библиотекой больших чисел
	gBig := big.NewInt(g)                   // передаем общее g в функцию для создания эземпляра для работы с библиотекой больших чисел
	pBig := big.NewInt(p)                   // передаем общее p в функцию для создания эземпляра для работы с библиотекой больших чисел
	itog := new(big.Int).Exp(a, gBig, pBig) // считаем g в степени a по модулю p
	return itog.Int64()                     // возвращаем результат в формате int64
}

func encrypt(text string, shift int64) string { // шифрование на шифре цезаря
	shift = shift % 26 // так как наше число может быть очень большим, делим его на 26(число букв в Англ Алфавите), чтобы работало корректно, шифр существует лишь для примера и далеко не самый эффективный, так как не в нем суть задания
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

func decrypt(text string, shift int64) string { // дешифровщик шифра Цезаря, опять же не так важно как само задание по получению одного и того же секретного числа всеми участниками
	shift = shift % 26
	if shift == 0 {
		shift += 2
	}
	return encrypt(text, 26-shift)
}
