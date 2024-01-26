package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func InitializeApp() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/scorecard", ScorecardHandler)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func ScorecardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		HandleScorecardPost(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/scorecard.html"))
	tmpl.Execute(w, nil)
}

func HandleScorecardPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}
	var scorecard Scorecard
	for i := 1; i <= 18; i++ {
		inputName := "hole" + fmt.Sprint(i)
		scoreStr := r.Form.Get(inputName)
		score, err := strconv.Atoi(scoreStr)
		if err != nil {
			http.Error(w, "Invalid score for "+inputName, http.StatusBadRequest)
			return
		}
		scorecard.Holes[i-1] = score
	}

	insertScore(scorecard)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
