package buscar_libros

import (
	"database/sql"
)

// Libro representa un libro en el sistema.
type Libro struct {
	ID             int    `json:"id"`
	Titulo         string `json:"titulo"`
	Autores        string `json:"autores"`
	ISBN           string `json:"isbn"`
	AñoPublicacion int    `json:"anio_publicacion"`
}

// BuscarLibros busca libros en la base de datos que coincidan con el término de búsqueda.
// Devuelve una lista de libros que coinciden con el término.
func BuscarLibros(db *sql.DB, termino string) ([]Libro, error) {
	// Consulta para buscar libros y sus autores que coincidan con el término
	query := `
		SELECT Libros.id, Libros.titulo, GROUP_CONCAT(Autores.nombre SEPARATOR ', '), Libros.isbn, Libros.año_publicacion 
		FROM Libros 
		JOIN Libros_Autores ON Libros.id = Libros_Autores.libro_id 
		JOIN Autores ON Libros_Autores.autor_id = Autores.id 
		WHERE Libros.titulo LIKE ? OR Autores.nombre LIKE ?
		GROUP BY Libros.id, Libros.titulo, Libros.isbn, Libros.año_publicacion;
	`

	rows, err := db.Query(query, "%"+termino+"%", "%"+termino+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []Libro
	for rows.Next() {
		var libro Libro
		var autores string
		if err := rows.Scan(&libro.ID, &libro.Titulo, &autores, &libro.ISBN, &libro.AñoPublicacion); err != nil {
			return nil, err
		}
		libro.Autores = autores
		libros = append(libros, libro)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return libros, nil
}
