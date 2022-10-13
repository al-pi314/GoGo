package nn

import (
	"encoding/json"
	"math"
	"math/rand"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

type Structure struct {
	InputNeurons         int
	HiddenNeuronsByLayer []int
	OutputNeurons        int
}

type MatDense struct {
	M *mat.Dense
}

type NeuralNetwork struct {
	Structure          Structure
	ActivationFuncName string
	activation         func(float64) float64

	WHiddenByLayer []MatDense
	BHiddenByLayer []MatDense

	WOut MatDense
	BOut MatDense
}

func (matDense MatDense) MarshalJSON() ([]byte, error) {
	marshalable := map[string]interface{}{}
	marshalable["rows"] = matDense.M.RawMatrix().Rows
	marshalable["cols"] = matDense.M.RawMatrix().Cols
	marshalable["data"] = matDense.M.RawMatrix().Data
	return json.Marshal(marshalable)
}

func (MatDense *MatDense) UnmarshalJSON(b []byte) error {
	marshalable := map[string]interface{}{}
	if err := json.Unmarshal(b, &marshalable); err != nil {
		return errors.Wrap(err, "custom type MatDense unmarshal error")
	}
	data := []float64{}
	unmarshaledData := marshalable["data"].([]interface{})
	for _, v := range unmarshaledData {
		data = append(data, v.(float64))
	}

	MatDense.M = mat.NewDense(marshalable["rows"].(int), marshalable["cols"].(int), data)
	return nil
}

func NewNeuralNetwork(nn NeuralNetwork) *NeuralNetwork {
	nn.activation = ActivationFunc(nn.ActivationFuncName)
	nn.WHiddenByLayer = []MatDense{}
	nn.BHiddenByLayer = []MatDense{}
	nn.WOut = MatDense{}
	nn.BOut = MatDense{}
	randomise := [][]float64{}

	prev_layer_size := nn.Structure.InputNeurons
	for _, curr_layer_size := range nn.Structure.HiddenNeuronsByLayer {
		currW := mat.NewDense(prev_layer_size, curr_layer_size, nil)
		nn.WHiddenByLayer = append(nn.WHiddenByLayer, MatDense{currW})
		randomise = append(randomise, currW.RawMatrix().Data)

		currB := mat.NewDense(1, curr_layer_size, nil)
		nn.BHiddenByLayer = append(nn.BHiddenByLayer, MatDense{currB})
		randomise = append(randomise, currB.RawMatrix().Data)

		prev_layer_size = curr_layer_size
	}

	nn.WOut.M = mat.NewDense(prev_layer_size, nn.Structure.OutputNeurons, nil)
	randomise = append(randomise, nn.WOut.M.RawMatrix().Data)
	nn.BOut.M = mat.NewDense(1, nn.Structure.OutputNeurons, nil)
	randomise = append(randomise, nn.BOut.M.RawMatrix().Data)

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
		return nn.activation(v)
	}

	prev_activation := input
	for i := range nn.Structure.HiddenNeuronsByLayer {
		hiddenLayerInput := new(mat.Dense)
		hiddenLayerInput.Mul(prev_activation, nn.WHiddenByLayer[i].M)
		hiddenLayerInput.Apply(addBaseFunction(nn.BHiddenByLayer[i].M), hiddenLayerInput)
		hiddenLayerActivations := new(mat.Dense)
		hiddenLayerActivations.Apply(applyReLU, hiddenLayerInput)

		prev_activation = hiddenLayerActivations
	}

	outputLayerInput := new(mat.Dense)
	outputLayerInput.Mul(prev_activation, nn.WOut.M)
	outputLayerInput.Apply(addBaseFunction(nn.BOut.M), outputLayerInput)
	output.Apply(applyActivation, outputLayerInput)

	return output
}

func (nn *NeuralNetwork) Mutate(rate float64) *NeuralNetwork {
	newNN := NeuralNetwork{
		Structure:          nn.Structure,
		ActivationFuncName: nn.ActivationFuncName,
		activation:         nn.activation,
		WHiddenByLayer:     []MatDense{},
		BHiddenByLayer:     []MatDense{},
		WOut:               MatDense{mat.NewDense(nn.WOut.M.RawMatrix().Rows, nn.WOut.M.RawMatrix().Cols, nil)},
		BOut:               MatDense{mat.NewDense(nn.BOut.M.RawMatrix().Rows, nn.BOut.M.RawMatrix().Cols, nil)},
	}

	mutateFunc := func(i int, j int, v float64) float64 {
		if rand.Float64() > rate {
			return -1 + 2*rand.Float64()
		}
		return v
	}

	for _, layer := range nn.WHiddenByLayer {
		newWeights := mat.NewDense(layer.M.RawMatrix().Rows, layer.M.RawMatrix().Cols, nil)
		newNN.WHiddenByLayer = append(newNN.WHiddenByLayer, MatDense{newWeights})
		newWeights.Apply(mutateFunc, layer.M)
	}

	for _, layer := range nn.BHiddenByLayer {
		newBiases := mat.NewDense(layer.M.RawMatrix().Rows, layer.M.RawMatrix().Cols, nil)
		newNN.BHiddenByLayer = append(newNN.BHiddenByLayer, MatDense{newBiases})
		newBiases.Apply(mutateFunc, layer.M)
	}

	nn.WOut.M.Apply(mutateFunc, newNN.WOut.M)
	nn.BOut.M.Apply(mutateFunc, newNN.BOut.M)

	return &newNN
}
