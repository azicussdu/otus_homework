package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func Top10(text string) []string {
	// тут буду считать количество слов в map
	wordsMap := make(map[string]uint16)

	words := strings.Fields(text)

	for _, word := range words {
		result := trimWord(word)
		if result != "-" {
			wordsMap[result]++
		}
	}

	var maxIndex uint16

	/* в countSlice я буду хранить слова как значения,
	а как индексы буду хранить количество повторении.

	[][]string - так как некоторые слова могут повторяться одинаковое количество раз
	*/
	countSlice := make([][]string, len(words)+1)

	for key, value := range wordsMap {
		if value > maxIndex {
			maxIndex = value
		}

		countSlice[value] = append(countSlice[value], key)
	}

	// k := 10 так как нужно вывести 10 топ встречающихся слов
	k := 10
	var resultSlice []string
	for i := maxIndex; i > 0; i-- {
		// сортируем каждый slice строк лексографический
		sort.Strings(countSlice[i])

		for j := 0; j < len(countSlice[i]); j++ {
			resultSlice = append(resultSlice, countSlice[i][j])
			k--
			if k == 0 {
				return resultSlice
			}
		}
	}

	return resultSlice
}

// функция для трима слова (удаляет по краям знаки кроме '-').
func trimWord(w string) string {
	return strings.TrimFunc(strings.ToLower(w), func(r rune) bool {
		// удаляем все что не буква и не дефис '-'
		// я так понял дефис удалять нельзя так как '-----' считается как слово

		return unicode.IsPunct(r) && r != '-'
	})
}
