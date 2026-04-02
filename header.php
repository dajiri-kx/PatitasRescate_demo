<?php
// Verificar si la sesión ya está activa antes de iniciarla
if (session_status() === PHP_SESSION_NONE) {
    session_start();
}

// Verificar si el cliente está logueado
$clienteLogueado = isset($_SESSION['cliente']['id_cliente']);
error_log("Cliente logueado: " . ($clienteLogueado ? 'Sí' : 'No'));
$nombre = isset($_SESSION['cliente']['nombre']) ? $_SESSION['cliente']['nombre'] : 'Usuario';
error_log("Nombre del cliente: " . $nombre);
?>

<header>
    <nav class="navbar navbar-expand-lg navbar-dark" style="background: linear-gradient(135deg, #ff8c00, #ff6b00); box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);">
        <div class="container d-flex align-items-center">
            <div>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarsExample07" aria-controls="navbarsExample07" aria-expanded="false" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <a class="navbar-brand ms-2 d-flex align-items-center" href="home.php">
                    <i class="bi bi-heart-pulse me-2"></i>
                    <span class="fw-bold">Patitas al rescate</span>
                </a>
            </div>

            <div class="col-md-3 text-end d-lg-none">
                <?php if ($clienteLogueado): ?>
                    <div class="dropdown text-end">
                        <button class="btn btn-outline-light dropdown-toggle" type="button" id="userDropdown" data-bs-toggle="dropdown" aria-expanded="false">
                            <?php echo htmlspecialchars($nombre); ?>
                        </button>
                        <ul class="dropdown-menu dropdown-menu-end" style="background-color: #ff8c00;">
                            <li><a class="dropdown-item text-white" href="dashboard.php">Ir al Dashboard</a></li>
                            <li><hr class="dropdown-divider"></li>
                            <li><a class="dropdown-item text-white" href="logout.php">Cerrar Sesión</a></li>
                        </ul>
                    </div>
                <?php else: ?>
                    <a href="login.php" type="button" class="btn btn-outline-light">Iniciar sesión</a>
                    <a href="registro.php" type="button" class="btn btn-light ms-2" style="color: #ff6b00;">Registrarse</a>
                <?php endif; ?>
            </div>

            <div class="collapse navbar-collapse" id="navbarsExample07">
                <ul class="navbar-nav me-auto mb-2 mb-lg-0">
                    <li class="nav-item">
                        <a class="nav-link active px-3 py-2 rounded" aria-current="page" href="home.php">Home</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link px-3 py-2 rounded" href="servicios.php">Servicios</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link px-3 py-2 rounded" href="ubicacion.php">Ubicación</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link px-3 py-2 rounded" href="contactenos.php">Contáctenos</a>
                    </li>
                </ul>

                <div class="col-md-3 text-end d-none d-lg-block">
                    <?php if ($clienteLogueado): ?>
                        <div class="dropdown text-end">
                            <button class="btn btn-outline-light dropdown-toggle" type="button" id="userDropdown" data-bs-toggle="dropdown" aria-expanded="false">
                                <?php echo htmlspecialchars($nombre); ?>
                            </button>
                            <ul class="dropdown-menu dropdown-menu-end" style="background-color: #ff8c00;">
                                <li><a class="dropdown-item text-white" href="dashboard.php">Ir al Dashboard</a></li>
                                <li><hr class="dropdown-divider"></li>
                                <li><a class="dropdown-item text-white" href="logout.php">Cerrar Sesión</a></li>
                            </ul>
                        </div>
                    <?php else: ?>
                        <a href="login.php" type="button" class="btn btn-outline-light">Iniciar sesión</a>
                        <a href="registro.php" type="button" class="btn btn-light ms-2" style="color: #ff6b00;">Registrarse</a>
                    <?php endif; ?>
                </div>
            </div>
        </div>
    </nav>
</header>