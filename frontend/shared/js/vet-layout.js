/* Vet layout — renders sidebar + topbar for veterinario portal */

function vetRequireAuth() {
    const user = Auth.get();
    if (!user) {
        window.location.replace(nav('/auth/login/'));
        document.body.style.display = 'none';
        return false;
    }
    if (user.rol !== 2) {
        const dest = user.rol === 0 ? '/admin/' : '/dashboard/';
        window.location.replace(nav(dest));
        document.body.style.display = 'none';
        return false;
    }
    return true;
}

function renderVetLayout(activeSection) {
    if (!vetRequireAuth()) return;

    const fb = CONFIG.FRONTEND_BASE;
    const nombre = Auth.getNombre();

    const sections = [
        { id: 'panel',    icon: 'bi-speedometer2',    label: 'Panel',      href: `${fb}/veterinario/` },
        { id: 'agenda',   icon: 'bi-calendar-check',  label: 'Mi Agenda',  href: `${fb}/veterinario/agenda/` },
    ];

    const sidebarLinks = sections.map(s => {
        const cls = s.id === activeSection ? ' active' : '';
        return `<a href="${s.href}" class="sidebar-link${cls}"><i class="bi ${s.icon}"></i><span>${s.label}</span></a>`;
    }).join('');

    document.body.classList.add('admin-body-page', 'vet-body-page');
    document.body.innerHTML = `
    <div class="admin-layout">
        <aside class="admin-sidebar">
            <a href="${fb}/veterinario/" class="sidebar-logo">
                <i class="bi bi-heart-pulse"></i>
                <span>Portal Veterinario</span>
            </a>
            <div class="sidebar-label">Mi espacio</div>
            ${sidebarLinks}
            <div class="sidebar-spacer"></div>
            <a href="#" class="sidebar-link logout" onclick="doVetLogout(); return false;">
                <i class="bi bi-box-arrow-left"></i><span>Cerrar sesión</span>
            </a>
        </aside>
        <div class="admin-main">
            <header class="admin-topbar">
                <div class="admin-badge" style="background:rgba(38,166,154,.15);color:#26a69a;">
                    <i class="bi bi-person-badge"></i> Veterinario
                </div>
                <div class="admin-user">
                    <span>${escapeHtml(nombre)}</span>
                    <div class="avatar" style="background:#26a69a;">Ve</div>
                </div>
            </header>
            <div class="admin-content" id="vet-content">
            </div>
        </div>
    </div>`;
}

async function doVetLogout() {
    try { await apiPost('/auth/logout', {}); } catch (_) {}
    Auth.clear();
    window.location.href = nav('/auth/login/');
}

function escapeHtml(str) {
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}
