<!-- contactenos.php -->
<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Patitas al Rescate - Contáctenos</title>
  <meta name="description" content="Ponte en contacto con nosotros en Patitas al Rescate para resolver tus dudas y programar citas.">
  <meta name="keywords" content="contacto, Patitas al Rescate, veterinaria, mascotas, agendar cita">
  <link rel="icon" href="favicon.ico">
  
  <!-- Preconexión a recursos externos -->
  <link rel="preconnect" href="https://cdn.jsdelivr.net">
  <link rel="preconnect" href="https://cdnjs.cloudflare.com">
  
  <!-- Bootstrap CSS -->
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
  <!-- Bootstrap Icons -->
  <link href="https://cdnjs.cloudflare.com/ajax/libs/bootstrap-icons/1.10.2/font/bootstrap-icons.min.css" rel="stylesheet">
  
  <style>
.contact-banner {
  background: linear-gradient(rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5)), 
              url('https://img.freepik.com/free-photo/adorable-white-bulldog-puppy-portrait-social-banner_53876-160763.jpg?t=st=1745534716~exp=1745538316~hmac=12be47b2d52f9db77437d59b39b61679c1110dfca142332cfeaae6887ba7ae01&w=1380') no-repeat center center;
  background-size: cover;
  height: 400px; /* Altura ajustada para mejor visualización */
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.7);
  position: relative;
  overflow: hidden;
}

    .contact-banner h1 {
      font-size: 3rem;
      margin: 0;
    }
    .contact-banner p {
      font-size: 1.25rem;
    }
    /* Estilo de la sección de contacto */
    .contact-section {
      background: linear-gradient(to bottom right, #f0f9ff, #e0f7fa);
      padding: 60px 0;
    }
    .contact-form, .contact-info {
      background: rgba(255, 255, 255, 0.98);
      padding: 30px;
      border-radius: 15px;
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
    }
    .contact-info i {
      font-size: 1.5rem;
    }
    .social-icons a {
      font-size: 2rem;
      margin-right: 15px;
      color: inherit;
      text-decoration: none;
    }
    .social-icons a:hover {
      color: #007bff;
    }
  </style>
</head>
<body>
  <?php include 'header.php'; ?>

  <!-- Banner de Contacto -->
  <section class="contact-banner text-center">
    <div>
      <h1>Contáctenos</h1>
      <p class="lead">¡Estamos aquí para ayudarte!</p>
    </div>
  </section>

  <!-- Sección de Contacto -->
  <section class="contact-section">
    <div class="container">
      <div class="row g-4">
        <!-- Formulario de Contacto -->
        <div class="col-lg-7">
          <div class="contact-form">
            <h2 class="mb-4">Envíanos un mensaje</h2>
            <form action="procesar_contacto.php" method="post">
              <div class="mb-3">
                <label for="nombre" class="form-label">Nombre Completo</label>
                <input type="text" class="form-control" id="nombre" name="nombre" placeholder="Tu nombre" required>
              </div>
              <div class="mb-3">
                <label for="email" class="form-label">Correo Electrónico</label>
                <input type="email" class="form-control" id="email" name="email" placeholder="tuemail@ejemplo.com" required>
              </div>
              <div class="mb-3">
                <label for="telefono" class="form-label">Teléfono</label>
                <input type="tel" class="form-control" id="telefono" name="telefono" placeholder="+506 0000-0000">
              </div>
              <div class="mb-3">
                <label for="mensaje" class="form-label">Mensaje</label>
                <textarea class="form-control" id="mensaje" name="mensaje" rows="5" placeholder="Escribe aquí tu mensaje" required></textarea>
              </div>
              <div class="d-grid">
                <button type="submit" class="btn btn-primary btn-lg">Enviar Mensaje</button>
              </div>
            </form>
          </div>
        </div>
        <!-- Información de Contacto -->
        <div class="col-lg-5">
          <div class="contact-info">
            <h2 class="mb-4">Información de Contacto</h2>
            <p><i class="bi bi-telephone-fill text-primary me-2"></i><strong>Teléfono:</strong> +506 1234-5678</p>
            <p><i class="bi bi-envelope-fill text-danger me-2"></i><strong>Email:</strong> info@patitasalrescate.com</p>
            <p><i class="bi bi-geo-alt-fill text-success me-2"></i><strong>Dirección:</strong> San José, Costa Rica - 100 metros norte del Parque Central</p>
            <p><i class="bi bi-clock-fill text-warning me-2"></i><strong>Horario:</strong> Lunes a Viernes: 8:00 AM - 6:00 PM</p>
            <hr>
            <h5 class="mb-3">Síguenos en Redes Sociales</h5>
            <div class="social-icons">
              <a href="#" aria-label="Facebook"><i class="bi bi-facebook"></i></a>
              <a href="#" aria-label="Instagram"><i class="bi bi-instagram"></i></a>
              <a href="#" aria-label="Twitter"><i class="bi bi-twitter"></i></a>
              <a href="#" aria-label="YouTube"><i class="bi bi-youtube"></i></a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <?php include 'footer.php'; ?>

  <!-- Bootstrap JS (con Popper) -->
  <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-NlNdpGZ1EnP9B2DT4lLGXSwET8fQLz+mDmI5FOFAe7I3rx+BTHDpldm9RN2IO6JR" crossorigin="anonymous"></script>
</body>
</html>
