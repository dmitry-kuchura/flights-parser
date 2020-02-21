package api

import (
	"net/http"

	u "skyup/utils"
)

func GetFlights(w http.ResponseWriter, r *http.Request) {
	collection := u.GetConnection()

	data, _ := u.FindMany(collection)

	resp := u.Message(true, "success")
	resp["data"] = data
	u.Respond(w, resp)
}
