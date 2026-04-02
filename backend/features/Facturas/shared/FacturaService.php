<?php

class FacturaService
{
    private PDO $conn;

    public function __construct(PDO $conn)
    {
        $this->conn = $conn;
    }

    public function obtenerPorCliente(int $idCliente): array
    {
        $query = "
            SELECT f.ID_FACTURA, f.FECHA_FACTURA, f.TOTAL,
                   CASE
                       WHEN c.FECHA_CITA < SYSDATE THEN 'Pagada'
                       ELSE 'Pendiente'
                   END AS ESTADO
            FROM CITAS_TABLAS.FACTURAS f
            JOIN CITAS_TABLAS.CITAS_SERVICIOS cs ON f.ID_FACTURA = cs.FACTURAS_ID_FACTURA
            JOIN CITAS_TABLAS.CITAS c ON cs.ID_CITA = c.ID_CITA
            WHERE c.ID_MASCOTA IN (
                SELECT m.ID_MASCOTA
                FROM USUARIOS_TABLAS.MASCOTAS m
                WHERE m.ID_CLIENTE = :id_cliente
            )
        ";
        $stmt = $this->conn->prepare($query);
        $stmt->bindParam(':id_cliente', $idCliente, PDO::PARAM_INT);
        $stmt->execute();
        return $stmt->fetchAll(PDO::FETCH_ASSOC);
    }
}
