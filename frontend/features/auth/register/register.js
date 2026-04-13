document.addEventListener('DOMContentLoaded', function () {
    const provinciaSelect = document.getElementById('provincia');
    const cantonSelect = document.getElementById('canton');
    const distritoSelect = document.getElementById('distrito');
    const alertBox = document.getElementById('alert-box');

    // --- Input filters: block non-digits on identificacion & telefono ---
    const idInput = document.getElementById('identificacion');
    const telInput = document.getElementById('telefono');
    [idInput, telInput].forEach(el => {
        el.addEventListener('input', () => { el.value = el.value.replace(/\D/g, ''); });
    });

    // --- Password requirements live feedback ---
    const pwInput = document.getElementById('password');
    const pwReq = document.getElementById('pw-requirements');
    const rules = [
        { test: v => v.length >= 8,          label: 'Mínimo 8 caracteres' },
        { test: v => /[A-Z]/.test(v),        label: 'Una letra mayúscula' },
        { test: v => /[a-z]/.test(v),        label: 'Una letra minúscula' },
        { test: v => /\d/.test(v),            label: 'Un número' },
        { test: v => /[^A-Za-z0-9]/.test(v), label: 'Un carácter especial (!@#$...)' },
    ];

    function renderPwRules() {
        const val = pwInput.value;
        pwReq.innerHTML = rules.map(r => {
            const ok = r.test(val);
            return `<span style="color:${ok ? '#198754' : '#6c757d'};">${ok ? '✓' : '○'} ${r.label}</span>`;
        }).join('<br>');
    }
    renderPwRules();
    pwInput.addEventListener('input', renderPwRules);

    function passwordIsValid() {
        return rules.every(r => r.test(pwInput.value));
    }

    // Bootstrap validation
    const forms = document.querySelectorAll('.needs-validation');
    forms.forEach(function (form) {
        form.addEventListener('submit', function (event) {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            }
            form.classList.add('was-validated');
        }, false);
    });

    // Load provinces
    fetch('https://api-geo-cr.vercel.app/provincias')
        .then(r => r.json())
        .then(data => {
            if (Array.isArray(data.data)) {
                data.data.forEach(p => {
                    const opt = document.createElement('option');
                    opt.value = p.idProvincia;
                    opt.textContent = p.descripcion;
                    provinciaSelect.appendChild(opt);
                });
            }
        });

    provinciaSelect.addEventListener('change', function () {
        const id = this.value;
        cantonSelect.innerHTML = '<option value="">Seleccione un cantón</option>';
        distritoSelect.innerHTML = '<option value="">Seleccione un distrito</option>';
        if (id) {
            fetch('https://api-geo-cr.vercel.app/provincias/' + id + '/cantones?limit=200')
                .then(r => r.json())
                .then(data => {
                    if (Array.isArray(data.data)) {
                        data.data.forEach(c => {
                            const opt = document.createElement('option');
                            opt.value = c.idCanton;
                            opt.textContent = c.descripcion;
                            cantonSelect.appendChild(opt);
                        });
                    }
                });
        }
    });

    cantonSelect.addEventListener('change', function () {
        const id = this.value;
        distritoSelect.innerHTML = '<option value="">Seleccione un distrito</option>';
        if (id) {
            fetch('https://api-geo-cr.vercel.app/cantones/' + id + '/distritos?limit=200')
                .then(r => r.json())
                .then(data => {
                    if (Array.isArray(data.data)) {
                        data.data.forEach(d => {
                            const opt = document.createElement('option');
                            opt.value = d.idDistrito;
                            opt.textContent = d.descripcion;
                            distritoSelect.appendChild(opt);
                        });
                    }
                });
        }
    });

    // Form submission
    document.getElementById('registerForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const form = e.target;
        form.classList.add('was-validated');
        if (!form.checkValidity()) return;

        // Custom checks beyond HTML5
        const errors = [];
        if (!/^\d{9}$/.test(idInput.value)) errors.push('La identificación debe ser exactamente 9 dígitos.');
        if (!/^\d{8}$/.test(telInput.value)) errors.push('El teléfono debe ser exactamente 8 dígitos.');
        if (!passwordIsValid()) errors.push('La contraseña no cumple todos los requisitos.');
        if (pwInput.value !== document.getElementById('confirmPassword').value) errors.push('Las contraseñas no coinciden.');

        if (errors.length) {
            alertBox.innerHTML = errors.map(m => `<div class="alert alert-danger py-2">${m}</div>`).join('');
            return;
        }

        alertBox.innerHTML = '';
        const btn = document.getElementById('btnSubmit');
        btn.disabled = true;
        btn.textContent = 'Registrando...';

        const body = {
            identificacion: document.getElementById('identificacion').value,
            nombre: document.getElementById('nombre').value,
            primerApellido: document.getElementById('primerApellido').value,
            correo: document.getElementById('correo').value,
            telefono: document.getElementById('telefono').value,
            password: document.getElementById('password').value,
            confirmPassword: document.getElementById('confirmPassword').value,
            direccionSennas: document.getElementById('direccionSennas').value,
        };

        try {
            await apiPost('/auth/register', body);
            window.location.href = nav('/auth/login/') + '?registered=1';
        } catch (err) {
            alertBox.innerHTML = `<div class="alert alert-danger">${escapeHtml(err.message)}</div>`;
        } finally {
            btn.disabled = false;
            btn.textContent = 'Registrar';
        }
    });
});
