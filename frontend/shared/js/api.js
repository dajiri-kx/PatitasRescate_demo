const CONFIG = {
    API_BASE: '/api',
    FRONTEND_BASE: '/frontend/features',
};

async function apiFetch(featurePath, options = {}) {
    const url = CONFIG.API_BASE + featurePath;
    const defaults = {
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
    };
    const opts = { ...defaults, ...options };

    const res = await fetch(url, opts);
    const data = await res.json();

    if (!res.ok) {
        if (res.status === 401) {
            Auth.clear();
            window.location.href = CONFIG.FRONTEND_BASE + '/auth/login/';
            return;
        }
        throw new Error(data.error || 'Error desconocido');
    }
    return data;
}

async function apiGet(featurePath) {
    return apiFetch(featurePath, { method: 'GET' });
}

async function apiPost(featurePath, body) {
    return apiFetch(featurePath, {
        method: 'POST',
        body: JSON.stringify(body),
    });
}

function nav(featurePath) {
    return CONFIG.FRONTEND_BASE + featurePath;
}
