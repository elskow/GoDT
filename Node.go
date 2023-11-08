package GoDT

import (
	"fmt"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"strconv"
	"strings"
)

type Data struct {
	rule         string
	feature      []string
	counter      map[string]int
	giniImpurity float64
	label        string
	nb           int
}

type Node struct {
	Y           []int
	X           dataframe.DataFrame
	bestFeature string
	bestValue   string
	depth       int
	maxDepth    int
	minDfSplit  int
	data        *Data
	left        *Node
	right       *Node
}

func NodeInit(X dataframe.DataFrame, Y []int, depth int, maxDepth int, minDfSplit int, rule string) *Node {
	if len(Y) == 0 && X.Nrow() == 0 {
		return nil
	}

	if len(Y) != X.Nrow() {
		panic("Y and X are not the same length")
	}

	generateCounter := count(Y)

	node := &Node{
		Y:     Y,
		X:     X,
		depth: depth,
		data: &Data{
			rule:         rule,
			feature:      X.Names(),
			counter:      generateCounter,
			giniImpurity: giniImpurity(generateCounter["0"], generateCounter["1"]),
			label:        maxCounts(generateCounter),
			nb:           len(Y),
		},
		left:        nil,
		right:       nil,
		bestFeature: "",
		bestValue:   "",
	}

	return node
}

func (node *Node) Split() (string, string) {
	dataFrame := node.X.Copy().Mutate(series.New(node.Y, series.Int, "Y"))

	giniBase := node.data.giniImpurity
	maxGain := 0.0
	bestFeature := ""
	bestValue := ""

	for _, feature := range node.data.feature {
		sorted := dataFrame.Arrange(dataframe.Sort(feature))
		xmean := meth(uniqueGotaSeries(sorted.Col(feature)).Float())

		for _, value := range xmean {
			leftCounter := countError(sorted.Filter(
				dataframe.F{Colname: feature, Comparator: series.Less, Comparando: value},
			).Col("Y").Int())

			rightCounter := countError(sorted.Filter(
				dataframe.F{Colname: feature, Comparator: series.Greater, Comparando: value},
			).Col("Y").Int())

			s0left, s1left, s0right, s1right := leftCounter["0"], leftCounter["1"], rightCounter["0"], rightCounter["1"]

			totalLeft := s0left + s1left
			totalRight := s0right + s1right

			weightLeft := float64(totalLeft) / float64(totalLeft+totalRight)
			weightRight := float64(totalRight) / float64(totalLeft+totalRight)

			weightedGini := weightLeft*giniImpurity(s0left, s1left) + weightRight*giniImpurity(s0right, s1right)

			giniGain := giniBase - weightedGini

			if giniGain >= maxGain {
				maxGain = giniGain
				bestFeature = feature
				bestValue = fmt.Sprint(value)
			}
		}
	}

	return bestFeature, bestValue
}

func (node *Node) sprout() {
	if node.depth < node.maxDepth && node.data.nb >= node.minDfSplit {
		node.bestFeature, node.bestValue = node.Split()

		dataFrame := node.X.Copy().Mutate(series.New(node.Y, series.Int, "Y"))

		if node.bestFeature != "" {
			panic("Node.bestFeature is empty")
		}

		leftDataFrame, rightDataFrame := dataFrame.Filter(
			dataframe.F{Colname: node.bestFeature, Comparator: series.LessEq, Comparando: node.bestValue},
		).Copy(), dataFrame.Filter(
			dataframe.F{Colname: node.bestFeature, Comparator: series.Greater, Comparando: node.bestValue},
		).Copy()

		leftLabel, err := leftDataFrame.Col("Y").Int()
		if err != nil {
			panic(err)
		}

		rightLabel, err := rightDataFrame.Col("Y").Int()
		if err != nil {
			panic(err)
		}

		node.left = NodeInit(leftDataFrame, leftLabel, node.depth+1, node.maxDepth, node.minDfSplit, fmt.Sprintf("%s <= %s", node.bestFeature, node.bestValue))

		if node.left != nil {
			node.left.sprout()
		}

		node.right = NodeInit(rightDataFrame, rightLabel, node.depth+1, node.maxDepth, node.minDfSplit, fmt.Sprintf("%s > %s", node.bestFeature, node.bestValue))

		if node.right != nil {
			node.right.sprout()
		}
	}
}

func (node *Node) print(padding int) {
	whitespace := strings.Repeat("-", padding*node.depth)
	fmt.Print(whitespace)

	if node.data.rule == "ROOT" {
		fmt.Println("ROOT")
	} else {
		fmt.Printf("Split on %s\n", node.data.rule)
	}

	fmt.Print(strings.Repeat(" ", padding*node.depth))
	fmt.Println("\t| Gini Impurity:", node.data.giniImpurity)
	fmt.Print(strings.Repeat(" ", padding*node.depth))
	fmt.Printf("\t| Class Distribution: %d/%d\n", node.data.counter["0"], node.data.counter["1"])
	fmt.Print(strings.Repeat(" ", padding*node.depth))
	fmt.Printf("\t| Predicted Class: %s\n", node.data.label)

	if node.left != nil {
		node.left.print(padding + 2)
	}

	if node.right != nil {
		node.right.print(padding + 2)
	}
}

func (node *Node) predict(X map[string]float64) string {
	if node.left == nil && node.right == nil {
		return node.data.label
	}

	parsed, err := strconv.ParseFloat(node.bestValue, 64)
	if err != nil {
		panic(err)
	}

	if X[node.bestFeature] <= parsed {
		return node.left.predict(X)
	} else {
		return node.right.predict(X)
	}
}
