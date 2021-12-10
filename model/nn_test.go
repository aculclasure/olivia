package model

import (
	"math/rand"
	"testing"
)

func TestCreateNeuralNetwork(t *testing.T) {
	nn := CreateNeuralNetwork(0.5, 10, 20, 12, 30)

	for i := 0; i < 3; i++ {
		if nn.Layers[i].Columns() != nn.Weights[i].Rows() {
			t.Errorf("Error with the creation of the RNN, the size of layers/weights is wrong.")
		}
	}
}

func TestNNFeedForward(t *testing.T) {
	nn := CreateNeuralNetwork(0.01, 1, 1)
	rand.Seed(5)
	output := nn.FeedForward([]float64{0.5})

	if output[0][0] != 0.3254432307595214 {
		t.Errorf("Error with the feed forward, the output is wrong.")
	}
}