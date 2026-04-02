<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';
require_once __DIR__ . '/../../../shared/database.php';
require_once __DIR__ . '/../shared/CitaService.php';

cors();
requireMethod('GET');
requireAuth();

$citaService = new CitaService($conn);

try {
    $servicios = $citaService->obtenerServicios();
    jsonResponse(['ok' => true, 'data' => $servicios]);
} catch (PDOException $e) {
    error_log('Error al obtener servicios: ' . $e->getMessage());
    jsonResponse(['ok' => false, 'error' => 'Error al obtener servicios.'], 500);
}
