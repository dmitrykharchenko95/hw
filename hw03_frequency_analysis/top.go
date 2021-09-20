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

	sort.Slice(inputSlice, func(i, j int) bool {
		return inputSlice[i] < inputSlice[j]
	})

	amountData := []StringAndAmount{}
	amount := 1

	for i := 0; i < len(inputSlice); {
		switch {
		case i == len(inputSlice)-1:
			amountData = append(amountData, StringAndAmount{inputSlice[i], amount})
			i++
		case inputSlice[i] == inputSlice[i+1]:
			i++
			amount++
			continue
		case inputSlice[i] != inputSlice[i+1]:
			amountData = append(amountData, StringAndAmount{inputSlice[i], amount})
			i++
			amount = 1
		}
	}
	sort.Slice(amountData, func(i, j int) bool {
		if amountData[i].Amount != amountData[j].Amount {
			return amountData[i].Amount > amountData[j].Amount
		}
		return amountData[i].String < amountData[j].String
	})

	var outputSlice []string
	num := 10
	if len(amountData) < 10 {
		num = len(amountData)
	}
	for i := 0; i < num; i++ {
		outputSlice = append(outputSlice, amountData[i].String)
	}

	return outputSlice
}
