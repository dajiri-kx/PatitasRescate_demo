document.addEventListener('DOMContentLoaded', async () => {
    Auth.requireAuth();

    const cliente = Auth.get();
    document.getElementById('nombre-usuario').textContent = cliente.nombre;

    // Set card links
    document.getElementById('card-mascotas').href = nav('/mis-mascotas/');
    document.getElementById('card-citas').href = nav('/mis-citas/');
    document.getElementById('card-facturas').href = nav('/mis-facturas/');
    document.getElementById('card-agendar').href = nav('/citas/agendar-cita/');

    // Load counts in parallel
    try {
        const [mascotasRes, citasRes, facturasRes] = await Promise.all([
            apiGet('/mascotas'),
            apiGet('/citas'),
            apiGet('/facturas'),
        ]);

        const mascotas = mascotasRes.data || [];
        const citas = citasRes.data || [];
        const facturas = facturasRes.data || [];

        document.getElementById('count-mascotas').textContent = mascotas.length;
        document.getElementById('count-citas').textContent = citas.filter(c => c.ESTADO === 'Activa').length;
        document.getElementById('count-facturas-pendientes').textContent = facturas.filter(f => f.ESTADO === 'Pendiente').length;
    } catch (err) {
        document.getElementById('alert-box').innerHTML =
            `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
    }
});
