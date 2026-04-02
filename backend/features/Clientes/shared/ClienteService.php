<?php

class ClienteService
{
    private PDO $conn;

    public function __construct(PDO $conn)
    {
        $this->conn = $conn;
    }

    public function obtenerPerfil(int $idCliente): ?array
    {
        $stmt = $this->conn->prepare(
            "SELECT ID_CLIENTE, NOMBRE, APELLIDO, CORREO, TELEFONO
             FROM USUARIOS_TABLAS.CLIENTES
             WHERE ID_CLIENTE = :id_cliente"
        );
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        $stmt->execute();
        $result = $stmt->fetch(PDO::FETCH_ASSOC);
        return $result ?: null;
    }

    public function actualizar(int $idCliente, string $nombre, string $correo, string $telefono): bool
    {
        $stmt = $this->conn->prepare(
            "UPDATE USUARIOS_TABLAS.CLIENTES
             SET NOMBRE = :nombre, CORREO = :correo, TELEFONO = :telefono
             WHERE ID_CLIENTE = :id_cliente"
        );
        $stmt->bindParam(':nombre', $nombre);
        $stmt->bindParam(':correo', $correo);
        $stmt->bindParam(':telefono', $telefono);
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        return $stmt->execute();
    }
}
