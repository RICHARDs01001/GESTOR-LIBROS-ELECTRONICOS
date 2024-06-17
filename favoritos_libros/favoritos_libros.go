package favoritos_libros

import (
	"database/sql"
)

// Favorito representa un libro favorito de un usuario.
type Favorito struct {
	UsuarioID int `json:"usuario_id"`
	LibroID   int `json:"libro_id"`
}

// ListarFavoritos devuelve una lista de libros favoritos de un usuario.
func ListarFavoritos(db *sql.DB, usuarioID string) ([]int, error) {
	rows, err := db.Query("SELECT libro_id FROM Favoritos WHERE usuario_id = ?", usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favoritos []int
	for rows.Next() {
		var libroID int
		if err := rows.Scan(&libroID); err != nil {
			return nil, err
		}
		favoritos = append(favoritos, libroID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favoritos, nil
}

// AgregarLibroAFavoritos agrega un libro a los favoritos de un usuario.
func AgregarLibroAFavoritos(db *sql.DB, usuarioID, libroID int) error {
	_, err := db.Exec("INSERT INTO Favoritos (usuario_id, libro_id) VALUES (?, ?)", usuarioID, libroID)
	return err
}
