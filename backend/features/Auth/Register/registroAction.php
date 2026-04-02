<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';
require_once __DIR__ . '/../../../shared/database.php';
require_once __DIR__ . '/../shared/AuthService.php';

cors();
requireMethod('POST');

$input = json_decode(file_get_contents('php://input'), true);

$identificacion = trim($input['identificacion'] ?? '');
$nombre = trim($input['nombre'] ?? '');
$apellido = trim($input['primerApellido'] ?? '');
$correo = trim($input['correo'] ?? '');
$telefono = trim($input['telefono'] ?? '');
$password = $input['password'] ?? '';
$confirmPassword = $input['confirmPassword'] ?? '';
$direccionSennas = trim($input['direccionSennas'] ?? '');

if ($password !== $confirmPassword) {
    jsonResponse(['ok' => false, 'error' => 'Las contraseñas no coinciden.'], 400);
}

if ($identificacion === '' || $nombre === '' || $apellido === '' || $correo === '' || $telefono === '' || $password === '') {
    jsonResponse(['ok' => false, 'error' => 'Todos los campos son obligatorios.'], 400);
}

if (!filter_var($correo, FILTER_VALIDATE_EMAIL)) {
    jsonResponse(['ok' => false, 'error' => 'Correo electrónico no válido.'], 400);
}

$auth = new AuthService($conn);

try {
    $auth->registrar($identificacion, $nombre, $apellido, $correo, $telefono, $direccionSennas, $password);
    jsonResponse(['ok' => true, 'message' => 'Cliente registrado exitosamente.']);
} catch (PDOException $e) {
    if (strpos($e->getMessage(), '20010') !== false) {
        jsonResponse(['ok' => false, 'error' => 'La cédula ya está registrada.'], 409);
    } elseif (strpos($e->getMessage(), '20011') !== false) {
        jsonResponse(['ok' => false, 'error' => 'El correo electrónico ya está registrado.'], 409);
    } else {
        error_log('Error al registrar cliente: ' . $e->getMessage());
        jsonResponse(['ok' => false, 'error' => 'No se pudo completar el registro. Intente más tarde.'], 500);
    }
}
