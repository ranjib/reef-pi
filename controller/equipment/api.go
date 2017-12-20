package equipment

import (
	"github.com/gorilla/mux"
	"github.com/reef-pi/reef-pi/controller/utils"
	"net/http"
)

//API
func (e *Controller) LoadAPI(r *mux.Router) {
	r.HandleFunc("/api/equipment/{id}", e.GetEquipment).Methods("GET")
	r.HandleFunc("/api/equipment", e.ListEquipments).Methods("GET")
	r.HandleFunc("/api/equipment", e.CreateEquipment).Methods("PUT")
	r.HandleFunc("/api/equipment/{id}", e.UpdateEquipment).Methods("POST")
	r.HandleFunc("/api/equipment/{id}", e.DeleteEquipment).Methods("DELETE")
}

func (c *Controller) GetEquipment(w http.ResponseWriter, r *http.Request) {
	fn := func(id string) (interface{}, error) {
		return c.Get(id)
	}
	utils.JSONGetResponse(fn, w, r)
}

func (c Controller) ListEquipments(w http.ResponseWriter, r *http.Request) {
	fn := func() (interface{}, error) {
		return c.List()
	}
	utils.JSONListResponse(fn, w, r)
}

func (c *Controller) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	var eq Equipment
	fn := func() error {
		return c.Create(eq)
	}
	utils.JSONCreateResponse(&eq, fn, w, r)
}

func (c *Controller) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	var eq Equipment
	fn := func(id string) error {
		return c.Update(id, eq)
	}
	utils.JSONUpdateResponse(&eq, fn, w, r)
}

func (c *Controller) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	fn := func(id string) error {
		return c.Delete(id)
	}
	utils.JSONDeleteResponse(fn, w, r)
}
