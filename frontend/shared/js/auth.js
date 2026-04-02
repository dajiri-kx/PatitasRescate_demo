const Auth = {
    KEY: 'patitas_user',

    save(cliente) {
        localStorage.setItem(this.KEY, JSON.stringify(cliente));
    },

    get() {
        const raw = localStorage.getItem(this.KEY);
        return raw ? JSON.parse(raw) : null;
    },

    clear() {
        localStorage.removeItem(this.KEY);
    },

    isLoggedIn() {
        return this.get() !== null;
    },

    requireAuth() {
        if (!this.isLoggedIn()) {
            window.location.href = nav('/auth/login/');
        }
    },

    getNombre() {
        const u = this.get();
        return u ? u.nombre : 'Usuario';
    },
};
