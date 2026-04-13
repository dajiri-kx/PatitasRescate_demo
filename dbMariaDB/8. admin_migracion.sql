-- Migración: Modo Administrador
-- Ejecutar después de 7. servicios_categorias.sql
USE patitas_rescate;

-- Roles: 0 = Admin, 1 = Cliente, 2 = Veterinario
ALTER TABLE USUARIOS
    ADD COLUMN ROL TINYINT NOT NULL DEFAULT 1;

-- Para promover un usuario registrado a administrador:
--   UPDATE USUARIOS SET ROL = 0 WHERE CORREO = 'correo_del_usuario@ejemplo.com';
