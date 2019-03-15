package datamodels

type Result struct {
	Status bool
	Code   string
	Msg    string
	Data   interface{}
}
