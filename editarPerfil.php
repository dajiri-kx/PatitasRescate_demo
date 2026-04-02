<?php
session_start();

if ($_SESSION['id_cliente'] == "") {
    header("Location: login.php");
    exit();
}

include('db.php');

$id_cliente = $_SESSION['id_cliente']; //usuario registrado

//solicitar los datos existentes del usuario
$infoCliente = "SELECT * FROM clientes WHERE id_cliente = $id_cliente";
$resultadoCliente = $conn->query($infoCliente);
$datosCliente = $resultadoCliente->fetch_assoc();

//si solicita actualizar datos por medio del form
if ($_SERVER["REQUEST_METHOD"] == "POST") {
    $nombre = $_POST['nombre'];
    $correo = $_POST['correo'];
    $telefono = $_POST['telefono'];

    $updateCliente = "UPDATE clientes SET nombre = '$nombre', correo = '$correo', telefono = '$telefono' WHERE id_cliente = $id_cliente";
    
    if ($conn->query($updateCliente) === TRUE) {
        echo "Actualización exitosa";
    } else {
        $alerta = "Error al actualizar: " . $conn->error;
    }
}
?>

<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <title>Editar Perfil</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <?php include 'header.php'; ?>

    <div class="container my-5">
        <h2>Editar Perfil</h2>
    </div>
</body>
</html>
