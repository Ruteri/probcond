package main

import (
	"fmt"
	"graphy/lib"
)

func main() {
	dag := &lib.DAG{}

	W := "w"
	K := "c"
	F := "f"

	C := "c"
	S := "s"

	negations := map[string]string{
		W: "nw",
		K: "wk",
		F: "ff",
		C: "cf",
		S: "sf",
	}

	// Add edges to the graph
	dag.AddEdge(W, C)
	dag.AddEdge(W, K)
	dag.AddEdge(W, F)
	dag.AddEdge(C, K)
	dag.AddEdge(C, F)
	dag.AddEdge(C, S)
	dag.AddEdge(K, F)
	dag.AddEdge(F, S)

	// Print the graph
	fmt.Print(dag.Format())
	questionnaire := lib.GenerateQuestionnaire(dag, []string{W}, negations)
	questionnaire.GatherAnswers()
	/*answers := []int{45, 90, 65, 70, 75, 10, 12, 60, 35, 8, 1}
	for i := range questionnaire.ProbConds {
		questionnaire.ProbConds[i].Answer = answers[i]
	}*/

	fmt.Print(questionnaire.Format())

	v := lib.CalculateProbability(S, questionnaire)
	fmt.Println("Probability: ", v)
}
