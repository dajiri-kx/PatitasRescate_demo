<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Patitas al Rescate - Ubicación</title>
  <meta name="description" content="Conoce la ubicación de nuestra clínica en San José, Costa Rica. Visítanos y descubre cómo llegar.">
  <meta name="keywords" content="ubicación, clínica, San José, Costa Rica, mapa">
  <link rel="icon" href="favicon.ico">
  
  <!-- Preconexión a recursos externos -->
  <link rel="preconnect" href="https://cdn.jsdelivr.net">
  <link rel="preconnect" href="https://cdnjs.cloudflare.com">

  <!-- Bootstrap CSS -->
  <link
    href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css"
    rel="stylesheet"
    integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH"
    crossorigin="anonymous"
  >
  <!-- Bootstrap Icons -->
  <link
    href="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-icons/1.10.2/font/bootstrap-icons.min.css"
    rel="stylesheet"
  >
  <style>
    /* Banner de Ubicación */
    .location-banner {
      background: url('https://img.freepik.com/premium-photo/concept-plastic-waste-pollution-from-plastic-waste-discarded-plastic-bags-black-backgroun_1066580-1557.jpg?w=1380') no-repeat center center;
      background-size: cover;
      height: 300px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.7);
    }
    .location-banner h1 {
      font-size: 3rem;
    }
    /* Fondo y estilo para el contenido principal */
    main {
      background: linear-gradient(to bottom right, #f0f9ff, #cbebff);
      padding: 40px 0;
    }
    .location-card {
      background: rgba(255, 255, 255, 0.95);
      padding: 30px;
      border-radius: 15px;
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
      margin-top: 50px;
    }
    ul li {
      font-size: 1.1rem;
      margin-bottom: 5px;
    }
    
    /* Estilos para el carrusel */
    .carousel-section {
      padding: 30px 0;
    }
    .carousel-img {
      height: 500px;
      object-fit: cover;
      width: 100%;
    }
    .carousel-control-prev-icon,
    .carousel-control-next-icon {
      background-color: rgba(0, 0, 0, 0.5);
      border-radius: 50%;
      padding: 20px;
    }
    .carousel-indicators button {
      background-color: #ff6b00;
    }
  </style>
</head>
<body>
  <!-- Incluir cabecera -->
  <?php include 'header.php'; ?>

  <!-- Banner de Ubicación -->
  <section class="location-banner">
    <div class="text-center">
      <h1>Encuéntranos</h1>
    </div>
  </section>

  <!-- Contenido principal -->
  <main>
    <div class="container">
      <!-- Carrusel de Imágenes -->
      <section class="carousel-section">
        <div id="clinicCarousel" class="carousel slide" data-bs-ride="carousel">
          <div class="carousel-indicators">
            <button type="button" data-bs-target="#clinicCarousel" data-bs-slide-to="0" class="active"></button>
            <button type="button" data-bs-target="#clinicCarousel" data-bs-slide-to="1"></button>
            <button type="button" data-bs-target="#clinicCarousel" data-bs-slide-to="2"></button>
            <button type="button" data-bs-target="#clinicCarousel" data-bs-slide-to="3"></button>
            <button type="button" data-bs-target="#clinicCarousel" data-bs-slide-to="4"></button>
          </div>
          <div class="carousel-inner rounded-3">
            <div class="carousel-item active">
              <img src="https://www.gaceta.unam.mx/wp-content/uploads/2022/08/220815-aca10-f2-ganan-terreno-mujeres-veterinaria.jpg" class="d-block w-100 carousel-img" alt="Profesionales veterinarios">
            </div>
            <div class="carousel-item">
              <img src="https://utn.ac.cr/sites/default/files/carreras/field_images/image2.jpeg" class="d-block w-100 carousel-img" alt="Instalaciones modernas">
            </div>
            <div class="carousel-item">
              <img src="https://www.vistazo.com/binrepository/768x576/0c0/0d0/none/12727/QAEM/campan-a-de-vacunacio-n-en-barrios-a-p_614493_20220919161917.jpg" class="d-block w-100 carousel-img" alt="Servicio de vacunación">
            </div>
            <div class="carousel-item">
              <img src="https://clinicaelpalau.es/wp-content/uploads/2018/09/perro-peluquero.jpg" class="d-block w-100 carousel-img" alt="Servicios estéticos">
            </div>
            <div class="carousel-item">
              <img src="https://okdiario.com/img/2024/10/23/perro-en-veterinario-635x358.jpg" class="d-block w-100 carousel-img" alt="Atención veterinaria">
            </div>
          </div>
          <button class="carousel-control-prev" type="button" data-bs-target="#clinicCarousel" data-bs-slide="prev">
            <span class="carousel-control-prev-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Anterior</span>
          </button>
          <button class="carousel-control-next" type="button" data-bs-target="#clinicCarousel" data-bs-slide="next">
            <span class="carousel-control-next-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Siguiente</span>
          </button>
        </div>
      </section>

      <!-- Sección de Ubicación (al final) -->
      <div class="row justify-content-center">
        <div class="col-md-8 col-lg-6">
          <div class="location-card text-center">
            <i class="bi bi-geo-alt display-4 text-success mb-3" aria-hidden="true"></i>
            <h2>Visítanos en Nuestra Clínica</h2>
            <p class="mb-4">
              Estamos ubicados en el corazón de San José, Costa Rica. ¡Te esperamos!
            </p>
            <ul class="list-unstyled mb-4">
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
  </main>
  <!-- Incluir pie de página -->
  <?php include 'footer.php'; ?>

  <!-- Bootstrap JS (con Popper) -->
  <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
</body>
</html>
