/* Admin layout — renders sidebar + topbar (mirroring ProyectoPyme _AdminLayout) */

function adminRequireAuth() {
    const user = Auth.get();
    if (!user) {
        window.location.replace(nav('/auth/login/'));
        document.body.style.display = 'none';
        return false;
    }
    if (user.rol !== 0) {
        const dest = user.rol === 2 ? '/veterinario/' : '/dashboard/';
        window.location.replace(nav(dest));
        document.body.style.display = 'none';
        return false;
    }
    return true;
}

function renderAdminLayout(activeSection) {
    if (!adminRequireAuth()) return;

    const fb = CONFIG.FRONTEND_BASE;
    const nombre = Auth.getNombre();

    const sections = [
        { id: 'panel',        icon: 'bi-speedometer2',     label: 'Panel',        href: `${fb}/admin/` },
        { id: 'servicios',    icon: 'bi-clipboard2-pulse', label: 'Servicios',    href: `${fb}/admin/servicios/` },
        { id: 'veterinarios', icon: 'bi-person-badge',     label: 'Veterinarios', href: `${fb}/admin/veterinarios/` },
        { id: 'clientes',     icon: 'bi-people',           label: 'Clientes',     href: `${fb}/admin/clientes/` },
        { id: 'citas',        icon: 'bi-calendar-check',   label: 'Citas',        href: `${fb}/admin/citas/` },
    ];

    const sidebarLinks = sections.map(s => {
        const cls = s.id === activeSection ? ' active' : '';
        return `<a href="${s.href}" class="sidebar-link${cls}"><i class="bi ${s.icon}"></i><span>${s.label}</span></a>`;
    }).join('');

    document.body.classList.add('admin-body-page');
    document.body.innerHTML = `
    <div class="admin-layout">
        <aside class="admin-sidebar">
            <a href="${fb}/admin/" class="sidebar-logo">
                <i class="bi bi-heart-pulse"></i>
                <span>Patitas Admin</span>
            </a>
            <div class="sidebar-label">Gestión</div>
            ${sidebarLinks}
            <div class="sidebar-spacer"></div>
            <a href="#" class="sidebar-link logout" onclick="doAdminLogout(); return false;">
                <i class="bi bi-box-arrow-left"></i><span>Cerrar sesión</span>
            </a>
        </aside>
        <div class="admin-main">
            <header class="admin-topbar">
                <div class="admin-badge">
                    <i class="bi bi-shield-lock"></i> Modo Administrador
                </div>
                <div class="admin-user">
                    <span>${escapeHtml(nombre)}</span>
                    <div class="avatar">Ad</div>
                </div>
            </header>
            <div class="admin-content" id="admin-content">
            </div>
        </div>
    </div>`;
}

async function doAdminLogout() {
    try { await apiPost('/auth/logout', {}); } catch (_) {}
    Auth.clear();
    window.location.href = nav('/auth/login/');
}

function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}

/* Helper: show modal */
function showModal(title, bodyHtml) {
    const existing = document.querySelector('.admin-modal-backdrop');
    if (existing) existing.remove();
    const backdrop = document.createElement('div');
    backdrop.className = 'admin-modal-backdrop';
    backdrop.innerHTML = `<div class="admin-modal"><h3>${title}</h3>${bodyHtml}</div>`;
    backdrop.addEventListener('click', e => { if (e.target === backdrop) backdrop.remove(); });
    document.body.appendChild(backdrop);
    return backdrop;
}

function closeModal() {
    const m = document.querySelector('.admin-modal-backdrop');
    if (m) m.remove();
}

/* Helper: show alert */
function showAlert(msg, type) {
    const el = document.getElementById('admin-alert');
    if (!el) return;
    el.className = `admin-alert ${type}`;
    el.textContent = msg;
    el.style.display = 'block';
    setTimeout(() => { el.style.display = 'none'; }, 4000);
}
