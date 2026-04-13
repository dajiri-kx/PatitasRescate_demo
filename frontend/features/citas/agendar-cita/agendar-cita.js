document.addEventListener('DOMContentLoaded', async () => {
    Auth.requireAuth();

    document.getElementById('link-dashboard').href = nav('/dashboard/');

    // Read category from URL params (from servicios page link)
    const params = new URLSearchParams(window.location.search);
    const categoriaParam = params.get('categoria') || '';

    // Category config for title customization
    const categoryConfig = {
        'Consulta':    { title: 'Agendar Consulta Veterinaria', subtitle: 'Evaluación médica completa para tu mascota', color: '#0d6efd' },
        'Vacunación':  { title: 'Agendar Vacunación', subtitle: 'Protege a tu mascota con vacunas al día', color: '#198754' },
        'Cirugía':     { title: 'Agendar Cirugía / Procedimiento', subtitle: 'Intervenciones con equipos especializados', color: '#dc3545' },
        'Estética':    { title: 'Agendar Servicio de Estética', subtitle: 'Baños, cortes y tratamientos de spa', color: '#ffc107' },
        'Diagnóstico': { title: 'Agendar Diagnóstico por Imágenes', subtitle: 'Radiografías, ecografías y exámenes', color: '#0dcaf0' },
    };

    const categorias = ['Consulta', 'Vacunación', 'Cirugía', 'Estética', 'Diagnóstico'];

    // Populate category dropdown
    const catSelect = document.getElementById('categoria');
    categorias.forEach(c => {
        const opt = document.createElement('option');
        opt.value = c;
        opt.textContent = c;
        catSelect.appendChild(opt);
    });

    function applyConfig(cat) {
        const config = categoryConfig[cat];
        if (config) {
            document.getElementById('form-title').textContent = config.title;
            document.getElementById('form-title').style.color = config.color;
            document.getElementById('form-subtitle').textContent = config.subtitle;
        } else {
            document.getElementById('form-title').textContent = 'Agendar Cita';
            document.getElementById('form-title').style.color = '#ffc107';
            document.getElementById('form-subtitle').textContent = '';
        }
    }

    // Load services for a selected category
    async function loadServicios(cat) {
        const section = document.getElementById('servicios-section');
        const container = document.getElementById('servicios-container');

        if (!cat) {
            section.style.display = 'none';
            container.innerHTML = '';
            return;
        }

        section.style.display = 'block';
        container.innerHTML = '<p class="text-muted">Cargando servicios...</p>';

        try {
            const res = await apiGet('/citas/servicios?categoria=' + encodeURIComponent(cat));
            if (!res.data.length) {
                container.innerHTML = '<p class="text-muted">No hay servicios disponibles en esta categoría.</p>';
            } else {
                container.innerHTML = res.data.map(s => `
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
        } catch (err) {
            container.innerHTML = `<p class="text-danger">Error: ${escapeHtml(err.message)}</p>`;
        }
    }

    // Category change handler
    catSelect.addEventListener('change', () => {
        applyConfig(catSelect.value);
        loadServicios(catSelect.value);
    });

    // If URL has ?categoria=, pre-select it
    if (categoriaParam && categorias.includes(categoriaParam)) {
        catSelect.value = categoriaParam;
        applyConfig(categoriaParam);
        loadServicios(categoriaParam);
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

    // Load mascotas and veterinarios
    try {
        const [mascotasRes, vetsRes] = await Promise.all([
            apiGet('/mascotas/nombres'),
            apiGet('/citas/veterinarios'),
        ]);

        const mascotaSelect = document.getElementById('id_mascota');
        mascotaSelect.innerHTML = '<option value="">-- Elige una mascota --</option>';
        mascotasRes.data.forEach(m => {
            const opt = document.createElement('option');
            opt.value = m.ID_MASCOTA;
            opt.textContent = m.NOMBRE;
            mascotaSelect.appendChild(opt);
        });

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
                try {
                    const checkoutRes = await apiPost('/checkout/crear-sesion', { id_factura: String(idFactura) });
                    window.location.href = checkoutRes.data.url;
                    return;
                } catch (checkoutErr) {
                    // Si falla Stripe, la cita se creó igual — mostrar modal
                    showResultModal(
                        'warning',
                        '¡Cita agendada!',
                        '<i class="bi bi-check-circle" style="font-size:48px;color:#198754;"></i>' +
                        '<p class="mt-3 mb-1" style="font-size:18px;font-weight:600;">Tu cita fue agendada exitosamente.</p>' +
                        '<p class="text-muted">No se pudo redirigir al pago en este momento.<br>Puedes completar el pago desde <strong>Mis Facturas</strong>.</p>'
                    );
                }
            } else {
                showResultModal(
                    'success',
                    '¡Cita agendada!',
                    '<i class="bi bi-calendar-check" style="font-size:48px;color:#198754;"></i>' +
                    '<p class="mt-3" style="font-size:18px;font-weight:600;">' + escapeHtml(res.data?.message || 'Cita agendada exitosamente.') + '</p>'
                );
            }
            document.getElementById('agendarForm').reset();
            document.getElementById('servicios-section').style.display = 'none';
        } catch (err) {
            alertBox.innerHTML = `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
        } finally {
            btn.disabled = false;
        }
    });

    function showResultModal(type, title, bodyHtml) {
        const header = document.getElementById('resultModalHeader');
        const titleEl = document.getElementById('resultModalTitle');
        const body = document.getElementById('resultModalBody');
        const footer = document.getElementById('resultModalFooter');
        const fb = CONFIG.FRONTEND_BASE;

        const colors = { success: '#198754', warning: '#ffc107' };
        header.style.background = colors[type] || '#198754';
        header.style.color = '#fff';
        titleEl.textContent = title;
        body.innerHTML = bodyHtml;

        const btns = type === 'warning'
            ? `<a href="${fb}/mis-facturas/" class="btn btn-warning"><i class="bi bi-receipt me-1"></i>Ir a Mis Facturas</a>
               <a href="${fb}/dashboard/" class="btn btn-outline-secondary">Volver al Dashboard</a>`
            : `<a href="${fb}/dashboard/" class="btn btn-success"><i class="bi bi-house-door me-1"></i>Volver al Dashboard</a>`;
        footer.innerHTML = btns;

        const modal = new bootstrap.Modal(document.getElementById('resultModal'));
        modal.show();
    }
});
