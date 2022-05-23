package signals

type Input interface {
	Value() bool
}

type Output interface {
	Set(bool)
}

type IO interface{
	Input
	Output
}