package main

import (
	"math/rand"
	"time"
)

type responses struct {
	modelData map[interface{}]interface{}
}

// responses fetches a random string based on a lookup
func newResponses(model map[interface{}]interface{}) *responses {
	rand.Seed(time.Now().UnixNano())
	r := &responses{
		modelData: model,
	}
	return r
}

// TODO: not be dumb
func (r responses) generateResponse(lookup string) string {
	return "placeholder"
	// fetch map and cast to array of strings
	// respMap := r.modelData[lookup]
	// respArr, valid := respMap.([]string)
	// if !valid {
	// 	log.Fatalf("Error: could not convert response map to array")
	// }
	// response := respArr[rand.Intn(len(respArr))]
	// return response
}
