<?php
include 'dbUsuariosTablas.php';
session_start();

if ($_SERVER['REQUEST_METHOD'] == 'POST') {
    $nombre = $_POST['nombre'];
    $especie = $_POST['especie']; // Nueva selección para tipo de mascota
    $raza = $_POST['raza'];
    $edad = $_POST['edad'];

    if (isset($_SESSION['cliente']['id_cliente'])) {
        $id_cliente = $_SESSION['cliente']['id_cliente'];

        try {
            $query = "INSERT INTO mascotas (nombre, especie, raza, meses, id_cliente) VALUES (:nombre, :especie, :raza, :edad, :id_cliente)";
            $stmt = $conn->prepare($query);
            $stmt->bindParam(':nombre', $nombre);
            $stmt->bindParam(':especie', $especie);
            $stmt->bindParam(':raza', $raza);
            $stmt->bindParam(':edad', $edad);
            $stmt->bindParam(':id_cliente', $id_cliente);

            if ($stmt->execute()) {
                $mensajeExito = "Mascota registrada con éxito.";
                header("Location: dashboard.php?mensaje=" . urlencode($mensajeExito));
                exit();
            } else {
                echo "Lo sentimos. Ha ocurrido un error al registrar la mascota.";
            }
        } catch (PDOException $e) {
            echo "Error: " . $e->getMessage();
        }
    }
}
?>

<!DOCTYPE html>
<html lang="es">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Agregar Mascota</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        body {
            background-color: #f8f9fa;
            font-family: 'Roboto', sans-serif;
        }
    </style>
</head>

<body>
    <?php include 'header.php'; ?>

    <div class="container py-5 px-5 my-5 shadow mx-auto rounded-3">
        <div class="form-container">

            <div class="container my-5">
                <div class="text-center mb-4">
                    <img src="https://img.freepik.com/premium-photo/close-up-australian-shepherd-dog-with-blue-pastel-background-dog-fashion-photo-generative-ai_796128-1443.jpg?w=1380" 
                    alt="Veterinaria"
                    class="img-fluid rounded shadow" style="max-height: 350px; object-fit: cover;">
                </div>
                <h1 class="text-center" style="color: #277e1c;">Agregar Mascota</h1>

                <form method="POST">
                    <div class="mb-3">
                        <label class="form-label">Nombre de la Mascota:</label>
                        <input type="text" name="nombre" class="form-control" placeholder="Ejemplo: Max" required>
                    </div>
                    <div class="mb-3">
                        <label class="form-label">Especie:</label>
                        <select name="especie" class="form-select" required>
                            <option value="Perro">Perro</option>
                            <option value="Gato">Gato</option>
                            <option value="Conejo">Conejo</option>
                            <option value="Hámster">Hamster</option>
                            <option value="Ave">Ave</option>
                            <option value="Caballo">Caballo</option>
                            <option value="Vaca">Vaca</option>
                            <option value="Oveja">Oveja</option>
                        </select>
                    </div>
                    <div class="mb-3">
                        <label class="form-label">Raza de la Mascota:</label>
                        <input type="text" name="raza" class="form-control" placeholder="Ejemplo: Labrador" required>
                    </div>
                    <div class="mb-3">
                        <label class="form-label">Edad de la Mascota en Meses:</label>
                        <input type="number" name="edad" class="form-control" min="0" placeholder="Ejemplo: 12"
                            required>
                    </div>
                    <div class="d-flex justify-content-between">
                        <button type="submit" class="btn btn-primary">Registrar Mascota</button>
                        <a href="dashboard.php" class="btn btn-secondary">Cancelar</a>
                    </div>
                </form>
            </div>
        </div>

        <?php include 'footer.php'; ?>
        <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"
            integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz"
            crossorigin="anonymous"></script>
</body>

</html>

<?php $conn = null; ?>