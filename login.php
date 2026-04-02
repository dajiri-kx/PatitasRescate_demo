<?php
session_start();
require 'db.php';

$error = '';
$success = '';

// Mostrar mensajes de éxito/error si existen
if (isset($_SESSION['message'])) {
    if (strpos($_SESSION['message'], 'Error') === 0) {
        $error = $_SESSION['message'];
    } else {
        $success = $_SESSION['message'];
    }
    unset($_SESSION['message']);
}

// Procesar el formulario de login
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    // Debug: Mostrar datos enviados desde el formulario en la consola
    error_log('Datos recibidos del formulario de login:');
    error_log(print_r($_POST, true));

    $username = trim($_POST['username']);
    $password = trim($_POST['password']);

    // Debug: Mostrar variables asignadas en la consola
    error_log('Variables asignadas:');
    error_log("Usuario: $username");
    error_log("Contraseña: $password");

    // Debug: Mostrar el correo ingresado después de aplicar trim y cualquier transformación
    error_log('Correo ingresado (después de trim): ' . $username);

    try {
        // Consulta para verificar las credenciales con JOIN
        $stmt = $conn->prepare(
            "SELECT u.ID_USUARIO, c.ID_CLIENTE, c.NOMBRE, c.APELLIDO, c.TELEFONO, u.CONTRASENA 
             FROM USUARIOS_TABLAS.USUARIOS u
             JOIN USUARIOS_TABLAS.CLIENTES c ON u.ID_CLIENTE = c.ID_CLIENTE
             WHERE u.CORREO = :username"
        );
        $stmt->bindParam(':username', $username);
        $stmt->execute();

        $user = $stmt->fetch(PDO::FETCH_ASSOC); // Obtener el registro directamente después de ejecutar la consulta

        error_log("Usuario encontrado: " . ($user ? 'Sí' : 'No'));
        error_log("ID_CLIENTE: " . ($user ? $user['ID_CLIENTE'] : 'No disponible'));

        if ($user) { // Si fetch devuelve un resultado
            // Depuración adicional: Mostrar el valor de $user antes de acceder a sus índices
            error_log('Valor de $user antes de acceder a índices: ' . print_r($user, true));

            // Verificar si $user es un array antes de acceder a sus índices
            if (is_array($user)) {
                // Debug: Mostrar hash almacenado y contraseña ingresada en el error_log
                error_log('Hash almacenado: ' . $user['CONTRASENA']);
                error_log('Contraseña ingresada: ' . $password);

                // Verificar la contraseña
                if (password_verify($password, $user['CONTRASENA'])) {
                    // Regenerar el ID de sesión para prevenir fijación de sesión
                    session_regenerate_id(true);

                    // Almacenar los datos del cliente en la sesión
                    $_SESSION['cliente'] = [
                        'id_cliente' => $user['ID_CLIENTE'],
                        'nombre' => $user['NOMBRE'],
                        'apellido' => $user['APELLIDO'],
                        'correo' => $username,
                        'telefono' => $user['TELEFONO']
                    ];

                    // Iniciar sesión
                    $_SESSION['logged_in'] = true;
                    $_SESSION['last_activity'] = time();

                    // Redirigir al dashboard o página principal
                    header("Location: dashboard.php");
                    exit();
                } else {
                    // Mensaje genérico para no revelar información
                    $error = "Credenciales incorrectas. Por favor intente nuevamente.";
                }
            } else {
                error_log('Error: $user no es un array. Valor actual: ' . print_r($user, true));
                $error = "Error interno. Por favor intente más tarde.";
            }
        } else {
            // Mismo mensaje aunque el usuario no exista (seguridad)
            $error = "Credenciales incorrectas. Por favor intente nuevamente.";
        }
    } catch (PDOException $e) {
        error_log("Error al iniciar sesión: " . $e->getMessage());
        $error = "Error al iniciar sesión. Por favor intente más tarde.";
    }
}


?>

<!DOCTYPE html>
<html lang="es">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Iniciar Sesión</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>

<body>
    <?php include 'header.php'; ?>
    <div class="container mt-5">
        <div class="row justify-content-center">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header text-center">
                        <h4>Iniciar Sesión</h4>
                    </div>
                    <div class="card-body">
                        <?php if (!empty($error)): ?>
                            <div class="alert alert-danger"> <?php echo $error; ?> </div>
                        <?php endif; ?>

                        <?php if (!empty($success)): ?>
                            <div class="alert alert-success"> <?php echo $success; ?> </div>
                        <?php endif; ?>

                        <form action="login.php" method="post">
                            <div class="mb-3">
                                <label for="username" class="form-label">Correo Electrónico</label>
                                <input type="email" class="form-control" id="username" name="username" required>
                            </div>
                            <div class="mb-3">
                                <label for="password" class="form-label">Contraseña</label>
                                <input type="password" class="form-control" id="password" name="password" required>
                            </div>
                            <div class="d-grid">
                                <button type="submit" class="btn btn-primary">Iniciar Sesión</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
</body>

</html>