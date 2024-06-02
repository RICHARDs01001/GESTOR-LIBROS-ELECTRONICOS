package resenas_libros

import (
	"database/sql"
	"fmt"
	"time"
)

// Reseña representa una reseña de un libro.
type Reseña struct {
	ID           int
	UsuarioID    int
	LibroID      int
	Calificacion int
	Texto        string
	Fecha        time.Time
}

// InsertarReseña inserta una nueva reseña en la base de datos.
func InsertarReseña(db *sql.DB, usuarioID, libroID, calificacion int, texto string) error {
	fecha := time.Now().Format("2006-01-02")
	query := "INSERT INTO Reseñas (usuario_id, libro_id, calificacion, texto, fecha) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, usuarioID, libroID, calificacion, texto, fecha)
	if err != nil {
		return fmt.Errorf("error al insertar la reseña: %v", err)
	}
	return nil
}

// ListarReseñas lista todas las reseñas de un libro.
func ListarReseñas(db *sql.DB, libroID int) ([]Reseña, error) {
	query := "SELECT id, usuario_id, libro_id, calificacion, texto, fecha FROM Reseñas WHERE libro_id = ?"
	rows, err := db.Query(query, libroID)
	if err != nil {
		return nil, fmt.Errorf("error al listar las reseñas: %v", err)
	}
	defer rows.Close()

	var reseñas []Reseña
	for rows.Next() {
		var reseña Reseña
		var fecha string
		if err := rows.Scan(&reseña.ID, &reseña.UsuarioID, &reseña.LibroID, &reseña.Calificacion, &reseña.Texto, &fecha); err != nil {
			return nil, fmt.Errorf("error al escanear la fila de reseña: %v", err)
		}
		reseña.Fecha, _ = time.Parse("2006-01-02", fecha)
		reseñas = append(reseñas, reseña)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre las filas de reseña: %v", err)
	}
	return reseñas, nil
}
