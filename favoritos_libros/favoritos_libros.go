package favoritos_libros

import (
	"database/sql"
	"fmt"
)

// Favorito representa un libro favorito del usuario.
type Favorito struct {
	UsuarioID      int
	LibroID        int
	Titulo         string
	Autor          string
	ISBN           string
	AñoPublicacion int
}

// AgregarLibroAFavoritos agrega un libro a la lista de favoritos del usuario.
// Verifica si el libro ya está en la lista de favoritos antes de agregarlo.
func AgregarLibroAFavoritos(db *sql.DB, usuarioID, libroID int) error {
	// Verificar si el libro ya está en la lista de favoritos
	existe, err := ExisteLibroEnFavoritos(db, usuarioID, libroID)
	if err != nil {
		return fmt.Errorf("error al verificar si el libro está en favoritos: %v", err)
	}

	if existe {
		return fmt.Errorf("el libro ya está en la lista de favoritos")
	}

	// Insertar el libro en la lista de favoritos
	_, err = db.Exec("INSERT INTO Favoritos (usuario_id, libro_id) VALUES (?, ?)", usuarioID, libroID)
	if err != nil {
		return fmt.Errorf("error al agregar libro a favoritos: %v", err)
	}
	fmt.Println("Libro agregado a favoritos.")
	return nil
}

// ExisteLibroEnFavoritos verifica si un libro ya está en la lista de favoritos del usuario.
// Devuelve verdadero si el libro ya está en la lista, de lo contrario, falso.
func ExisteLibroEnFavoritos(db *sql.DB, usuarioID, libroID int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Favoritos WHERE usuario_id = ? AND libro_id = ?", usuarioID, libroID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error al verificar si el libro está en favoritos: %v", err)
	}
	return count > 0, nil
}

// ListarFavoritos devuelve la lista de libros favoritos del usuario.
// Consulta la base de datos y devuelve una lista de libros favoritos.
func ListarFavoritos(db *sql.DB, usuarioID int) ([]Favorito, error) {
	var favoritos []Favorito

	// Consulta para obtener la lista de libros favoritos del usuario
	rows, err := db.Query(`
		SELECT Favoritos.usuario_id, Favoritos.libro_id, Libros.titulo, Autores.nombre, Libros.isbn, Libros.año_publicacion 
		FROM Favoritos 
		JOIN Libros ON Favoritos.libro_id = Libros.id 
		JOIN Libros_Autores ON Libros.id = Libros_Autores.libro_id 
		JOIN Autores ON Libros_Autores.autor_id = Autores.id 
		WHERE Favoritos.usuario_id = ?`, usuarioID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener lista de favoritos: %v", err)
	}
	defer rows.Close()

	// Escanear los resultados y agregarlos a la lista de favoritos
	for rows.Next() {
		var favorito Favorito
		err := rows.Scan(&favorito.UsuarioID, &favorito.LibroID, &favorito.Titulo, &favorito.Autor, &favorito.ISBN, &favorito.AñoPublicacion)
		if err != nil {
			return nil, fmt.Errorf("error al escanear favorito: %v", err)
		}
		favoritos = append(favoritos, favorito)
	}

	return favoritos, nil
}
