package daemon

import (
	"github.com/gorilla/mux"
	"github.com/reef-pi/reef-pi/controller/settings"
	"github.com/reef-pi/reef-pi/controller/utils"
	"log"
	"net/http"
	"os"
)

var DefaultCredentials = utils.Credentials{
	User:     "reef-pi",
	Password: "reef-pi",
}

func (r *ReefPi) API() error {
	creds, err := r.a.GetCredentials()
	if err != nil {
		log.Println("ERROR: Failed to load credentials. Error", err)
		if err := r.store.Update(Bucket, "credentials", DefaultCredentials); err != nil {
			return err
		}
		creds = DefaultCredentials
	}
	err, router := startAPIServer(r.settings.Address, creds, r.settings.HTTPS)
	if err != nil {
		return err
	}
	r.AuthenticatedAPI(router)
	r.UnAuthenticatedAPI(router)
	if r.settings.Prometheus {
		r.prometheus()
	}
	if os.Getenv("REEF_PI_LIST_API") == "1" {
		utils.SummarizeAPI()
	}
	return nil
}

func (r *ReefPi) UnAuthenticatedAPI(router *mux.Router) {
	router.HandleFunc("/auth/signin", r.a.SignIn).Methods("POST")
	router.HandleFunc("/auth/signout", r.a.SignOut).Methods("GET")
}

// Authenticated API using the BasicAuth middleware
func (r *ReefPi) AuthenticatedAPI(router *mux.Router) {
	http.Handle("/api/", r.a.Authenticate(router.ServeHTTP))

	router.HandleFunc("/api/capabilities", r.GetCapabilities).Methods("GET")
	for _, sController := range r.subsystems {
		sController.LoadAPI(router)
	}
	utils.APIDoc(router.HandleFunc("/api/settings", r.GetSettings).Methods("GET"), nil, &settings.DefaultSettings)
	router.HandleFunc("/api/settings", r.UpdateSettings).Methods("POST")
	router.HandleFunc("/api/credentials", r.a.UpdateCredentials).Methods("POST")
	router.HandleFunc("/api/telemetry", r.telemetry.GetConfig).Methods("GET")
	router.HandleFunc("/api/telemetry", r.telemetry.UpdateConfig).Methods("POST")
	router.HandleFunc("/api/telemetry/test_message", r.telemetry.SendTestMessage).Methods("POST")
	router.HandleFunc("/api/errors/clear", r.clearErrors).Methods("DELETE")
	router.HandleFunc("/api/errors/{id}", r.deleteError).Methods("DELETE")
	router.HandleFunc("/api/errors/{id}", r.getError).Methods("GET")
	router.HandleFunc("/api/errors", r.listErrors).Methods("GET")
	router.HandleFunc("/api/me", r.a.Me).Methods("GET")
	if r.h != nil {
		router.HandleFunc("/api/health_stats", r.h.GetStats).Methods("GET")
	}
	r.dm.LoadAPI(router)
	for _, sController := range r.subsystems {
		sController.LoadAPI(router)
	}
	if r.settings.Capabilities.Dashboard {
		router.HandleFunc("/api/dashboard", r.GetDashboard).Methods("GET")
		router.HandleFunc("/api/dashboard", r.UpdateDashboard).Methods("POST")
	}
}

func startAPIServer(address string, creds utils.Credentials, https bool) (error, *mux.Router) {
	assets := http.FileServer(http.Dir("ui/assets"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ui/home.html")
	})
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ui/favicon.ico")
	})
	router := mux.NewRouter()
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))
	images := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", images))
	http.Handle("/auth/", router)
	if https {
		if err := utils.GenerateCerts(); err != nil {
			return err, nil
		}
		go func() {
			log.Printf("Starting https server at: %s\n", address)
			if err := http.ListenAndServeTLS(address, "server.crt", "server.key", nil); err != nil {
				log.Println("ERROR: Failed to run https server. Error:", err)
			}
		}()
	} else {
		go func() {
			log.Printf("Starting http server at: %s\n", address)
			if err := http.ListenAndServe(address, nil); err != nil {
				log.Println("ERROR: Failed to run http server. Error:", err)
			}
		}()
	}
	return nil, router
}
