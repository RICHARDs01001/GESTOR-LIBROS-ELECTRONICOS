package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"gestor_libros/autenticacion"
	"gestor_libros/buscar_libros"
	"gestor_libros/favoritos_libros"
	"gestor_libros/resenas_libros"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {
	var err error
	// Conectar a la base de datos (modifica la cadena de conexión según tus credenciales)
	db, err = sql.Open("")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	r := gin.Default()

	r.Static("/statics", "./HTML/statics")
	r.LoadHTMLGlob("HTML/templates/*")

	r.GET("/", serveIndex)
	r.GET("/principal", servePrincipal)

	r.POST("/register", handleRegister)
	r.POST("/login", handleLogin)
	r.GET("/read/:bookID", handleReadBook)

	// Rutas adicionales
	r.GET("/search", searchBooks)
	r.GET("/reviews/:bookID", listReviews)
	r.POST("/reviews", addReview)
	r.GET("/favorites", listFavorites)
	r.POST("/favorites", addFavorite)

	r.Run(":8080") // Inicia el servidor en el puerto 8080
}

func serveIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func servePrincipal(c *gin.Context) {
	c.HTML(http.StatusOK, "principal.html", nil)
}

func handleRegister(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos inválidos"})
		return
	}

	// Generar el hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al generar el hash de la contraseña"})
		return
	}

	repo := &autenticacion.UsuarioRepositorySQL{DB: db}
	err = autenticacion.Registrarse(repo, user.Username, user.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al registrar el usuario"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario registrado con éxito"})
}

func handleLogin(c *gin.Context) {
	var creds autenticacion.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos inválidos"})
		return
	}

	log.Println("Intentando iniciar sesión con correo:", creds.CorreoElectronico)

	repo := &autenticacion.UsuarioRepositorySQL{DB: db}
	user, err := autenticacion.IniciarSesion(repo, creds.CorreoElectronico, creds.Contraseña)
	if err != nil {
		log.Println("Error en IniciarSesion:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Usuario o contraseña incorrectos"})
		return
	}

	log.Println("Inicio de sesión exitoso para usuario:", user.Nombre)
	c.JSON(http.StatusOK, gin.H{"message": "Usuario autenticado con éxito", "user": user})
}

func handleReadBook(c *gin.Context) {
	bookID := c.Param("bookID")
	bookPath := fmt.Sprintf("./LIBROS/%s.pdf", bookID)

	_, err := os.Stat(bookPath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Libro no encontrado"})
		return
	}

	c.File(bookPath)
}

func searchBooks(c *gin.Context) {
	term := c.Query("term")
	libros, err := buscar_libros.BuscarLibros(db, term)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al buscar libros"})
		return
	}
	c.JSON(http.StatusOK, libros)
}

func listFavorites(c *gin.Context) {
	userID := c.Query("userID")
	favoritos, err := favoritos_libros.ListarFavoritos(db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al listar favoritos"})
		return
	}
	c.JSON(http.StatusOK, favoritos)
}

func addFavorite(c *gin.Context) {
	var fav favoritos_libros.Favorito
	if err := c.ShouldBindJSON(&fav); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos inválidos"})
		return
	}

	err := favoritos_libros.AgregarLibroAFavoritos(db, fav.UsuarioID, fav.LibroID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al agregar a favoritos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Libro agregado a favoritos"})
}

func listReviews(c *gin.Context) {
	bookID := c.Param("bookID")
	resenas, err := resenas_libros.ListarReseñas(db, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al listar reseñas"})
		return
	}
	c.JSON(http.StatusOK, resenas)
}

func addReview(c *gin.Context) {
	var resena resenas_libros.Reseña
	if err := c.ShouldBindJSON(&resena); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos inválidos"})
		return
	}

	err := resenas_libros.InsertarReseña(db, resena.UsuarioID, resena.LibroID, resena.Calificacion, resena.Texto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al agregar reseña"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reseña agregada con éxito"})
}
