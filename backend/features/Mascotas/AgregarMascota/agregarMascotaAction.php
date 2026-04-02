<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';
require_once __DIR__ . '/../../../shared/database.php';
require_once __DIR__ . '/../shared/MascotaService.php';

cors();
requireMethod('POST');
$cliente = requireAuth();

$input = json_decode(file_get_contents('php://input'), true);

$idCliente = (int) $cliente['id_cliente'];
$nombre = trim($input['nombre'] ?? '');
$especie = trim($input['especie'] ?? '');
$raza = trim($input['raza'] ?? '');
$edad = filter_var($input['edad'] ?? null, FILTER_VALIDATE_INT);

if ($nombre === '' || $especie === '' || $raza === '' || $edad === false || $edad === null) {
    jsonResponse(['ok' => false, 'error' => 'Todos los campos son obligatorios.'], 400);
}

$especiesValidas = ['Perro', 'Gato', 'Conejo', 'Hámster', 'Ave', 'Caballo', 'Vaca', 'Oveja'];
if (!in_array($especie, $especiesValidas, true)) {
    jsonResponse(['ok' => false, 'error' => 'Especie no válida.'], 400);
}

$mascotaService = new MascotaService($conn);

try {
    if ($mascotaService->agregar($nombre, $especie, $raza, $edad, $idCliente)) {
        jsonResponse(['ok' => true, 'message' => 'Mascota registrada con éxito.']);
    } else {
        jsonResponse(['ok' => false, 'error' => 'Error al registrar la mascota.'], 500);
    }
} catch (PDOException $e) {
    error_log('Error al agregar mascota: ' . $e->getMessage());
    jsonResponse(['ok' => false, 'error' => 'No se pudo registrar la mascota. Intente más tarde.'], 500);
}
