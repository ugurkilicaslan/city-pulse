const BASE = 'http://localhost:3649/api/1.0';

/* ─── API Helper ─── */
async function api(path, options = {}) {
    const token = localStorage.getItem('cp_token');
    const headers = {
        'Content-Type': 'application/json',
        ...(token ? { 'Authorization': 'Bearer ' + token } : {}),
        ...(options.headers || {})
    };

    try {
        const r = await fetch(BASE + path, { ...options, headers });
        if (r.status === 401) {
            localStorage.removeItem('cp_token');
            localStorage.removeItem('cp_user');
            window.location.href = '/ui/auth.html';
            return null;
        }
        if (!r.ok) {
            const errData = await r.json().catch(() => ({}));
            throw new Error(errData.error || `HTTP ${r.status}`);
        }
        return await r.json();
    } catch (err) {
        console.error('API Hatası:', path, err);
        throw err;
    }
}

/* ─── Analytics Tracker ─── */
function trackEvent(type, section = '', metadata = {}) {
    try {
        api('/analytics/track', {
            method: 'POST',
            body: JSON.stringify({
                type: type,
                page: window.location.pathname,
                section: section,
                metadata: metadata
            })
        }).catch(() => {});
    } catch(e) {}
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
let currentBookmarkFilter = 'all';
let cachedBookmarks = [];

/* ─── Navigation ─── */
function navigateTo(sectionId) {
    document.querySelectorAll('.section').forEach(s => s.classList.add('hidden'));
    const target = document.getElementById('section-' + sectionId);
    if (target) target.classList.remove('hidden');

    document.querySelectorAll('.nav-tab').forEach(t => {
        t.classList.toggle('active', t.dataset.section === sectionId);
    });

    const hasCtrl = ['exchange', 'news', 'weather'].includes(sectionId);
    const ctrl = document.getElementById('ctrlGroup');
    if (ctrl) ctrl.style.display = hasCtrl ? 'flex' : 'none';

    currentSection = sectionId;
    window.location.hash = sectionId;

    trackEvent('section_view', sectionId);

    if (sectionId === 'bookmarks') {
        loadBookmarks();
    } else if (sectionId === 'alerts') {
        loadAlerts();
    } else if (cachedData) {
        renderSection(sectionId, cachedData);
    }
}

function onContextChange() {
    currentCity = document.getElementById('citySelect').value;
    currentLang = document.getElementById('langSelect').value;
    cachedData = null;
    loadSnapshot();
}

function refresh() {
    cachedData = null;
    trackEvent('data_refresh', currentSection);
    loadSnapshot();
    if (currentSection === 'bookmarks') loadBookmarks();
    if (currentSection === 'alerts') loadAlerts();
}

function setLoading(on) {
    const btn = document.getElementById('refreshBtn');
    if (!btn) return;
    btn.disabled = on;
    const ico = btn.querySelector('.spin-ico');
    if (ico) ico.style.animation = on ? 'spin .7s linear infinite' : '';
}

/* ─── Main Data Load ─── */
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
        const [snap, weather, crypto] = await Promise.all([
            api(`/city/snapshot?city=${encodeURIComponent(city)}&lang=${lang}`),
            api(`/weather/current?city=${encodeURIComponent(city)}`),
            api(`/crypto/prices`)
        ]);
        if (!snap) return;

        cachedData = { ...snap, weather, crypto };
        const pt = document.getElementById('pageTime');
        if (pt) pt.textContent = '⏱ ' + new Date().toLocaleTimeString('tr-TR');

        renderSection(currentSection, cachedData);

    } catch (err) {
        showError(err.message);
    } finally {
        setLoading(false);
    }
}

function renderSection(sec, data) {
    if (!data) return;
    switch (sec) {
        case 'exchange': renderExchange(data.exchange); break;
        case 'news':     renderNews(data.news);         break;
        case 'games':    renderGames(data.games);       break;
        case 'nasa':     renderNasa(data.nasa);         break;
        case 'github':   renderGithub(data.github);     break;
        case 'weather':  renderWeather(data.weather);   break;
        case 'crypto':   renderCrypto(data.crypto);     break;
    }
}

/* ─── Render: Exchange ─── */
function renderExchange(ex) {
    const el = document.getElementById('xgrid');
    if (!el) return;
    if (!ex || !ex.rates || !ex.rates.length) {
        el.innerHTML = '<div class="empty">Döviz verisi bulunamadı</div>';
        return;
    }
    el.innerHTML = ex.rates.map(r => `
        <div class="xcard">
            <div class="xcard-pair">${r.pair}</div>
            <div class="xcard-rate">${r.rate ? r.rate.toFixed(4) : '—'}</div>
            <div class="xcard-sub">Baz: ${ex.base || 'TRY'}</div>
        </div>`).join('');
}

/* ─── Render: News ─── */
function renderNews(news) {
    const el = document.getElementById('nlist');
    const badge = document.getElementById('newsBadge');
    if (badge) badge.textContent = (news && news.articles) ? `${news.articles.length} Haber` : '';
    if (!el) return;
    if (!news || !news.articles || !news.articles.length) {
        el.innerHTML = '<div class="empty" style="grid-column:1/-1">Haber bulunamadı</div>';
        return;
    }
    el.innerHTML = news.articles.map(a => `
        <div class="news-card">
            ${a.image ? `<img class="news-img" src="${a.image}" alt="" onerror="this.style.display='none'">` : ''}
            <div class="news-body">
                <div class="news-src">${a.source || ''} • ${relTime(a.publishedAt)}</div>
                <div class="news-title">${a.title || ''}</div>
                <div class="news-desc">${a.description || ''}</div>
                <div style="display:flex; justify-content:space-between; align-items:center; margin-top:0.8rem;">
                    <a class="news-link" href="${a.url}" target="_blank" rel="noopener">Haberi Oku →</a>
                    <button class="bookmark-btn" onclick="saveBookmark('news', '${encodeURIComponent(a.title || '')}', '${encodeURIComponent(a.url || '')}', '${encodeURIComponent(a.source || '')}')">
                        ★ Kaydet
                    </button>
                </div>
            </div>
        </div>`).join('');
}

/* ─── Render: Games ─── */
function renderGames(games) {
    const el = document.getElementById('glist');
    const badge = document.getElementById('gamesBadge');
    if (badge) badge.textContent = (games && games.deals) ? `${games.deals.length} Fırsat` : '';
    if (!el) return;
    if (!games || !games.deals || !games.deals.length) {
        el.innerHTML = '<div class="empty" style="grid-column:1/-1">Oyun fırsatı bulunamadı</div>';
        return;
    }
    el.innerHTML = games.deals.map(g => `
        <div class="game-card">
            ${g.thumb ? `<img class="game-thumb" src="${g.thumb}" alt="" onerror="this.style.display='none'">` : ''}
            <div class="game-body">
                <div class="game-title">${g.title || ''}</div>
                <div class="game-prices">
                    <span class="game-sale">$${g.salePrice || '0'}</span>
                    ${g.normalPrice ? `<span class="game-norm">$${g.normalPrice}</span>` : ''}
                    ${g.savings ? `<span class="game-savings">${pct(g.savings)}</span>` : ''}
                </div>
                <div style="display:flex; justify-content:space-between; align-items:center; margin-top:0.6rem;">
                    <a class="game-deal-link" href="https://www.cheapshark.com/redirect?dealID=${g.dealID}" target="_blank" rel="noopener">Fırsata Git →</a>
                    <button class="bookmark-btn" onclick="saveBookmark('game', '${encodeURIComponent(g.title || '')}', 'https://www.cheapshark.com/redirect?dealID=${g.dealID}', '${encodeURIComponent(g.salePrice || '')}')">
                        ★ Kaydet
                    </button>
                </div>
            </div>
        </div>`).join('');
}

/* ─── Render: NASA ─── */
function renderNasa(apod) {
    const el = document.getElementById('nasaContent');
    if (!el) return;
    if (!apod || !apod.url) {
        el.innerHTML = '<div class="empty">NASA verisi alınamadı</div>';
        return;
    }

    const isVideo = apod.mediaType === 'video';
    const mediaHtml = isVideo
        ? `<iframe class="nasa-video" src="${apod.url}" frameborder="0" allowfullscreen></iframe>`
        : `<img class="nasa-img" src="${apod.hdUrl || apod.url}" alt="${apod.title||''}" loading="lazy">`;

    el.innerHTML = `
        <div class="nasa-page">
            <div class="nasa-media">${mediaHtml}</div>
            <div class="nasa-info">
                <div class="nasa-badge">🚀 Astronomy Picture of the Day</div>
                <h2 class="nasa-title">${apod.title||''}</h2>
                <div class="nasa-date">📅 ${apod.date||''}</div>
                <p class="nasa-desc">${apod.explanation||''}</p>
                <div style="display:flex; gap:1rem; align-items:center; margin-top:1rem;">
                    <a href="${apod.url}" target="_blank" rel="noopener"
                       style="display:inline-flex;align-items:center;gap:.4rem;padding:.6rem 1.2rem;background:var(--primary);color:var(--dark);font-family:'Aldrich',sans-serif;font-size:.75rem;text-transform:uppercase;text-decoration:none;border-radius:7px;border:2px solid var(--dark);box-shadow:2px 2px 0 var(--dark);font-weight:bold;">
                       🔍 Tam Boyut Görüntüle
                    </a>
                    <button class="bookmark-btn" style="padding:0.6rem 1.2rem; font-size:0.75rem;" onclick="saveBookmark('other', '${encodeURIComponent(apod.title || 'NASA APOD')}', '${encodeURIComponent(apod.url)}', '${encodeURIComponent(apod.date)}')">
                        ★ Favorilere Ekle
                    </button>
                </div>
            </div>
        </div>`;
}

/* ─── Render: GitHub ─── */
function renderGithub(repos) {
    const el = document.getElementById('ghgrid');
    if (!el) return;
    if (!repos || !repos.length) {
        el.innerHTML = '<div class="empty" style="grid-column:1/-1">Repo bulunamadı</div>';
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
                <span class="gh-stars">★ ${(r.stars||0).toLocaleString()}</span>
            </div>
        </a>`).join('');
}

/* ─── Render: Weather ─── */
function renderWeather(w) {
    const el = document.getElementById('weatherContent');
    const cityLbl = document.getElementById('cityLabelWeather');
    if (cityLbl && w) cityLbl.textContent = w.city || 'Şehir';
    if (!el) return;
    if (!w || w.error) {
        el.innerHTML = '<div class="empty">Hava durumu alınamadı</div>';
        return;
    }
    
    el.innerHTML = `
        <div style="background:var(--card); border:1px solid var(--bdr); border-radius:12px; padding:2rem; width:100%; max-width:400px; text-align:center; box-shadow:0 8px 30px rgba(0,0,0,0.15);">
            <div style="font-size:4rem; margin-bottom:1rem; animation: geo-float 3s infinite alternate;">${w.emoji}</div>
            <div style="font-family:'Aldrich',sans-serif; font-size:2.5rem; color:var(--t1); font-weight:bold; margin-bottom:0.5rem;">${w.temperature}°C</div>
            <div style="font-size:1.1rem; color:var(--primary); font-weight:600; margin-bottom:1.5rem; text-transform:uppercase;">${w.condition}</div>
            <div style="display:flex; justify-content:space-around; border-top:1px solid var(--bdr); padding-top:1.5rem;">
                <div>
                    <div style="font-size:0.75rem; color:var(--t3); text-transform:uppercase; margin-bottom:0.3rem;">Rüzgar</div>
                    <div style="font-size:1rem; color:var(--t1); font-weight:bold;">${w.windSpeed} km/s</div>
                </div>
                <div>
                    <div style="font-size:0.75rem; color:var(--t3); text-transform:uppercase; margin-bottom:0.3rem;">Vakit</div>
                    <div style="font-size:1rem; color:var(--t1); font-weight:bold;">${w.isDay ? 'Gündüz ☀️' : 'Gece 🌙'}</div>
                </div>
            </div>
        </div>`;
}

/* ─── Render: Crypto ─── */
function renderCrypto(data) {
    const el = document.getElementById('cryptolist');
    if (!el) return;
    if (!data || !data.coins || !data.coins.length) {
        el.innerHTML = '<div class="empty" style="grid-column:1/-1">Kripto verisi alınamadı</div>';
        return;
    }
    
    el.innerHTML = data.coins.map(c => {
        const isUp = c.change24hUsd >= 0;
        const clr = isUp ? 'var(--mint)' : '#fb7185';
        const sign = isUp ? '+' : '';
        return `
        <div style="background:var(--card); border:1px solid var(--bdr); border-radius:10px; padding:1.2rem; display:flex; flex-direction:column; gap:0.8rem; border-left:4px solid var(--primary);">
            <div style="display:flex; align-items:center; gap:0.8rem;">
                <img src="${c.logoUrl}" style="width:32px; height:32px; border-radius:50%; background:#fff; padding:2px;">
                <div style="flex:1;">
                    <div style="font-family:'Aldrich',sans-serif; font-size:1rem; color:var(--t1); font-weight:bold;">${c.name}</div>
                    <div style="font-size:0.7rem; color:var(--t3); text-transform:uppercase;">${c.symbol}</div>
                </div>
                <button class="bookmark-btn" onclick="saveBookmark('crypto', '${encodeURIComponent(c.name)}', 'https://coingecko.com', '$${c.priceUsd}')">
                    ★
                </button>
            </div>
            <div style="display:flex; align-items:flex-end; justify-content:space-between; margin-top:0.5rem;">
                <div style="font-family:'Aldrich',sans-serif; font-size:1.4rem; font-weight:bold; color:var(--t1);">$${c.priceUsd.toLocaleString(undefined, {minimumFractionDigits:2, maximumFractionDigits:2})}</div>
                <div style="font-family:'Inter',sans-serif; font-size:0.85rem; font-weight:600; color:${clr}; background:rgba(0,0,0,0.2); padding:0.2rem 0.5rem; border-radius:6px;">${sign}${c.change24hUsd.toFixed(2)}%</div>
            </div>
        </div>`;
    }).join('');
}

/* ─── Bookmarks (Favoriler) ─── */
async function saveBookmark(type, encodedTitle, encodedUrl, extra) {
    const title = decodeURIComponent(encodedTitle);
    const url = decodeURIComponent(encodedUrl);
    try {
        await api('/me/bookmarks', {
            method: 'POST',
            body: JSON.stringify({
                type: type,
                title: title,
                url: url,
                metadata: { extra: decodeURIComponent(extra || '') }
            })
        });
        trackEvent('bookmark_add', type, { title });
        alert(`"${title}" başarıyla favorilere eklendi!`);
    } catch (err) {
        alert('Favori eklenirken hata: ' + err.message);
    }
}

async function loadBookmarks(filter = currentBookmarkFilter) {
    currentBookmarkFilter = filter;
    const el = document.getElementById('bookmarkList');
    const badge = document.getElementById('bookmarkCountBadge');
    if (!el) return;
    el.innerHTML = '<div class="sk sk-c" style="grid-column:1/-1;"></div>';

    try {
        const query = filter && filter !== 'all' ? `?type=${filter}` : '';
        const res = await api(`/me/bookmarks${query}`);
        cachedBookmarks = res.bookmarks || [];
        if (badge) badge.textContent = `${cachedBookmarks.length} Kayıt`;

        if (!cachedBookmarks.length) {
            el.innerHTML = '<div class="empty" style="grid-column:1/-1;">Henüz kaydedilmiş yer iminiz yok. Haber, oyun veya kripto kartlarındaki "★ Kaydet" butonunu kullanarak ekleyebilirsiniz.</div>';
            return;
        }

        el.innerHTML = cachedBookmarks.map(b => `
            <div class="news-card" style="border-left:3px solid var(--primary);">
                <div class="news-body">
                    <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:0.5rem;">
                        <span class="col-badge" style="text-transform:uppercase;">${b.type}</span>
                        <button onclick="deleteBookmark('${b.id}')" style="background:rgba(255,0,0,0.2); color:#ff6b6b; border:1px solid rgba(255,0,0,0.4); border-radius:6px; padding:0.2rem 0.5rem; font-size:0.75rem; cursor:pointer;">
                            Sil ✕
                        </button>
                    </div>
                    <div class="news-title" style="font-size:1rem;">${b.title}</div>
                    <div style="margin-top:1rem;">
                        <a class="news-link" href="${b.url}" target="_blank" rel="noopener">Kaynağa Git →</a>
                    </div>
                </div>
            </div>`).join('');
    } catch (err) {
        el.innerHTML = `<div class="empty" style="grid-column:1/-1;">Yer imleri yüklenemedi: ${err.message}</div>`;
    }
}

function filterBookmarks(filter, btn) {
    if (btn) {
        btn.parentElement.querySelectorAll('button').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
    }
    loadBookmarks(filter);
}

async function deleteBookmark(id) {
    if (!confirm('Bu yer imini silmek istediğinize emin misiniz?')) return;
    try {
        await api(`/me/bookmarks/${id}`, { method: 'DELETE' });
        loadBookmarks();
    } catch (err) {
        alert('Silme hatası: ' + err.message);
    }
}

/* ─── Price Alerts (Fiyat Alarmları) ─── */
async function loadAlerts() {
    const el = document.getElementById('alertsList');
    const badge = document.getElementById('alertCountBadge');
    if (!el) return;
    el.innerHTML = '<div class="sk sk-x" style="grid-column:1/-1;"></div>';

    try {
        const res = await api('/me/alerts');
        const alerts = res.alerts || [];
        if (badge) badge.textContent = `${alerts.length} Alarm`;

        if (!alerts.length) {
            el.innerHTML = '<div class="empty" style="grid-column:1/-1;">Henüz kurulmuş fiyat alarmınız bulunmuyor.</div>';
            return;
        }

        el.innerHTML = alerts.map(a => `
            <div style="background:var(--card); border:1px solid var(--bdr); border-radius:10px; padding:1.2rem; display:flex; justify-content:space-between; align-items:center;">
                <div>
                    <div style="font-family:'Aldrich',sans-serif; font-size:1.1rem; color:var(--primary); font-weight:bold;">${a.pair}</div>
                    <div style="font-size:0.85rem; color:var(--t2); margin-top:0.2rem;">
                        Hedef: <strong>${a.condition === 'above' ? '≥' : '≤'} ${a.targetRate}</strong>
                    </div>
                    ${a.note ? `<div style="font-size:0.75rem; color:var(--t3); margin-top:0.3rem;">Not: ${a.note}</div>` : ''}
                    <div style="font-size:0.72rem; color:${a.triggered ? '#fb7185' : 'var(--mint)'}; margin-top:0.4rem;">
                        ${a.triggered ? '● Tetiklendi' : '● Aktif Takipte'}
                    </div>
                </div>
                <button onclick="deleteAlert('${a.id}')" style="background:rgba(255,0,0,0.2); color:#ff6b6b; border:1px solid rgba(255,0,0,0.4); border-radius:6px; padding:0.4rem 0.8rem; font-size:0.8rem; cursor:pointer;">
                    Sil ✕
                </button>
            </div>`).join('');
    } catch (err) {
        el.innerHTML = `<div class="empty" style="grid-column:1/-1;">Alarmlar yüklenemedi: ${err.message}</div>`;
    }
}

async function handleCreateAlert(e) {
    if (e) e.preventDefault();
    const pair = document.getElementById('alertPair').value;
    const condition = document.getElementById('alertCondition').value;
    const targetRate = parseFloat(document.getElementById('alertTargetRate').value);
    const note = document.getElementById('alertNote').value.trim();

    if (!targetRate || isNaN(targetRate)) {
        alert('Lütfen geçerli bir hedef fiyat girin.');
        return;
    }

    try {
        await api('/me/alerts', {
            method: 'POST',
            body: JSON.stringify({ pair, condition, targetRate, note })
        });
        trackEvent('alert_create', 'alerts', { pair, targetRate });
        document.getElementById('alertTargetRate').value = '';
        document.getElementById('alertNote').value = '';
        loadAlerts();
    } catch (err) {
        alert('Alarm kurulamadı: ' + err.message);
    }
}

async function deleteAlert(id) {
    if (!confirm('Bu alarmı silmek istediğinize emin misiniz?')) return;
    try {
        await api(`/me/alerts/${id}`, { method: 'DELETE' });
        loadAlerts();
    } catch (err) {
        alert('Silme hatası: ' + err.message);
    }
}

/* ─── Profile & Preferences ─── */
async function openProfileModal() {
    const modal = document.getElementById('profileModal');
    if (!modal) return;
    modal.style.display = 'flex';
    modal.classList.remove('hidden');

    try {
        const res = await api('/me');
        if (res) {
            document.getElementById('profUsername').textContent = res.username || '-';
            document.getElementById('profEmail').textContent = res.email || '-';
            document.getElementById('profBookmarkCount').textContent = res.bookmark_count != null ? res.bookmark_count : '0';
            if (res.preferences) {
                if (res.preferences.defaultCity) document.getElementById('prefDefaultCity').value = res.preferences.defaultCity;
                if (res.preferences.defaultLang) document.getElementById('prefDefaultLang').value = res.preferences.defaultLang;
            }
        }
    } catch(e) {
        console.error('Profil yükleme hatası:', e);
    }
}

function closeProfileModal() {
    const modal = document.getElementById('profileModal');
    if (modal) {
        modal.style.display = 'none';
        modal.classList.add('hidden');
    }
    const msgEl = document.getElementById('prefMsg');
    if (msgEl) msgEl.textContent = '';
}

async function handleSavePreferences(e) {
    if (e) {
        e.preventDefault();
        e.stopPropagation();
    }
    const defaultCity = document.getElementById('prefDefaultCity').value;
    const defaultLang = document.getElementById('prefDefaultLang').value;
    const msgEl = document.getElementById('prefMsg');

    try {
        await api('/me/preferences', {
            method: 'PUT',
            body: JSON.stringify({ defaultCity, defaultLang })
        });
        if (msgEl) {
            msgEl.textContent = 'Tercihleriniz başarıyla güncellendi!';
            msgEl.style.color = 'var(--mint)';
        }
        const citySelect = document.getElementById('citySelect');
        if (citySelect && defaultCity) {
            citySelect.value = defaultCity;
            onContextChange();
        }
        setTimeout(() => closeProfileModal(), 900);
    } catch (err) {
        if (msgEl) {
            msgEl.textContent = 'Hata: ' + err.message;
            msgEl.style.color = '#fb7185';
        }
    }
    return false;
}

async function handleResetPreferences() {
    if (!confirm('Tercihlerinizi sıfırlamak istiyor musunuz?')) return;
    try {
        await api('/me/preferences', { method: 'DELETE' });
        openProfileModal();
    } catch(e) {}
}

/* ─── Error state ─── */
function showError(msg) {
    const el = document.getElementById('xgrid');
    if (el) el.innerHTML = `
        <div class="empty" style="grid-column:1/-1">
            <span class="empty-ico">⚠️</span>
            <strong>Bağlantı hatası</strong>
            <code>${msg}</code>
            Servis çalışıyor mu? → <code>go run ./main.go --dev</code>
        </div>`;
}

/* ─── Global Exports for Inline HTML Handlers ─── */
window.openProfileModal = openProfileModal;
window.closeProfileModal = closeProfileModal;
window.handleSavePreferences = handleSavePreferences;
window.handleResetPreferences = handleResetPreferences;
window.saveBookmark = saveBookmark;
window.deleteBookmark = deleteBookmark;
window.filterBookmarks = filterBookmarks;
window.handleCreateAlert = handleCreateAlert;
window.deleteAlert = deleteAlert;
window.navigateTo = navigateTo;
window.onContextChange = onContextChange;
window.refresh = refresh;

/* ─── Init ─── */
document.addEventListener('DOMContentLoaded', () => {
    const user = JSON.parse(localStorage.getItem('cp_user') || '{}');
    const elUser = document.getElementById('navUser');
    if (elUser && user.username) {
        elUser.innerHTML = `👤 <span>${user.username}</span>`;
    }

    const hash = window.location.hash.replace('#', '') || 'exchange';
    const validSections = ['exchange', 'news', 'games', 'nasa', 'github', 'weather', 'crypto', 'bookmarks', 'alerts'];
    const initSection = validSections.includes(hash) ? hash : 'exchange';

    navigateTo(initSection);
    loadSnapshot();
});