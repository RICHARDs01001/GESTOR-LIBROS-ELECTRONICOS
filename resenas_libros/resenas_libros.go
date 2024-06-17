package resenas_libros

import (
	"database/sql"
	"time"
)

// Reseña representa una reseña de un libro realizada por un usuario.
type Reseña struct {
	ID           int       `json:"id"`
	UsuarioID    int       `json:"usuario_id"`
	LibroID      int       `json:"libro_id"`
	Calificacion int       `json:"calificacion"`
	Texto        string    `json:"texto"`
	Fecha        time.Time `json:"fecha"`
}

// ListarReseñas devuelve una lista de reseñas de un libro.
func ListarReseñas(db *sql.DB, libroID string) ([]Reseña, error) {
	rows, err := db.Query("SELECT id, usuario_id, libro_id, calificacion, texto, fecha FROM Reseñas WHERE libro_id = ?", libroID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reseñas []Reseña
	for rows.Next() {
		var reseña Reseña
		if err := rows.Scan(&reseña.ID, &reseña.UsuarioID, &reseña.LibroID, &reseña.Calificacion, &reseña.Texto, &reseña.Fecha); err != nil {
			return nil, err
		}
		reseñas = append(reseñas, reseña)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reseñas, nil
}

// InsertarReseña inserta una nueva reseña en la base de datos.
func InsertarReseña(db *sql.DB, usuarioID, libroID, calificacion int, texto string) error {
	_, err := db.Exec("INSERT INTO Reseñas (usuario_id, libro_id, calificacion, texto, fecha) VALUES (?, ?, ?, ?, ?)", usuarioID, libroID, calificacion, texto, time.Now())
	return err
}
