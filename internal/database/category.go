// Package database
package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *Category {
	return &Category{db: db}
}

func (c *Category) Create(name string, description string) (Category, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("INSERT INTO category (id, name, description) VALUES ($1, $2, $3)",
		id, name, description)
	if err != nil {
		return Category{}, err
	}

	return Category{ID: id, Name: name, Description: description}, nil
}

func (c *Category) FindAll() ([]Category, error) {
	rows, err := c.db.Query("SELECT * FROM category")
	if err != nil {
		return nil, err
	}

	var categories []Category

	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (c *Category) FindByID(ID string) (Category, error) {
	var category Category

	row := c.db.QueryRow("SELECT * FROM category WHERE id = ?", ID)
	if err := row.Scan(&category.ID, &category.Name, &category.Description); err != nil {
		return Category{}, err
	}

	return category, nil
}

func (c *Category) FindByCourseID(ID string) (Category, error) {
	var category Category

	row := c.db.QueryRow("SELECT c.id, c.name, c.description FROM category c JOIN course co ON co.category_id = c.id WHERE co.id = $1", ID)
	if err := row.Scan(&category.ID, &category.Name, &category.Description); err != nil {
		return Category{}, err
	}

	return category, nil
}
