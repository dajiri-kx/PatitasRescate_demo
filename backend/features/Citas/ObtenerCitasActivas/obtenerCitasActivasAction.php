<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';
require_once __DIR__ . '/../../../shared/database.php';
require_once __DIR__ . '/../shared/CitaService.php';

cors();
requireMethod('GET');
$cliente = requireAuth();

$citaService = new CitaService($conn);

try {
    $citas = $citaService->obtenerCitasActivasPorCliente((int) $cliente['id_cliente']);
    jsonResponse(['ok' => true, 'data' => $citas]);
} catch (PDOException $e) {
    error_log('Error al obtener citas activas: ' . $e->getMessage());
    jsonResponse(['ok' => false, 'error' => 'Error al obtener citas activas.'], 500);
}
