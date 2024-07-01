package lector_libros

import (
	"database/sql"
)

type Libro struct {
	ID        int    `json:"id"`
	Contenido string `json:"contenido"`
}

func LeerLibro(db *sql.DB, libroID int) (*Libro, error) {
	var libro Libro
	err := db.QueryRow("SELECT id, contenido FROM Libros WHERE id = ?", libroID).Scan(&libro.ID, &libro.Contenido)
	if err != nil {
		return nil, err
	}
	return &libro, nil
}
