package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type StringAndAmount struct {
	String string
	Amount int
}

func Top10(inputString string) []string {
	if len(inputString) == 0 {
		return []string{}
	}
	inputString = strings.ToLower(inputString)
	inputString = strings.ReplaceAll(inputString, "- ", "")
	inputSlice := strings.FieldsFunc(inputString, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-'
	})

	inputMap := make(map[string]int)

	for _, v := range inputSlice {
		inputMap[v]++
	}

	fullOutputSlice := []StringAndAmount{}
	for key, val := range inputMap {
		fullOutputSlice = append(fullOutputSlice, StringAndAmount{key, val})
	}

	sort.Slice(fullOutputSlice, func(i, j int) bool {
		if fullOutputSlice[i].Amount != fullOutputSlice[j].Amount {
			return fullOutputSlice[i].Amount > fullOutputSlice[j].Amount
		}
		return fullOutputSlice[i].String < fullOutputSlice[j].String
	})

	var outputSlice []string
	num := 10
	if len(fullOutputSlice) < 10 {
		num = len(fullOutputSlice)
	}
	for i := 0; i < num; i++ {
		outputSlice = append(outputSlice, fullOutputSlice[i].String)
	}

	return outputSlice
}
