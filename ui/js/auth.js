const BASE = 'http://localhost:3649/api/1.0';

(function() {
    const token = localStorage.getItem('cp_token');
    if (token) {
        try {
            const payload = JSON.parse(atob(token.split('.')[1]));
            if (payload.exp * 1000 > Date.now()) {
                window.location.href = '/ui/dashboard.html';
                return;
            }
        } catch(e) {}
        localStorage.removeItem('cp_token');
        localStorage.removeItem('cp_user');
    }
})();

let activeTab = 'login';

function switchTab(tab) {
    activeTab = tab;
    document.getElementById('loginForm').classList.toggle('hidden', tab !== 'login');
    document.getElementById('registerForm').classList.toggle('hidden', tab !== 'register');
    document.getElementById('tabLogin').classList.toggle('tab-active', tab === 'login');
    document.getElementById('tabRegister').classList.toggle('tab-active', tab === 'register');
    clearMessages();
}

function showMsg(id, text, isError) {
    const el = document.getElementById(id);
    el.textContent = text;
    el.className = 'form-msg ' + (isError ? 'form-msg-error' : 'form-msg-ok');
    el.style.display = 'block';
}

function clearMessages() {
    ['loginMsg', 'registerMsg'].forEach(id => {
        const el = document.getElementById(id);
        el.style.display = 'none';
        el.textContent = '';
    });
}

async function handleLogin(e) {
    e.preventDefault();
    const btn = document.getElementById('loginBtn');
    btn.disabled = true;
    btn.textContent = 'Giriş yapılıyor...';

    const email = document.getElementById('loginEmail').value.trim();
    const password = document.getElementById('loginPassword').value;

    try {
        const res = await fetch(`${BASE}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        const data = await res.json();

        if (!res.ok) {
            showMsg('loginMsg', data.error || 'Giriş başarısız', true);
            return;
        }

        localStorage.setItem('cp_token', data.token);
        localStorage.setItem('cp_user', JSON.stringify(data.user));
        window.location.href = '/ui/dashboard.html';

    } catch(err) {
        showMsg('loginMsg', 'Sunucuya bağlanılamadı', true);
    } finally {
        btn.disabled = false;
        btn.textContent = 'Giriş Yap';
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const btn = document.getElementById('registerBtn');
    btn.disabled = true;
    btn.textContent = 'Kayıt yapılıyor...';

    const email = document.getElementById('regEmail').value.trim();
    const username = document.getElementById('regUsername').value.trim();
    const password = document.getElementById('regPassword').value;
    const confirm = document.getElementById('regConfirm').value;

    if (password !== confirm) {
        showMsg('registerMsg', 'Şifreler eşleşmiyor', true);
        btn.disabled = false;
        btn.textContent = 'Kayıt Ol';
        return;
    }

    if (password.length < 6) {
        showMsg('registerMsg', 'Şifre en az 6 karakter olmalı', true);
        btn.disabled = false;
        btn.textContent = 'Kayıt Ol';
        return;
    }

    try {
        const res = await fetch(`${BASE}/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, username, password })
        });

        const data = await res.json();

        if (!res.ok) {
            showMsg('registerMsg', data.error || 'Kayıt başarısız', true);
            return;
        }

        showMsg('registerMsg', 'Kayıt başarılı! Giriş yapılıyor...', false);
        setTimeout(() => {
            document.getElementById('loginEmail').value = email;
            switchTab('login');
        }, 1200);

    } catch(err) {
        showMsg('registerMsg', 'Sunucuya bağlanılamadı', true);
    } finally {
        btn.disabled = false;
        btn.textContent = 'Kayıt Ol';
    }
}

document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('loginFormEl').addEventListener('submit', handleLogin);
    document.getElementById('registerFormEl').addEventListener('submit', handleRegister);
});
