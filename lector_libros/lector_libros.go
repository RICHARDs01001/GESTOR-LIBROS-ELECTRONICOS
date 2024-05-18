package lector_libros

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
)

// ObtenerRutaPdf obtiene la ruta del archivo PDF de un libro dado su ID.
// Consulta la base de datos para obtener la ruta.
func ObtenerRutaPdf(db *sql.DB, libroID int) (string, error) {
	var rutaPdf string
	err := db.QueryRow("SELECT ruta_pdf FROM Libros WHERE id = ?", libroID).Scan(&rutaPdf)
	if err != nil {
		return "", fmt.Errorf("error al obtener ruta del PDF: %v", err)
	}
	return rutaPdf, nil
}

// LeerLibro abre un archivo PDF utilizando el navegador Microsoft Edge.
// Verifica si el archivo existe antes de intentar abrirlo.
func LeerLibro(rutaPdf string) error {
	// Verificar si el archivo existe
	if _, err := os.Stat(rutaPdf); os.IsNotExist(err) {
		return fmt.Errorf("el archivo %s no existe", rutaPdf)
	}

	// Comando para abrir el archivo PDF en Microsoft Edge
	cmd := exec.Command(`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`, rutaPdf)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo %s: %v", rutaPdf, err)
	}

	fmt.Printf("Abriendo el libro %s...\n", rutaPdf)
	return nil
}
