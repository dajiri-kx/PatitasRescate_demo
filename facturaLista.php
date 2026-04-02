<!-- facturaLista.php -->
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Lista de Facturas - Patitas al Rescate</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <?php include 'header.php'; ?>

    <div class="container mt-5">
        <h1 class="text-center mb-4">Lista de Facturas</h1>
        <div class="table-responsive">
            <table class="table table-dark table-striped table-bordered text-center">
                <thead class="table-primary">
                    <tr>
                        <th>#</th>
                        <th>Fecha</th>
                        <th>Servicio</th>
                        <th>Veterinario</th>
                        <th>Precio (₡)</th>
                        <th>Acciones</th> <!-- no hace nada --> 
                    </tr>
                </thead>
                <tbody>
                    <?php
                    // placeholders
                    $facturas = [
                        ['fecha' => '2025-03-10', 'servicio' => 'Consulta Veterinaria', 'veterinario' => 'Dra. María Rodríguez', 'precio' => 25000],
                        ['fecha' => '2025-03-12', 'servicio' => 'Baño y Corte de Pelo', 'veterinario' => 'Dr. Juan Pérez', 'precio' => 15000],
                        ['fecha' => '2025-03-15', 'servicio' => 'Vacunación', 'veterinario' => 'Dra. Laura Gómez', 'precio' => 10000],
                    ];
                    
                    foreach ($facturas as $index => $factura) {
                        echo "<tr>";
                        echo "<td>" . ($index + 1) . "</td>";
                        echo "<td>{$factura['fecha']}</td>";
                        echo "<td>{$factura['servicio']}</td>";
                        echo "<td>{$factura['veterinario']}</td>";
                        echo "<td>₡" . number_format($factura['precio'], 0, ',', '.') . "</td>";
                        echo "<td>">
                             "<a href='factura.php?id=" . ($index + 1) . "' class='btn btn-primary btn-sm'>Ver Detalle</a> "
                            . "<button class='btn btn-danger btn-sm ms-2'>Eliminar</button>"
                        . "</td>";
                        echo "</tr>";
                    }
                    ?>
                </tbody>
            </table>
        </div>
    </div>

    <?php include 'footer.php'; ?>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz" crossorigin="anonymous"></script>
</body>
</html>