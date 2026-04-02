<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';
require_once __DIR__ . '/../../../shared/database.php';
require_once __DIR__ . '/../shared/AuthService.php';

cors();
requireMethod('POST');

$input = json_decode(file_get_contents('php://input'), true);
$username = trim($input['username'] ?? '');
$password = $input['password'] ?? '';

if ($username === '' || $password === '') {
    jsonResponse(['ok' => false, 'error' => 'Todos los campos son obligatorios.'], 400);
}

$auth = new AuthService($conn);

try {
    $cliente = $auth->login($username, $password);

    if ($cliente) {
        session_regenerate_id(true);
        $_SESSION['cliente'] = $cliente;
        $_SESSION['logged_in'] = true;
        $_SESSION['last_activity'] = time();
        jsonResponse(['ok' => true, 'data' => $cliente]);
    } else {
        jsonResponse(['ok' => false, 'error' => 'Credenciales incorrectas.'], 401);
    }
} catch (PDOException $e) {
    error_log('Error al iniciar sesión: ' . $e->getMessage());
    jsonResponse(['ok' => false, 'error' => 'Error al iniciar sesión. Intente más tarde.'], 500);
}
