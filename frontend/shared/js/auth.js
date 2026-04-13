/*
   auth.js — Estado de autenticación en el lado del cliente (localStorage).

   ESTADO DUAL DE AUTENTICACIÓN:
   El backend mantiene la sesión real con una cookie HttpOnly (segura, no accesible
   desde JS). Pero el frontend necesita saber el nombre y rol del usuario para
   renderizar la navbar, mostrar badges (Admin/Vet), y decidir qué links mostrar.

   SOLUCIÓN: Después del login exitoso, el backend retorna los datos del usuario
   en el JSON response, y el frontend los guarda en localStorage bajo 'patitas_user'.

   IMPORTANTE: localStorage NO es la fuente de verdad de seguridad.
   - La cookie HttpOnly es lo que el backend verifica en cada petición.
   - localStorage solo se usa para UI (nombre, rol para badges).
   - Si alguien manipula localStorage, no puede hacer nada malicioso porque
     el backend siempre verifica la cookie real.

   OBJETO ALMACENADO (ejemplo):
   { nombre: "Juan", apellido: "Pérez", correo: "juan@...", rol: 1, ... }
*/
const Auth = {
    KEY: 'patitas_user',

    // Guardar datos del usuario después del login exitoso.
    save(cliente) {
        localStorage.setItem(this.KEY, JSON.stringify(cliente));
    },

    // Obtener datos del usuario o null si no está logueado.
    get() {
        const raw = localStorage.getItem(this.KEY);
        return raw ? JSON.parse(raw) : null;
    },

    // Limpiar estado local (logout o sesión expirada).
    clear() {
        localStorage.removeItem(this.KEY);
    },

    // ¿Hay un usuario guardado en localStorage?
    isLoggedIn() {
        return this.get() !== null;
    },

    // Guard de navegación: redirige a login si no está logueado.
    // Se llama al inicio de páginas protegidas (dashboard, mascotas, citas, etc.).
    requireAuth() {
        if (!this.isLoggedIn()) {
            window.location.href = nav('/auth/login/');
        }
    },

    // Obtener nombre para mostrar en la navbar.
    getNombre() {
        const u = this.get();
        return u ? u.nombre : 'Usuario';
    },
};
