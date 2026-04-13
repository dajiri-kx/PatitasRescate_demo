document.addEventListener('DOMContentLoaded', async () => {
    Auth.requireAuth();

    try {
        const res = await apiGet('/facturas');
        renderFacturas(res.data);
    } catch (err) {
        document.getElementById('alert-box').innerHTML =
            `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
    }
});

function renderFacturas(facturas) {
    const container = document.getElementById('facturas-container');
    if (!facturas || facturas.length === 0) {
        container.innerHTML = '<div class="alert alert-info mb-0">No hay facturas registradas.</div>';
        return;
    }
    container.innerHTML = `
        <div class="table-responsive">
            <table class="table table-striped table-hover">
                <thead><tr><th>ID</th><th>Fecha</th><th>Total</th><th>Estado</th><th>Acción</th></tr></thead>
                <tbody>
                    ${facturas.map(f => {
                        const badgeClass = f.ESTADO === 'Pagada' ? 'bg-success' : 'bg-warning text-dark';
                        const total = Number(f.TOTAL).toLocaleString('es-CR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
                        const pagarBtn = f.ESTADO === 'Pendiente'
                            ? `<button class="btn btn-success btn-sm" onclick="iniciarPago('${f.ID_FACTURA}', this)">
                                   <i class="bi bi-credit-card"></i> Pagar
                               </button>`
                            : `<span class="text-success"><i class="bi bi-check-circle-fill"></i></span>`;
                        return `
                            <tr>
                                <td>${escapeHtml(String(f.ID_FACTURA))}</td>
                                <td>${escapeHtml(f.FECHA_FACTURA)}</td>
                                <td><strong>₡${total}</strong></td>
                                <td><span class="badge ${badgeClass}">${escapeHtml(f.ESTADO)}</span></td>
                                <td>${pagarBtn}</td>
                            </tr>`;
                    }).join('')}
                </tbody>
            </table>
        </div>`;
}

async function iniciarPago(idFactura, btn) {
    btn.disabled = true;
    btn.innerHTML = '<span class="spinner-border spinner-border-sm"></span> Procesando...';
    try {
        const res = await apiPost('/checkout/crear-sesion', { id_factura: String(idFactura) });
        window.location.href = res.data.url;
    } catch (err) {
        btn.disabled = false;
        btn.innerHTML = '<i class="bi bi-credit-card"></i> Pagar';
        alert('Error al iniciar el pago: ' + err.message);
    }
}
