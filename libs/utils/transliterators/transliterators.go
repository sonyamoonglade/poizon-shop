package transliterators

import "unicode/utf8"

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

	"Ъ": "#",
	"ъ": "@",

	"Ы": "!",
	"ы": "+",

	"Ю": "]",
	"ю": "[",

	"Я": "}",
	"я": "{",
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

	"#": "Ъ",
	"@": "ъ",

	"!": "Ы",
	"+": "ы",

	"]": "Ю",
	"[": "ю",

	"}": "Я",
	"{": "я",
}

func Encode(strs []string) []string {
	out := make([]string, 0, len(strs))
	for _, sentence := range strs {
		if onlyEnglishOrNumberLetters(sentence) {
			out = append(out, sentence)
			continue
		}
		sub := make([]byte, 0, len(sentence))
		for _, ch := range []rune(sentence) {
			if dch, ok := dict[string(ch)]; ok {
				sub = append(sub, []byte(dch)...)
			} else {
				sub = utf8.AppendRune(sub, ch)
			}
		}
		out = append(out, string(sub))
	}

	return out
}

func Decode(strs []string) []string {
	out := make([]string, 0, len(strs))
	for _, sentence := range strs {
		if onlyEnglishOrNumberLetters(sentence) {
			out = append(out, sentence)
			continue
		}
		sub := make([]byte, 0, len(sentence))
		for _, ch := range []rune(sentence) {
			if dch, ok := reversedDict[string(ch)]; ok {
				sub = append(sub, []byte(dch)...)
			} else {
				sub = utf8.AppendRune(sub, ch)
			}
		}
		out = append(out, string(sub))
	}

	return out
}

func onlyEnglishOrNumberLetters(sentence string) bool {

	for _, ch := range sentence {
		if !(('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('0' <= ch && ch <= '9') || ch == ' ') {
			return false
		}
	}
	return true
}
