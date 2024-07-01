package resenas_libros

import (
	"database/sql"
)

type Reseña struct {
	ID           int    `json:"id"`
	LibroID      int    `json:"libro_id"`
	UsuarioID    int    `json:"usuario_id"`
	Calificacion int    `json:"calificacion"`
	Texto        string `json:"texto"`
}

func ListarReseñas(db *sql.DB, libroID string) ([]Reseña, error) {
	rows, err := db.Query("SELECT id, libro_id, usuario_id, calificacion, texto FROM Reseñas WHERE libro_id = ?", libroID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reseñas []Reseña
	for rows.Next() {
		var reseña Reseña
		if err := rows.Scan(&reseña.ID, &reseña.LibroID, &reseña.UsuarioID, &reseña.Calificacion, &reseña.Texto); err != nil {
			return nil, err
		}
		reseñas = append(reseñas, reseña)
	}

	return reseñas, nil
}

func InsertarReseña(db *sql.DB, usuarioID, libroID, calificacion int, texto string) error {
	_, err := db.Exec("INSERT INTO Reseñas (libro_id, usuario_id, calificacion, texto) VALUES (?, ?, ?, ?)", libroID, usuarioID, calificacion, texto)
	if err != nil {
		return err
	}
	return nil
}
