package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBayesSingleNode(t *testing.T) {
	dag := &DAG{}

	A := "A"
	negations := map[string]string{
		A: "!A",
	}

	dag.AddNode(A)

	questionnaire := GenerateQuestionnaire(dag, nil, negations, nil)
	require.Equal(t, 1, len(questionnaire.ProbConds))

	questionnaire.ProbConds[0].Answer = 50
	require.Equal(t, 50, CalculateProbability(A, questionnaire))

	questionnaire.ProbConds[0].Answer = 33
	require.Equal(t, 33, CalculateProbability(A, questionnaire))
}

func TestBayesSingleNodeWithExperiment(t *testing.T) {
	dag := &DAG{}

	A := "A"
	negations := map[string]string{
		A: "!A",
	}

	dag.AddNode(A)

	experiments := map[string][]string{
		"A": []string{
			"E1",
			"E2",
		},
	}
	questionnaire := GenerateQuestionnaire(dag, nil, negations, experiments)
	require.Equal(t, 1, len(questionnaire.ProbConds))
	require.Equal(t, 2, len(questionnaire.Experiments))

	questionnaire.ProbConds[0].Answer = 33
	questionnaire.Experiments[0].AnswerIfTrue = 50
	questionnaire.Experiments[0].AnswerIfFalse = 50
	questionnaire.Experiments[1].AnswerIfTrue = 50
	questionnaire.Experiments[1].AnswerIfFalse = 50
	require.Equal(t, 33, CalculateProbability(A, questionnaire))

	questionnaire.Experiments[0].AnswerIfTrue = 70
	questionnaire.Experiments[0].AnswerIfFalse = 50
	require.Equal(t, 41, CalculateProbability(A, questionnaire))

	questionnaire.Experiments[1].AnswerIfTrue = 20
	questionnaire.Experiments[1].AnswerIfFalse = 30
	require.Equal(t, 31, CalculateProbability(A, questionnaire))
}

func TestBayesSingleEdge(t *testing.T) {
	dag := &DAG{}

	A := "A"
	B := "B"
	negations := map[string]string{
		A: "!A",
		B: "!B",
	}

	dag.AddEdge(A, B)

	questionnaire := GenerateQuestionnaire(dag, nil, negations, nil)

	require.Equal(t, `P(A (~!A) |  | ) = 0
P(B (~!B) | A | ) = 0
P(B (~!B) |  | !A) = 0`, questionnaire.Format())

	require.Equal(t, 3, len(questionnaire.ProbConds))

	questionnaire.ProbConds[0].Answer = 13
	questionnaire.ProbConds[1].Answer = 76
	questionnaire.ProbConds[2].Answer = 41

	require.Equal(t, questionnaire.ProbConds[0].Answer, CalculateProbability(A, questionnaire))

	require.Equal(t, 46, CalculateProbability(B, questionnaire))
}

func TestBayesSingleEdgeWGiven(t *testing.T) {
	dag := &DAG{}

	A := "A"
	B := "B"
	negations := map[string]string{
		A: "!A",
		B: "!B",
	}

	dag.AddEdge(A, B)

	questionnaire := GenerateQuestionnaire(dag, []string{A}, negations, nil)
	require.Equal(t, `P(B (~!B) | A | ) = 0`, questionnaire.Format())

	require.Equal(t, 100, CalculateProbability(A, questionnaire))

	questionnaire.ProbConds[0].Answer = 76
	require.Equal(t, questionnaire.ProbConds[0].Answer, CalculateProbability(B, questionnaire))
}

func TestBayesSimplishWGiven(t *testing.T) {
	dag := &DAG{}

	A := "A"
	B := "B"
	C := "C"
	D := "D"
	E := "E"
	negations := map[string]string{
		A: "!A",
		B: "!B",
		C: "!C",
		D: "!D",
		E: "!E",
	}

	experiments := map[string][]string{
		E: []string{"EX1", "EX2"},
		C: []string{"EX3"},
	}

	dag.AddEdge(A, C)
	dag.AddEdge(A, D)
	dag.AddEdge(B, C)
	dag.AddEdge(B, E)
	dag.AddEdge(C, D)
	dag.AddEdge(D, E)

	questionnaire := GenerateQuestionnaire(dag, []string{A, B}, negations, experiments)
	require.Equal(t, `P(C (~!C) | A, B | ) = 0
P(D (~!D) | A, C | ) = 0
P(D (~!D) | A | !C) = 0
P(E (~!E) | B, D | ) = 0
P(E (~!E) | B | !D) = 0
P(EX3 | C) = 0
P(EX3 | ~!C) = 0
P(EX1 | E) = 0
P(EX1 | ~!E) = 0
P(EX2 | E) = 0
P(EX2 | ~!E) = 0`, questionnaire.Format())

	questionnaire.ProbConds[0].Answer = 76
	questionnaire.ProbConds[1].Answer = 13
	questionnaire.ProbConds[2].Answer = 34
	questionnaire.ProbConds[3].Answer = 31
	questionnaire.ProbConds[4].Answer = 54

	questionnaire.Experiments[0].AnswerIfTrue = 50
	questionnaire.Experiments[0].AnswerIfFalse = 50
	questionnaire.Experiments[1].AnswerIfTrue = 50
	questionnaire.Experiments[1].AnswerIfFalse = 50
	questionnaire.Experiments[2].AnswerIfTrue = 50
	questionnaire.Experiments[2].AnswerIfFalse = 50

	require.Equal(t, 100, CalculateProbability(A, questionnaire))
	require.Equal(t, 100, CalculateProbability(B, questionnaire))
	require.Equal(t, questionnaire.ProbConds[0].Answer, CalculateProbability(C, questionnaire))
	require.Equal(t, 18, CalculateProbability(D, questionnaire))
	require.Equal(t, 50, CalculateProbability(E, questionnaire))

	questionnaire.Experiments[0].AnswerIfTrue = 66
	questionnaire.Experiments[0].AnswerIfFalse = 46

	require.Equal(t, 100, CalculateProbability(A, questionnaire))
	require.Equal(t, 100, CalculateProbability(B, questionnaire))
	require.Equal(t, 82, CalculateProbability(C, questionnaire))
	require.Equal(t, 19, CalculateProbability(D, questionnaire))
	require.Equal(t, 50, CalculateProbability(E, questionnaire))

	questionnaire.Experiments[1].AnswerIfTrue = 70
	questionnaire.Experiments[1].AnswerIfFalse = 50
	questionnaire.Experiments[2].AnswerIfTrue = 30
	questionnaire.Experiments[2].AnswerIfFalse = 20

	require.Equal(t, 100, CalculateProbability(A, questionnaire))
	require.Equal(t, 100, CalculateProbability(B, questionnaire))
	require.Equal(t, 82, CalculateProbability(C, questionnaire))
	require.Equal(t, 19, CalculateProbability(D, questionnaire))
	require.Equal(t, 67, CalculateProbability(E, questionnaire))
}

func TestBayesNaive(t *testing.T) {
	dag := &DAG{}

	W := "w"
	K := "k"
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

	questionnaire := GenerateQuestionnaire(dag, []string{W}, negations, nil)
	for i := range questionnaire.ProbConds {
		questionnaire.ProbConds[i].Answer = 50
	}

	if v := CalculateProbability(S, questionnaire); v != 50 {
		t.Error("probability", v, "expected", 50)
	}

	for i := range questionnaire.ProbConds {
		questionnaire.ProbConds[i].Answer = 33
	}

	if v := CalculateProbability(S, questionnaire); v != 33 {
		t.Error("probability", v, "expected", 33)
	}

	require.Equal(t, `P(c (~cf) | w | ) = 33
P(k (~wk) | w, c | ) = 33
P(k (~wk) | w | cf) = 33
P(f (~ff) | w, c, k | ) = 33
P(f (~ff) | w, c | wk) = 33
P(f (~ff) | w, k | cf) = 33
P(f (~ff) | w | cf, wk) = 33
P(s (~sf) | c, f | ) = 33
P(s (~sf) | c | ff) = 33
P(s (~sf) | f | cf) = 33
P(s (~sf) |  | cf, ff) = 33`, questionnaire.Format())

	answers := []int{45, 90, 65, 70, 75, 10, 12, 60, 35, 8, 1}
	for i := range questionnaire.ProbConds {
		questionnaire.ProbConds[i].Answer = answers[i]
	}

	require.Equal(t, 22, CalculateProbability(S, questionnaire))
}
