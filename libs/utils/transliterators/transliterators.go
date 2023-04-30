package transliterators

import (
	"unicode/utf8"
	"unsafe"

	fn "github.com/sonyamoonglade/go_func"
)

var dict = map[string]string{
	"А": "A",
	"а": "a",

	"Б": "B",
	"б": "b",

	"В": "V",
	"в": "v",

	"Г": "G",
	"г": "g",

	"Д": "D",
	"д": "d",

	"Е": "E",
	"е": "e",

	"Ё": ">",
	"ё": "<",

	"Ж": "J",
	"ж": "j",

	"З": "Z",
	"з": "z",

	"И": "I",
	"и": "i",

	"Й": "~",
	"й": "*",

	"К": "K",
	"к": "k",

	"Л": "L",
	"л": "l",

	"М": "M",
	"м": "m",

	"Н": "N",
	"н": "n",

	"О": "O",
	"о": "o",

	"П": "P",
	"п": "p",

	"Р": "R",
	"р": "r",

	"С": "S",
	"с": "s",

	"Т": "T",
	"т": "t",

	"У": "U",
	"у": "u",

	"Ф": "F",
	"ф": "f",

	"Х": "H",
	"х": "h",

	"Ч": ")",
	"ч": "(",

	"Ш": "&",
	"ш": "?",

	"Щ": "$",
	"щ": "%",

	"ц": "#",
	"Ц": "!",

	"ъ": "@",
	"ы": "+",

	"Ю": "]",
	"ю": "[",

	"Я": "}",
	"я": "{",

	"ь": ".",

	"э": ",",
}

var reversedDict = map[string]string{
	"A": "А",
	"a": "а",

	"B": "Б",
	"b": "б",

	"V": "В",
	"v": "в",

	"G": "Г",
	"g": "г",

	"D": "Д",
	"d": "д",

	"E": "Е",
	"e": "е",

	">": "Ё",
	"<": "ё",

	"J": "Ж",
	"j": "ж",

	"Z": "З",
	"z": "з",

	"I": "И",
	"i": "и",

	"~": "Й",
	"*": "й",

	"K": "К",
	"k": "к",

	"L": "Л",
	"l": "л",

	"M": "М",
	"m": "м",

	"N": "Н",
	"n": "н",

	"O": "О",
	"o": "о",

	"P": "П",
	"p": "п",

	"R": "Р",
	"r": "р",

	"S": "С",
	"s": "с",

	"T": "Т",
	"t": "т",

	"U": "У",
	"u": "у",

	"F": "Ф",
	"f": "ф",

	"H": "Х",
	"h": "х",

	")": "Ч",
	"(": "ч",

	"&": "Ш",
	"?": "ш",

	"$": "Щ",
	"%": "щ",

	"#": "ц",
	"!": "Ц",

	"@": "ъ",

	"+": "ы",

	"]": "Ю",
	"[": "ю",

	"}": "Я",
	"{": "я",

	".": "ь",

	",": "э",
}

const (
	noEncPrefixByte = '='
	noEncPrefixStr  = "="
	spaceSeparator  = " "
)

func Encode(values []string) []string {
	return fn.
		Of(values).
		Map(func(value string, _ int) string {
			return fn.
				OfSplitString(value, spaceSeparator).
				Map(func(word string, _ int) string {
					if isNumerical(word) {
						return word
					}
					if isEnglishOrNumerical(word) {
						return noEncPrefixStr + word
					}
					return fn.Compose2(byteSliceToStr, encodeToAscii)(word)
				}).
				Join(spaceSeparator)
		}).
		Result
}

func Decode(encoded []string) []string {
	return fn.
		Of(encoded).
		Map(func(value string, _ int) string {
			return fn.
				OfSplitString(value, spaceSeparator).
				Map(func(word string, _ int) string {
					if isNumerical(word) {
						return word
					}
					if isNotEncoded(word) {
						return word[1:]
					}
					return fn.Compose2(byteSliceToStr, decodeFromAscii)(word)
				}).
				Join(spaceSeparator)
		}).
		Result
}

func isEnglishOrNumerical(word string) bool {
	for _, ch := range word {
		if !(('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('0' <= ch && ch <= '9')) {
			return false
		}
	}
	return true
}

func isNumerical(word string) bool {
	for _, ch := range word {
		if !('0' <= ch && ch <= '9') {
			return false
		}
	}
	return true
}

func encodeToAscii(s string) []byte {
	return fn.Reduce(func(acc []byte, ch rune, _ int) []byte {
		if encoded, ok := dict[string(ch)]; ok {
			return append(acc, []byte(encoded)...)
		}
		return utf8.AppendRune(acc, ch)
	}, []rune(s), make([]byte, 0, len(s)))
}

func decodeFromAscii(s string) []byte {
	return fn.Reduce(func(acc []byte, ch rune, _ int) []byte {
		if decoded, ok := reversedDict[string(ch)]; ok {
			return append(acc, []byte(decoded)...)
		}
		return utf8.AppendRune(acc, ch)
	}, []rune(s), make([]byte, 0, len(s)))
}

func byteSliceToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func isNotEncoded(s string) bool {
	return s[0] == noEncPrefixByte
}
