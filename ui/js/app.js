const BASE = 'http://localhost:3649/api/1.0';
const KEY  = 'local-dev-key';

async function api(path) {
    const r = await fetch(BASE + path, { headers: { 'X-API-Key': KEY } });
    if (!r.ok) throw new Error(`HTTP ${r.status}`);
    return r.json();
}

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

/* ─── Render: Exchange ─── */
function renderExchange(ex) {
    if (!ex) {
        document.getElementById('xgrid').innerHTML =
            '<div class="empty" style="grid-column:1/-1"><span class="empty-ico">💱</span>Döviz verisi alınamadı</div>';
        return;
    }
    const r = ex.rates || {};
    // base = USD → r.TRY, r.EUR, r.GBP, r.CHF, r.JPY
    const cards = [
        { flag:'🇺🇸', pair:'USD / TRY', desc:'1 Dolar kaç TL',   v: r.TRY },
        { flag:'🇪🇺', pair:'EUR / TRY', desc:'1 Euro kaç TL',    v: r.TRY && r.EUR ? r.TRY / r.EUR : null },
        { flag:'🇨🇭', pair:'CHF / TRY', desc:'1 Frank kaç TL',   v: r.TRY && r.CHF ? r.TRY / r.CHF : null },
        { flag:'⚖️',  pair:'EUR / USD', desc:'1 Euro kaç Dolar', v: r.EUR ? 1 / r.EUR : null },
    ];
    document.getElementById('xgrid').innerHTML = cards.map(c => `
        <div class="xcard">
            <div class="xcard-pair">${c.flag} ${c.pair}</div>
            <div class="xcard-rate">${c.v ? c.v.toFixed(4) : '—'}</div>
            <div class="xcard-desc">${c.desc}</div>
            <div class="xcard-date">📅 ${ex.date || ''}</div>
        </div>`).join('');
}

/* ─── Render: News ─── */
function renderNews(articles) {
    const el = document.getElementById('nlist');
    if (!articles || !articles.length) {
        el.innerHTML = '<div class="empty"><span class="empty-ico">📭</span>Haber bulunamadı</div>';
        return;
    }
    document.getElementById('newsBadge').textContent = articles.length + ' haber';
    el.innerHTML = articles.map(a => `
        <a class="ncard" href="${a.url||'#'}" target="_blank" rel="noopener">
            ${a.image ? `<img class="nimg" src="${a.image}" alt="" loading="lazy" onerror="this.style.display='none'">` : ''}
            <div class="nbody">
                <div class="ntitle">${a.title||''}</div>
                <div class="nmeta">
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
    if (!deals || !deals.length) {
        el.innerHTML = '<div class="empty"><span class="empty-ico">🎮</span>Şu an fırsat bulunamadı</div>';
        return;
    }
    document.getElementById('gamesBadge').textContent = deals.length + ' fırsat';
    el.innerHTML = deals.map(d => `
        <div class="gcard">
            ${d.thumb ? `<img class="gimg" src="${d.thumb}" alt="" loading="lazy" onerror="this.style.display='none'">` : '<div class="gimg"></div>'}
            <div class="gbody">
                <div class="gtitle">${d.title||'—'}</div>
                <div class="gprices">
                    <span class="gprice-s">$${parseFloat(d.salePrice||0).toFixed(2)}</span>
                    ${d.normalPrice ? `<span class="gprice-o">$${parseFloat(d.normalPrice).toFixed(2)}</span>` : ''}
                    ${d.savingsPercent ? `<span class="gsave">${pct(d.savingsPercent)}</span>` : ''}
                </div>
                ${d.steamRating ? `<div class="grating">⭐ ${d.steamRating}</div>` : ''}
            </div>
        </div>`).join('');
}

/* ─── Render: NASA ─── */
function renderNasa(apod) {
    const el = document.getElementById('nasaBox');
    if (!apod) {
        el.innerHTML = '<div class="empty"><span class="empty-ico">🌌</span>NASA verisi alınamadı</div>';
        return;
    }
    const media = apod.mediaType === 'video'
        ? `<iframe src="${apod.url}" style="width:100%; height:200px; border:none; border-radius:8px; margin-bottom:1rem;"></iframe>`
        : `<img src="${apod.url}" alt="${apod.title}" style="width:100%; height:auto; max-height:250px; object-fit:cover; border-radius:8px; margin-bottom:1rem;">`;
        
    el.innerHTML = `
        ${media}
        <h3 style="font-size: .9rem; font-weight: 600; margin-bottom: .4rem;">${apod.title}</h3>
        <p style="font-size: .75rem; color: var(--t2); line-height: 1.5; display: -webkit-box; -webkit-line-clamp: 4; -webkit-box-orient: vertical; overflow: hidden;">${apod.explanation}</p>
        <div style="font-size: .65rem; color: var(--t3); margin-top: .6rem;">🗓️ ${apod.date}</div>
    `;
}

/* ─── Render: GitHub ─── */
function renderGithub(repos) {
    const el = document.getElementById('ghlist');
    if (!repos || !repos.length) {
        el.innerHTML = '<div class="empty"><span class="empty-ico">🐙</span>Repo bulunamadı</div>';
        return;
    }
    el.innerHTML = repos.map(r => `
        <a class="ncard" href="${r.url}" target="_blank" rel="noopener">
            <img class="nimg" src="${r.avatarURL}" alt="" loading="lazy" onerror="this.style.display='none'">
            <div class="nbody">
                <div class="ntitle">${r.fullName}</div>
                <p style="font-size: .7rem; color: var(--t2); margin-bottom: .4rem; display: -webkit-box; -webkit-line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden;">${r.description || 'Açıklama yok'}</p>
                <div class="nmeta">
                    <span class="nsource">⭐ ${r.stars}</span>
                    ${r.language ? `<div class="ndot"></div><span class="ntime">💻 ${r.language}</span>` : ''}
                </div>
            </div>
        </a>`).join('');
}

/* ─── Loading state ─── */
function setLoading(on) {
    const btn = document.getElementById('refreshBtn');
    btn.disabled = on;
    btn.classList.toggle('loading', on);
}

/* ─── Main ─── */
async function loadData() {
    setLoading(true);
    const cityEl = document.getElementById('citySelect');
    const city   = cityEl.value;
    const lang   = document.getElementById('langSelect').value;
    document.getElementById('cityLabel').textContent =
        cityEl.options[cityEl.selectedIndex].text;

    try {
        const snap = await api(`/city/snapshot?city=${encodeURIComponent(city)}&lang=${lang}`);
        document.getElementById('pageTime').textContent =
            '🕐 ' + new Date().toLocaleTimeString('tr-TR');
        renderExchange(snap.exchange);
        renderNews(snap.news);
        renderGames(snap.gameDeals);
        renderNasa(snap.nasaApod);
        renderGithub(snap.githubTrend);
    } catch (e) {
        console.error(e);
        document.getElementById('xgrid').innerHTML = `
            <div class="empty" style="grid-column:1/-1">
                <span class="empty-ico">⚡</span>
                <strong>Servise ulaşılamıyor</strong>
                <code>${e.message}</code>
                Servis çalışıyor mu? →&nbsp;<code>go run ./main.go --dev</code>
            </div>`;
        document.getElementById('nlist').innerHTML = '';
        document.getElementById('glist').innerHTML = '';
        document.getElementById('nasaBox').innerHTML = '';
        document.getElementById('ghlist').innerHTML = '';
    } finally {
        setLoading(false);
    }
}

// Initial load
document.addEventListener("DOMContentLoaded", () => {
    loadData();
});
