package buscar_libros

import (
	"database/sql"
)

// Libro representa un libro en el sistema.
type Libro struct {
	ID             int
	Titulo         string
	Autores        string
	ISBN           string
	AñoPublicacion int
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
		GROUP BY Libros.id
	`
	rows, err := db.Query(query, "%"+termino+"%", "%"+termino+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libros []Libro
	// Escanear los resultados y agregarlos a la lista de libros
	for rows.Next() {
		var libro Libro
		err := rows.Scan(&libro.ID, &libro.Titulo, &libro.Autores, &libro.ISBN, &libro.AñoPublicacion)
		if err != nil {
			return nil, err
		}
		libros = append(libros, libro)
	}

	return libros, nil
}
