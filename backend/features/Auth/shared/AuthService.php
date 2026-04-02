<?php

class AuthService
{
    private PDO $conn;

    public function __construct(PDO $conn)
    {
        $this->conn = $conn;
    }

    public function login(string $correo, string $password): ?array
    {
        $stmt = $this->conn->prepare(
            "SELECT u.ID_USUARIO, c.ID_CLIENTE, c.NOMBRE, c.APELLIDO, c.TELEFONO, u.CONTRASENA
             FROM USUARIOS_TABLAS.USUARIOS u
             JOIN USUARIOS_TABLAS.CLIENTES c ON u.ID_CLIENTE = c.ID_CLIENTE
             WHERE u.CORREO = :correo"
        );
        $stmt->bindParam(':correo', $correo);
        $stmt->execute();
        $user = $stmt->fetch(PDO::FETCH_ASSOC);

        if ($user && is_array($user) && password_verify($password, $user['CONTRASENA'])) {
            return [
                'id_cliente' => $user['ID_CLIENTE'],
                'nombre'     => $user['NOMBRE'],
                'apellido'   => $user['APELLIDO'],
                'correo'     => $correo,
                'telefono'   => $user['TELEFONO'],
            ];
        }

        return null;
    }

    public function registrar(
        string $identificacion,
        string $nombre,
        string $apellido,
        string $correo,
        string $telefono,
        string $direccion,
        string $password
    ): int {
        $stmt = $this->conn->prepare(
            "BEGIN registrarCliente(:identificacion, :nombre, :apellido, :correo, :telefono, :direccion, :password, :idCliente); END;"
        );
        $stmt->bindParam(':identificacion', $identificacion);
        $stmt->bindParam(':nombre', $nombre);
        $stmt->bindParam(':apellido', $apellido);
        $stmt->bindParam(':correo', $correo);
        $stmt->bindParam(':telefono', $telefono);
        $stmt->bindParam(':direccion', $direccion);
        $hash = password_hash($password, PASSWORD_DEFAULT);
        $stmt->bindParam(':password', $hash);

        $idCliente = null;
        $stmt->bindParam(':idCliente', $idCliente, PDO::PARAM_INT | PDO::PARAM_INPUT_OUTPUT, 32);

        $stmt->execute();
        return (int) $idCliente;
    }
}
