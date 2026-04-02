<?php

class CitaService
{
    private PDO $conn;

    public function __construct(PDO $conn)
    {
        $this->conn = $conn;
    }

    public function obtenerCitasPorCliente(int $idCliente): array
    {
        $query = "
            SELECT C.ID_CITA, C.FECHA_CITA, C.ESTADO, M.NOMBRE AS MASCOTA, V.NOMBRE AS VETERINARIO
            FROM CITAS_TABLAS.CITAS C
            JOIN USUARIOS_TABLAS.MASCOTAS M ON C.ID_MASCOTA = M.ID_MASCOTA
            JOIN USUARIOS_TABLAS.VETERINARIOS V ON C.ID_VETERINARIO = V.ID_VETERINARIO
            WHERE M.ID_CLIENTE = :id_cliente
        ";
        $stmt = $this->conn->prepare($query);
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        $stmt->execute();
        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }

    public function obtenerVeterinarios(): array
    {
        $stmt = $this->conn->prepare("SELECT ID_VETERINARIO, NOMBRE FROM usuarios_tablas.veterinarios");
        $stmt->execute();
        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }

    public function obtenerServicios(): array
    {
        $stmt = $this->conn->prepare("SELECT ID_SERVICIO, NOMBRE_SERVICIO, DESCRIPCION FROM servicios_tablas.SERVICIOS");
        $stmt->execute();
        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }

    public function agendarCita(int $idCliente, int $idMascota, int $idVeterinario, string $fechaCita, string $serviciosList): void
    {
        $sql = "BEGIN agendarCita(:id_cliente, :id_mascota, :id_veterinario, TO_DATE(:fecha_cita, 'YYYY-MM-DD HH24:MI'), :servicios); END;";
        $stmt = $this->conn->prepare($sql);
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        $stmt->bindParam(':id_mascota', $idMascota, PDO::PARAM_INT);
        $stmt->bindParam(':id_veterinario', $idVeterinario, PDO::PARAM_INT);
        $stmt->bindParam(':fecha_cita', $fechaCita);
        $stmt->bindParam(':servicios', $serviciosList);
        $stmt->execute();
    }

    public function cancelarCita(int $idCita, int $idCliente): bool
    {
        $sql = "DELETE FROM CITAS_TABLAS.CITAS WHERE ID_CITA = :id_cita AND ID_MASCOTA IN (
            SELECT ID_MASCOTA FROM USUARIOS_TABLAS.MASCOTAS WHERE ID_CLIENTE = :id_cliente
        )";
        $stmt = $this->conn->prepare($sql);
        $stmt->bindParam(':id_cita', $idCita, PDO::PARAM_INT);
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        return $stmt->execute();
    }

    public function obtenerCitasActivasPorCliente(int $idCliente): array
    {
        $query = "
            SELECT C.ID_CITA, C.FECHA_CITA, M.NOMBRE AS MASCOTA
            FROM CITAS_TABLAS.CITAS C
            JOIN USUARIOS_TABLAS.MASCOTAS M ON C.ID_MASCOTA = M.ID_MASCOTA
            WHERE M.ID_CLIENTE = :id_cliente AND C.ESTADO = 'Activa'
            ORDER BY C.FECHA_CITA
        ";
        $stmt = $this->conn->prepare($query);
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        $stmt->execute();
        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }
}
