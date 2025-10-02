package req

import (
	"net/http"
	pkg "purple/links/pkg/res"
)

func HandleBody[T any](w http.ResponseWriter, req *http.Request) (*T, error) {
	body, err := Decode[T](req.Body)
	if err != nil {
		pkg.Json(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	err = IsValid(body)
	if err != nil {
		pkg.Json(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &body, nil
}

// OV20pt8w3Q5K8D9BrXyw
