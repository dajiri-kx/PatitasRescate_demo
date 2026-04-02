<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';
require_once __DIR__ . '/../../../shared/database.php';
require_once __DIR__ . '/../shared/MascotaService.php';

cors();
requireMethod('GET');
$cliente = requireAuth();

$mascotaService = new MascotaService($conn);

try {
    $mascotas = $mascotaService->obtenerPorCliente((int) $cliente['id_cliente']);
    jsonResponse(['ok' => true, 'data' => $mascotas]);
} catch (PDOException $e) {
    error_log('Error al obtener mascotas: ' . $e->getMessage());
    jsonResponse(['ok' => false, 'error' => 'Error al obtener mascotas.'], 500);
}
