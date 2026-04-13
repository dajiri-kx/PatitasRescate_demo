document.addEventListener('DOMContentLoaded', async () => {
    Auth.requireAuth();
    document.getElementById('link-agendar').href = nav('/citas/agendar-cita/');

    try {
        const res = await apiGet('/citas');
        renderCitas(res.data);
    } catch (err) {
        document.getElementById('alert-box').innerHTML =
            `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
    }
});

function renderCitas(citas) {
    const tbody = document.getElementById('citas-body');
    if (!citas || citas.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="text-center"><div class="alert alert-info mb-0">No hay citas agendadas. ¡Agenda tu primera cita!</div></td></tr>';
        return;
    }
    tbody.innerHTML = citas.map(c => {
        const badgeClass = c.ESTADO === 'Activa' ? 'bg-warning text-dark' : (c.ESTADO === 'Completada' ? 'bg-success' : 'bg-secondary');
        const acciones = c.ESTADO === 'Activa'
            ? `<a href="${nav('/citas/cancelar-cita/')}" class="btn btn-outline-danger btn-sm"><i class="bi bi-x-circle"></i> Cancelar</a>`
            : '';
        return `
            <tr>
                <td>${escapeHtml(String(c.ID_CITA))}</td>
                <td>${escapeHtml(c.MASCOTA)}</td>
                <td>${escapeHtml(c.VETERINARIO || '')}</td>
                <td>${escapeHtml(c.FECHA_CITA)}</td>
                <td><span class="badge ${badgeClass}">${escapeHtml(c.ESTADO)}</span></td>
                <td>${acciones}</td>
            </tr>`;
    }).join('');
}
