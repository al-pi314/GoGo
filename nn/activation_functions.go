package nn

import "math"

var activationFunc = map[string]func(float64) float64{
	"SIGMOID": sigmoid,
}

func ActivationFunc(name string) func(float64) float64 {
	if f, ok := activationFunc[name]; ok {
		return f
	}
	return sigmoid
}

func sigmoid(v float64) float64 {
	return 1 / (1.0 + math.Exp(-v))
}
