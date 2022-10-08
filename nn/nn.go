package nn

import (
	"math"
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

type Structure struct {
	InputNeurons         int
	HiddenNeuronsByLayer []int
	OutputNeurons        int
}

type NeuralNetwork struct {
	Structure Structure

	wHiddenByLayer []*mat.Dense
	bHiddenByLayer []*mat.Dense

	wOut *mat.Dense
	bOut *mat.Dense
}

func NewNeuralNetwork(s Structure) *NeuralNetwork {
	nn := NeuralNetwork{
		Structure: s,
	}

	nn.wHiddenByLayer = []*mat.Dense{}
	nn.bHiddenByLayer = []*mat.Dense{}
	randomise := [][]float64{}

	prev_layer_size := s.InputNeurons
	for _, curr_layer_size := range s.HiddenNeuronsByLayer {
		currW := mat.NewDense(prev_layer_size, curr_layer_size, nil)
		nn.wHiddenByLayer = append(nn.wHiddenByLayer, currW)
		randomise = append(randomise, currW.RawMatrix().Data)

		currB := mat.NewDense(1, curr_layer_size, nil)
		nn.bHiddenByLayer = append(nn.bHiddenByLayer, currB)
		randomise = append(randomise, currB.RawMatrix().Data)

		prev_layer_size = curr_layer_size
	}

	nn.wOut = mat.NewDense(prev_layer_size, s.OutputNeurons, nil)
	randomise = append(randomise, nn.wOut.RawMatrix().Data)
	nn.bOut = mat.NewDense(1, s.OutputNeurons, nil)
	randomise = append(randomise, nn.bOut.RawMatrix().Data)

	for _, param := range randomise {
		for i := range param {
			param[i] = -1 + 2*rand.Float64()
		}
	}

	return &nn
}

func (nn *NeuralNetwork) Predict(input *mat.Dense) *mat.Dense {
	output := new(mat.Dense)

	// helper functions
	addBaseFunction := func(source *mat.Dense) func(int, int, float64) float64 {
		return func(_, col int, v float64) float64 { return v + source.At(0, col) }
	}
	applyReLU := func(_, _ int, v float64) float64 {
		return math.Max(0, v)
	}
	applySigmoid := func(_, _ int, v float64) float64 {
		return 1.0 / (1.0 + math.Exp(-v))
	}

	prev_activation := input
	for i := range nn.Structure.HiddenNeuronsByLayer {
		hiddenLayerInput := new(mat.Dense)
		hiddenLayerInput.Mul(prev_activation, nn.wHiddenByLayer[i])
		hiddenLayerInput.Apply(addBaseFunction(nn.bHiddenByLayer[i]), hiddenLayerInput)
		hiddenLayerActivations := new(mat.Dense)
		hiddenLayerActivations.Apply(applyReLU, hiddenLayerInput)

		prev_activation = hiddenLayerActivations
	}

	outputLayerInput := new(mat.Dense)
	outputLayerInput.Mul(prev_activation, nn.wOut)
	outputLayerInput.Apply(addBaseFunction(nn.bOut), outputLayerInput)
	output.Apply(applySigmoid, outputLayerInput)

	return output
}
