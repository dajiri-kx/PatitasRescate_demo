document.addEventListener('DOMContentLoaded', async () => {
    Auth.requireAuth();
    document.getElementById('link-agregar').href = nav('/mascotas/agregar-mascota/');

    try {
        const res = await apiGet('/mascotas');
        renderMascotas(res.data);
    } catch (err) {
        document.getElementById('alert-box').innerHTML =
            `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
    }
});

function renderMascotas(mascotas) {
    const tbody = document.getElementById('mascotas-body');
    if (!mascotas || mascotas.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" class="text-center"><div class="alert alert-info mb-0">No hay mascotas registradas. ¡Agrega tu primera mascota!</div></td></tr>';
        return;
    }
    tbody.innerHTML = mascotas.map(m => `
        <tr>
            <td>${escapeHtml(m.NOMBRE_MASCOTA || 'No disponible')}</td>
            <td>${escapeHtml(m.ESPECIE || 'No disponible')}</td>
            <td>${escapeHtml(m.RAZA || 'No disponible')}</td>
            <td>${escapeHtml(String(m.MESES || 'No disponible'))}</td>
            <td>
                <a href="${nav('/proximamente/')}?msg=Editar mascota próximamente." class="btn btn-outline-primary btn-sm"><i class="bi bi-pencil"></i></a>
                <button class="btn btn-outline-danger btn-sm" onclick="confirmarEliminacion('${escapeHtml(m.NOMBRE_MASCOTA || '')}')"><i class="bi bi-trash"></i></button>
            </td>
        </tr>
    `).join('');
}

function confirmarEliminacion(nombre) {
    if (confirm('¿Estás seguro de que quieres eliminar a ' + nombre + '? Esta acción no se puede deshacer.')) {
        alert('Funcionalidad de eliminar próximamente.');
    }
}
