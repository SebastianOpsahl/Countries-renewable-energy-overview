package handlers

import (
	"fmt"
	"net/http"

	"groupXX/structures"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	//ensures that the data in the response is text of html type so the user can see it
	w.Header().Set("content-type", "text/html")

	//offer information for redirection to other paths, because that is the reason they got the default handler
	output := `<html>
<head>
<style>
  body {
	background-color: #4b5563;
	color: #ffffff;
	font-family: calibri;
  }
  a{
	color: #ffffff;
	background-color: #353740;
	text-decoration: none;
	padding: 5px 10px;
    border-radius: 4px;
  }
</style>
</head>
<body>
<h2>This is an webservice to give information about the percentage of renewables used in countries around the world</h2>
You are now at root level, choose what you want to do:<br><br>`

  	//redirections
	output += "<a href=\"" + structures.RENEWABLECURRENT_PATH + "\">" + "Search for current renewables data in given country/countries</a><br><br>"

	output += "<a href=\"" + structures.RENEWABLEHISTORY_PATH + "\">" + "Search for history renewables data in given country/countries</a><br><br>"

	output += "<a href=\"" + structures.NOTIFICATIONS_PATH + "\">" + "Notification endpoint</a><br><br>"

	output += "<a href=\"" + structures.STATUS_PATH + "\">" + "Status page for the API's</a><br><br><br>" + "For information about our service and how to use it: <br><br>"

	output += "<a href=\"" + structures.INFO_PATH + "\">" + "Information page</a><br><br>"

	output += `</body></html>`

	//check for error when printing out
	_, err := fmt.Fprintf(w, "%v", output)

	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}