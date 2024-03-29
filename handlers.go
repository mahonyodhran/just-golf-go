package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type AddScoreTemplateData struct {
	Golfers []Golfer
	Courses []Course
}

func inc(i int) int {
	return i + 1
}

func InitializeApp() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/scorecard", ScorecardHandler)
	http.HandleFunc("/scorecards", ScorecardsHandler)
	http.HandleFunc("/course", CourseHandler)
	http.HandleFunc("/golfer", GolferHandler)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func ScorecardsHandler(w http.ResponseWriter, r *http.Request) {
	scorecards, err := getScorecards()
	if err != nil {
		log.Fatal("Error getting scorecards: ", err)
	}
	// TODO - Clean up this, very awkwarad for some reason
	// Read the template file
	templateFile, err := os.ReadFile("templates/scorecard/scorecards.html")
	if err != nil {
		log.Fatal("Error reading template file: ", err)
		return
	}

	// Convert the file content to a string
	templateContent := string(templateFile)

	// Parse the HTML templates
	tmpl := template.New("scorecards").Funcs(template.FuncMap{"inc": inc})
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		log.Fatal("Error parsing template: ", err)
		return
	}

	err = tmpl.Execute(w, scorecards)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}
}

func ScorecardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		HandleScorecardPost(w, r)
		return
	}

	golfers, err := getGolfers()
	if err != nil {
		log.Fatal("Error getting golfers: ", err)
	}

	courses, err := getCourses()
	if err != nil {
		log.Fatal("Error getting courses: ", err)
	}

	data := AddScoreTemplateData{
		Golfers: golfers,
		Courses: courses,
	}

	tmpl := template.Must(template.ParseFiles("templates/scorecard/add-scorecard.html"))
	tmpl.Execute(w, data)
}

func HandleScorecardPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}
	// TODO - Break this up into helper validation methods
	var scorecard Scorecard
	var golferIDStr = r.Form.Get("golferID")
	scorecard.GolferID, err = strconv.Atoi(golferIDStr)
	if err != nil {
		http.Error(w, "Invalid golfer ID", http.StatusBadRequest)
		return
	}

	courseIDStr := r.Form.Get("courseID")
	scorecard.CourseID, err = strconv.Atoi(courseIDStr)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

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

func CourseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		HandleCoursePost(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/course/add-course.html"))
	tmpl.Execute(w, nil)
}

func HandleCoursePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}
	var courseName = r.Form.Get("courseName")
	if err != nil {
		http.Error(w, "Invalid Course Name", http.StatusBadRequest)
		return
	}

	insertCourse(courseName)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GolferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		HandleGolferPost(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/golfer/add-golfer.html"))
	tmpl.Execute(w, nil)
}

func HandleGolferPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		return
	}
	var first_name = r.Form.Get("firstName")
	if err != nil {
		http.Error(w, "Invalid First Name", http.StatusBadRequest)
		return
	}
	var last_name = r.Form.Get("lastName")
	if err != nil {
		http.Error(w, "Invalid Last Name", http.StatusBadRequest)
		return
	}
	var index_str = r.Form.Get("index")
	index, err := strconv.Atoi(index_str)
	if err != nil {
		http.Error(w, "Invalid Index", http.StatusBadRequest)
		return
	}

	insertGolfer(first_name, last_name, index)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
