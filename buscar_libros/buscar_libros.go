package buscar_libros

import (
	"database/sql"
)

type Libro struct {
	ID         int    `json:"id"`
	Titulo     string `json:"titulo"`
	Autor      string `json:"autor"`
	CoverImage string `json:"coverImage"`
	PdfPath    string `json:"pdfPath"`
}

func BuscarLibros(db *sql.DB, termino string) ([]Libro, error) {
	rows, err := db.Query("SELECT id, titulo, autor, cover_image, pdf_path FROM Libros WHERE titulo LIKE ?", "%"+termino+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var libro Libro
		if err := rows.Scan(&libro.ID, &libro.Titulo, &libro.Autor, &libro.CoverImage, &libro.PdfPath); err != nil {
			return nil, err
		}
		libros = append(libros, libro)
	}

	return libros, nil
}
