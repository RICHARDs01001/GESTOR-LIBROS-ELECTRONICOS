package favoritos_libros

import (
	"database/sql"
)

type Favorito struct {
	ID        int `json:"id"`
	UsuarioID int `json:"usuario_id"`
	LibroID   int `json:"libro_id"`
}

func ListarFavoritos(db *sql.DB, usuarioID string) ([]Favorito, error) {
	rows, err := db.Query("SELECT id, usuario_id, libro_id FROM Favoritos WHERE usuario_id = ?", usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favoritos []Favorito
	for rows.Next() {
		var favorito Favorito
		if err := rows.Scan(&favorito.ID, &favorito.UsuarioID, &favorito.LibroID); err != nil {
			return nil, err
		}
		favoritos = append(favoritos, favorito)
	}

	return favoritos, nil
}

func AgregarLibroAFavoritos(db *sql.DB, usuarioID, libroID int) error {
	_, err := db.Exec("INSERT INTO Favoritos (usuario_id, libro_id) VALUES (?, ?)", usuarioID, libroID)
	if err != nil {
		return err
	}
	return nil
}
