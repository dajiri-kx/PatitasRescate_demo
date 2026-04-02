<?php
session_start();
require_once __DIR__ . '/../../../shared/middleware.php';

cors();
requireMethod('POST');

session_unset();
session_destroy();
jsonResponse(['ok' => true, 'message' => 'Sesión cerrada.']);
