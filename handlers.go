package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

func inc(i int) int {
	return i + 1
}

func InitializeApp() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/scorecard", ScorecardHandler)
	http.HandleFunc("/scorecards", ScorecardsHandler)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func ScorecardsHandler(w http.ResponseWriter, r *http.Request) {
	scorecards, err := getScorecards()
	if err != nil {
		fmt.Println(err) // TODO
	}

	// TODO - Clean up this, very awkwarad for some reason
	// Read the template file
	templateFile, err := os.ReadFile("templates/scorecards.html")
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return
	}

	// Convert the file content to a string
	templateContent := string(templateFile)

	// Parse the HTML template
	tmpl := template.New("scorecards").Funcs(template.FuncMap{"inc": inc})
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	fmt.Println("Parsed Template:", tmpl.DefinedTemplates())

	err = tmpl.Execute(w, scorecards)
	if err != nil {
		fmt.Println(err) // TODO: Handle the error appropriately
	}
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
