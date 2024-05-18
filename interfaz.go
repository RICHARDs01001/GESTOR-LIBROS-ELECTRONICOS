package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gestor_libros/autenticacion"
	"gestor_libros/buscar_libros"
	"gestor_libros/favoritos_libros"
	libros "gestor_libros/lector_libros"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Conectar a la base de datos
	db, err := sql.Open("mysql", "INGRESARUSUARIO Y CLAVE @/gestor_lib_electr")
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Bienvenido a nuestra aplicación de libros electronicos!")

	for {
		// Mostrar opciones al usuario
		fmt.Println("Por favor, elige una opción:")
		fmt.Println("1. Iniciar sesión")
		fmt.Println("2. Registrarse")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			iniciarSesion(reader, db)
		case "2":
			registrarse(reader, db)
		default:
			fmt.Println("Opción no válida. Por favor, elige 1 para iniciar sesión o 2 para registrarte.")
		}
	}
}

func iniciarSesion(reader *bufio.Reader, db *sql.DB) {
	// Solicitar credenciales de usuario
	fmt.Print("Ingresa tu correo electrónico: ")
	correoElectronico, _ := reader.ReadString('\n')
	correoElectronico = strings.TrimSpace(correoElectronico)

	fmt.Print("Ingresa tu contraseña: ")
	contraseña, _ := reader.ReadString('\n')
	contraseña = strings.TrimSpace(contraseña)

	// Intentar iniciar sesión
	usuario, err := autenticacion.IniciarSesion(db, correoElectronico, contraseña)
	if err != nil {
		fmt.Printf("Error al iniciar sesión: %v\n", err)
		return
	}
	fmt.Printf("Bienvenido, %s!\n", usuario.Nombre)

	// Mostrar menú principal
	mostrarMenuUsuario(reader, db, usuario.ID)
}

func registrarse(reader *bufio.Reader, db *sql.DB) {
	// Solicitar datos de registro
	fmt.Print("Ingresa tu nombre: ")
	nombre, _ := reader.ReadString('\n')
	nombre = strings.TrimSpace(nombre)

	fmt.Print("Ingresa tu correo electrónico: ")
	correoElectronico, _ := reader.ReadString('\n')
	correoElectronico = strings.TrimSpace(correoElectronico)

	fmt.Print("Ingresa tu contraseña: ")
	contraseña, _ := reader.ReadString('\n')
	contraseña = strings.TrimSpace(contraseña)

	// Intentar registrarse
	err := autenticacion.Registrarse(db, nombre, correoElectronico, contraseña)
	if err != nil {
		fmt.Printf("Error al registrarse: %v\n", err)
		return
	}

	fmt.Println("Usuario registrado con éxito. Por favor, inicia sesión.")
}

func mostrarMenuUsuario(reader *bufio.Reader, db *sql.DB, userID int) {
	for {
		// Mostrar opciones del menú principal
		fmt.Println("Por favor, elige una opción:")
		fmt.Println("1. Buscar libros")
		fmt.Println("2. Ver libros favoritos")
		fmt.Println("3. Ver todos los libros")
		fmt.Println("4. Cerrar sesión")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			buscarLibros(reader, db, userID)
		case "2":
			verFavoritos(reader, db, userID)
		case "3":
			verTodosLosLibros(reader, db, userID)
		case "4":
			fmt.Println("Has cerrado sesión exitosamente.")
			return
		default:
			fmt.Println("Opción no válida. Por favor, elige 1 para buscar libros, 2 para ver tus libros favoritos, 3 para ver todos los libros o 4 para cerrar sesión.")
		}
	}
}

func buscarLibros(reader *bufio.Reader, db *sql.DB, userID int) {
	// Solicitar término de búsqueda
	fmt.Print("Ingresa tu búsqueda: ")
	busqueda, _ := reader.ReadString('\n')
	busqueda = strings.TrimSpace(busqueda)

	// Buscar libros
	librosEncontrados, err := buscar_libros.BuscarLibros(db, busqueda)
	if err != nil {
		fmt.Printf("Error al buscar libros: %v\n", err)
		return
	}

	// Mostrar resultados
	for i, libro := range librosEncontrados {
		fmt.Printf("%d. Título: %s, Autores: %s, ISBN: %s, Año de Publicación: %d\n", i+1, libro.Titulo, libro.Autores, libro.ISBN, libro.AñoPublicacion)
	}

	fmt.Println("Por favor, elige una opción:")
	fmt.Println("1. Elige un libro")
	fmt.Println("2. Realizar nueva búsqueda")

	opcion, _ := reader.ReadString('\n')
	opcion = strings.TrimSpace(opcion)

	switch opcion {
	case "1":
		fmt.Print("Ingresa el número del libro que deseas elegir: ")
		numeroLibroStr, _ := reader.ReadString('\n')
		numeroLibro, err := strconv.Atoi(strings.TrimSpace(numeroLibroStr))
		if err != nil || numeroLibro < 1 || numeroLibro > len(librosEncontrados) {
			fmt.Println("Número de libro no válido. Por favor, intenta de nuevo.")
		} else {
			libroElegido := librosEncontrados[numeroLibro-1]
			fmt.Printf("Has elegido el libro: %s\n", libroElegido.Titulo)
			mostrarMenuLibro(reader, db, userID, libroElegido.ID)
		}
	case "2":
		buscarLibros(reader, db, userID)
	default:
		fmt.Println("Opción no válida. Por favor, elige 1 para elegir un libro o 2 para realizar una nueva búsqueda.")
	}
}

func mostrarMenuLibro(reader *bufio.Reader, db *sql.DB, userID, libroID int) {
	for {
		// Mostrar opciones para el libro elegido
		fmt.Println("Por favor, elige una opción:")
		fmt.Println("1. Leer libro")
		fmt.Println("2. Agregar a favoritos")
		fmt.Println("3. Regresar al menú principal")

		opcionLibro, _ := reader.ReadString('\n')
		opcionLibro = strings.TrimSpace(opcionLibro)

		switch opcionLibro {
		case "1":
			// Obtener y leer libro
			rutaPdf, err := libros.ObtenerRutaPdf(db, libroID)
			if err != nil {
				fmt.Printf("Error al obtener la ruta del PDF: %v\n", err)
				return
			}

			err = libros.LeerLibro(rutaPdf)
			if err != nil {
				fmt.Printf("Error al leer el libro: %v\n", err)
				return
			}
		case "2":
			// Agregar libro a favoritos
			err := favoritos_libros.AgregarLibroAFavoritos(db, userID, libroID)
			if err != nil {
				fmt.Printf("Error al agregar el libro a favoritos: %v\n", err)
				return
			}
			fmt.Println("Libro agregado a favoritos exitosamente.")
		case "3":
			return
		default:
			fmt.Println("Opción no válida. Por favor, elige 1 para leer libros, 2 para agregar a favoritos o 3 para regresar al menú principal.")
		}
	}
}

func verFavoritos(reader *bufio.Reader, db *sql.DB, userID int) {
	// Listar libros favoritos
	favoritos, err := favoritos_libros.ListarFavoritos(db, userID)
	if err != nil {
		fmt.Printf("Error al listar favoritos: %v\n", err)
		return
	}

	// Mostrar favoritos
	for i, favorito := range favoritos {
		fmt.Printf("%d. Título: %s, Autor: %s, ISBN: %s, Año de Publicación: %d\n", i+1, favorito.Titulo, favorito.Autor, favorito.ISBN, favorito.AñoPublicacion)
	}

	fmt.Println("Por favor, elige una opción:")
	fmt.Println("1. Elige un libro para leer")
	fmt.Println("2. Regresar al menú principal")

	opcion, _ := reader.ReadString('\n')
	opcion = strings.TrimSpace(opcion)

	switch opcion {
	case "1":
		fmt.Print("Ingresa el número del libro que deseas leer: ")
		numeroLibroStr, _ := reader.ReadString('\n')
		numeroLibro, err := strconv.Atoi(strings.TrimSpace(numeroLibroStr))
		if err != nil || numeroLibro < 1 || numeroLibro > len(favoritos) {
			fmt.Println("Número de libro no válido. Por favor, intenta de nuevo.")
		} else {
			libroElegido := favoritos[numeroLibro-1]
			rutaPdf, err := libros.ObtenerRutaPdf(db, libroElegido.LibroID)
			if err != nil {
				fmt.Printf("Error al obtener la ruta del PDF: %v\n", err)
				return
			}

			err = libros.LeerLibro(rutaPdf)
			if err != nil {
				fmt.Printf("Error al leer el libro: %v\n", err)
				return
			}
		}
	case "2":
		return
	default:
		fmt.Println("Opción no válida. Por favor, elige 1 para leer un libro o 2 para regresar al menú principal.")
	}
}

func verTodosLosLibros(reader *bufio.Reader, db *sql.DB, userID int) {
	// Buscar todos los libros (sin término de búsqueda)
	librosEncontrados, err := buscar_libros.BuscarLibros(db, "")
	if err != nil {
		fmt.Printf("Error al buscar libros: %v\n", err)
		return
	}

	// Mostrar libros encontrados
	for i, libro := range librosEncontrados {
		fmt.Printf("%d. Título: %s, Autores: %s, ISBN: %s, Año de Publicación: %d\n", i+1, libro.Titulo, libro.Autores, libro.ISBN, libro.AñoPublicacion)
	}

	fmt.Println("Por favor, elige una opción:")
	fmt.Println("1. Elige un libro")
	fmt.Println("2. Regresar al menú principal")

	opcion, _ := reader.ReadString('\n')
	opcion = strings.TrimSpace(opcion)

	switch opcion {
	case "1":
		fmt.Print("Ingresa el número del libro que deseas elegir: ")
		numeroLibroStr, _ := reader.ReadString('\n')
		numeroLibro, err := strconv.Atoi(strings.TrimSpace(numeroLibroStr))
		if err != nil || numeroLibro < 1 || numeroLibro > len(librosEncontrados) {
			fmt.Println("Número de libro no válido. Por favor, intenta de nuevo.")
		} else {
			libroElegido := librosEncontrados[numeroLibro-1]
			fmt.Printf("Has elegido el libro: %s\n", libroElegido.Titulo)
			mostrarMenuLibro(reader, db, userID, libroElegido.ID)
		}
	case "2":
		return
	default:
		fmt.Println("Opción no válida. Por favor, elige 1 para elegir un libro o 2 para regresar al menú principal.")
	}
}
