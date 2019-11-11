package util

import (
	"math/rand"
	"time"
)

func RandomBooleanByRatio(ratio float64) bool{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	flagArray := make([]bool, 10)
	for i, _ := range flagArray{
		if i < int(ratio*10){
			flagArray[i] = true
		}
	}
	randomIndex := r.Intn(10)
	result := flagArray[randomIndex]
	return result
}