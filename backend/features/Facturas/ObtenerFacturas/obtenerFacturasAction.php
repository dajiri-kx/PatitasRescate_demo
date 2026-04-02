<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';
require_once __DIR__ . '/../../../shared/database.php';
require_once __DIR__ . '/../shared/FacturaService.php';

cors();
requireMethod('GET');
$cliente = requireAuth();

$facturaService = new FacturaService($conn);

try {
    $facturas = $facturaService->obtenerPorCliente((int) $cliente['id_cliente']);
    jsonResponse(['ok' => true, 'data' => $facturas]);
} catch (PDOException $e) {
    error_log('Error al obtener facturas: ' . $e->getMessage());
    jsonResponse(['ok' => false, 'error' => 'Error al obtener facturas.'], 500);
}
