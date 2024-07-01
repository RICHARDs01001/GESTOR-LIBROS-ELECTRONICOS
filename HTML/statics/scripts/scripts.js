// Obtén los modales
var loginModal = document.getElementById('loginModal');
var registerModal = document.getElementById('registerModal');

// Obtén los botones que abren los modales
var loginBtn = document.getElementById('loginBtn');
var registerBtn = document.getElementById('registerBtn');

// Obtén los elementos <span> que cierran los modales
var closeBtns = document.getElementsByClassName('close');

// Cuando el usuario hace clic en el botón de inicio de sesión, abre el modal de inicio de sesión
loginBtn.onclick = function(event) {
    event.preventDefault();
    loginModal.style.display = 'block';
}

// Cuando el usuario hace clic en el botón de registro, abre el modal de registro
registerBtn.onclick = function(event) {
    event.preventDefault();
    registerModal.style.display = 'block';
}

// Cuando el usuario hace clic en <span> (x), cierra el modal
for (var i = 0; i < closeBtns.length; i++) {
    closeBtns[i].onclick = function() {
        this.parentElement.parentElement.parentElement.style.display = 'none';
    }
}

// Cuando el usuario hace clic en cualquier lugar fuera del modal, cierra el modal
window.onclick = function(event) {
    if (event.target == loginModal) {
        loginModal.style.display = 'none';
    } else if (event.target == registerModal) {
        registerModal.style.display = 'none';
    }
}

// Manejar registro
document.querySelector('.register input[type="submit"]').addEventListener('click', function(event) {
    event.preventDefault();
    const username = document.querySelector('#registerModal #username').value;
    const email = document.querySelector('#registerModal #email').value;
    const password = document.querySelector('#registerModal #password').value;

    fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username: username,
            email: email,
            password: password,
        }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.message === 'Usuario registrado con éxito') {
            window.location.href = '/principal';
        } else {
            alert(data.message);
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });
});

// Manejar inicio de sesión
document.querySelector('.login input[type="submit"]').addEventListener('click', function(event) {
    event.preventDefault();
    const email = document.querySelector('#loginModal #loginEmail').value;
    const password = document.querySelector('#loginModal #loginPassword').value;

    fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            email: email,
            password: password,
        }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.message === 'Usuario autenticado con éxito') {
            window.location.href = '/principal';
        } else {
            alert(data.message);
        }
    })
    .catch((error) => {
        console.error('Error:', error);
    });
});

// Funciones para abrir y cerrar el modal
function openModal(modalId) {
    document.getElementById(modalId).style.display = "block";
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = "none";
}

// Cargar reseñas y abrir modal de reseñas
function openReviewsModal(bookId) {
    fetch(`/reviews/${bookId}`)
        .then(response => response.json())
        .then(data => {
            const reviewsContent = document.getElementById('reviewsContent');
            reviewsContent.innerHTML = '';
            data.forEach(review => {
                reviewsContent.innerHTML += `<p>${review.texto} - Calificación: ${review.calificacion}</p>`;
            });
            document.getElementById('addReviewForm').onsubmit = function(event) {
                event.preventDefault();
                addReview(bookId);
            };
            openModal('reviewsModal');
        })
        .catch(error => console.error('Error:', error));
}

// Agregar reseña
function addReview(bookId) {
    const reviewText = document.getElementById('reviewText').value;
    const reviewRating = document.getElementById('reviewRating').value;

    fetch('/reviews', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            libro_id: bookId,
            texto: reviewText,
            calificacion: reviewRating,
        }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.message === 'Reseña agregada con éxito') {
                openReviewsModal(bookId); // Recargar las reseñas
            } else {
                alert(data.message);
            }
        })
        .catch(error => console.error('Error:', error));
}

// Agregar a favoritos
function addToFavorites(bookId) {
    fetch('/favorites', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            libro_id: bookId,
        }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.message === 'Libro agregado a favoritos') {
                alert('Libro agregado a favoritos');
            } else {
                alert(data.message);
            }
        })
        .catch(error => console.error('Error:', error));
}

// Cargar favoritos y abrir modal de favoritos
function loadFavorites() {
    fetch('/favorites')
        .then(response => response.json())
        .then(data => {
            const favoritesContent = document.getElementById('favoritesContent');
            favoritesContent.innerHTML = '';
            data.forEach(book => {
                favoritesContent.innerHTML += `
                    <div class="book">
                        <img src="${book.coverImage}" alt="Portada del libro ${book.titulo}">
                        <div class="book-options">
                            <ul>
                                <li><a href="${book.pdfPath}" target="_blank">Leer libro</a></li>
                                <li><a href="#" onclick="openReviewsModal(${book.id})">Ver reseñas</a></li>
                                <li><a href="#" onclick="addToFavorites(${book.id})">+ Favoritos</a></li>
                            </ul>
                        </div>
                    </div>
                `;
            });
            openModal('favoritesModal');
        })
        .catch(error => console.error('Error:', error));
}

// Salir
function logout() {
    window.location.href = '/';
}

// Búsqueda de libros
document.getElementById('searchButton').addEventListener('click', function() {
    const searchTerm = document.getElementById('searchInput').value;
    fetch(`/search?term=${searchTerm}`)
        .then(response => response.json())
        .then(data => {
            const bookGrid = document.getElementById('bookGrid');
            bookGrid.innerHTML = '';
            data.forEach(book => {
                bookGrid.innerHTML += `
                    <div class="book">
                        <img src="${book.coverImage}" alt="Portada del libro ${book.titulo}">
                        <div class="book-options">
                            <ul>
                                <li><a href="${book.pdfPath}" target="_blank">Leer libro</a></li>
                                <li><a href="#" onclick="openReviewsModal(${book.id})">Ver reseñas</a></li>
                                <li><a href="#" onclick="addToFavorites(${book.id})">+ Favoritos</a></li>
                            </ul>
                        </div>
                    </div>
                `;
            });
        })
        .catch(error => console.error('Error:', error));
});




