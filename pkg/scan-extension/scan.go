package scanextension

import (
	"unicode"
	"unicode/utf8"
)

func ScanWordsOnly(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Валидными символами(те, которые считаются единым словом) являются буквы и цифры
	// Лично я бы цифры убрал, но в приложенном тесте слово со слитной цифрой считалось единым словом
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			break
		}
	}
	// Сканит до тех пор пока не возникнет любой символ кроме буквы или цифры
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return i + width, data[start:i], nil
		}
	}
	// EOF - конец, можно возвращать весь результат
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// Обработка дальше
	return start, nil, nil
}
