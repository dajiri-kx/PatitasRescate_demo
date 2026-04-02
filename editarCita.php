<?php
include('db.php');
session_start();

if (!isset($_SESSION['id_cliente'])) {
    header("Location: login.php");
    exit();
}

$id_cliente = $_SESSION['id_cliente'];
$id_cita = $_GET['id'] ?? null;

if (!$id_cita) {
    header("Location: dashboard.php");
    exit();
}

// Obtener información de la cita
try {
    $sql = "SELECT * FROM citas WHERE id_cita = :id_cita AND id_cliente = :id_cliente";
    $stmt = $conn->prepare($sql);
    $stmt->bindParam(':id_cita', $id_cita, PDO::PARAM_INT);
    $stmt->bindParam(':id_cliente', $id_cliente, PDO::PARAM_INT);
    $stmt->execute();
    $cita = $stmt->fetch(PDO::FETCH_ASSOC);

    if (!$cita) {
        header("Location: dashboard.php");
        exit();
    }
} catch (PDOException $e) {
    echo "Error: " . $e->getMessage();
    exit();
}

// Obtener mascotas del cliente
try {
    $sql_mascotas = "SELECT * FROM mascotas WHERE id_cliente = :id_cliente";
    $stmt_mascotas = $conn->prepare($sql_mascotas);
    $stmt_mascotas->bindParam(':id_cliente', $id_cliente, PDO::PARAM_INT);
    $stmt_mascotas->execute();
    $mascotas = $stmt_mascotas->fetchAll(PDO::FETCH_ASSOC);
} catch (PDOException $e) {
    echo "Error: " . $e->getMessage();
    exit();
}

// Procesar actualización
if ($_SERVER["REQUEST_METHOD"] == "POST") {
    $id_mascota = $_POST['id_mascota'];
    $fecha = $_POST['fecha'];
    $hora = $_POST['hora'];

    try {
        $update = "UPDATE citas SET id_mascota = :id_mascota, fecha = :fecha, hora = :hora WHERE id_cita = :id_cita AND id_cliente = :id_cliente";
        $stmt = $conn->prepare($update);
        $stmt->bindParam(':id_mascota', $id_mascota, PDO::PARAM_INT);
        $stmt->bindParam(':fecha', $fecha);
        $stmt->bindParam(':hora', $hora);
        $stmt->bindParam(':id_cita', $id_cita, PDO::PARAM_INT);
        $stmt->bindParam(':id_cliente', $id_cliente, PDO::PARAM_INT);

        if ($stmt->execute()) {
            $mensaje = "Cita actualizada con éxito.";
            header("Location: dashboard.php?success=1");
            exit();
        } else {
            $error = "Error al actualizar la cita.";
        }
    } catch (PDOException $e) {
        $error = "Error: " . $e->getMessage();
    }
}
?>

<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Editar Cita - Patitas al Rescate</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <?php include 'header.php'; ?>

    <div class="container my-5">
        <h1 class="text-center">Editar Cita</h1>

        <?php if (isset($error)): ?>
            <div class="alert alert-danger"><?php echo $error; ?></div>
        <?php endif; ?>

        <form method="POST">
            <div class="mb-3">
                <label for="id_mascota" class="form-label">Mascota:</label>
                <select name="id_mascota" id="id_mascota" class="form-select" required>
                    <?php foreach ($mascotas as $mascota): ?>
                        <option value="<?php echo $mascota['id_mascota']; ?>" 
                            <?php echo ($mascota['id_mascota'] == $cita['id_mascota']) ? 'selected' : ''; ?>>
                            <?php echo $mascota['nombre']; ?>
                        </option>
                    <?php endforeach; ?>
                </select>
            </div>

            <div class="mb-3">
                <label for="fecha" class="form-label">Fecha:</label>
                <input type="date" name="fecha" id="fecha" class="form-control" 
                       value="<?php echo $cita['fecha']; ?>" required>
            </div>

            <div class="mb-3">
                <label for="hora" class="form-label">Hora:</label>
                <input type="time" name="hora" id="hora" class="form-control" 
                       value="<?php echo $cita['hora']; ?>" required min="08:00" max="18:00">
            </div>

            <div class="d-grid gap-2 d-md-flex justify-content-md-end">
                <a href="dashboard.php" class="btn btn-secondary me-md-2">Cancelar</a>
                <button type="submit" class="btn btn-primary">Guardar Cambios</button>
            </div>
        </form>
    </div>

    <?php include 'footer.php'; ?>
</body>
</html>

<?php $conn = null; ?>