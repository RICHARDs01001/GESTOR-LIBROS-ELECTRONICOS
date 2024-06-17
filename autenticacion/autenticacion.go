package autenticacion

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Usuario representa un usuario en el sistema.
type Usuario struct {
	ID                int    `json:"id"`
	Nombre            string `json:"nombre"`
	CorreoElectronico string `json:"correo_electronico"`
}

type Credentials struct {
	Nombre            string `json:"nombre,omitempty"`
	CorreoElectronico string `json:"correo_electronico"`
	Contraseña        string `json:"contraseña"`
}

type UsuarioRepository interface {
	FindUserByEmail(email string) (*Usuario, string, error)
	SaveUser(user *Usuario, hashedPassword string) error
}

type UsuarioRepositorySQL struct {
	DB *sql.DB
}

func (repo *UsuarioRepositorySQL) FindUserByEmail(email string) (*Usuario, string, error) {
	var usuario Usuario
	var hashedPassword string
	err := repo.DB.QueryRow("SELECT id, nombre, contraseña FROM Usuarios WHERE correo_electronico = ?", email).Scan(&usuario.ID, &usuario.Nombre, &hashedPassword)
	if err != nil {
		return nil, "", fmt.Errorf("error al obtener usuario: %v", err)
	}
	return &usuario, hashedPassword, nil
}

func (repo *UsuarioRepositorySQL) SaveUser(user *Usuario, hashedPassword string) error {
	_, err := repo.DB.Exec("INSERT INTO Usuarios (nombre, correo_electronico, contraseña) VALUES (?, ?, ?)", user.Nombre, user.CorreoElectronico, hashedPassword)
	if err != nil {
		return fmt.Errorf("error al registrar usuario: %v", err)
	}
	return nil
}

// IniciarSesion permite a un usuario iniciar sesión en el sistema.
// Comprueba las credenciales proporcionadas y devuelve el usuario si son correctas.
func IniciarSesion(repo UsuarioRepository, correoElectronico, contraseña string) (*Usuario, error) {
	usuario, hashedPassword, err := repo.FindUserByEmail(correoElectronico)
	if err != nil {
		return nil, err
	}

	// Comparar la contraseña hash con la proporcionada
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(contraseña))
	if err != nil {
		return nil, fmt.Errorf("contraseña incorrecta: %v", err)
	}

	// Si la contraseña es correcta, devolver el usuario
	return usuario, nil
}

// Registrarse permite a un nuevo usuario registrarse en el sistema.
// Cifra la contraseña y guarda los detalles del usuario en la base de datos.
func Registrarse(repo UsuarioRepository, nombre, correoElectronico, contraseña string) error {
	// Generar un hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(contraseña), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al encriptar la contraseña: %v", err)
	}

	// Convertir hashedPassword a string
	hashedPasswordStr := string(hashedPassword)

	// Crear un nuevo usuario y guardarlo en la base de datos
	user := &Usuario{
		Nombre:            nombre,
		CorreoElectronico: correoElectronico,
	}
	err = repo.SaveUser(user, hashedPasswordStr)
	if err != nil {
		return err
	}

	return nil
}
