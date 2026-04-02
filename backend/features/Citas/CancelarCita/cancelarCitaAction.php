<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';
require_once __DIR__ . '/../../../shared/database.php';
require_once __DIR__ . '/../shared/CitaService.php';

cors();
requireMethod('POST');
$cliente = requireAuth();

$input = json_decode(file_get_contents('php://input'), true);
$idCliente = (int) $cliente['id_cliente'];
$idCita = filter_var($input['id_cita'] ?? null, FILTER_VALIDATE_INT);

if (!$idCita) {
    jsonResponse(['ok' => false, 'error' => 'Seleccione una cita válida.'], 400);
}

$citaService = new CitaService($conn);

try {
    if ($citaService->cancelarCita($idCita, $idCliente)) {
        jsonResponse(['ok' => true, 'message' => 'Cita cancelada correctamente.']);
    } else {
        jsonResponse(['ok' => false, 'error' => 'No se pudo cancelar la cita.'], 422);
    }
} catch (PDOException $e) {
    error_log('Error al cancelar cita: ' . $e->getMessage());
    jsonResponse(['ok' => false, 'error' => 'No se pudo cancelar la cita. Intente más tarde.'], 500);
}
