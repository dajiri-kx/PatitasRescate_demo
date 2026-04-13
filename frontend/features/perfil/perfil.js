document.addEventListener('DOMContentLoaded', () => {
    Auth.requireAuth();
    document.getElementById('link-editar').href = nav('/proximamente/') + '?msg=Editar perfil próximamente.';

    const cliente = Auth.get();
    const el = document.getElementById('perfil-info');
    el.innerHTML = `
        <div class="profile-field">
            <span class="profile-label"><i class="bi bi-person-badge me-2"></i>ID Cliente</span>
            <span class="profile-value">${escapeHtml(String(cliente.id_cliente))}</span>
        </div>
        <div class="profile-field">
            <span class="profile-label"><i class="bi bi-person me-2"></i>Nombre</span>
            <span class="profile-value">${escapeHtml(cliente.nombre)} ${escapeHtml(cliente.apellido)}</span>
        </div>
        <div class="profile-field">
            <span class="profile-label"><i class="bi bi-envelope me-2"></i>Correo</span>
            <span class="profile-value">${escapeHtml(cliente.correo)}</span>
        </div>
        <div class="profile-field">
            <span class="profile-label"><i class="bi bi-telephone me-2"></i>Teléfono</span>
            <span class="profile-value">${escapeHtml(cliente.telefono)}</span>
        </div>
    `;
});
