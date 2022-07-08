package http

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/guregu/null.v3"
)

type CreateBookRequestPayload struct {
	Id          null.Int      `json:"id,omitempty"`
	BookName    string        `json:"book_name,omitempty"`
	AuthorName  AuthorDetails `json:"author_name,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description null.String   `json:"description,omitempty"`
	Publisher   string        `json:"publisher,omitempty"`
	Price       int           `json:"price,omitempty"`
	EditionType null.Int      `json:"edition,omitempty"`
}

type AuthorDetails struct {
	AuthorFirstName string `json:"author_first_name,omitempty"`
	AuthorLastName  string `json:"author_last_name,omitempty"`
}

func (cbrp CreateBookRequestPayload) Validation() error {
	return validation.ValidateStruct(&cbrp,
		validation.Field(&cbrp.BookName, validation.Required),
		validation.Field(&cbrp.Publisher, validation.Required),
		validation.Field(&cbrp.AuthorName, validation.Required),
		validation.Field(&cbrp.Title, validation.Required),
		validation.Field(&cbrp.Price, validation.When(cbrp.Price > 0, validation.Required).Else(validation.Nil.Error("Price is required Fields"))),
		validation.Field(&cbrp.EditionType),
	)
}

func (ad AuthorDetails) Validate() error {
	return validation.ValidateStruct(&ad,
		validation.Field(&ad.AuthorFirstName, validation.Required),
		validation.Field(&ad.AuthorLastName, validation.Required),
	)
}
