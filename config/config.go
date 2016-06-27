package config

var config configger

func init() {
	config = newTomlAdapter()
}

type configger interface {
	Decode(interface{}) error
}

func Decode(v interface{}) error {
	return config.Decode(v)
}
