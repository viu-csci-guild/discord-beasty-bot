package responses

import (
	"log"
	"math/rand"
	"time"
)

type Responses struct {
	modelData map[interface{}]interface{}
}

// responses fetches a random string based on a lookup
func NewResponses(model map[interface{}]interface{}) *Responses {
	rand.Seed(time.Now().UnixNano())
	r := &Responses{
		modelData: model,
	}
	return r
}

// returns response string randomly selected based on lookup
func (r Responses) GenerateResponse(lookup string) string {
	// fetch map and cast to array of strings
	respMap, valid := r.modelData[lookup].([]interface{})
	if !valid {
		log.Fatalf("Error: could not type response as array of interfaces")
	}
	respArr := make([]string, 0, len(respMap))
	for _, v := range respMap {
		respArr = append(respArr, v.(string))
	}
	response := respArr[rand.Intn(len(respArr))]
	return response
}
