package lector_libros

import (
	"database/sql"
	"fmt"
)

// Libro representa un libro completo que puede ser leído por el usuario.
type Libro struct {
	ID             int    `json:"id"`
	Titulo         string `json:"titulo"`
	RutaPDF        string `json:"ruta_pdf"`
	AñoPublicacion int    `json:"año_publicacion"`
}

// LeerLibro permite al usuario obtener la URL del PDF del libro.
func LeerLibro(db *sql.DB, libroID int) (*Libro, error) {
	var libro Libro

	// Consultar la base de datos para obtener la ruta del PDF del libro
	err := db.QueryRow("SELECT id, titulo, ruta_pdf, año_publicacion FROM Libros WHERE id = ?", libroID).Scan(&libro.ID, &libro.Titulo, &libro.RutaPDF, &libro.AñoPublicacion)
	if err != nil {
		return nil, fmt.Errorf("error al obtener el libro: %v", err)
	}

	return &libro, nil
}
