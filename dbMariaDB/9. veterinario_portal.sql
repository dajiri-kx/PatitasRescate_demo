-- =====================================================
-- 9. Portal Veterinario — vincular USUARIOS con VETERINARIOS
-- =====================================================

-- Agregar columna nullable que vincula un usuario con un veterinario
ALTER TABLE USUARIOS ADD COLUMN ID_VETERINARIO INT NULL;
ALTER TABLE USUARIOS ADD CONSTRAINT fk_usuarios_veterinarios
    FOREIGN KEY (ID_VETERINARIO) REFERENCES VETERINARIOS(ID_VETERINARIO);
