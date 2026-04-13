/*
   api.js — Wrapper centralizado para todas las llamadas HTTP al backend.

   PATRÓN:
   Todos los .js de features usan apiGet() y apiPost() en vez de fetch() directo.
   Esto centraliza: base URL, headers, credentials, y manejo de 401.

   FLUJO DE DATOS:
   Feature JS → apiGet('/citas/mis-citas') → apiFetch('/citas/mis-citas', {method:'GET'})
   → fetch('http://localhost:8080/api/citas/mis-citas', {credentials:'include'})
   → Backend responde JSON {ok:true, data:...} o {ok:false, error:"..."}
   → apiFetch retorna el objeto parseado o lanza Error.

   MANEJO DE 401 (sesión expirada):
   Si el backend responde 401, verificamos si el usuario ESTABA logueado:
   - Sí → limpiamos localStorage y redirigimos a login (sesión expiró).
   - No → lanzamos error normal (intento sin autenticar, e.g., login fallido).
   Esto evita redirecciones infinitas cuando un login falla con 401.

   CREDENTIALS: 'include' envía la cookie de sesión en cada petición.
   Sin esto, gorilla/sessions no recibiría la cookie y toda petición sería anónima.
*/
const CONFIG = {
    API_BASE: '/api',
    FRONTEND_BASE: '/frontend/features',
};

// apiFetch es el fetch centralizado. Todos los endpoints pasan por aquí.
async function apiFetch(featurePath, options = {}) {
    const url = CONFIG.API_BASE + featurePath;
    const defaults = {
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',  // Envía la cookie de sesión siempre.
    };
    const opts = { ...defaults, ...options };

    const res = await fetch(url, opts);
    const data = await res.json();

    if (!res.ok) {
        if (res.status === 401) {
            // ¿Estaba logueado antes de esta petición?
            const wasLoggedIn = Auth.isLoggedIn();
            Auth.clear();
            if (wasLoggedIn) {
                // Sesión expiró → redirigir a login.
                window.location.href = CONFIG.FRONTEND_BASE + '/auth/login/';
                return;
            }
            // No estaba logueado → error normal (e.g., login incorrecto).
            throw new Error(data.error || 'Credenciales incorrectas');
        }
        throw new Error(data.error || 'Error desconocido');
    }
    return data;
}

// apiGet — shortcut para peticiones GET.
async function apiGet(featurePath) {
    return apiFetch(featurePath, { method: 'GET' });
}

// apiPost — shortcut para peticiones POST con body JSON.
async function apiPost(featurePath, body) {
    return apiFetch(featurePath, {
        method: 'POST',
        body: JSON.stringify(body),
    });
}

// nav — helper para construir URLs del frontend sin repetir el prefijo.
function nav(featurePath) {
    return CONFIG.FRONTEND_BASE + featurePath;
}
