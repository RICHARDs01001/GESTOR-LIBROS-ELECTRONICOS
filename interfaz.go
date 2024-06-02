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
	"gestor_libros/lector_libros"
	"gestor_libros/resenas_libros"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:001002003@Rsg@/gestor_lib_electr")
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Bienvenido a nuestra aplicación de libros electrónicos!")

	for {
		fmt.Println("Por favor, elige una opción:")
		fmt.Println("1. Iniciar sesión")
		fmt.Println("2. Registrarse")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			continue
		}
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
	fmt.Print("Ingresa tu correo electrónico: ")
	correoElectronico, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	correoElectronico = strings.TrimSpace(correoElectronico)

	fmt.Print("Ingresa tu contraseña: ")
	contraseña, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	contraseña = strings.TrimSpace(contraseña)

	usuario, err := autenticacion.IniciarSesion(db, correoElectronico, contraseña)
	if err != nil {
		fmt.Printf("Error al iniciar sesión: %v\n", err)
		return
	}
	fmt.Printf("Bienvenido, %s!\n", usuario.Nombre)

	mostrarMenuUsuario(reader, db, usuario.ID)
}

func registrarse(reader *bufio.Reader, db *sql.DB) {
	fmt.Print("Ingresa tu nombre: ")
	nombre, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	nombre = strings.TrimSpace(nombre)

	fmt.Print("Ingresa tu correo electrónico: ")
	correoElectronico, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	correoElectronico = strings.TrimSpace(correoElectronico)

	fmt.Print("Ingresa tu contraseña: ")
	contraseña, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	contraseña = strings.TrimSpace(contraseña)

	err = autenticacion.Registrarse(db, nombre, correoElectronico, contraseña)
	if err != nil {
		fmt.Printf("Error al registrarse: %v\n", err)
		return
	}

	fmt.Println("Usuario registrado con éxito. Por favor, inicia sesión.")
}

func mostrarMenuUsuario(reader *bufio.Reader, db *sql.DB, usuarioID int) {
	for {
		fmt.Println("Por favor, elige una opción:")
		fmt.Println("1. Buscar libros")
		fmt.Println("2. Ver libros favoritos")
		fmt.Println("3. Salir")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			continue
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			buscarLibros(reader, db, usuarioID)
		case "2":
			verFavoritos(reader, db, usuarioID)
		case "3":
			fmt.Println("Saliendo...")
			return
		default:
			fmt.Println("Opción no válida. Por favor, elige una opción del 1 al 3.")
		}
	}
}

func buscarLibros(reader *bufio.Reader, db *sql.DB, usuarioID int) {
	fmt.Print("Ingresa el término de búsqueda: ")
	termino, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	termino = strings.TrimSpace(termino)

	librosEncontrados, err := buscar_libros.BuscarLibros(db, termino)
	if err != nil {
		fmt.Printf("Error al buscar libros: %v\n", err)
		return
	}

	if len(librosEncontrados) == 0 {
		fmt.Println("No se encontraron libros con ese término.")
		return
	}

	for i, libro := range librosEncontrados {
		fmt.Printf("%d. %s - %s\n", i+1, libro.Titulo, libro.Autores)
	}

	fmt.Print("Elige el número del libro para ver más opciones o presiona Enter para volver: ")
	numeroLibroStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	numeroLibroStr = strings.TrimSpace(numeroLibroStr)
	if numeroLibroStr == "" {
		return
	}

	numeroLibro, err := strconv.Atoi(numeroLibroStr)
	if err != nil || numeroLibro < 1 || numeroLibro > len(librosEncontrados) {
		fmt.Println("Número de libro no válido. Por favor, intenta de nuevo.")
		return
	}

	libroElegido := librosEncontrados[numeroLibro-1]
	mostrarMenuLibro(reader, db, usuarioID, libroElegido.ID)
}

func mostrarMenuLibro(reader *bufio.Reader, db *sql.DB, usuarioID, libroID int) {
	for {
		fmt.Println("Por favor, elige una opción:")
		fmt.Println("1. Leer libro")
		fmt.Println("2. Agregar a favoritos")
		fmt.Println("3. Ver reseñas")
		fmt.Println("4. Agregar reseña")
		fmt.Println("5. Volver")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			continue
		}
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			rutaPdf, err := lector_libros.ObtenerRutaPdf(db, libroID)
			if err != nil {
				fmt.Printf("Error al obtener la ruta del PDF: %v\n", err)
				continue
			}
			err = lector_libros.LeerLibro(rutaPdf)
			if err != nil {
				fmt.Printf("Error al leer el libro: %v\n", err)
			}
		case "2":
			err := favoritos_libros.AgregarLibroAFavoritos(db, usuarioID, libroID)
			if err != nil {
				fmt.Printf("Error al agregar libro a favoritos: %v\n", err)
			} else {
				fmt.Println("Libro agregado a favoritos.")
			}
		case "3":
			verReseñas(db, libroID)
		case "4":
			agregarReseña(reader, db, usuarioID, libroID)
		case "5":
			return
		default:
			fmt.Println("Opción no válida. Por favor, elige una opción del 1 al 5.")
		}
	}
}

func verFavoritos(reader *bufio.Reader, db *sql.DB, usuarioID int) {
	favoritos, err := favoritos_libros.ListarFavoritos(db, usuarioID)
	if err != nil {
		fmt.Printf("Error al obtener lista de favoritos: %v\n", err)
		return
	}

	if len(favoritos) == 0 {
		fmt.Println("No tienes libros favoritos.")
		return
	}

	for i, favorito := range favoritos {
		fmt.Printf("%d. %s - %s\n", i+1, favorito.Titulo, favorito.Autor)
	}

	fmt.Print("Elige el número del libro para ver más opciones o presiona Enter para volver: ")
	numeroLibroStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	numeroLibroStr = strings.TrimSpace(numeroLibroStr)
	if numeroLibroStr == "" {
		return
	}

	numeroLibro, err := strconv.Atoi(numeroLibroStr)
	if err != nil || numeroLibro < 1 || numeroLibro > len(favoritos) {
		fmt.Println("Número de libro no válido. Por favor, intenta de nuevo.")
		return
	}

	libroElegido := favoritos[numeroLibro-1]
	mostrarMenuLibro(reader, db, usuarioID, libroElegido.LibroID)
}

func agregarReseña(reader *bufio.Reader, db *sql.DB, usuarioID, libroID int) {
	fmt.Print("Ingresa el título de tu reseña: ")
	titulo, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	titulo = strings.TrimSpace(titulo)

	fmt.Print("Ingresa el contenido de tu reseña: ")
	contenido, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}
	contenido = strings.TrimSpace(contenido)

	if titulo == "" || contenido == "" {
		fmt.Println("El título y el contenido de la reseña no pueden estar vacíos.")
		return
	}

	// Aquí puedes añadir la lógica para obtener la calificación del usuario si es necesario.
	// Por ahora, estableceremos una calificación fija (por ejemplo, 5)
	calificacion := 5

	err = resenas_libros.InsertarReseña(db, usuarioID, libroID, calificacion, contenido)
	if err != nil {
		fmt.Printf("Error al agregar reseña: %v\n", err)
	} else {
		fmt.Println("Reseña agregada con éxito.")
	}
}
func verReseñas(db *sql.DB, libroID int) {
	reseñas, err := resenas_libros.ListarReseñas(db, libroID)
	if err != nil {
		fmt.Printf("Error al obtener las reseñas: %v\n", err)
		return
	}

	if len(reseñas) == 0 {
		fmt.Println("No hay reseñas para este libro.")
		return
	}

	for _, reseña := range reseñas {
		fmt.Printf("Usuario %d - Calificación: %d\n", reseña.UsuarioID, reseña.Calificacion)
		fmt.Printf("Fecha: %s\n", reseña.Fecha.Format("2006-01-02"))
		fmt.Printf("Texto: %s\n", reseña.Texto)
		fmt.Println("-------------------------------")
	}
}
