document.addEventListener('DOMContentLoaded', async () => {
    Auth.requireAuth();

    document.getElementById('link-dashboard').href = nav('/dashboard/');

    // Read category from URL params
    const params = new URLSearchParams(window.location.search);
    const categoria = params.get('categoria') || '';

    // Category-specific UI configuration
    const categoryConfig = {
        'Consulta':    { title: 'Agendar Consulta Veterinaria', subtitle: 'Evaluación médica completa para tu mascota', icon: 'bi-stethoscope', color: '#0d6efd' },
        'Vacunación':  { title: 'Agendar Vacunación', subtitle: 'Protege a tu mascota con vacunas al día', icon: 'bi-shield-check', color: '#198754' },
        'Cirugía':     { title: 'Agendar Cirugía / Procedimiento', subtitle: 'Intervenciones con equipos especializados', icon: 'bi-heart-pulse', color: '#dc3545' },
        'Estética':    { title: 'Agendar Servicio de Estética', subtitle: 'Baños, cortes y tratamientos de spa', icon: 'bi-scissors', color: '#ffc107' },
        'Diagnóstico': { title: 'Agendar Diagnóstico por Imágenes', subtitle: 'Radiografías, ecografías y exámenes', icon: 'bi-camera', color: '#0dcaf0' },
    };

    const config = categoryConfig[categoria];
    if (config) {
        document.getElementById('form-title').textContent = config.title;
        document.getElementById('form-title').style.color = config.color;
        document.getElementById('form-subtitle').textContent = config.subtitle;
    }

    // Generate time slots
    const horaSelect = document.getElementById('hora');
    for (let h = 8; h <= 17; h++) {
        for (const min of ['00', '30']) {
            const time = String(h).padStart(2, '0') + ':' + min;
            const opt = document.createElement('option');
            opt.value = time;
            opt.textContent = time;
            horaSelect.appendChild(opt);
        }
    }

    // Build services API URL
    const serviciosUrl = categoria
        ? '/citas/servicios?categoria=' + encodeURIComponent(categoria)
        : '/citas/servicios';

    // Load form data in parallel
    try {
        const [mascotasRes, serviciosRes, vetsRes] = await Promise.all([
            apiGet('/mascotas/nombres'),
            apiGet(serviciosUrl),
            apiGet('/citas/veterinarios'),
        ]);

        // Populate mascotas
        const mascotaSelect = document.getElementById('id_mascota');
        mascotaSelect.innerHTML = '<option value="">-- Elige una mascota --</option>';
        mascotasRes.data.forEach(m => {
            const opt = document.createElement('option');
            opt.value = m.ID_MASCOTA;
            opt.textContent = m.NOMBRE;
            mascotaSelect.appendChild(opt);
        });

        // Populate servicios with price and description
        const serviciosContainer = document.getElementById('servicios-container');
        if (serviciosRes.data.length === 0) {
            serviciosContainer.innerHTML = '<p class="text-muted">No hay servicios disponibles en esta categoría.</p>';
        } else {
            serviciosContainer.innerHTML = serviciosRes.data.map(s => `
                <div class="form-check border rounded p-2 mb-2">
                    <input class="form-check-input" type="checkbox" name="servicio" value="${s.ID_SERVICIO}" id="servicio_${s.ID_SERVICIO}">
                    <label class="form-check-label w-100" for="servicio_${s.ID_SERVICIO}">
                        <strong>${escapeHtml(s.NOMBRE_SERVICIO)}</strong>
                        <span class="text-success float-end">₡${Number(s.PRECIO).toLocaleString('es-CR')}</span>
                        <br><small class="text-muted">${escapeHtml(s.DESCRIPCION)}</small>
                    </label>
                </div>
            `).join('');
        }

        // Populate veterinarios
        const vetSelect = document.getElementById('veterinario');
        vetSelect.innerHTML = '<option value="">-- Elige un veterinario --</option>';
        vetsRes.data.forEach(v => {
            const opt = document.createElement('option');
            opt.value = v.ID_VETERINARIO;
            opt.textContent = v.NOMBRE;
            vetSelect.appendChild(opt);
        });
    } catch (err) {
        document.getElementById('alert-box').innerHTML =
            `<div class="alert alert-danger">Error al cargar datos: ${escapeHtml(err.message)}</div>`;
    }

    // Form submit
    document.getElementById('agendarForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const alertBox = document.getElementById('alert-box');
        alertBox.innerHTML = '';
        const btn = document.getElementById('btnSubmit');
        btn.disabled = true;

        const serviciosChecked = Array.from(document.querySelectorAll('input[name="servicio"]:checked')).map(c => c.value);

        try {
            const res = await apiPost('/citas/agendar', {
                id_mascota: document.getElementById('id_mascota').value,
                fecha: document.getElementById('fecha').value,
                hora: document.getElementById('hora').value,
                servicio: serviciosChecked,
                veterinario: document.getElementById('veterinario').value,
            });

            // Intentar redirigir al checkout de Stripe
            const idFactura = res.data && res.data.id_factura;
            if (idFactura) {
                alertBox.innerHTML = `<div class="alert alert-info">Cita agendada. Redirigiendo al pago...</div>`;
                try {
                    const checkoutRes = await apiPost('/checkout/crear-sesion', { id_factura: String(idFactura) });
                    window.location.href = checkoutRes.data.url;
                    return;
                } catch (checkoutErr) {
                    // Si falla Stripe, igual la cita se creó
                    alertBox.innerHTML = `<div class="alert alert-warning">Cita agendada exitosamente. No se pudo redirigir al pago: ${escapeHtml(checkoutErr.message)}. Puedes pagar desde Mis Facturas.</div>`;
                }
            } else {
                alertBox.innerHTML = `<div class="alert alert-success">${escapeHtml(res.data?.message || res.message || 'Cita agendada exitosamente.')}</div>`;
            }
            document.getElementById('agendarForm').reset();
        } catch (err) {
            alertBox.innerHTML = `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
        } finally {
            btn.disabled = false;
        }
    });
});
