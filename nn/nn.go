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
	Structure      Structure
	ActivationFunc func(float64) float64

	wHiddenByLayer []*mat.Dense
	bHiddenByLayer []*mat.Dense

	wOut *mat.Dense
	bOut *mat.Dense
}

func NewNeuralNetwork(nn NeuralNetwork) *NeuralNetwork {
	nn.wHiddenByLayer = []*mat.Dense{}
	nn.bHiddenByLayer = []*mat.Dense{}
	randomise := [][]float64{}

	prev_layer_size := nn.Structure.InputNeurons
	for _, curr_layer_size := range nn.Structure.HiddenNeuronsByLayer {
		currW := mat.NewDense(prev_layer_size, curr_layer_size, nil)
		nn.wHiddenByLayer = append(nn.wHiddenByLayer, currW)
		randomise = append(randomise, currW.RawMatrix().Data)

		currB := mat.NewDense(1, curr_layer_size, nil)
		nn.bHiddenByLayer = append(nn.bHiddenByLayer, currB)
		randomise = append(randomise, currB.RawMatrix().Data)

		prev_layer_size = curr_layer_size
	}

	nn.wOut = mat.NewDense(prev_layer_size, nn.Structure.OutputNeurons, nil)
	randomise = append(randomise, nn.wOut.RawMatrix().Data)
	nn.bOut = mat.NewDense(1, nn.Structure.OutputNeurons, nil)
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
	applyActivation := func(_, _ int, v float64) float64 {
		return nn.ActivationFunc(v)
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
	output.Apply(applyActivation, outputLayerInput)

	return output
}

func (nn *NeuralNetwork) Mutate(rate float64) *NeuralNetwork {
	newNN := NeuralNetwork{
		Structure:      nn.Structure,
		ActivationFunc: nn.ActivationFunc,
		wHiddenByLayer: []*mat.Dense{},
		bHiddenByLayer: []*mat.Dense{},
		wOut:           &mat.Dense{},
		bOut:           &mat.Dense{},
	}

	mutateFunc := func(i int, j int, v float64) float64 {
		if rand.Float64() > rate {
			return -1 + 2*rand.Float64()
		}
		return v
	}

	for _, layer := range nn.wHiddenByLayer {
		newWeights := mat.Dense{}
		newNN.wHiddenByLayer = append(newNN.wHiddenByLayer, &newWeights)
		layer.Apply(mutateFunc, &newWeights)
	}

	for _, layer := range nn.bHiddenByLayer {
		newBiases := mat.Dense{}
		newNN.bHiddenByLayer = append(newNN.bHiddenByLayer, &newBiases)
		layer.Apply(mutateFunc, &newBiases)
	}

	nn.wOut.Apply(mutateFunc, newNN.wOut)
	nn.bOut.Apply(mutateFunc, newNN.bOut)

	return &newNN
}
