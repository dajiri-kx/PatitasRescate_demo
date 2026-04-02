<!-- servicios.php -->
<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Patitas al Rescate - Servicios Especializados</title>
  <meta name="description" content="Descubre los servicios especializados en Patitas al Rescate para el cuidado y bienestar de tus mascotas.">
  <meta name="keywords" content="servicios, veterinaria, estética, mascotas, salud, bienestar">
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
    /* Banner de servicios */
    .services-banner {
      background: url('https://img.freepik.com/premium-photo/husky-dog-portrait-beautiful-photo-selective-focus_73944-25165.jpg?w=1380') no-repeat center center;
      background-size: cover;
      height: 450px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      text-shadow: 2px 2px 4px rgba(0,0,0,0.6);
    }
    /* Animación al pasar el mouse por las tarjetas */
    .service-card {
      transition: transform 0.3s;
    }
    .service-card:hover {
      transform: translateY(-5px);
    }
    .icon-large {
      font-size: 3rem;
    }
    /* Fondos diferenciados para categorías */
    .bg-vet {
      background-color: #e3f2fd; /* azul claro */
    }
    .bg-estetic {
      background-color: #fff3e0; /* naranja suave */
    }
    .bg-special {
      background-color: #e8f5e9; /* verde suave */
    }
  </style>
</head>
<body>
  <!-- Incluir cabecera -->
  <?php include 'header.php'; ?>

  <!-- Banner Hero para Servicios -->
  <section class="services-banner">
    <div class="text-center">
      <h1>Servicios Especializados</h1>
      <p class="lead">Cuidamos de tus mascotas con profesionalismo y pasión</p>
      <a href="agendarCita.php" class="btn btn-lg btn-primary mt-3">Agenda tu Cita</a>
    </div>
  </section>

  <!-- Sección de introducción -->
  <section class="py-5">
    <div class="container">
      <div class="row align-items-center">
        <div class="col-lg-6">
          <h2 class="mb-3">Nuestros Servicios</h2>
          <p>
            En <strong>Patitas al Rescate</strong> ofrecemos un amplio abanico de servicios para el cuidado integral de tus mascotas. Contamos con profesionales capacitados y tecnología de punta para garantizar tratamientos de calidad, ya sea en salud o estética.
          </p>
        </div>
        <div class="col-lg-6">
          <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRh8REzUuc7h-4zIK_Vw6aMpahX7FM82jUuDg&s" alt="Servicios para mascotas" class="img-fluid rounded">
        </div>
      </div>
    </div>
  </section>

  <!-- Sección de tarjetas de servicios -->
  <section class="py-5">
    <div class="container">
      <div class="row g-4">
        <!-- Servicio 1: Consulta Veterinaria -->
        <div class="col-md-6 col-lg-4">
          <div class="card service-card bg-vet h-100 shadow-sm">
            <div class="card-body text-center">
              <i class="bi bi-stethoscope icon-large text-primary mb-3" aria-hidden="true"></i>
              <h4 class="card-title">Consulta Veterinaria</h4>
              <p class="card-text">Evaluación médica completa, diagnósticos y tratamientos personalizados para la salud de tus mascotas.</p>
              <a href="agendarCita.php" class="btn btn-primary">Agendar Cita</a>
            </div>
          </div>
        </div>

        <!-- Servicio 2: Vacunación -->
        <div class="col-md-6 col-lg-4">
          <div class="card service-card bg-vet h-100 shadow-sm">
            <div class="card-body text-center">
              <i class="bi bi-shield-check icon-large text-success mb-3" aria-hidden="true"></i>
              <h4 class="card-title">Vacunación</h4>
              <p class="card-text">Programas de vacunación completos para proteger a tus mascotas en cada etapa de su vida.</p>
              <a href="agendarCita.php" class="btn btn-success">Agendar Cita</a>
            </div>
          </div>
        </div>

        <!-- Servicio 3: Cirugías y Procedimientos -->
        <div class="col-md-6 col-lg-4">
          <div class="card service-card bg-vet h-100 shadow-sm">
            <div class="card-body text-center">
              <i class="bi bi-heart-pulse icon-large text-danger mb-3" aria-hidden="true"></i>
              <h4 class="card-title">Cirugías &amp; Procedimientos</h4>
              <p class="card-text">Intervenciones quirúrgicas y tratamientos avanzados en un ambiente seguro y profesional.</p>
              <a href="agendarCita.php" class="btn btn-danger">Agendar Cita</a>
            </div>
          </div>
        </div>

        <!-- Servicio 4: Estética y Spa -->
        <div class="col-md-6 col-lg-4">
          <div class="card service-card bg-estetic h-100 shadow-sm">
            <div class="card-body text-center">
              <i class="bi bi-scissors icon-large text-warning mb-3" aria-hidden="true"></i>
              <h4 class="card-title">Estética y Spa</h4>
              <p class="card-text">Baños, cortes, tratamientos de spa y más para que tu mascota luzca radiante.</p>
              <a href="agendarCita.php" class="btn btn-warning">Agendar Cita</a>
            </div>
          </div>
        </div>

        <!-- Servicio 5: Diagnóstico por Imágenes -->
        <div class="col-md-6 col-lg-4">
          <div class="card service-card bg-special h-100 shadow-sm">
            <div class="card-body text-center">
              <i class="bi bi-camera icon-large text-info mb-3" aria-hidden="true"></i>
              <h4 class="card-title">Diagnóstico por Imágenes</h4>
              <p class="card-text">Utilizamos tecnología de punta para ofrecer diagnósticos precisos y confiables.</p>
              <a href="agendarCita.php" class="btn btn-info">Agendar Cita</a>
            </div>
          </div>
        </div>

        <!-- Servicio 6: Atención de Emergencia -->
        <div class="col-md-6 col-lg-4">
          <div class="card service-card bg-special h-100 shadow-sm">
            <div class="card-body text-center">
              <i class="bi bi-exclamation-triangle icon-large text-warning mb-3" aria-hidden="true"></i>
              <h4 class="card-title">Atención de Emergencia</h4>
              <p class="card-text">Servicio 24/7 para atender emergencias y brindar asistencia inmediata.</p>
              <a href="agendarCita.php" class="btn btn-warning">Agendar Cita</a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Sección de beneficios -->
  <section class="py-5 bg-light">
    <div class="container">
      <h2 class="text-center mb-4">¿Por qué elegirnos?</h2>
      <div class="row text-center">
        <div class="col-md-4 mb-4">
          <i class="bi bi-award icon-large text-primary mb-2" aria-hidden="true"></i>
          <h5>Calidad y Experiencia</h5>
          <p>Años de experiencia y un equipo profesional comprometido con el bienestar de tu mascota.</p>
        </div>
        <div class="col-md-4 mb-4">
          <i class="bi bi-hand-thumbs-up icon-large text-success mb-2" aria-hidden="true"></i>
          <h5>Atención Personalizada</h5>
          <p>Cada mascota recibe un trato único, adaptado a sus necesidades.</p>
        </div>
        <div class="col-md-4 mb-4">
          <i class="bi bi-chat-dots icon-large text-info mb-2" aria-hidden="true"></i>
          <h5>Soporte Continuo</h5>
          <p>Estamos siempre disponibles para acompañarte en el cuidado de tu mascota.</p>
        </div>
      </div>
    </div>
  </section>

  <!-- Incluir pie de página -->
  <?php include 'footer.php'; ?>

  <!-- Bootstrap JS (con Popper) -->
  <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
</body>
</html>
