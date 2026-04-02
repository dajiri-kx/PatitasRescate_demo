<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';

cors();
requireMethod('POST');

$input = json_decode(file_get_contents('php://input'), true);

$nombre = trim($input['nombre'] ?? '');
$email = trim($input['email'] ?? '');
$telefono = trim($input['telefono'] ?? '');
$mensaje = trim($input['mensaje'] ?? '');

if ($nombre === '' || $email === '' || $mensaje === '') {
    jsonResponse(['ok' => false, 'error' => 'Complete los campos obligatorios.'], 400);
}

if (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
    jsonResponse(['ok' => false, 'error' => 'Correo electrónico no válido.'], 400);
}

error_log("Mensaje de contacto - Nombre: $nombre, Email: $email, Tel: $telefono, Msg: $mensaje");
jsonResponse(['ok' => true, 'message' => '¡Mensaje enviado exitosamente!']);
