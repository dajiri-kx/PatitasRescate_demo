async function loadComponent(elementId, url) {
    const el = document.getElementById(elementId);
    if (!el) return;
    try {
        const res = await fetch(url);
        el.innerHTML = await res.text();
    } catch (e) {
        console.error('Error cargando componente:', url, e);
    }
}

function renderHeader() {
    const loggedIn = Auth.isLoggedIn();
    const nombre = Auth.getNombre();
    const fb = CONFIG.FRONTEND_BASE;

    /* --- Nav links depend on auth state --- */
    const navLinks = loggedIn
        ? `<li class="nav-item"><a class="nav-link px-3 py-2 rounded" href="${fb}/dashboard/"><i class="bi bi-house-door me-1"></i>Inicio</a></li>
           <li class="nav-item"><a class="nav-link px-3 py-2 rounded" href="${fb}/mis-citas/"><i class="bi bi-calendar-check me-1"></i>Mis Citas</a></li>
           <li class="nav-item"><a class="nav-link px-3 py-2 rounded" href="${fb}/mis-mascotas/"><i class="bi bi-heart me-1"></i>Mis Mascotas</a></li>
           <li class="nav-item"><a class="nav-link px-3 py-2 rounded" href="${fb}/mis-facturas/"><i class="bi bi-receipt me-1"></i>Mis Facturas</a></li>
           <li class="nav-item"><a class="nav-link px-3 py-2 rounded" href="${fb}/servicios/"><i class="bi bi-clipboard2-pulse me-1"></i>Servicios</a></li>`
        : `<li class="nav-item"><a class="nav-link px-3 py-2 rounded" href="${fb}/servicios/">Servicios</a></li>
           <li class="nav-item"><a class="nav-link px-3 py-2 rounded" href="${fb}/ubicacion/">Ubicación</a></li>
           <li class="nav-item"><a class="nav-link px-3 py-2 rounded" href="${fb}/contactenos/">Contáctenos</a></li>`;

    /* --- Right‑side auth / user block --- */
    const userDropdown = (id) => `
        <div class="dropdown text-end">
            <button class="btn btn-outline-light dropdown-toggle" type="button" id="${id}" data-bs-toggle="dropdown" aria-expanded="false">
                <i class="bi bi-person-circle me-1"></i>${escapeHtml(nombre)}
            </button>
            <ul class="dropdown-menu dropdown-menu-end" style="background-color:#ff8c00;">
                <li><a class="dropdown-item text-white" href="${fb}/perfil/"><i class="bi bi-person me-1"></i>Mi Perfil</a></li>
                <li><hr class="dropdown-divider"></li>
                <li><a class="dropdown-item text-white" href="#" onclick="doLogout()"><i class="bi bi-box-arrow-right me-1"></i>Cerrar Sesión</a></li>
            </ul>
        </div>`;

    const authBlockMobile = loggedIn
        ? userDropdown('userDropdownMobile')
        : `<a href="${fb}/auth/login/" class="btn btn-outline-light">Iniciar sesión</a>
           <a href="${fb}/auth/register/" class="btn btn-light ms-2" style="color:#ff6b00;">Registrarse</a>`;

    const authBlockDesktop = loggedIn
        ? userDropdown('userDropdownDesktop')
        : `<a href="${fb}/auth/login/" class="btn btn-outline-light">Iniciar sesión</a>
           <a href="${fb}/auth/register/" class="btn btn-light ms-2" style="color:#ff6b00;">Registrarse</a>`;

    document.getElementById('app-header').innerHTML = `
    <nav class="navbar navbar-expand-lg navbar-dark" style="background:linear-gradient(135deg,#ff8c00,#ff6b00);box-shadow:0 2px 10px rgba(0,0,0,0.1);">
        <div class="container d-flex align-items-center">
            <div>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarsMain" aria-controls="navbarsMain" aria-expanded="false" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <a class="navbar-brand ms-2 d-flex align-items-center" href="${loggedIn ? fb + '/dashboard/' : fb + '/home/'}">
                    <i class="bi bi-heart-pulse me-2"></i>
                    <span class="fw-bold">Patitas al rescate</span>
                </a>
            </div>
            <div class="col-md-3 text-end d-lg-none">${authBlockMobile}</div>
            <div class="collapse navbar-collapse" id="navbarsMain">
                <ul class="navbar-nav me-auto mb-2 mb-lg-0">
                    ${navLinks}
                </ul>
                <div class="col-md-3 text-end d-none d-lg-block">${authBlockDesktop}</div>
            </div>
        </div>
    </nav>`;
}

function renderFooter() {
    document.getElementById('app-footer').innerHTML = `
    <div class="container">
        <footer class="py-2 py-md-5">
            <div class="d-flex flex-column flex-sm-row justify-content-between pt-2 border-top">
                <p>&copy; 2025 Patitas al rescate. Todos los derechos reservados.</p>
                <ul class="list-unstyled d-flex">
                    <li class="ms-3"><a class="link-body-emphasis" href="#"><i class="bi bi-twitter"></i></a></li>
                    <li class="ms-3"><a class="link-body-emphasis" href="#"><i class="bi bi-instagram"></i></a></li>
                    <li class="ms-3"><a class="link-body-emphasis" href="#"><i class="bi bi-facebook"></i></a></li>
                    <li class="ms-3"><a class="link-body-emphasis" href="#"><i class="bi bi-tiktok"></i></a></li>
                </ul>
            </div>
        </footer>
    </div>`;
}

async function doLogout() {
    try {
        await apiPost('/auth/logout', {});
    } catch (_) { /* ignore */ }
    Auth.clear();
    window.location.href = nav('/auth/login/');
}

function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

function initLayout() {
    renderHeader();
    renderFooter();
}

document.addEventListener('DOMContentLoaded', initLayout);
