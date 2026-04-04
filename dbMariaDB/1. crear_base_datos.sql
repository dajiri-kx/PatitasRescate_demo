-- Crear la base de datos y el usuario de aplicación
-- Orden de ejecución: 1 → 2 → 3 → 4

CREATE DATABASE IF NOT EXISTS patitas_rescate
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_general_ci;

-- Crear usuario de aplicación
CREATE USER IF NOT EXISTS 'Progra_PAR'@'%' IDENTIFIED BY 'PrograPAR_2026';

-- Otorgar permisos sobre la base de datos
GRANT SELECT, INSERT, UPDATE, DELETE ON patitas_rescate.* TO 'Progra_PAR'@'%';
FLUSH PRIVILEGES;

USE patitas_rescate;
