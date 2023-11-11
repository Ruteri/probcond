package lib

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ProbCond struct {
	InQuestion string
	QNegation  string
	Conditions []string
	Negations  []string
	Answer     int
}

func (probCond *ProbCond) Format() string {
	return fmt.Sprintf("P(%s (~%s) | %s | %s) = %d", probCond.InQuestion, probCond.QNegation, strings.Join(probCond.Conditions, ", "), strings.Join(probCond.Negations, ", "), probCond.Answer)
}

type Experiment struct {
	InQuestion         string
	Hypothesis         string
	HypothesisNegation string
	AnswerIfTrue       int // p (e|q)
	AnswerIfFalse      int // p (e|~q)
}

func (e *Experiment) Format() string {
	return fmt.Sprintf("P(%s | %s) = %d\nP(%s | ~%s) = %d", e.InQuestion, e.Hypothesis, e.AnswerIfTrue, e.InQuestion, e.HypothesisNegation, e.AnswerIfFalse)
}

type Questionnaire struct {
	ProbConds   []*ProbCond
	Experiments []*Experiment
}

func (q *Questionnaire) Format() string {
	return strings.Join(
		append(
			Transform(q.ProbConds, func(pc *ProbCond) string { return pc.Format() }),
			Transform(q.Experiments, (*Experiment).Format)...,
		),
		"\n",
	)
}

func (q *Questionnaire) Print() {
	fmt.Printf(q.Format())
}

func GenerateQuestionnaire(dag *DAG, given []string, negations map[string]string, experiments map[string][]string) *Questionnaire {
	q := &Questionnaire{}

	dag.Traverse(func(node *Node) bool {
		nodeParents := dag.NodeParents(node.Value)
		q.ProbConds = append(q.ProbConds, MakeProbConds(dag, node, []*ProbCond{}, nodeParents, given, negations)...)
		if relevantExperiments, found := experiments[node.Value]; found {
			q.Experiments = append(q.Experiments, Transform(relevantExperiments, func(experiment string) *Experiment {
				hypothesisNegation, found := negations[node.Value]
				if !found {
					hypothesisNegation = "not " + node.Value
				}
				return &Experiment{
					InQuestion:         experiment,
					Hypothesis:         node.Value,
					HypothesisNegation: hypothesisNegation,
				}
			})...)
		}
		return true
	})

	return q
}

func MakeProbConds(dag *DAG, node *Node, curConds []*ProbCond, unprocessedParents []*Node, given []string, negations map[string]string) []*ProbCond {
	isGiven := len(Filter(given, func(s *string) bool { return *s == node.Value })) > 0
	if isGiven {
		return curConds
	}

	qNegation, found := negations[node.Value]
	if !found {
		qNegation = "not " + node.Value
	}

	if dag.IsRoot(node) {
		return append(curConds, &ProbCond{
			InQuestion: node.Value,
			QNegation:  qNegation,
			Conditions: []string{},
			Negations:  []string{},
		})
	}

	if len(unprocessedParents) == 0 {
		return curConds
	}

	parentGiven := len(Filter(given, func(s *string) bool { return *s == unprocessedParents[0].Value })) > 0

	conds := []*ProbCond{}

	pNegation, found := negations[unprocessedParents[0].Value]
	if !found {
		pNegation = "not " + unprocessedParents[0].Value
	}

	for _, cond := range curConds {
		conds = append(conds, &ProbCond{
			InQuestion: node.Value,
			QNegation:  qNegation,
			Conditions: append(cond.Conditions, unprocessedParents[0].Value),
			Negations:  cond.Negations,
		})

		if !parentGiven {
			conds = append(conds, &ProbCond{
				InQuestion: node.Value,
				QNegation:  qNegation,
				Conditions: cond.Conditions,
				Negations:  append(cond.Negations, pNegation),
			})
		}
	}

	if len(curConds) == 0 {
		conds = append(conds, &ProbCond{
			InQuestion: node.Value,
			QNegation:  qNegation,
			Conditions: []string{unprocessedParents[0].Value},
			Negations:  nil,
		})

		if !parentGiven {
			conds = append(conds, &ProbCond{
				InQuestion: node.Value,
				QNegation:  qNegation,
				Conditions: nil,
				Negations:  []string{pNegation},
			})
		}
	}

	return MakeProbConds(dag, node, conds, unprocessedParents[1:], given, negations)
}

func (q *Questionnaire) GatherAnswers() {
	for _, pc := range q.ProbConds {
		v, err := GetAnswer(pc)
		if err != nil {
			panic(err.Error())
		}
		pc.Answer = v
	}

	for _, ex := range q.Experiments {
		v, err := GetAnswer(&ProbCond{
			InQuestion: ex.InQuestion,
			Conditions: []string{ex.Hypothesis},
		})
		if err != nil {
			panic(err.Error())
		}
		ex.AnswerIfTrue = v

		v, err = GetAnswer(&ProbCond{
			InQuestion: ex.InQuestion,
			Negations:  []string{ex.HypothesisNegation},
		})
		if err != nil {
			panic(err.Error())
		}
		ex.AnswerIfFalse = v
	}
}

func (pc *ProbCond) FormatAsQuestion() []string {
	condStrings := []string{fmt.Sprintf("What is the probability that %s given", pc.InQuestion)}
	condStrings = append(condStrings, pc.Conditions...)
	condStrings = append(condStrings, pc.Negations...)

	return condStrings
}

func GetAnswer(pc *ProbCond) (int, error) {
	fmt.Printf("%s? (0-100): ", strings.Join(pc.FormatAsQuestion(), "\n  - and "))

	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
		return 0, err
	}

	// remove the delimeter from the string
	input = strings.TrimSuffix(input, "\n")
	return strconv.Atoi(input)
}
