package http

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/guregu/null.v3"
)

type SearchBookRequestPayload struct {
	BookName        null.String `json:"book_name"`
	AuthorFirstName null.String `json:"author_first_name"`
	AuthorLastName  null.String `json:"author_last_name"`
	Title           null.String `json:"title"`
	Description     null.String `json:"description"`
	Publisher       null.String `json:"publisher"`
	StartingPrice   null.Int    `json:"starting_price"`
	EndingPrice     null.Int    `json:"ending_price"`
	PageSize        int         `json:"page_size"`
	PageNumber      int         `json:"page_number"`
}

func (sbrp SearchBookRequestPayload) Validate() error {
	return validation.ValidateStruct(&sbrp,
		validation.Field(sbrp.BookName),
		validation.Field(sbrp.AuthorFirstName),
		validation.Field(sbrp.AuthorLastName),
		validation.Field(sbrp.Title),
		validation.Field(sbrp.Description),
		validation.Field(sbrp.Publisher),
		validation.Field(sbrp.StartingPrice),
		validation.Field(sbrp.EndingPrice),
		validation.Field(sbrp.PageSize, validation.Required),
		validation.Field(sbrp.PageNumber, validation.Required))
}
