/*
auth/service.go — Lógica de negocio para autenticación (login y registro).

FLUJO DE DATOS — LOGIN:
 1. handlers.go loginHandler() recibe {username, password} del frontend.
 2. Llama a Login(ctx, correo, password).
 3. Login hace un JOIN USUARIOS + CLIENTES para obtener los datos del usuario
    por correo. También obtiene ID_VETERINARIO (0 si no es veterinario).
 4. bcrypt.CompareHashAndPassword verifica la contraseña contra el hash almacenado.
 5. Si es correcto, retorna *ClienteData con todos los campos del usuario.
 6. Si el correo no existe o la contraseña es incorrecta, retorna nil (sin error)
    → el handler responde 401 "Credenciales incorrectas".

FLUJO DE DATOS — REGISTRO:
 1. handlers.go registerHandler() recibe los datos del formulario.
 2. Llama a Registrar(ctx, identificacion, nombre, ..., password).
 3. Registrar usa una TRANSACCIÓN para garantizar consistencia:
    a. bcrypt.GenerateFromPassword genera el hash de la contraseña.
    b. Verifica que la cédula no exista ya en CLIENTES.
    c. Verifica que el correo no exista ya en USUARIOS.
    d. INSERT INTO CLIENTES → obtiene ID_CLIENTE auto-generado.
    e. INSERT INTO USUARIOS con el ID_CLIENTE obtenido.
    f. COMMIT si todo fue exitoso, o ROLLBACK automático en error.
 4. Retorna el ID_CLIENTE del nuevo registro.

TABLA DE RELACIONES:
CLIENTES (1) ←→ (1) USUARIOS — Todo cliente tiene exactamente un usuario.
USUARIOS.ROL default = 1 (Cliente). Un admin puede cambiar el ROL manualmente.
USUARIOS.ID_VETERINARIO es opcional — vincula con VETERINARIOS para rol=2.
*/
package auth

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// ClienteData es el struct que retorna Login() y se envía al frontend como JSON.
// También se usa para construir la ClienteSession que se guarda en la cookie.
type ClienteData struct {
	IDCliente     int64  `json:"id_cliente"`
	IDVeterinario int64  `json:"id_veterinario"`
	Nombre        string `json:"nombre"`
	Apellido      string `json:"apellido"`
	Correo        string `json:"correo"`
	Telefono      string `json:"telefono"`
	Rol           int    `json:"rol"`
}

// Login busca el usuario por correo y verifica la contraseña con bcrypt.
// Retorna nil, nil si las credenciales son incorrectas (no es un "error" del sistema).
// Retorna nil, error si hay un problema de base de datos.
func (s *AuthService) Login(ctx context.Context, correo, password string) (*ClienteData, error) {
	// JOIN USUARIOS + CLIENTES: USUARIOS tiene el correo/contraseña/rol,
	// CLIENTES tiene los datos personales (nombre, apellido, teléfono).
	// IFNULL(u.ID_VETERINARIO, 0): si no es veterinario, retorna 0.
	row := s.db.QueryRowContext(ctx,
		`SELECT c.ID_CLIENTE, c.NOMBRE, c.APELLIDO, c.TELEFONO, u.CONTRASENA, u.ROL, IFNULL(u.ID_VETERINARIO, 0)
		 FROM USUARIOS u
		 JOIN CLIENTES c ON u.ID_CLIENTE = c.ID_CLIENTE
		 WHERE u.CORREO = ?`,
		correo,
	)

	var idCliente, idVet int64
	var nombre, apellido, telefono, hash string
	var rol int
	if err := row.Scan(&idCliente, &nombre, &apellido, &telefono, &hash, &rol, &idVet); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Correo no registrado — no es un error de sistema
		}
		return nil, err
	}

	// bcrypt compara el password en texto plano con el hash almacenado.
	// Si no coincide, retorna error → lo tratamos como credenciales incorrectas.
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return nil, nil
	}

	return &ClienteData{
		IDCliente:     idCliente,
		IDVeterinario: idVet,
		Nombre:        nombre,
		Apellido:      apellido,
		Correo:        correo,
		Telefono:      telefono,
		Rol:           rol,
	}, nil
}

// Registrar crea un nuevo cliente + usuario en una transacción atómica.
// Si algo falla, el defer tx.Rollback() deshace todo automáticamente.
// Retorna el ID_CLIENTE generado por el auto-increment.
func (s *AuthService) Registrar(ctx context.Context, identificacion, nombre, apellido, correo, telefono, direccion, password string) (int64, error) {
	// Generar hash bcrypt antes de abrir la transacción (operación costosa en CPU)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // No-op si ya se hizo Commit

	// Paso 1: Verificar que la cédula no exista (evita duplicados en CLIENTES)
	var count int
	err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM CLIENTES WHERE DIDENTIDAD_CLIENTE = ?`, identificacion).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("ORA-20010: La cédula o documento de identidad ya está registrada")
	}

	// Paso 2: Verificar que el correo no exista (evita duplicados en USUARIOS)
	err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM USUARIOS WHERE CORREO = ?`, correo).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("ORA-20010: El correo electrónico ya está registrado")
	}

	// Paso 3: Insertar en CLIENTES (datos personales)
	// LastInsertId() obtiene el ID auto-generado por MariaDB
	result, err := tx.ExecContext(ctx,
		`INSERT INTO CLIENTES (DIDENTIDAD_CLIENTE, NOMBRE, APELLIDO, EMAIL, TELEFONO, DIRECCION, FECHA_REGISTRO)
		 VALUES (?, ?, ?, ?, ?, ?, NOW())`,
		identificacion, nombre, apellido, correo, telefono, direccion,
	)
	if err != nil {
		return 0, err
	}

	idCliente, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Paso 4: Insertar en USUARIOS (credenciales) vinculado al cliente
	// ROL no se especifica → usa el DEFAULT de la columna (1 = Cliente)
	_, err = tx.ExecContext(ctx,
		`INSERT INTO USUARIOS (ID_CLIENTE, CORREO, CONTRASENA) VALUES (?, ?, ?)`,
		idCliente, correo, string(hash),
	)
	if err != nil {
		return 0, err
	}

	// Paso 5: Commit — si llegamos aquí, todo fue exitoso
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return idCliente, nil
}
