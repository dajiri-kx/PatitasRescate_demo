<!-- home.php -->
<!DOCTYPE html>
<html lang="es" data-bs-theme="auto">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Patitas al rescate - Inicio</title>
    <meta name="description" content="Página principal de Patitas al rescate, ofreciendo servicios veterinarios y estéticos para tus mascotas.">
    <meta name="keywords" content="veterinaria, mascotas, estética, cuidado de mascotas, agendar cita">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-icons/1.10.2/font/bootstrap-icons.min.css" rel="stylesheet">
    <style>
        .banner {
            background: url('https://img.freepik.com/foto-gratis/vista-frontal-hermoso-perro-espacio-copia_23-2148786560.jpg?t=st=1741669395~exp=1741672995~hmac=14b1d4b75eeb86cd027c3840dc72b34813358280e29c8b310ddf01488929d360&w=1380') no-repeat center center;
            background-size: cover;
            height: 400px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
        }
        .testimonios {
            background-color: #f8f9fa;
            padding: 50px 0;
        }
        .testimonio-card {
            background: white;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
    </style>
</head>

<body>
    <header>
        <?php include 'header.php'; ?>
    </header>

    <!-- Banner con slogan e imágenes -->
    <section class="banner">
        <div class="text-center">
            <h1>Bienvenidos a Patitas al Rescate</h1>
            <p class="lead">Cuidamos de tus mascotas como si fueran nuestras</p>
            <a href="agendarCita.php" class="btn btn-primary btn-lg">Agendar Cita</a>
        </div>
    </section>

   <!-- Sección de Servicios -->
<section class="container my-5">
    <div class="row">
        <!-- Servicios Veterinarios -->
        <div class="col-md-6 mb-4">
            <div class="card h-100 shadow">
                <div class="card-body text-center">
                    <i class="bi bi-heart-pulse display-4 text-primary mb-3"></i> <!-- Ícono de servicios veterinarios -->
                    <h2 class="card-title">Servicios Veterinarios</h2>
                    <p class="card-text">
                        Ofrecemos una amplia gama de servicios veterinarios para asegurar la salud y bienestar de tus mascotas. Desde consultas generales hasta cirugías especializadas.
                    </p>
                    <ul class="list-unstyled">
                        <li><i class="bi bi-check-circle text-success"></i> Consultas generales</li>
                        <li><i class="bi bi-check-circle text-success"></i> Vacunación</li>
                        <li><i class="bi bi-check-circle text-success"></i> Cirugías</li>
                        <li><i class="bi bi-check-circle text-success"></i> Diagnóstico por imágenes</li>
                        <li><i class="bi bi-check-circle text-success"></i> Hospitalización</li>
                    </ul>
                    <a href="#" class="btn btn-primary">Más información</a>
                </div>
            </div>
        </div>

        <!-- Servicios Estéticos -->
        <div class="col-md-6 mb-4">
            <div class="card h-100 shadow">
                <div class="card-body text-center">
                    <i class="bi bi-scissors display-4 text-warning mb-3"></i> <!-- Ícono de servicios estéticos -->
                    <h2 class="card-title">Servicios Estéticos</h2>
                    <p class="card-text">
                        Nuestros servicios estéticos incluyen baños, cortes de pelo, uñas y más, para que tu mascota no solo esté saludable, sino también radiante.
                    </p>
                    <ul class="list-unstyled">
                        <li><i class="bi bi-check-circle text-success"></i> Baños y secado profesional</li>
                        <li><i class="bi bi-check-circle text-success"></i> Cortes de pelo a medida</li>
                        <li><i class="bi bi-check-circle text-success"></i> Corte de uñas</li>
                        <li><i class="bi bi-check-circle text-success"></i> Limpieza dental</li>
                        <li><i class="bi bi-check-circle text-success"></i> Perfumería para mascotas</li>
                    </ul>
                    <a href="#" class="btn btn-warning">Más información</a>
                </div>
            </div>
        </div>
    </div>
</section>

   <!-- Información de contacto y ubicación -->
<section class="contacto py-5 bg-white">
    <div class="container">
        <div class="row">
            <!-- Información de contacto -->
            <div class="col-md-6 mb-4">
                <div class="card h-100 shadow-sm border-0">
                    <div class="card-body text-center">
                        <i class="bi bi-telephone display-4 text-primary mb-3"></i>
                        <h2 class="card-title">Contacto</h2>
                        <p class="card-text">
                            Estamos aquí para ayudarte. ¡Contáctanos!
                        </p>
                        <ul class="list-unstyled">
                            <li><i class="bi bi-envelope me-2"></i>info@patitasalrescate.com</li>
                            <li><i class="bi bi-phone me-2"></i>+506 1234 5678</li>
                            <li><i class="bi bi-clock me-2"></i>Lunes a Viernes: 8:00 AM - 6:00 PM</li>
                        </ul>
                        <a href="#" class="btn btn-primary">Enviar mensaje</a>
                    </div>
                </div>
            </div>

            <!-- Ubicación -->
            <div class="col-md-6 mb-4">
                <div class="card h-100 shadow-sm border-0">
                    <div class="card-body text-center">
                        <i class="bi bi-geo-alt display-4 text-success mb-3"></i>
                        <h2 class="card-title">Ubicación</h2>
                        <p class="card-text">
                            Visítanos en nuestra clínica.
                        </p>
                        <ul class="list-unstyled">
                            <li><i class="bi bi-house-door me-2"></i>San José, Costa Rica</li>
                            <li><i class="bi bi-signpost me-2"></i>100 metros norte del Parque Central</li>
                        </ul>
                        <div class="ratio ratio-16x9">
                            <iframe src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3929.982734693637!2d-84.077288685227!3d9.93536587651835!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x8fa0e3f5c3b5c5b5%3A0x5c5c5c5c5c5c5c5c!2sParque%20Central%2C%20San%20Jos%C3%A9%2C%20Costa%20Rica!5e0!3m2!1ses!2scr!4v1641234567890!5m2!1ses!2scr" style="border:0;" allowfullscreen="" loading="lazy"></iframe>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</section>

   <!-- Testimonios de clientes -->
<section class="testimonios py-5 bg-light">
    <div class="container">
        <h2 class="text-center mb-5">Calificaciones y reseñas</h2>
        <div class="row row-cols-1 row-cols-md-3 g-4">
            <!-- Testimonio 1 -->
            <div class="col">
                <div class="card h-100 shadow-sm border-0">
                    <div class="card-body text-center">
                        <img src="https://static-00.iconduck.com/assets.00/profile-default-icon-2048x2045-u3j7s5nj.png" class="rounded-circle mb-3" alt="Cliente 1" style="width: 80px; height: 80px;">
                        <h5 class="card-title">María Gómez</h5>
                        <p class="card-text">
                            "Excelente servicio, mi perro fue tratado con mucho cuidado y profesionalismo."
                        </p>
                        <div class="text-warning">
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Testimonio 2 -->
            <div class="col">
                <div class="card h-100 shadow-sm border-0">
                    <div class="card-body text-center">
                        <img src="https://static-00.iconduck.com/assets.00/profile-default-icon-2048x2045-u3j7s5nj.png" class="rounded-circle mb-3" alt="Cliente 2" style="width: 80px; height: 80px;">
                        <h5 class="card-title">Juan Pérez</h5>
                        <p class="card-text">
                            "Muy contento con el corte de pelo de mi gato, quedó hermoso."
                        </p>
                        <div class="text-warning">
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-half"></i>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Testimonio 3 -->
            <div class="col">
                <div class="card h-100 shadow-sm border-0">
                    <div class="card-body text-center">
                        <img src="https://static-00.iconduck.com/assets.00/profile-default-icon-2048x2045-u3j7s5nj.png" class="rounded-circle mb-3" alt="Cliente 3" style="width: 80px; height: 80px;">
                        <h5 class="card-title">Ana Rodríguez</h5>
                        <p class="card-text">
                            "Recomiendo totalmente Patitas al Rescate, siempre atentos y amables."
                        </p>
                        <div class="text-warning">
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                            <i class="bi bi-star-fill"></i>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</section>

    <footer>
        <?php include 'footer.php'; ?>
    </footer>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
</body>

</html>