package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/analytics/internal/ports"
)

type SumReactionResponse struct {
	Sum int `json:"sum" example:"3"`
}

type SumReactionRequest struct {
	ObjectID uint32 `json:"object_id" example:"123"`
}

type SumReaction struct {
	Statistics ports.Statistics
}

// @Summary Get summary time of Reaction
// @Description Get summary time of Reaction
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /analytics/sum-reaction [get]
func (i SumReaction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := logrus.WithField("RAPI", "/analytics/sum-reaction")

	var err error
	var req SumReactionRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("error parse body request: ", err)
		http.Error(w, "{}", http.StatusBadRequest)
		return
	}

	var response SumReactionResponse
	response.Sum, err = i.Statistics.GetSumReaction(req.ObjectID)
	if err != nil {
		log.Error("error Statistics.GetSumReaction: ", err)
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
