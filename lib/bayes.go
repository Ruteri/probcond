package lib

import "fmt"

func CalculateProbability(node string, q *Questionnaire) int {
	cache := make(map[string]float64)
	return int(CalculateNodeProb(node, q, &cache, nil)*100 + 0.5)
}

func CalculateNodeProb(node string, q *Questionnaire, cache *map[string]float64, visited *map[string]struct{}) float64 {
	if cv, found := (*cache)[node]; found {
		return cv
	}

	if visited == nil {
		visitedMap := make(map[string]struct{})
		visited = &visitedMap
	}

	if _, found := (*visited)[node]; found {
		fmt.Println("error: cyclic graph! aborting")
		panic("abort! cyclic graph")
	}

	(*visited)[node] = struct{}{}

	ret := 0.0

	var nodeQs []*ProbCond = Filter(q.ProbConds, func(pc **ProbCond) bool { return (*pc).InQuestion == node || (*pc).QNegation == node })
	for _, nq := range nodeQs {
		cc := float64(nq.Answer) / 100
		for _, cq := range nq.Conditions {
			cc = cc * CalculateNodeProb(cq, q, cache, visited)
		}
		for _, cq := range nq.Negations {
			cc = cc * (1 - CalculateNodeProb(cq, q, cache, visited))
		}
		fmt.Println("calculated", *nq, cc)
		ret += cc
	}

	if len(nodeQs) == 0 {
		ret = 1.0
	}

	fmt.Println("calculated node", node, ret)

	// Take experiment data into account (bayes rule)
	ret = Accumulate(
		ret,
		Filter(q.Experiments, func(ex **Experiment) bool { return (*ex).Hypothesis == node }),
		func(cR float64, ex **Experiment) float64 {
			ansIfTrue := float64((*ex).AnswerIfTrue) / 100
			ansIfFalse := float64((*ex).AnswerIfFalse) / 100
			return (ansIfTrue * cR) / ((ansIfTrue * cR) + (ansIfFalse * (1 - cR)))
		},
	)

	(*cache)[node] = ret
	return ret
}

func Accumulate[Acc any, T any](a Acc, ts []T, cb func(Acc, *T) Acc) Acc {
	for _, t := range ts {
		a = cb(a, &t)
	}
	return a
}

func Filter[T any](ts []T, cb func(*T) bool) []T {
	rs := []T{}
	for _, t := range ts {
		if cb(&t) {
			rs = append(rs, t)
		}
	}
	return rs
}
