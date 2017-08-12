package randomforest

import (
	"log"
	"math"
	"sort"
)

type Node struct {
	left      *Node
	right     *Node
	feature   int
	label     float64
	gain      float64
	threshold float64
	depth     float64
}

func NewNode() *Node {
	n := &Node{}
	n.gain = 0.0
	n.threshold = 0.0
	return n
}

func (n *Node) Grow(features [][]float64, labels []float64, fun Criterion) {
	uniq := unique(labels)
	size := len(features)
	featureN := len(features[0])

	// 全データが同一クラスの場合は終了
	if len(uniq) == 1 {
		n.label = labels[0]
		return
	}

	selectedFeatures := selectRandomFeatures(featureN, 2)

	if len(features) == 0 {
		return
	}

	n.getNodeLabel(labels)

	nodeImpurity := fun(labels)

	for _, featureCol := range selectedFeatures {
		level := unique(selectCol(features, featureCol))
		sort.Float64s(level)
		thresholds := getThresholds(level)

		// search best split
		for i := 0; i < len(thresholds); i++ {
			_, _, lhsl, rhsl := divide(features, labels, featureCol, thresholds[i])

			leftRatio := float64(len(lhsl)) / float64(size)
			rightRatio := float64(len(rhsl)) / float64(size)

			informationGain := nodeImpurity - (fun(lhsl)*leftRatio + fun(rhsl)*rightRatio)

			if informationGain > n.gain {
				n.gain = informationGain
				n.feature = featureCol
				n.threshold = thresholds[i]
			}
		}
	}

	log.Println("size:", size, "depth:", n.depth, "label:", n.label, "feature:", n.feature, "th:", n.threshold, "gain:", n.gain)

	// partition not found
	if n.gain == 0.0 {
		return
	}

	lhs, rhs, lhsl, rhsl := divide(features, labels, n.feature, n.threshold)

	n.left = NewNode()
	n.left.depth = n.depth + 1
	n.left.Grow(lhs, lhsl, fun)

	n.right = NewNode()
	n.right.depth = n.depth + 1
	n.right.Grow(rhs, rhsl, fun)
}

func (n *Node) getNodeLabel(labels []float64) {
	uniq := unique(labels)

	cnt := 0
	for i := 0; i < len(uniq); i++ {
		if cnt < count(labels, uniq[i]) {
			cnt = count(labels, uniq[i])
			n.label = uniq[i]
		}
	}
}

func (n *Node) Predict(feature []float64) float64 {
	if n.left != nil && n.right != nil {
		if feature[n.feature] <= n.threshold {
			return n.left.Predict(feature)
		} else {
			return n.right.Predict(feature)
		}
	}

	return n.label
}

func gini(labels []float64) float64 {
	n := len(labels)
	classes := unique(labels)
	gini := 1.0

	for i := 0; i < len(classes); i++ {
		cnt := 0
		for j := 0; j < n; j++ {
			if labels[j] == classes[i] {
				cnt++
			}
		}
		gini -= math.Pow(float64(cnt)/float64(n), 2.0)
	}

	return gini
}

func entropy() {}

func error() {}

func mse() {}