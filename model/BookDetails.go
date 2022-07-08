package model

import "gopkg.in/guregu/null.v3"

type Book struct {
	Id              int      `gorm:"column:id" gorm:"primaryKey" json:"id"`
	Title           string   `gorm:"column:title" json:"title"`
	BookName        string   `gorm:"column:book_name" json:"book_name"`
	Description     string   `gorm:"column:description" json:"description"`
	Publisher       string   `gorm:"column:publisher" json:"publisher"`
	Price           int      `gorm:"column:price" json:"price"`
	EditionType     null.Int `gorm:"column:edition" json:"edition_type"`
	AuthorFirstName string   `gorm:"column:author_first_name" json:"author_first_name"`
	AuthorLastName  string   `gorm:"column:author_last_name" json:"author_last_name"`
}
