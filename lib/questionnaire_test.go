package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuestionnaireSingleNode(t *testing.T) {
	dag := &DAG{}

	A := "A"
	negations := map[string]string{
		A: "!A",
	}

	dag.AddNode(A)

	questionnaire := GenerateQuestionnaire(dag, nil, negations)

	require.Equal(t, `P(A (~!A) |  | ) = 0`, questionnaire.Format())
}

func TestQuestionnaireSingleEdge(t *testing.T) {
	dag := &DAG{}

	A := "A"
	B := "B"
	negations := map[string]string{
		A: "!A",
		B: "!B",
	}

	dag.AddEdge(A, B)

	questionnaire := GenerateQuestionnaire(dag, nil, negations)

	require.Equal(t, `P(A (~!A) |  | ) = 0
P(B (~!B) | A | ) = 0
P(B (~!B) |  | !A) = 0`, questionnaire.Format())
}

func TestQuestionnaireSingleEdgeWGiven(t *testing.T) {
	dag := &DAG{}

	A := "A"
	B := "B"
	negations := map[string]string{
		A: "!A",
		B: "!B",
	}

	dag.AddEdge(A, B)

	questionnaire := GenerateQuestionnaire(dag, []string{A}, negations)

	require.Equal(t, `P(B (~!B) | A | ) = 0`, questionnaire.Format())
}

func TestQuestionnaire(t *testing.T) {
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

	questionnaire := GenerateQuestionnaire(dag, []string{W}, negations)
	require.Equal(t, `P(c (~cf) | w | ) = 0
P(k (~wk) | w, c | ) = 0
P(k (~wk) | w | cf) = 0
P(f (~ff) | w, c, k | ) = 0
P(f (~ff) | w, c | wk) = 0
P(f (~ff) | w, k | cf) = 0
P(f (~ff) | w | cf, wk) = 0
P(s (~sf) | c, f | ) = 0
P(s (~sf) | c | ff) = 0
P(s (~sf) | f | cf) = 0
P(s (~sf) |  | cf, ff) = 0`, questionnaire.Format())
}
