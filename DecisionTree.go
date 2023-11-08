package GoDT

import "github.com/go-gota/gota/dataframe"

type DecisionTree struct {
	Root *Node
}

func TreeInit(X dataframe.DataFrame, Y []int, maxDepth int, minDfSplit int) *DecisionTree {
	return &DecisionTree{
		Root: NodeInit(X, Y, 0, maxDepth, minDfSplit, "ROOT"),
	}
}

func (tree *DecisionTree) Sprout() {
	tree.Root.sprout()
}

func (tree *DecisionTree) Predict(X dataframe.DataFrame) []string {
	features := tree.Root.data.feature

	var predictions []string

	x, _ := X.Dims()

	for i := 0; i < x; i++ {
		nmap := make(map[string]float64)
		for _, feature := range features {
			nmap[feature] = X.Col(feature).Elem(i).Float()
		}
		predictions = append(predictions, tree.Root.predict(nmap))
	}

	return predictions
}

func (tree *DecisionTree) Print() {
	tree.Root.print(1)
}
