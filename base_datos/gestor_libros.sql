DROP DATABASE IF EXISTS gestor_lib_electr;
CREATE DATABASE gestor_lib_electr;
USE gestor_lib_electr;

CREATE TABLE Autores (
    id INT AUTO_INCREMENT,
    nombre VARCHAR(100),
    PRIMARY KEY (id)
);

CREATE TABLE Editoriales (
    id INT AUTO_INCREMENT,
    nombre VARCHAR(100),
    PRIMARY KEY (id)
);

CREATE TABLE Libros (
    id INT AUTO_INCREMENT,
    titulo VARCHAR(200),
    isbn VARCHAR(20),
    año_publicacion YEAR,
    numero_paginas INT,
    ruta_pdf VARCHAR(255), 
    editorial_id INT,
    PRIMARY KEY (id),
    FOREIGN KEY (editorial_id) REFERENCES Editoriales(id)
);



CREATE TABLE Usuarios (
    id INT AUTO_INCREMENT,
    nombre VARCHAR(100),
    correo_electronico VARCHAR(100),
    contraseña VARCHAR(100),
    PRIMARY KEY (id)
);

CREATE TABLE Favoritos (
    usuario_id INT,
    libro_id INT,
    PRIMARY KEY (usuario_id, libro_id),
    FOREIGN KEY (usuario_id) REFERENCES Usuarios(id),
    FOREIGN KEY (libro_id) REFERENCES Libros(id)
);

CREATE TABLE Reseñas (
    id INT AUTO_INCREMENT,
    usuario_id INT,
    libro_id INT,
    calificacion INT,
    texto TEXT,
    fecha DATE,
    PRIMARY KEY (id),
    FOREIGN KEY (usuario_id) REFERENCES Usuarios(id),
    FOREIGN KEY (libro_id) REFERENCES Libros(id)
);

CREATE TABLE Libros_Autores (
    libro_id INT,
    autor_id INT,
    PRIMARY KEY (libro_id, autor_id),
    FOREIGN KEY (libro_id) REFERENCES Libros(id),
    FOREIGN KEY (autor_id) REFERENCES Autores(id)
);



