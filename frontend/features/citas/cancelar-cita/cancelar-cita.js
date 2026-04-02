document.addEventListener('DOMContentLoaded', async () => {
    Auth.requireAuth();

    document.getElementById('link-dashboard').href = nav('/dashboard/');

    try {
        const res = await apiGet('/Citas/ObtenerCitasActivas/obtenerCitasActivasAction.php');
        const select = document.getElementById('id_cita');

        if (!res.data || res.data.length === 0) {
            select.innerHTML = '<option value="">No tienes citas activas</option>';
            document.getElementById('btnSubmit').disabled = true;
        } else {
            select.innerHTML = '<option value="">-- Selecciona una cita --</option>';
            res.data.forEach(c => {
                const opt = document.createElement('option');
                opt.value = c.ID_CITA;
                opt.textContent = c.FECHA_CITA + ' con ' + c.MASCOTA;
                select.appendChild(opt);
            });
        }
    } catch (err) {
        document.getElementById('alert-box').innerHTML =
            `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
    }

    document.getElementById('cancelarForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const alertBox = document.getElementById('alert-box');
        alertBox.innerHTML = '';
        const btn = document.getElementById('btnSubmit');
        btn.disabled = true;

        try {
            const res = await apiPost('/Citas/CancelarCita/cancelarCitaAction.php', {
                id_cita: document.getElementById('id_cita').value,
            });
            alertBox.innerHTML = `<div class="alert alert-success">${escapeHtml(res.message)}</div>`;
            // Reload the list
            setTimeout(() => window.location.reload(), 1500);
        } catch (err) {
            alertBox.innerHTML = `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
        } finally {
            btn.disabled = false;
        }
    });
});
