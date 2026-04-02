<!-- home.php -->
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Gracias por tu compra - Patitas al Rescate</title>

    <!-- CSS externo primero -->
    <link rel="stylesheet" href="https://2-22-4-dot-lead-pages.appspot.com/static/lp918/min/default_thank_you.css">
    <!-- Bootstrap después -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
    <!-- Font Awesome -->
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css">
    <!-- Fuentes -->
    <link href="https://fonts.googleapis.com/css?family=Lato:300,400|Montserrat:700" rel="stylesheet">

    <!-- Sobrescribir estilos para evitar colores no deseados -->
    <style>
        /* Corregir colores del header y footer */
        header, header * {
            color: #212529 !important;
        }

        footer, footer * {
            color: #212529 !important;
        }

        .nav-link,
        .navbar-brand,
        .dropdown-item {
            color: #212529 !important;
        }

        .btn-outline-primary,
        .btn-outline-primary:hover {
            color: #0d6efd !important;
            border-color: #0d6efd !important;
        }

        .link-body-emphasis {
            color: #212529 !important;
        }

        /* Forzar color azul del botón "Descargar Factura" */
        .btn.btn-primary {
            background-color: #0d6efd !important;
            border-color: #0d6efd !important;
            color: #fff !important;
        }

        /* Cambiar el botón "Registrarse" a blanco */
        a.btn.btn-secondary,
        button.btn.btn-secondary {
            background-color: #ffffff !important;
            color: #212529 !important;
            border-color: #ced4da !important;
        }

        a.btn.btn-secondary:hover,
        button.btn.btn-secondary:hover {
            background-color: #f8f9fa !important;
        }
    </style>
</head>
<body>
    <?php include 'header.php'; ?>

    <div class="container mt-5 mb-5 text-center text-dark">
        <i class="fas fa-check fa-5x text-success mb-4"></i>
        <h1 class="site-header__title">¡GRACIAS POR TU COMPRA!</h1>
        <p class="lead mt-3">
            Tu pago ha sido procesado con éxito. Apreciamos tu confianza en <strong>Patitas al Rescate</strong>.
        </p>
        <a href="descargar_factura.php" class="btn btn-primary mt-3">
            <i class="fas fa-download me-2"></i> Descargar Factura
        </a>
    </div>

    <?php include 'footer.php'; ?>

    <!-- Scripts de Bootstrap -->
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
</body>
</html>