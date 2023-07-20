package requests

import (
	"encoding/json"
	"github.com/Swapica/swapica-svc/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewGetCommissionEstimateRequest(r *http.Request) (resources.GetCommissionEstimateRequest, error) {
	request := struct {
		Data resources.GetCommissionEstimateRequest
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data, errors.Wrap(err, "failed to decode request")
	}

	return request.Data, nil
}
