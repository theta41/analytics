package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/analytics/internal/ports"
)

type AcceptedResponse struct {
	Count int `json:"count" example:"2"`
}

type CountAcceptedTask struct {
	Statistics ports.Statistics
}

// @Summary Get count accepted task
// @Description Get count accepted task
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /analytics/count-accepted [get]
func (i CountAcceptedTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := logrus.WithField("RAPI", "/analytics/count-accepted")

	var response AcceptedResponse
	var err error
	response.Count, err = i.Statistics.GetCountAcceptedTask()
	if err != nil {
		log.Error("error Statistics.GetCountAcceptedTask: ", err)
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(response)
	if err != nil {
		log.Error("error marshal response: ", err)
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		log.Error("error write response: ", err)
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}
}
