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
$idMascota = filter_var($input['id_mascota'] ?? null, FILTER_VALIDATE_INT);
$fecha = trim($input['fecha'] ?? '');
$hora = trim($input['hora'] ?? '');
$servicio = $input['servicio'] ?? [];
$veterinario = filter_var($input['veterinario'] ?? null, FILTER_VALIDATE_INT);

if (!$idMascota || !$fecha || !$hora || empty($servicio) || !$veterinario) {
    jsonResponse(['ok' => false, 'error' => 'Todos los campos son obligatorios.'], 400);
}

if (!preg_match('/^\d{4}-\d{2}-\d{2}$/', $fecha)) {
    jsonResponse(['ok' => false, 'error' => 'Formato de fecha inválido.'], 400);
}

if (!preg_match('/^\d{2}:\d{2}$/', $hora)) {
    jsonResponse(['ok' => false, 'error' => 'Formato de hora inválido.'], 400);
}

$serviciosLimpios = array_filter($servicio, 'is_numeric');
if (empty($serviciosLimpios)) {
    jsonResponse(['ok' => false, 'error' => 'Seleccione al menos un servicio válido.'], 400);
}

$fechaCita = $fecha . ' ' . $hora;
$serviciosList = implode(',', $serviciosLimpios);

$citaService = new CitaService($conn);

try {
    $citaService->agendarCita($idCliente, $idMascota, $veterinario, $fechaCita, $serviciosList);
    jsonResponse(['ok' => true, 'message' => 'Cita agendada con éxito.']);
} catch (PDOException $e) {
    $msg = $e->getMessage();
    $errores = [
        'ORA-20001' => 'El cliente no existe.',
        'ORA-20002' => 'La mascota no existe.',
        'ORA-20003' => 'La mascota no pertenece al cliente.',
        'ORA-20004' => 'La mascota ya tiene una cita activa a la misma hora.',
        'ORA-20005' => 'El veterinario no existe.',
        'ORA-20006' => 'El veterinario no está disponible en esa fecha y hora.',
        'ORA-20007' => 'El servicio solicitado no es válido.',
        'ORA-20008' => 'No hay suficiente stock para uno de los productos asociados al servicio.',
    ];

    $avisoMsg = 'Error inesperado. Intente más tarde.';
    foreach ($errores as $code => $texto) {
        if (strpos($msg, $code) !== false) {
            $avisoMsg = $texto;
            break;
        }
    }
    error_log('Error al agendar cita: ' . $msg);
    jsonResponse(['ok' => false, 'error' => $avisoMsg], 422);
}
