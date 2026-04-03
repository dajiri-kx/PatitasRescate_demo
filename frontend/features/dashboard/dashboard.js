document.addEventListener('DOMContentLoaded', async () => {
    Auth.requireAuth();

    const proximamenteUrl = nav('/proximamente/');
    document.getElementById('link-editar-perfil').href = proximamenteUrl;
    document.getElementById('link-agregar-mascota').href = nav('/mascotas/agregar-mascota/');
    document.getElementById('link-agendar-cita').href = nav('/citas/agendar-cita/');

    const cliente = Auth.get();

    // Render profile
    const perfilEl = document.getElementById('perfil-info');
    perfilEl.innerHTML = `
        <p><strong>ID Cliente:</strong> ${escapeHtml(String(cliente.id_cliente))}</p>
        <p><strong>Nombre:</strong> ${escapeHtml(cliente.nombre)}</p>
        <p><strong>Apellido:</strong> ${escapeHtml(cliente.apellido)}</p>
        <p><strong>Correo:</strong> ${escapeHtml(cliente.correo)}</p>
        <p><strong>Teléfono:</strong> ${escapeHtml(cliente.telefono)}</p>
    `;

    // Load data in parallel
    try {
        const [mascotasRes, citasRes, facturasRes] = await Promise.all([
            apiGet('/mascotas'),
            apiGet('/citas'),
            apiGet('/facturas'),
        ]);

        renderMascotas(mascotasRes.data);
        renderCitas(citasRes.data);
        renderFacturas(facturasRes.data);
    } catch (err) {
        document.getElementById('alert-box').innerHTML =
            `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
    }
});

function renderMascotas(mascotas) {
    const tbody = document.getElementById('mascotas-body');
    if (!mascotas || mascotas.length === 0) {
        tbody.innerHTML = '<tr><td colspan="4" class="text-center"><div class="alert alert-info">No hay mascotas asociadas. Por favor registre sus mascotas.</div></td></tr>';
        return;
    }
    tbody.innerHTML = mascotas.map(m => `
        <tr>
            <td>${escapeHtml(m.NOMBRE_MASCOTA || 'No disponible')}</td>
            <td>${escapeHtml(m.RAZA || 'No disponible')}</td>
            <td>${escapeHtml(String(m.MESES || 'No disponible'))}</td>
            <td>
                <a href="${nav('/proximamente/')}?msg=Editar mascota próximamente." class="btn btn-primary btn-sm">Editar</a>
                <button class="btn btn-danger btn-sm" onclick="confirmarEliminacion('${escapeHtml(m.NOMBRE_MASCOTA || '')}')">Eliminar</button>
            </td>
        </tr>
    `).join('');
}

function renderCitas(citas) {
    const tbody = document.getElementById('citas-body');
    if (!citas || citas.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" class="text-center"><div class="alert alert-info">No hay citas agendadas. Por favor agende una cita.</div></td></tr>';
        return;
    }
    tbody.innerHTML = citas.map(c => {
        const badgeClass = c.ESTADO === 'Activa' ? 'bg-warning' : (c.ESTADO === 'Completada' ? 'bg-success' : 'bg-secondary');
        const acciones = c.ESTADO === 'Activa'
            ? `<a href="${nav('/proximamente/')}?msg=Editar cita próximamente." class="btn btn-warning btn-sm">Editar</a>
               <a href="${nav('/citas/cancelar-cita/')}" class="btn btn-danger btn-sm">Cancelar</a>`
            : '';
        return `
            <tr>
                <td>${escapeHtml(String(c.ID_CITA))}</td>
                <td>${escapeHtml(c.MASCOTA)}</td>
                <td>${escapeHtml(c.FECHA_CITA)}</td>
                <td><span class="badge ${badgeClass}">${escapeHtml(c.ESTADO)}</span></td>
                <td>${acciones}</td>
            </tr>`;
    }).join('');
}

function renderFacturas(facturas) {
    const container = document.getElementById('facturas-container');
    if (!facturas || facturas.length === 0) {
        container.innerHTML = '<div class="alert alert-info">No hay facturas registradas para las mascotas del usuario.</div>';
        return;
    }
    container.innerHTML = `
        <div class="table-responsive">
            <table class="table table-striped">
                <thead><tr><th>ID Factura</th><th>Fecha</th><th>Total</th><th>Estado</th></tr></thead>
                <tbody>
                    ${facturas.map(f => {
                        const badgeClass = f.ESTADO === 'Pagada' ? 'bg-success' : 'bg-warning';
                        const total = Number(f.TOTAL).toLocaleString('es-CR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
                        return `
                            <tr>
                                <td>${escapeHtml(String(f.ID_FACTURA))}</td>
                                <td>${escapeHtml(f.FECHA_FACTURA)}</td>
                                <td>₡${total}</td>
                                <td><span class="badge ${badgeClass}">${escapeHtml(f.ESTADO)}</span></td>
                            </tr>`;
                    }).join('')}
                </tbody>
            </table>
        </div>`;
}

function confirmarEliminacion(nombre) {
    if (confirm('¿Estás seguro de que quieres eliminar a ' + nombre + '? Esta acción no se puede deshacer.')) {
        alert('Funcionalidad de eliminar próximamente.');
    }
}
