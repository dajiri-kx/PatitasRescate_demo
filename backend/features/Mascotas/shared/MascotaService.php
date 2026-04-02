<?php

class MascotaService
{
    private PDO $conn;

    public function __construct(PDO $conn)
    {
        $this->conn = $conn;
    }

    public function obtenerPorCliente(int $idCliente): array
    {
        $query = "
            SELECT
                m.id_mascota AS ID_MASCOTA,
                m.nombre AS NOMBRE_MASCOTA,
                m.especie AS ESPECIE,
                m.raza AS RAZA,
                m.meses AS MESES,
                c.nombre AS NOMBRE_CLIENTE,
                c.apellido AS APELLIDO_CLIENTE
            FROM usuarios_tablas.mascotas m
            JOIN usuarios_tablas.clientes c ON m.id_cliente = c.id_cliente
            WHERE m.id_cliente = :id_cliente
        ";
        $stmt = $this->conn->prepare($query);
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        $stmt->execute();
        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }

    public function obtenerNombresPorCliente(int $idCliente): array
    {
        $query = "SELECT ID_MASCOTA, NOMBRE FROM usuarios_tablas.mascotas WHERE id_cliente = :id_cliente";
        $stmt = $this->conn->prepare($query);
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        $stmt->execute();
        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }

    public function agregar(string $nombre, string $especie, string $raza, int $edad, int $idCliente): bool
    {
        $query = "INSERT INTO usuarios_tablas.mascotas (nombre, especie, raza, meses, id_cliente)
                  VALUES (:nombre, :especie, :raza, :edad, :id_cliente)";
        $stmt = $this->conn->prepare($query);
        $stmt->bindParam(':nombre', $nombre);
        $stmt->bindParam(':especie', $especie);
        $stmt->bindParam(':raza', $raza);
        $stmt->bindParam(':edad', $edad, PDO::PARAM_INT);
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        return $stmt->execute();
    }
}
