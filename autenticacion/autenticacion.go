package autenticacion

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Usuario struct {
	ID                int    `json:"id"`
	Nombre            string `json:"nombre"`
	CorreoElectronico string `json:"correo_electronico"`
}

type Credentials struct {
	Nombre            string `json:"nombre"`
	CorreoElectronico string `json:"correo_electronico"`
	Contraseña        string `json:"contraseña"`
}

type UsuarioRepositorySQL struct {
	DB *sql.DB
}

func IniciarSesion(repo *UsuarioRepositorySQL, correoElectronico, contraseña string) (*Usuario, error) {
	var usuario Usuario
	var hashedPassword string

	err := repo.DB.QueryRow("SELECT id, nombre, contraseña FROM Usuarios WHERE correo_electronico = ?", correoElectronico).Scan(&usuario.ID, &usuario.Nombre, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(contraseña))
	if err != nil {
		return nil, errors.New("contraseña incorrecta")
	}
	return &usuario, nil
}

func Registrarse(repo *UsuarioRepositorySQL, nombre, correoElectronico, hashedPassword string) error {
	_, err := repo.DB.Exec("INSERT INTO Usuarios (nombre, correo_electronico, contraseña) VALUES (?, ?, ?)", nombre, correoElectronico, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}
