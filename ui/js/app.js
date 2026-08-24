const BASE = 'http://localhost:3649/api/1.0';

/* ─── API helper ─── */
async function api(path) {
    const token = localStorage.getItem('cp_token');
    const headers = token ? { 'Authorization': 'Bearer ' + token } : {};
    const r = await fetch(BASE + path, { headers });
    if (r.status === 401) {
        localStorage.removeItem('cp_token');
        localStorage.removeItem('cp_user');
        window.location.href = '/ui/auth.html';
        return;
    }
    if (!r.ok) throw new Error(`HTTP ${r.status}`);
    return r.json();
}

/* ─── Helpers ─── */
function relTime(iso) {
    if (!iso) return '';
    const m = Math.floor((Date.now() - new Date(iso)) / 60000);
    if (m < 2)    return 'az önce';
    if (m < 60)   return `${m}dk önce`;
    if (m < 1440) return `${Math.floor(m/60)}sa önce`;
    return new Date(iso).toLocaleDateString('tr-TR');
}

function pct(s) {
    const n = parseFloat(s);
    return isNaN(n) ? '' : `-%${Math.round(n)}`;
}

function getLangColor(lang) {
    const colors = {
        JavaScript:'#f1e05a', TypeScript:'#3178c6', Python:'#3572A5',
        Go:'#00ADD8', Rust:'#dea584', Java:'#b07219', 'C++':'#f34b7d',
        C:'#555555', Ruby:'#701516', Swift:'#F05138', Kotlin:'#A97BFF',
        PHP:'#4F5D95', CSS:'#563d7c', HTML:'#e34c26', Shell:'#89e051'
    };
    return colors[lang] || '#9575cd';
}

/* ─── State ─── */
let currentSection = 'exchange';
let cachedData = null;
let currentCity = 'istanbul';
let currentLang = 'tr';

/* ─── Navigation ─── */
function navigateTo(sectionId) {
    // Hide all sections
    document.querySelectorAll('.section').forEach(s => s.classList.add('hidden'));
    const target = document.getElementById('section-' + sectionId);
    if (target) target.classList.remove('hidden');

    // Update tabs
    document.querySelectorAll('.nav-tab').forEach(t => {
        t.classList.toggle('active', t.dataset.section === sectionId);
    });

    // Show/hide city+lang controls (only relevant for exchange+news)
    const hasCtrl = ['exchange', 'news'].includes(sectionId);
    const ctrl = document.getElementById('ctrlGroup');
    if (ctrl) ctrl.style.display = hasCtrl ? 'flex' : 'none';

    currentSection = sectionId;
    window.location.hash = sectionId;

    // Render with cached data if available
    if (cachedData) {
        renderSection(sectionId, cachedData);
    }
}

function onContextChange() {
    currentCity = document.getElementById('citySelect').value;
    currentLang = document.getElementById('langSelect').value;
    cachedData = null; // invalidate cache
    loadSnapshot();
}

function refresh() {
    cachedData = null;
    loadSnapshot();
}

/* ─── Loading state ─── */
function setLoading(on) {
    const btn = document.getElementById('refreshBtn');
    if (!btn) return;
    btn.disabled = on;
    const ico = btn.querySelector('.spin-ico');
    if (ico) ico.style.animation = on ? 'spin .7s linear infinite' : '';
}

/* ─── Main data load ─── */
async function loadSnapshot() {
    setLoading(true);
    const city = document.getElementById('citySelect')?.value || currentCity;
    const lang = document.getElementById('langSelect')?.value || currentLang;

    const cityLabel = document.getElementById('cityLabel');
    const cityLabelNews = document.getElementById('cityLabelNews');
    const cityEl = document.getElementById('citySelect');
    if (cityEl && cityLabel) {
        cityLabel.textContent = cityEl.options[cityEl.selectedIndex].text;
        if (cityLabelNews) cityLabelNews.textContent = cityEl.options[cityEl.selectedIndex].text;
    }

    try {
        const snap = await api(`/city/snapshot?city=${encodeURIComponent(city)}&lang=${lang}`);
        if (!snap) return;

        cachedData = snap;
        const pt = document.getElementById('pageTime');
        if (pt) pt.textContent = '🕐 ' + new Date().toLocaleTimeString('tr-TR');

        renderSection(currentSection, snap);

    } catch(e) {
        console.error(e);
        showError(`Servise ulaşılamıyor: ${e.message}`);
    } finally {
        setLoading(false);
    }
}

/* ─── Section router ─── */
function renderSection(section, snap) {
    switch(section) {
        case 'exchange': renderExchange(snap.exchange); break;
        case 'news':     renderNews(snap.news); break;
        case 'games':    renderGames(snap.gameDeals); break;
        case 'nasa':     renderNasa(snap.nasaApod); break;
        case 'github':   renderGithub(snap.githubTrend); break;
    }
}

/* ─── Render: Exchange ─── */
function renderExchange(ex) {
    const el = document.getElementById('xgrid');
    if (!el) return;
    if (!ex) {
        el.innerHTML = '<div class="empty" style="grid-column:1/-1"><span class="empty-ico">💱</span>Döviz verisi alınamadı</div>';
        return;
    }
    const r = ex.rates || {};
    const cards = [
        { flag:'🇺🇸', pair:'USD / TRY', desc:'1 Dolar = kaç TL',    v: r.TRY },
        { flag:'🇪🇺', pair:'EUR / TRY', desc:'1 Euro = kaç TL',     v: r.TRY && r.EUR ? r.TRY / r.EUR : null },
        { flag:'🇨🇭', pair:'CHF / TRY', desc:'1 Frank = kaç TL',    v: r.TRY && r.CHF ? r.TRY / r.CHF : null },
        { flag:'🇬🇧', pair:'GBP / TRY', desc:'1 Sterlin = kaç TL',  v: r.TRY && r.GBP ? r.TRY / r.GBP : null },
        { flag:'🇯🇵', pair:'JPY / TRY', desc:'100 Yen = kaç TL',    v: r.TRY && r.JPY ? (r.TRY / r.JPY) * 100 : null },
        { flag:'⚖️',  pair:'EUR / USD', desc:'1 Euro = kaç Dolar',  v: r.EUR ? 1 / r.EUR : null },
        { flag:'🪙',  pair:'XAU / USD', desc:'Altın (ons)',          v: r.XAU ? 1 / r.XAU : null },
        { flag:'💎',  pair:'BTC / USD', desc:'1 Bitcoin kaç Dolar',  v: r.BTC ? 1 / r.BTC : null },
    ];
    el.innerHTML = cards.map(c => `
        <div class="xcard">
            <div class="xcard-pair">${c.flag} ${c.pair}</div>
            <div class="xcard-rate">${c.v ? c.v.toFixed(c.pair.includes('JPY') ? 2 : 4) : '—'}</div>
            <div class="xcard-desc">${c.desc}</div>
            <div class="xcard-date">📅 ${ex.date || ''}</div>
        </div>`).join('');
}

/* ─── Render: News ─── */
function renderNews(articles) {
    const el = document.getElementById('nlist');
    if (!el) return;
    if (!articles || !articles.length) {
        el.innerHTML = '<div class="empty" style="grid-column:1/-1"><span class="empty-ico">📭</span>Haber bulunamadı</div>';
        return;
    }
    const badge = document.getElementById('newsBadge');
    if (badge) badge.textContent = articles.length + ' haber';

    el.innerHTML = articles.map(a => `
        <a class="ncard-full" href="${a.url||'#'}" target="_blank" rel="noopener">
            ${a.image ? `<img class="ncard-full-img" src="${a.image}" alt="" loading="lazy" onerror="this.style.display='none'">` : ''}
            <div class="ncard-full-body">
                <div class="ncard-full-title">${a.title||''}</div>
                <div class="ncard-full-meta">
                    <span class="nsource">${a.source?.name||'Kaynak'}</span>
                    <div class="ndot"></div>
                    <span class="ntime">${relTime(a.publishedAt)}</span>
                </div>
            </div>
        </a>`).join('');
}

/* ─── Render: Games ─── */
function renderGames(deals) {
    const el = document.getElementById('glist');
    if (!el) return;
    if (!deals || !deals.length) {
        el.innerHTML = '<div class="empty" style="grid-column:1/-1"><span class="empty-ico">🎮</span>Şu an fırsat bulunamadı</div>';
        return;
    }
    const badge = document.getElementById('gamesBadge');
    if (badge) badge.textContent = deals.length + ' fırsat';

    el.innerHTML = deals.map(d => `
        <div class="gcard">
            ${d.thumb ? `<img class="gimg" src="${d.thumb}" alt="" loading="lazy" onerror="this.style.display='none'" style="width:100%;height:120px;object-fit:cover;border-radius:6px;margin-bottom:0.8rem;background:var(--bg2);">` : ''}
            <div class="gtitle" style="font-size:.9rem;font-weight:600;margin-bottom:0.5rem;">${d.title||'—'}</div>
            <div class="gprices">
                <span class="gprice-s">$${parseFloat(d.salePrice||0).toFixed(2)}</span>
                ${d.normalPrice ? `<span class="gprice-o">$${parseFloat(d.normalPrice).toFixed(2)}</span>` : ''}
                ${d.savingsPercent ? `<span class="gsave">${pct(d.savingsPercent)}</span>` : ''}
            </div>
            ${d.steamRating ? `<div class="grating" style="margin-top:0.4rem;">⭐ ${d.steamRating}</div>` : ''}
        </div>`).join('');
}

/* ─── Render: NASA ─── */
function renderNasa(apod) {
    const el = document.getElementById('nasaContent');
    if (!el) return;
    if (!apod) {
        el.innerHTML = '<div class="empty"><span class="empty-ico">🌌</span>NASA verisi alınamadı.<br>API anahtarı geçerli mi?</div>';
        return;
    }

    const isVideo = apod.mediaType === 'video';
    const mediaHtml = isVideo
        ? `<iframe src="${apod.url}" style="width:100%;height:100%;min-height:450px;border:none;" allowfullscreen></iframe>`
        : `<img src="${apod.url}" alt="${apod.title}" style="width:100%;height:100%;max-height:75vh;object-fit:contain;" loading="lazy">`;

    el.innerHTML = `
        <div class="nasa-page">
            <div class="nasa-media">${mediaHtml}</div>
            <div class="nasa-info">
                <div class="nasa-badge">🚀 Astronomy Picture of the Day</div>
                <h2 class="nasa-title">${apod.title||''}</h2>
                <div class="nasa-date">🗓️ ${apod.date||''}</div>
                <p class="nasa-desc">${apod.explanation||''}</p>
                ${apod.copyright ? `<div class="nasa-credit">📷 © ${apod.copyright}</div>` : ''}
                <a href="${apod.url}" target="_blank" rel="noopener"
                   style="display:inline-flex;align-items:center;gap:.4rem;padding:.6rem 1.2rem;background:var(--primary);color:var(--dark);font-family:'Aldrich',sans-serif;font-size:.75rem;text-transform:uppercase;text-decoration:none;border-radius:7px;border:2px solid var(--dark);box-shadow:2px 2px 0 var(--dark);transition:transform .15s,box-shadow .15s;font-weight:bold;"
                   onmouseover="this.style.transform='translate(-2px,-2px)';this.style.boxShadow='4px 4px 0 var(--dark)'"
                   onmouseout="this.style.transform='';this.style.boxShadow='2px 2px 0 var(--dark)'">
                   🔗 Tam Boyut Görüntüle
                </a>
            </div>
        </div>`;
}

/* ─── Render: GitHub ─── */
function renderGithub(repos) {
    const el = document.getElementById('ghgrid');
    if (!el) return;
    if (!repos || !repos.length) {
        el.innerHTML = '<div class="empty" style="grid-column:1/-1"><span class="empty-ico">🐙</span>Repo bulunamadı</div>';
        return;
    }
    el.innerHTML = repos.map(r => `
        <a class="gh-card" href="${r.url||'#'}" target="_blank" rel="noopener">
            <div class="gh-card-head">
                <img class="gh-avatar" src="${r.avatarURL||''}" alt="" onerror="this.style.display='none'">
                <span class="gh-name">${r.fullName||''}</span>
            </div>
            <p class="gh-desc">${r.description||'Açıklama yok'}</p>
            <div class="gh-footer">
                ${r.language ? `
                <div class="gh-lang">
                    <span class="gh-lang-dot" style="background:${getLangColor(r.language)};"></span>
                    <span>${r.language}</span>
                </div>` : ''}
                <span class="gh-stars">⭐ ${(r.stars||0).toLocaleString()}</span>
            </div>
        </a>`).join('');
}

/* ─── Error state ─── */
function showError(msg) {
    const el = document.getElementById('xgrid');
    if (el) el.innerHTML = `
        <div class="empty" style="grid-column:1/-1">
            <span class="empty-ico">⚡</span>
            <strong>Bağlantı hatası</strong>
            <code>${msg}</code>
            Servis çalışıyor mu? → <code>go run ./main.go --dev</code>
        </div>`;
}

/* ─── Init ─── */
document.addEventListener('DOMContentLoaded', () => {
    // Read hash for initial section
    const hash = window.location.hash.replace('#', '') || 'exchange';
    const validSections = ['exchange', 'news', 'games', 'nasa', 'github'];
    const initSection = validSections.includes(hash) ? hash : 'exchange';

    // Navigate without triggering data load yet
    document.querySelectorAll('.section').forEach(s => s.classList.add('hidden'));
    const target = document.getElementById('section-' + initSection);
    if (target) target.classList.remove('hidden');
    document.querySelectorAll('.nav-tab').forEach(t => {
        t.classList.toggle('active', t.dataset.section === initSection);
    });
    currentSection = initSection;
    const hasCtrl = ['exchange', 'news'].includes(initSection);
    const ctrl = document.getElementById('ctrlGroup');
    if (ctrl) ctrl.style.display = hasCtrl ? 'flex' : 'none';

    // Load data
    loadSnapshot();
});
