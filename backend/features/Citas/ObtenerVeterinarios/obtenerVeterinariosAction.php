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
    $veterinarios = $citaService->obtenerVeterinarios();
    jsonResponse(['ok' => true, 'data' => $veterinarios]);
} catch (PDOException $e) {
    error_log('Error al obtener veterinarios: ' . $e->getMessage());
    jsonResponse(['ok' => false, 'error' => 'Error al obtener veterinarios.'], 500);
}
