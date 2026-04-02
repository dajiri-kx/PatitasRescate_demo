<?php
include('db.php');
session_start();

// Validación básica de sesión
if (!isset($_SESSION['cliente']['id_cliente'])) {
    header("Location: login.php");
    exit();
}

$id_cliente = $_SESSION['cliente']['id_cliente'];
error_log("ID Cliente: " . $id_cliente);

$aviso = "";

// Preparar datos para las opciones del formulario
$mascotas = [];
$veterinarios = [];
$servicios = [];

try {
    $queryMascotas = "SELECT ID_MASCOTA, NOMBRE FROM usuarios_tablas.mascotas WHERE id_cliente = :id_cliente";
    $stmtMascotas = $conn->prepare($queryMascotas);
    $stmtMascotas->bindParam(':id_cliente', $id_cliente, PDO::PARAM_INT);
    $stmtMascotas->execute();
    $mascotas = $stmtMascotas->fetchAll(PDO::FETCH_ASSOC);
    error_log("Mascotas: " . json_encode($mascotas));

    $queryVeterinarios = "SELECT ID_VETERINARIO, NOMBRE FROM usuarios_tablas.veterinarios";
    $stmtVeterinarios = $conn->prepare($queryVeterinarios);
    $stmtVeterinarios->execute();
    $veterinarios = $stmtVeterinarios->fetchAll(PDO::FETCH_ASSOC);
    error_log("Veterinarios: " . json_encode($veterinarios));

    $queryServicios = "SELECT ID_SERVICIO, NOMBRE_SERVICIO, DESCRIPCION FROM servicios_tablas.SERVICIOS";
    $stmtServicios = $conn->prepare($queryServicios);
    $stmtServicios->execute();
    $servicios = $stmtServicios->fetchAll(PDO::FETCH_ASSOC);
    error_log("Servicios: " . json_encode($servicios));

} catch (PDOException $e) {
    error_log("Error al preparar datos del formulario: " . $e->getMessage());
}

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $id_mascota = $_POST['id_mascota'] ?? null;
    $fecha = $_POST['fecha'] ?? null;
    $hora = $_POST['hora'] ?? null;
    $servicio = $_POST['servicio'] ?? null;
    $veterinario = $_POST['veterinario'] ?? null;

    if ($id_mascota && $fecha && $hora && $servicio && $veterinario) {
        try {
            $fecha_cita = $fecha . ' ' . $hora; // Combinar fecha y hora

            // Log para verificar los datos enviados al procedimiento
            error_log("Datos enviados al procedimiento agendarCita:");
            error_log("ID Cliente: " . $id_cliente);
            error_log("ID Mascota: " . $id_mascota);
            error_log("ID Veterinario: " . $veterinario);
            error_log("Fecha Cita: " . $fecha_cita);
            error_log("Servicios: " . implode(',', $servicio));

            // Convertir los servicios seleccionados en una cadena separada por comas
            $servicios_list = implode(',', $servicio);

            // Llamar al procedimiento almacenado agendarCita
            $sql = "BEGIN agendarCita(:id_cliente, :id_mascota, :id_veterinario, TO_DATE(:fecha_cita, 'YYYY-MM-DD HH24:MI'), :servicios); END;";
            $stmt = $conn->prepare($sql);

            $stmt->bindParam(':id_cliente', $id_cliente, PDO::PARAM_INT);
            $stmt->bindParam(':id_mascota', $id_mascota, PDO::PARAM_INT);
            $stmt->bindParam(':id_veterinario', $veterinario, PDO::PARAM_INT);
            $stmt->bindParam(':fecha_cita', $fecha_cita);
            $stmt->bindParam(':servicios', $servicios_list);

            $stmt->execute();

            $aviso = "Cita agendada con éxito.";
        } catch (PDOException $e) {
            // Preparar mensajes de error amigables
            if (strpos($e->getMessage(), 'ORA-20001') !== false) {
                $aviso = "Error: El cliente no existe.";
            } elseif (strpos($e->getMessage(), 'ORA-20002') !== false) {
                $aviso = "Error: La mascota no existe.";
            } elseif (strpos($e->getMessage(), 'ORA-20003') !== false) {
                $aviso = "Error: La mascota no pertenece al cliente.";
            } elseif (strpos($e->getMessage(), 'ORA-20004') !== false) {
                $aviso = "Error: La mascota ya tiene una cita activa a la misma hora. La cita anterior ha sido cancelada.";
            } elseif (strpos($e->getMessage(), 'ORA-20005') !== false) {
                $aviso = "Error: El veterinario no existe.";
            } elseif (strpos($e->getMessage(), 'ORA-20006') !== false) {
                $aviso = "Error: El veterinario no está disponible en esa fecha y hora.";
            } elseif (strpos($e->getMessage(), 'ORA-20007') !== false) {
                $aviso = "Error: El servicio solicitado no es válido.";
            } elseif (strpos($e->getMessage(), 'ORA-20008') !== false) {
                $aviso = "Error: No hay suficiente stock para uno de los productos asociados al servicio.";
            } else {
                $aviso = "Error inesperado: " . $e->getMessage();
            }
        }
    } else {
        $aviso = "Todos los campos son obligatorios.";
    }
}
?>

<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>Agendar Cita</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>

<?php include 'header.php'; ?>

<div class="container my-5">
<div class="text-center mb-4">
    <img src="https://okvet.co/wp-content/uploads/2020/06/que-es-una-veterinaria.jpg" 
         alt="Veterinaria" 
         class="img-fluid rounded shadow" 
         style="max-height: 250px; object-fit: cover;">
</div>

<h1 class="text-center" style="color: #ffc107;">Agendar Cita</h1>

<?php if (!empty($aviso)) : ?>
    <div class="alert alert-info text-center">
        <?= htmlspecialchars($aviso) ?>
    </div>
<?php endif; ?>

<form action="agendarCita.php" method="POST">
    <!-- Mascotas -->
    <div class="mb-3">
        <label for="id_mascota" class="form-label">Selecciona la Mascota:</label>
        <select name="id_mascota" id="id_mascota" class="form-select" required>
            <option value="">-- Elige una mascota --</option>
            <?php foreach ($mascotas as $mascota) : ?>
                <option value="<?= htmlspecialchars($mascota['ID_MASCOTA']) ?>">
                    <?= htmlspecialchars($mascota['NOMBRE']) ?>
                </option>
            <?php endforeach; ?>
        </select>
    </div>

    <!-- Fecha -->
    <div class="mb-3">
        <label for="fecha" class="form-label">Fecha:</label>
        <input type="date" name="fecha" id="fecha" class="form-control" required>
    </div>

    <!-- Hora -->
    <div class="mb-3">
        <label for="hora" class="form-label">Hora:</label>
        <select name="hora" id="hora" class="form-select" required>
            <?php
            for ($h = 8; $h <= 17; $h++) {
                foreach (['00', '30'] as $min) {
                    $time = sprintf("%02d:%s", $h, $min);
                    echo "<option value='$time'>$time</option>";
                }
            }
            ?>
        </select>
    </div>

    <!-- Servicios -->
    <div class="mb-3">
        <label for="servicios" class="form-label">Selecciona los Servicios:</label>
        <?php foreach ($servicios as $serv) : ?>
            <div class="form-check">
                <input class="form-check-input" type="checkbox" name="servicio[]" value="<?= htmlspecialchars($serv['ID_SERVICIO']) ?>" id="servicio_<?= htmlspecialchars($serv['ID_SERVICIO']) ?>">
                <label class="form-check-label" for="servicio_<?= htmlspecialchars($serv['ID_SERVICIO']) ?>">
                    <?= htmlspecialchars($serv['NOMBRE_SERVICIO']) ?>
                </label>
            </div>
        <?php endforeach; ?>
    </div>

    <!-- Veterinarios -->
    <div class="mb-3">
        <label for="veterinario" class="form-label">Selecciona el Veterinario:</label>
        <select name="veterinario" id="veterinario" class="form-select" required>
            <?php foreach ($veterinarios as $v) : ?>
                <option value="<?= htmlspecialchars($v['ID_VETERINARIO']) ?>">
                    <?= htmlspecialchars($v['NOMBRE']) ?>
                </option>
            <?php endforeach; ?>
        </select>
    </div>

    <button type="submit" class="btn btn-primary">Agendar Cita</button>
    <a href="dashboard.php" class="btn btn-link mt-3">Volver al Dashboard</a>
</form>
</div>

<?php include 'footer.php'; ?>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" crossorigin="anonymous"></script>
<script>
    const fechaInput = document.getElementById('fecha');
    const hoy = new Date();
    fechaInput.min = hoy.toISOString().split('T')[0];

    fechaInput.addEventListener('change', () => {
        const dia = new Date(fechaInput.value).getDay();
        if (dia === 0 || dia === 6) {
            alert("Por favor selecciona un día entre lunes y viernes.");
            fechaInput.value = '';
        }
    });
</script>
</body>
</html>

<?php $conn = null; ?>
