package product

import (
	"errors"
	"net/http"
	"purple/links/configs"
	"purple/links/pkg/middleware"
	"purple/links/pkg/req"
	pkg "purple/links/pkg/res"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	*ProductRepository
	*configs.Config
}

type ProductHandler struct {
	*ProductRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
	}
	router.Handle("POST /products", middleware.AuthMiddleware(handler.Create(), deps.Config))
	router.HandleFunc("GET /products/{id}", handler.Read())
	router.Handle("PATCH /products/{id}", middleware.AuthMiddleware(handler.Update(), deps.Config))
	router.Handle("DELETE /products/{id}", middleware.AuthMiddleware(handler.Delete(), deps.Config))
}

func (h *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ProductCreateRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := body.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		product := NewProduct(body)
		newProduct, err := h.ProductRepository.Create(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pkg.Json(w, newProduct, http.StatusCreated)
	}
}

func (h *ProductHandler) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product, err := h.ProductRepository.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}
		pkg.Json(w, product, http.StatusOK)
	}
}

func (h *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body, err := req.HandleBody[ProductUpdateRequest](w, r)
		if err != nil {
			return
		}

		if err := body.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		product, err := h.ProductRepository.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body.ApplyUpdates(product)

		product, err = h.ProductRepository.Update(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pkg.Json(w, product, http.StatusOK)
	}
}

func (h *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.ProductRepository.Delete(uint64(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		pkg.Json(w, "Deleted", http.StatusNoContent)
	}
}
