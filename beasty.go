package beasty

type beasty struct {
	token        string
	watchChannel string
	discordbot   struct{}
	config       struct{}
	responses    struct{}
	storage      struct{}
}

func NewBeasty(t string, w string) *beasty {
	b := &beasty{
		token:        t,
		watchChannel: w,
	}
	return b
}

func (b beasty) Start() {
	return
}
