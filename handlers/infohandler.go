package handlers

import (
	"fmt"
	"net/http"
	"github.com/gomarkdown/markdown"
	"io/ioutil"
	
	"groupXX/structures"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	//takes method of the request, if GET then forward to function, else write info to user
	switch r.Method {
	case http.MethodGet:
		InfoGetHandler(w, r)
	default:
		http.Error(w, "This service only offers the <a href=\"/diag\">/diag endpoint</a> that shows comprehensive HTTP "+
			"information for any received request. Please redirect any request to that endpoint.", http.StatusInternalServerError)
		return
	}
}

func InfoGetHandler(w http.ResponseWriter, r *http.Request) {
	//ensures that the data in the response is text of html type so the user can see it
	w.Header().Set("content-type", "text/html")

	//read the contents of your Markdown file
	mdFileContent, err := ioutil.ReadFile("./README.md")
	if err != nil {
		http.Error(w, "Error when reading markdown file", http.StatusInternalServerError)
		return
	}

	// Convert the Markdown content to HTML
	htmlContent := markdown.ToHTML(mdFileContent, nil, nil)

	//uses css to style the UI for the information pages more user friendly, with color and links as buttons
	output := `<html>
<head>
<style>
  body { 
    background-color: #4b5563;
	color: #ffffff;
	font-family: calibri;
	overflow-x: hidden;
  }
  a{
	color: #ffffff;
	background-color: #353740;
	text-decoration: none;
	padding: 2px 5px;
    border-radius: 5px;
  }
</style>
</head>
<body>`
  	output += "<a href=\"" + structures.DEFAULT_PATH + "\">" + "Go back to default handler (main menu)</a><br>"

	//makes a string out of the html content
	output += string(htmlContent)

	output += `</body></html>`

	//check for error when printing
	_, err = fmt.Fprintf(w, "%v", output)

	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}