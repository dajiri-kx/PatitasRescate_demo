document.addEventListener('DOMContentLoaded', () => {
    Auth.requireAuth();

    document.getElementById('link-dashboard').href = nav('/dashboard/');

    document.getElementById('mascotaForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const alertBox = document.getElementById('alert-box');
        alertBox.innerHTML = '';
        const btn = document.getElementById('btnSubmit');
        btn.disabled = true;
        btn.textContent = 'Registrando...';

        try {
            const res = await apiPost('/Mascotas/AgregarMascota/agregarMascotaAction.php', {
                nombre: document.getElementById('nombre').value,
                especie: document.getElementById('especie').value,
                raza: document.getElementById('raza').value,
                edad: parseInt(document.getElementById('edad').value, 10),
            });
            alertBox.innerHTML = `<div class="alert alert-success">${escapeHtml(res.message)}</div>`;
            setTimeout(() => {
                window.location.href = nav('/dashboard/');
            }, 1500);
        } catch (err) {
            alertBox.innerHTML = `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
        } finally {
            btn.disabled = false;
            btn.textContent = 'Registrar Mascota';
        }
    });
});
