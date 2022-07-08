package repository

import (
	"RestAPI/model"
	request "RestAPI/order/http"
	"RestAPI/response"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type BookOrderDB struct {
	DB *gorm.DB
}

func (db BookOrderDB) CreateBookRepo(req request.CreateBookRequestPayload) model.Book {
	bookDetails := model.Book{}
	bookDetails.BookName = req.BookName
	bookDetails.Title = req.Title
	bookDetails.AuthorFirstName = fmt.Sprintf(req.AuthorName.AuthorFirstName)
	bookDetails.AuthorLastName = fmt.Sprintf(req.AuthorName.AuthorLastName)
	bookDetails.Description = req.Description.ValueOrZero()
	bookDetails.EditionType = req.EditionType
	bookDetails.Price = req.Price
	bookDetails.Publisher = req.Publisher

	if res := db.DB.Create(&bookDetails); res.Error != nil {
		fmt.Println(res.Error)
	}
	return bookDetails
}

func (db BookOrderDB) GetBookDataIfExist(rw http.ResponseWriter, req request.CreateBookRequestPayload) model.Book {
	getExistingBook := model.Book{}
	bookDetails := model.Book{
		BookName:        req.BookName,
		Title:           req.Title,
		AuthorFirstName: fmt.Sprintf(req.AuthorName.AuthorFirstName),
		AuthorLastName:  fmt.Sprintf(req.AuthorName.AuthorLastName),
		Description:     req.Description.ValueOrZero(),
		EditionType:     req.EditionType,
		Price:           req.Price,
		Publisher:       req.Publisher,
	}

	defer func() {
		if rec := recover(); rec != nil {
			log.Println("panic recovered in Book Repo method GetBookDataIfExist():", rec)
			http.Error(rw, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}()

	if res := db.DB.Where(fmt.Sprint("book_name = '", bookDetails.BookName, "' AND author_first_name = '",
		bookDetails.AuthorFirstName, "' AND author_last_name = '", bookDetails.AuthorLastName, "' AND title = '", bookDetails.Title, "'")).Find(&getExistingBook); res.Error != nil {
		log.Fatalln(res.Error)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(&response.GlobalResponse{Status: 500, Message: "Internal Server Error", Data: res.Error})
	}
	return getExistingBook
}

func (db BookOrderDB) FindBookById(rw http.ResponseWriter, req request.CreateBookRequestPayload) model.Book {
	getExistingBook := model.Book{}
	defer func() {
		if rec := recover(); rec != nil {
			http.Error(rw, http.StatusText(500), http.StatusInternalServerError)
			log.Println("error while findBookById(): ", rec)
			return
		}
	}()
	if res := db.DB.Where(fmt.Sprint("id != ", int(req.Id.ValueOrZero()), " AND book_name = '", req.BookName, "' AND author_first_name = '",
		req.AuthorName.AuthorFirstName, "' AND author_last_name = '", req.AuthorName.AuthorLastName, "' AND title = '", req.Title, "'")).Find(&getExistingBook); res.RowsAffected > 0 {
		log.Println("Data exist with above parameter: ", res.Error)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(&response.GlobalResponse{Status: 500, Message: "Data Already Exist with different id: ", Data: request.CreateBookRequestPayload{
			BookName: getExistingBook.BookName,
			Title:    getExistingBook.Title,
			AuthorName: request.AuthorDetails{
				AuthorFirstName: getExistingBook.BookName,
				AuthorLastName:  getExistingBook.BookName,
			},
		}})
		return model.Book{}
	}
	log.Println(getExistingBook)
	getExistingBook.Id = int(req.Id.ValueOrZero())
	if res := db.DB.Find(&getExistingBook); res.Error != nil {
		log.Println("DB Error: ", res.Error)
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(&response.GlobalResponse{Status: 500, Message: "Internal server error", Data: res.Error})
		return model.Book{}
	}
	return getExistingBook
}

func (db BookOrderDB) FindBookByFilter(rw http.ResponseWriter, request request.SearchBookRequestPayload) []model.Book {
	resp := make([]model.Book, 0)
	offSetData := 0
	if request.PageNumber > 1 {
		offSetData = (request.PageNumber - 1) * request.PageSize
	}
	query := ""
	if request.BookName.NullString.Valid {
		query += fmt.Sprint(" book_name LIKE '%", request.BookName.ValueOrZero(), "%' AND")
	}
	if request.AuthorFirstName.NullString.Valid {
		query += fmt.Sprint(" author_first_name LIKE '%", request.AuthorFirstName.ValueOrZero(), "%' AND")
	}
	if request.AuthorLastName.NullString.Valid {
		query += fmt.Sprint(" author_last_name LIKE '%", request.AuthorLastName.ValueOrZero(), "%' AND")
	}
	if request.Title.NullString.Valid {
		query += fmt.Sprint(" title LIKE '%", request.Title.ValueOrZero(), "%' AND")
	}
	if request.Description.NullString.Valid {
		query += fmt.Sprint(" description LIKE '%", request.Description.ValueOrZero(), "%' AND")
	}
	if request.Publisher.NullString.Valid {
		query += fmt.Sprint(" publisher LIKE '%", request.Publisher.ValueOrZero(), "%' AND")
	}
	if request.StartingPrice.Valid && request.StartingPrice.ValueOrZero() > 0 && request.EndingPrice.Valid && request.EndingPrice.ValueOrZero() > 0 {
		query += fmt.Sprintf(" price >= %d && price <= %d AND", request.StartingPrice.ValueOrZero(), request.EndingPrice.ValueOrZero())
	}
	if len(query) > 0 {
		query = query[:len(query)-3]
	}
	if res := db.DB.Offset(offSetData).Limit(request.PageSize).Where(query).Find(&resp); res.Error != nil {
		log.Println("Error in DB Query")
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(&response.GlobalResponse{
			Status:  500,
			Message: "Internal Server Error!",
			Data:    res.Error,
		})
	}
	return resp
}

func (db BookOrderDB) UpdateBookDetails(rw http.ResponseWriter, req request.CreateBookRequestPayload) model.Book {
	book := model.Book{
		Id:              int(req.Id.ValueOrZero()),
		BookName:        req.BookName,
		Title:           req.Title,
		AuthorFirstName: fmt.Sprintf(req.AuthorName.AuthorFirstName),
		AuthorLastName:  fmt.Sprintf(req.AuthorName.AuthorLastName),
		Description:     req.Description.ValueOrZero(),
		EditionType:     req.EditionType,
		Price:           req.Price,
		Publisher:       req.Publisher,
	}
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("error while updating Book Details in repo: ", rec)
			http.Error(rw, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}()
	if err := db.DB.Updates(&book); err.Error != nil {
		return model.Book{}
	}
	return book
}
