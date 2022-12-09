package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/swapica-svc/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func NewCreateMatchRequest(r *http.Request) (resources.CreateMatchRequest, error) {
	request := struct {
		Data resources.CreateMatchRequest
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data, errors.Wrap(err, "failed to decode request")
	}

	return request.Data, nil
}
