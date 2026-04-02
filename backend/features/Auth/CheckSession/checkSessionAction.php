<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';

cors();
requireMethod('GET');

if (isset($_SESSION['cliente']['id_cliente'])) {
    jsonResponse(['ok' => true, 'data' => $_SESSION['cliente']]);
} else {
    jsonResponse(['ok' => false, 'error' => 'No autenticado.'], 401);
}
