package beasty

type responses struct {
	data struct{}
}

func newResponses() *responses {
	r := &responses{}
	return r
}

func (r responses) generateResponse(lookup string) string {
	return ""
}
