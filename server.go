package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var tpl *template.Template
var LocalHost = "8080"

func main () {
	exec.Command("open", "http://localhost:"+LocalHost).Start()
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("css"))))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// we parse all the files in the static directory
	// if the error is not a nil, the program panics (template.Must)
	tpl = template.Must(template.ParseGlob("static/*.html"))

    // HandleFunc takes a string (specifying a server's resource path), and a function that will handle requests for that path
    // everytime we receive a request for a URL ending in "/" --> call the homePage to generate a response
    http.HandleFunc("/", homePage)
    http.HandleFunc("/ascii-art", GetData)

    // Start the web server and specify the port to listen for incoming requests
    // We pass the string "localhost:8080" to the server, which will cause it to accept requests only from your own machine on port 8080 
    // Port 8080 is usually used for web servers (that’s the port that web browsers make HTTP requests to by default)
    // ListenAndServe listens for browser requests and responds to them.
    err := http.ListenAndServe(":"+LocalHost, nil) // The nil value in the 2nd arg means that requests will be handled using functions set up via homePage.
    // ListenAndServe will run forever, unless it encounters an error. If it does, it will return that error, which we log before the program exits
    log.Fatal(err)
}

// A net/http handler func handles browser requests at a certain path
// http.ResponseWriter writes data to the browser response, w is a value for updating the response sent to the browser
// A pointer to an http.Request value represents a request from the browser
func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
	} else {
		tpl.ExecuteTemplate(w,"index.html", nil)
	}
}

func GetData(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
    	userInput := r.FormValue("inputtext")
		// for _, v := range userInput {
		// 	if !(v >= 32 && v <= 126) {
		// 		http.Error(w, "ERROR-400\nBad request!", http.StatusBadRequest)
		// 		return
		// 	}
		// }
		
		userFont := r.FormValue("Banner")
		var checkbanner int
		if userFont == "standard.txt" || userFont == "thinkertoy.txt" || userFont == "shadow.txt" || userFont == "" {
			checkbanner = 1
		} else {
			checkbanner = 0
		}
		if checkbanner == 0 {
			http.Error(w, "500 Internal Server Error",http.StatusInternalServerError)
			return
		}
	
		if userInput == "" || userFont == "" {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		result := PrintAscii(userInput, userFont)
		
		tpl.ExecuteTemplate(w, "result.html", result)

	default:
		fmt.Fprintf(w, "Only POST methods supported")
	}
}

func PrintAscii(text, font string) string {
	var file = []byte{}

	fonts := os.ReadFile(font)
	file, err := fonts
	if err != nil {
		return ""
	}
	// if font == "standard" {
	// 	file, _ = os.ReadFile("standard.txt")
	// } else if font == "shadow" {
	// 	file, _ = os.ReadFile("shadow.txt")
	// } else {
	// 	file, _ = os.ReadFile("thinkertoy.txt")
	// }

	temp := strings.ReplaceAll(string(file), "\r", "")
	fontFile := strings.Split(temp, "\n")
	

	tempInput := strings.ReplaceAll(text, "\r\n", "\\n")
	input := strings.Split(tempInput, "\\n")

	mainSlice := [][]rune{}

	for _, v := range input {
		runeHolder := []rune(v) // [104,102,33,123]
		mainSlice = append(mainSlice, runeHolder) // [  [104,102,33,123]   ]
	}

	str := ""

	for _, slice := range mainSlice{
		if len(slice) == 0 {
			str += "\n"
		} else {
			for i := 0; i < 8; i++ {
				for _, r := range slice{
					str += fontFile[int((r-32)*9+1)+i]
				}
				str += "\n"
			}
		}
	
	}

	return str
}
