package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"gestor_libros/autenticacion"
	"gestor_libros/buscar_libros"
	"gestor_libros/favoritos_libros"
	"gestor_libros/lector_libros"
	"gestor_libros/resenas_libros"
)

func main() {
	db, err := sql.Open("mysql", "")
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer db.Close()

	usuarioRepo := &autenticacion.UsuarioRepositorySQL{DB: db}
	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		var credentials autenticacion.Credentials
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}
		usuario, err := autenticacion.IniciarSesion(usuarioRepo, credentials.CorreoElectronico, credentials.Contraseña)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, usuario)
	})

	router.POST("/register", func(c *gin.Context) {
		var credentials autenticacion.Credentials
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}
		err := autenticacion.Registrarse(usuarioRepo, credentials.Nombre, credentials.CorreoElectronico, credentials.Contraseña)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Usuario registrado con éxito"})
	})

	router.GET("/books", func(c *gin.Context) {
		termino := c.Query("search")
		libros, err := buscar_libros.BuscarLibros(db, termino)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, libros)
	})

	router.GET("/favorites/:userID", func(c *gin.Context) {
		usuarioID := c.Param("userID")
		favoritos, err := favoritos_libros.ListarFavoritos(db, usuarioID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, favoritos)
	})

	router.POST("/favorites", func(c *gin.Context) {
		var favorito favoritos_libros.Favorito
		if err := c.ShouldBindJSON(&favorito); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}
		err := favoritos_libros.AgregarLibroAFavoritos(db, favorito.UsuarioID, favorito.LibroID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Libro agregado a favoritos"})
	})

	router.GET("/reviews/:bookID", func(c *gin.Context) {
		libroID := c.Param("bookID")
		reseñas, err := resenas_libros.ListarReseñas(db, libroID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, reseñas)
	})

	router.POST("/reviews", func(c *gin.Context) {
		var reseña resenas_libros.Reseña
		if err := c.ShouldBindJSON(&reseña); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}
		err := resenas_libros.InsertarReseña(db, reseña.UsuarioID, reseña.LibroID, reseña.Calificacion, reseña.Texto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Reseña agregada con éxito"})
	})

	router.GET("/read/:bookID", func(c *gin.Context) {
		libroID, err := strconv.Atoi(c.Param("bookID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de libro inválido"})
			return
		}
		libro, err := lector_libros.LeerLibro(db, libroID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, libro)
	})

	router.Run(":8080")
}
