package hunter

type HunterOpt struct{}

type HunterObj struct{}

func New(opt HunterOpt) *HunterObj {
	obj := HunterObj{}

	return &obj
}
