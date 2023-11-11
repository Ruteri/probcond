package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"graphy/lib"
)

type DAGData struct {
	Edges []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"edges"`
	Nodes []struct {
		Value    string `json:"value"`
		Negation string `json:"negation"`
	} `json:"nodes"`
	Experiments [][2]string `json:"experiments"`
	Given       []string    `json:"given"`
}

type QuestionnaireData struct {
	Nodes []struct {
		Value string
	}
	Questionnaire lib.Questionnaire
}

type Result struct {
	InQuestion string
	Result     int
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", http.FileServer(http.Dir("static")).ServeHTTP).Methods("GET")
	r.HandleFunc("/dag.js", http.FileServer(http.Dir("static")).ServeHTTP).Methods("GET")
	r.HandleFunc("/graph.js", http.FileServer(http.Dir("static")).ServeHTTP).Methods("GET")
	r.HandleFunc("/questionnaire.js", http.FileServer(http.Dir("static")).ServeHTTP).Methods("GET")
	r.HandleFunc("/dag", handleDAG).Methods("POST")
	r.HandleFunc("/questionnaire", handleQuestionnaire).Methods("POST")

	fmt.Println("Server is running on http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func handleDAG(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var dagData DAGData
	err = json.Unmarshal(bodyBytes, &dagData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dag := &lib.DAG{}

	negations := make(map[string]string)
	for _, node := range dagData.Nodes {
		if node.Negation != "" {
			negations[node.Value] = node.Negation
		}
		dag.AddNode(node.Value)
	}

	for _, edge := range dagData.Edges {
		dag.AddEdge(edge.Src, edge.Dst)
	}

	experimentsMap := make(map[string][]string)
	for _, ex := range dagData.Experiments {
		experimentsMap[ex[0]] = append(experimentsMap[ex[0]], ex[1])
	}

	questionnaire := lib.GenerateQuestionnaire(dag, dagData.Given, negations, experimentsMap)
	log.Print(questionnaire)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questionnaire)
}

func handleQuestionnaire(w http.ResponseWriter, r *http.Request) {
	var questionnaireData QuestionnaireData
	err := json.NewDecoder(r.Body).Decode(&questionnaireData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var results []Result
	cache := make(map[string]float64)
	for _, node := range questionnaireData.Nodes {
		prob := lib.CalculateNodeProb(node.Value, &questionnaireData.Questionnaire, &cache, nil)
		result := Result{
			InQuestion: node.Value,
			Result:     int(prob*100 + 0.5),
		}
		results = append(results, result)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
