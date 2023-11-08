package GoDT

import (
	"fmt"
	"github.com/go-gota/gota/series"
	"math"
	"sort"
)

type Counter map[string]int

func count(Y []int) Counter {
	counter := make(Counter)

	for _, y := range Y {
		counter[fmt.Sprint(y)]++
	}

	return counter
}

func maxCounts(counter Counter) string {
	keys := make([]string, 0, len(counter))

	for key := range counter {
		keys = append(keys, key)
	}

	if len(keys) <= 1 {
		return keys[0]
	}

	sort.SliceStable(keys, func(i int, j int) bool {
		return counter[keys[i]] > counter[keys[j]]
	})

	return keys[0]
}

func countError(Y []int, err error) Counter {
	counter := make(Counter)

	for _, _counter := range Y {
		counter[fmt.Sprint(_counter)]++
	}

	return counter
}

func giniImpurity(s0 int, s1 int) float64 {
	if s0+s1 == 0 {
		return 0.0
	}

	prob0 := float64(s0) / float64(s0+s1)
	prob1 := float64(s1) / float64(s0+s1)

	return 1.0 - (math.Pow(prob0, 2) + math.Pow(prob1, 2))
}

func setFromList(list []string) (set []string) {
	track := make(map[string]bool)

	for _, item := range list {
		if _, ok := track[item]; !ok {
			track[item] = true
			set = append(set, item)
		}
	}

	return
}

func uniqueGotaSeries(s series.Series) series.Series {
	return series.New(setFromList(s.Records()), s.Type(), s.Name)
}

func meth(col []float64) []float64 {
	var methed []float64

	for i := 0; i < len(col)-1; i++ {
		methed = append(methed, (col[i]+col[i+1])/2)
	}

	return methed
}
