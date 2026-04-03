document.addEventListener('DOMContentLoaded', async () => {
    Auth.requireAuth();

    document.getElementById('link-dashboard').href = nav('/dashboard/');

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

    // Load form data in parallel
    try {
        const [mascotasRes, serviciosRes, vetsRes] = await Promise.all([
            apiGet('/mascotas/nombres'),
            apiGet('/citas/servicios'),
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

        // Populate servicios
        const serviciosContainer = document.getElementById('servicios-container');
        serviciosContainer.innerHTML = serviciosRes.data.map(s => `
            <div class="form-check">
                <input class="form-check-input" type="checkbox" name="servicio" value="${s.ID_SERVICIO}" id="servicio_${s.ID_SERVICIO}">
                <label class="form-check-label" for="servicio_${s.ID_SERVICIO}">${escapeHtml(s.NOMBRE_SERVICIO)}</label>
            </div>
        `).join('');

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
            alertBox.innerHTML = `<div class="alert alert-success">${escapeHtml(res.message)}</div>`;
            document.getElementById('agendarForm').reset();
        } catch (err) {
            alertBox.innerHTML = `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
        } finally {
            btn.disabled = false;
        }
    });
});
