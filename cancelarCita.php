<?php
include('db.php');
session_start();
$id_cliente = $_SESSION['id_cliente'];

if ($_SERVER['REQUEST_METHOD'] == "POST") {
    $id_cita = $_POST['id_cita'];

    try {
        $cancelar = "DELETE FROM citas WHERE id_cita = :id_cita AND id_cliente = :id_cliente";
        $stmt = $conn->prepare($cancelar);
        $stmt->bindParam(':id_cita', $id_cita, PDO::PARAM_INT);
        $stmt->bindParam(':id_cliente', $id_cliente, PDO::PARAM_INT);

        if ($stmt->execute()) {
            $aviso = "Cita cancelada correctamente.";
        } else {
            $aviso = "Error. No se pudo cancelar la cita.";
        }
    } catch (PDOException $e) {
        $aviso = "Error: " . $e->getMessage();
    }
}
?>

<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cancelar Cita</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>

    <?php include 'header.php'; ?>

    <div class="container my-5">
        <h1 class="text-center">Cancelar Cita</h1>

        <?php if (isset($aviso)) : ?>
            <div class="alert alert-info text-center">
                <?= $aviso ?>
            </div>
        <?php endif; ?>

        <form action="cancelarCita.php" method="POST">
            <div class="mb-3">
                <label for="id_cita" class="form-label">Selecciona la cita a cancelar:</label>
                <select name="id_cita" id="id_cita" class="form-select" required>
                    <?php
                    try {
                        $citas = "SELECT c.id_cita, c.fecha, c.hora, m.nombre AS mascota
                                  FROM citas c
                                  JOIN mascotas m ON c.id_mascota = m.id_mascota
                                  WHERE c.id_cliente = :id_cliente
                                  ORDER BY c.fecha, c.hora";
                        $stmt = $conn->prepare($citas);
                        $stmt->bindParam(':id_cliente', $id_cliente, PDO::PARAM_INT);
                        $stmt->execute();
                        $resultadoCitas = $stmt->fetchAll(PDO::FETCH_ASSOC);

                        if (empty($resultadoCitas)) {
                            echo "<option value=''>No tienes citas activas</option>";
                        } else {
                            foreach ($resultadoCitas as $row) {
                                $texto = $row['fecha'] . " a las " . $row['hora'] . " con " . $row['mascota'];
                                echo "<option value='".$row['id_cita']."'>".$texto."</option>";
                            }
                        }
                    } catch (PDOException $e) {
                        echo "<option value=''>Error al cargar citas</option>";
                    }
                    ?>
                </select>
            </div>

            <button type="submit" class="btn btn-danger">Cancelar Cita</button>
        </form>

        <a href="dashboard.php" class="btn btn-link mt-3">Volver al Dashboard</a>
    </div>

    <?php include 'footer.php'; ?>
</body>
</html>

<?php $conn = null; ?>
