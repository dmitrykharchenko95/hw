package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
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
	punctuationMarks := []string{"!", ",", ".", ":", "`", " -", `"`, "(", ")", " —"}
	for _, val := range punctuationMarks {
		inputString = strings.ReplaceAll(inputString, val, "")
	}
	inputSlice := strings.Fields(inputString)
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

	fmt.Println(amountData)

	var outputSlice []string
	var num int
	if len(amountData) < 10 {
		num = len(amountData)
	} else {
		num = 10
	}
	for i := 0; i < num; i++ {
		outputSlice = append(outputSlice, amountData[i].String)
	}

	return outputSlice
}
