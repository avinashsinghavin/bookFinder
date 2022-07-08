package order

import (
	requrstPayload "RestAPI/order/http"
	"RestAPI/repository"
	"RestAPI/response"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type BookOrderHttpExtension struct {
	Router    *mux.Router
	DB        *gorm.DB
	BookOrder repository.BookOrderDB
}

func (bohe *BookOrderHttpExtension) RegisterRoutes() {
	bohe.Router.HandleFunc("/", bohe.Greeting()).Methods("GET", "POST")
	bohe.Router.HandleFunc("/addBooks", bohe.AddBooks()).Methods("POST", "PUT")
	bohe.Router.HandleFunc("/searchBook", bohe.SearchBooks()).Methods("POST")
}

func (bohe BookOrderHttpExtension) AddBooks() func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(rw, http.StatusText(500), http.StatusInternalServerError)
			}
		}()
		switch req.Method {
		case http.MethodPost:
			bohe.createBook(rw, req)
		case http.MethodPut:
			bohe.updateBookDetails(rw, req)
		}
	}
}

func (boht BookOrderHttpExtension) createBook(rw http.ResponseWriter, req *http.Request) {
	log.Println(req)
	var httpRequestPayLoad requrstPayload.CreateBookRequestPayload
	data := json.NewDecoder(req.Body)
	if err := data.Decode(&httpRequestPayLoad); err != nil {
		fmt.Println(err)
		http.Error(rw, "Error while decoding request payload", http.StatusBadRequest)
	}

	if err := httpRequestPayLoad.Validation(); err != nil {
		fmt.Println(httpRequestPayLoad.Validation())
		http.Error(rw, fmt.Sprintf("Validation Exception #{httpRequestPayLoad.Validation()}"), http.StatusBadRequest)
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			log.Println("panic recovered in HandleSupplyOrderStateUpdate")
			http.Error(rw, http.StatusText(500), http.StatusInternalServerError)
		}
	}()

	rw.Header().Set("Content-Type", "application/json")
	var globalResponse response.GlobalResponse
	boht.BookOrder.DB = boht.DB
	getBookDetails := boht.BookOrder.GetBookDataIfExist(rw, httpRequestPayLoad)
	if getBookDetails.Id != 0 {
		json.NewEncoder(rw).Encode(globalResponse.GetGlobalResponse(409, "Data Exist", fmt.Sprint("data already exist with id ", getBookDetails.Id)))
		return
	}
	savedData := boht.BookOrder.CreateBookRepo(httpRequestPayLoad)
	globalResponse = globalResponse.GetGlobalResponse(201, "Created Successfully", savedData)
	json.NewEncoder(rw).Encode(globalResponse)
}

func (boht BookOrderHttpExtension) updateBookDetails(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	var httpRequestPayLoad requrstPayload.CreateBookRequestPayload
	boht.BookOrder.DB = boht.DB
	if headerContentType := req.Header.Get("Content-Type"); headerContentType != "application/json" {
		http.Error(rw, "Content Type is not application/json", http.StatusBadRequest)
		return
	}

	data := json.NewDecoder(req.Body)
	if err := data.Decode(&httpRequestPayLoad); err != nil {
		http.Error(rw, "Error while decoding request payload", http.StatusBadRequest)
		return
	}
	if httpRequestPayLoad.Id.ValueOrZero() == 0 {
		json.NewEncoder(rw).Encode(response.GlobalResponse{Status: 422, Message: "Id is required for updating the book"})
		return
	}

	getBook := boht.BookOrder.FindBookById(rw, httpRequestPayLoad)

	if getBook.Id == 0 {
		return
	}
	getBook = boht.BookOrder.UpdateBookDetails(rw, httpRequestPayLoad)
	json.NewEncoder(rw).Encode(response.GlobalResponse{Status: 204, Message: "Data updated successfully", Data: getBook})
}

func (bohe BookOrderHttpExtension) SearchBooks() func(rw http.ResponseWriter, re *http.Request) {
	return func(rw http.ResponseWriter, re *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		request := requrstPayload.SearchBookRequestPayload{}
		data := json.NewDecoder(re.Body)
		if err := data.Decode(&request); err != nil {
			http.Error(rw, "Error while decoding request payload", http.StatusBadRequest)
			return
		}
		defer func() {
			if err := recover(); err != nil {
				json.NewEncoder(rw).Encode(response.GlobalResponse{
					Status: 500, Message: "Internal Server Error", Data: err,
				})
			}
		}()
		bohe.BookOrder.DB = bohe.DB
		listOfBooks := bohe.BookOrder.FindBookByFilter(rw, request)
		if len(listOfBooks) == 0 {
			json.NewEncoder(rw).Encode(response.GlobalResponse{
				Status:  200,
				Message: "Success!",
				Data:    "Data match not found!",
			})
		}
		json.NewEncoder(rw).Encode(response.GlobalResponse{
			Status:  200,
			Message: "Success!",
			Data:    listOfBooks,
		})
	}
}

func (boht BookOrderHttpExtension) Greeting() func(res http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode("Welcome to Go Lang! server hosted on port 8000! This is Book Details Finder.")
	}
}
