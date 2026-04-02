<?php
$host = getenv('DB_HOST') ?: 'localhost';
$port = getenv('DB_PORT') ?: '1521';
$service_name = getenv('DB_SERVICE') ?: 'xe';
$db_user = getenv('DB_USER') ?: 'Progra_PAR';
$db_password = getenv('DB_PASS') ?: 'Progra_PAR';

try {
    $dsn = "oci:dbname=(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=$host)(PORT=$port))(CONNECT_DATA=(SERVICE_NAME=$service_name)))";
    $conn = new PDO($dsn, $db_user, $db_password);
    $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
} catch (PDOException $e) {
    error_log("Error de conexión: " . $e->getMessage());
    http_response_code(500);
    echo json_encode(['ok' => false, 'error' => 'Error de conexión a la base de datos.']);
    exit();
}
