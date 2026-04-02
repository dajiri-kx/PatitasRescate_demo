document.getElementById('contactForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const alertBox = document.getElementById('alert-box');
    alertBox.innerHTML = '';
    const btn = document.getElementById('btnSubmit');
    btn.disabled = true;
    btn.textContent = 'Enviando...';

    try {
        await apiPost('/Contacto/EnviarContacto/contactoAction.php', {
            nombre: document.getElementById('nombre').value,
            email: document.getElementById('email').value,
            telefono: document.getElementById('telefono').value,
            mensaje: document.getElementById('mensaje').value,
        });
        document.getElementById('success-box').style.display = 'block';
        document.getElementById('contactForm').reset();
    } catch (err) {
        alertBox.innerHTML = `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
    } finally {
        btn.disabled = false;
        btn.textContent = 'Enviar Mensaje';
    }
});
