# ⚡ City Pulse — Real-Time City Intelligence & Dashboard

Modern, katmanlı Go mimarisi ve zenginleştirilmiş koyu temalı web arayüzü ile gerçek zamanlı şehir, finans ve teknoloji verilerini bir araya getiren kapsamlı bir gösterge paneli (dashboard).

![City Pulse Preview](/ui/images/logo.png)

---

## 🚀 Temel Özellikler

- **💰 Canlı Döviz Kurları:** Popüler döviz çiftleri (USD/TRY, EUR/TRY, GBP/TRY vb.) anlık takip ve kur dönüşüm desteği.
- **📰 Günün Haberleri:** GNews API entegrasyonuyla şehir ve kategori bazlı filtrelenebilir son dakika haberleri.
- **🎮 Oyun İndirimleri & Fırsatlar:** CheapShark API ile PC oyunlarında aktif indirimler, tasarruf oranları ve fırsat yönlendirmeleri.
- **🚀 NASA APOD:** NASA Günün Astronomik Karesi (HD fotoğraf & video açıklamalarıyla).
- **⭐ GitHub Trending:** Günün popüler açık kaynak yazılım depoları, yıldız sayıları ve dillerine göre listeleme.
- **⛅ Canlı Hava Durumu:** Open-Meteo entegrasyonu ile rüzgar hızı, gece/gündüz durumu ve sıcaklık takibi.
- **🪙 Canlı Kripto Piyasası:** CoinGecko API ile popüler kripto paraların 24 saatlik değişimleri ve anlık fiyatları.
- **🔖 Favorilerim (Bookmarks):** İlgini çeken haberleri, oyunları ve kripto paraları tek tıkla veritabanına kaydetme ve yönetme.
- **🔔 Fiyat Alarmları (Price Alerts):** Döviz veya kripto varlıklar için hedef fiyat belirleme ve durum takibi.
- **👤 Profil & Tercihler:** Varsayılan şehir ve dil tercihlerini kişiselleştirme.

---

## 🏛️ Mimari & Altyapı (Architecture)

Proje, kurumsal standartlarda **katmanlı (layered/hexagonal) mimari** ile geliştirilmiştir:

```text
city-pulse/
├── app/
│   ├── auth/          # JWT authentication & API Key middleware
│   ├── clients/       # Harici servis HTTP client'ları (Weather, Crypto, NASA, News, Exchange, Game, GitHub)
│   ├── daos/          # MongoDB Driver v2 veri erişim katmanı (Data Access Objects)
│   ├── dtos/          # Data Transfer Objects
│   ├── handlers/      # Gin HTTP Handler katmanı (REST API endpoints)
│   ├── middleware/    # Rate limiter, Panic recovery, Request logger, CORS
│   ├── models/        # BSON / JSON veri modelleri
│   ├── routes/        # Router tanımları ve dependency injection
│   └── services/      # İş mantığı (Business Logic Layer)
├── internal/
│   ├── apierror/      # Standart API hata yönetimi
│   ├── cache/         # Thread-safe in-memory TTL cache mekanizması
│   ├── config/        # Viper ile ortam yapılandırması (.json / .env)
│   ├── metrics/       # Endpoint performans ve istek sayaçları
│   └── ratelimit/     # Token-bucket algoritmalı hız sınırlayıcı
└── ui/                # Modern Glassmorphism & Aldrich temalı SPA arayüzü
```

---

## 🛠️ Kurulum ve Çalıştırma

### Gereksinimler
- [Go](https://go.dev/) (1.22+)
- [MongoDB](https://www.mongodb.com/) (Yerel veya Atlas)

### 1. Depoyu Klonlayın
```bash
git clone https://github.com/KULLANICI_ADIN/city-pulse.git
cd city-pulse
```

### 2. Yapılandırma
`config/local-config.example.json` dosyasını `config/local-config.json` olarak kopyalayın:
```bash
cp config/local-config.example.json config/local-config.json
```

### 3. Uygulamayı Başlatın
```powershell
go run ./main.go --dev
```

Tarayıcınızdan `http://localhost:3649` adresine giderek uygulamaya erişebilirsiniz.

---

## 📜 Lisans

Bu proje [MIT Lisansı](LICENSE) altında lisanslanmıştır.