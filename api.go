package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)


type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer{
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

// type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// // Implementing the http.Handler interface for HandlerFunc
// func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	err := f(w, r)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }



func (s *APIServer) Run () {
	router := mux.NewRouter()

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET","POST","DELETE"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	router.Use(cors)

	router.HandleFunc("/member", s.HandleMember).Methods("GET","POST","DELETE")
	router.HandleFunc("/member/id/{id}", s.HandleGetMemberByID).Methods("GET","POST","DELETE")
	router.HandleFunc("/member/tech/{tech}", s.HandleGetMemberByTech).Methods("GET")
	router.HandleFunc("/member/delete/{id}", s.HandleDeleteMember).Methods("DELETE")
	//router.HandleFunc("/member/update/{id}", s.HandleUpdateMember).Methods("POST")

	
	log.Println("Json API server running on port",s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal(err)
	}

}

func (s *APIServer) HandleMember (w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		s.HandleGetMembers(w,r)
	case http.MethodPost:
		s.HandleCreateMember(w,r)
	case http.MethodDelete:
		s.HandleDeleteMember(w,r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)			
	}
	
}



func (s *APIServer) HandleGetMembers(w http.ResponseWriter, r *http.Request) error {
	members, err := s.store.GetMembers()
	if err != nil {
			fmt.Println("HandleAPi", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(members); err != nil {
		// Log the JSON encoding error for debugging purposes
		fmt.Println("Error encoding JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
}

	return nil
}





func (s *APIServer) HandleGetMemberByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	
	// Parse the ID into a UUID type
	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	member, err := s.store.GetMemberByID(parsedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(member); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}






func (s *APIServer) HandleGetMemberByTech(w http.ResponseWriter, r *http.Request) {
	tech := mux.Vars(r)["tech"]
	
	fmt.Println("bhencho abhi aya nhi Hadler mein")
	

	results, err := s.store.GetMembersByTech(tech)
	if err != nil {
			
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("err bich mein hain bhai")
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("ending mein error hain")
			return
	}

}



func (s *APIServer) HandleDeleteMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	parsedID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.store.DeleteMember(parsedID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}



// func (s *APIServer) HandleUpdateMember (w http.ResponseWriter, r *http.Request) {

// 	createMemberReq := new(CreateMemberRequest)



// }




func (s *APIServer) HandleCreateMember (w http.ResponseWriter, r *http.Request) {

	createMemberReq := new(CreateMemberRequest)

	if err := json.NewDecoder(r.Body).Decode(createMemberReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//log.Printf("Decoded request: %+v\n", createMemberReq)

	newMember := NewMember(
		createMemberReq.FirstName,
		createMemberReq.LastName,
		createMemberReq.Tech,
		createMemberReq.About,
		createMemberReq.Discord,
		createMemberReq.Linkedin,
	)

	if err := s.store.CreateMember(newMember); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//log.Printf("New member details: %+v\n", newMember)

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newMember); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}




