package autenticacion

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Usuario representa un usuario en el sistema.
type Usuario struct {
	ID                int
	Nombre            string
	CorreoElectronico string
}

// IniciarSesion permite a un usuario iniciar sesión en el sistema.
// Comprueba las credenciales proporcionadas y devuelve el usuario si son correctas.
func IniciarSesion(db *sql.DB, correoElectronico, contraseña string) (*Usuario, error) {
	var usuario Usuario
	var hashedPassword string

	// Consultar la base de datos para obtener el usuario y la contraseña hash
	err := db.QueryRow("SELECT id, nombre, contraseña FROM Usuarios WHERE correo_electronico = ?", correoElectronico).Scan(&usuario.ID, &usuario.Nombre, &hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario: %v", err)
	}

	// Comparar la contraseña hash con la proporcionada
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(contraseña))
	if err != nil {
		return nil, fmt.Errorf("contraseña incorrecta: %v", err)
	}

	// Si la contraseña es correcta, devolver el usuario
	return &usuario, nil
}

// Registrarse permite a un nuevo usuario registrarse en el sistema.
// Cifra la contraseña y guarda los detalles del usuario en la base de datos.
func Registrarse(db *sql.DB, nombre, correoElectronico, contraseña string) error {
	// Generar un hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(contraseña), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al encriptar la contraseña: %v", err)
	}

	// Insertar el nuevo usuario en la base de datos
	_, err = db.Exec("INSERT INTO Usuarios (nombre, correo_electronico, contraseña) VALUES (?, ?, ?)", nombre, correoElectronico, hashedPassword)
	if err != nil {
		return fmt.Errorf("error al registrar usuario: %v", err)
	}

	return nil
}
