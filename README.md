# Gestor de Libros

Este proyecto es una aplicación de gestión de libros desarrollada en Go. Permite a los usuarios iniciar sesión, buscar libros, agregar libros a favoritos y leer libros en formato PDF.

## Características

- **Inicio de sesión de usuarios**
- **Registro de nuevos usuarios**
- **Búsqueda de libros por título o autor**
- **Agregado de libros a la lista de favoritos**
- **Visualización de la lista de favoritos**
- **Lectura de libros en formato PDF**

## Requisitos

- Go (>=1.13)
- MySQL (>=5.7)

## Instalación

1. **Clona el repositorio:**

   ```bash
   git clone https://github.com/tu_usuario/gestor-libros.git
   
2. Instala las dependencias:
go mod download

3. Configuración de la base de datos: En el archivo interfaz_main.go, asegúrate de ingresar el nombre de usuario y contraseña de tu propia base de datos MySQL en las siguientes líneas:

db, err := sql.Open("mysql", "tu_usuario_mysql:tu_contraseña_mysql@tcp(tu_host_mysql:tu_puerto_mysql)/gestor_lib_electr")

4. Configuración de los datos de los libros: En el archivo .sql, asegúrate de insertar la información de los libros o documentos en PDF de acuerdo a cada usuario que desee usar el código. Asegúrate de incluir la ruta de los PDF.

5. Configuración del lector de PDF: En el archivo lector_libros.go, en caso de que prefieras cambiar la forma y la fuente desde la que se abre el PDF, asegúrate de editar la ruta del programa en la línea correspondiente:

cmd := exec.Command(`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`, rutaPdf)

                                  Uso

Inicia sesión si ya tienes una cuenta, o regístrate si eres un nuevo usuario.

Explora la biblioteca buscando libros por título o autor.

Agrega libros a tu lista de favoritos para acceder a ellos fácilmente más tarde.

Lee tus libros favoritos en formato PDF directamente desde la aplicación.


