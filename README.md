# YouTube Video Özetleyici (Summarify)

Bu proje, YouTube videolarını otomatik olarak özetleyen bir web uygulamasıdır. Video URL'sini girdikten sonra yapay zeka kullanarak videonun içeriğini özetler ve bu özeti farklı dillerde sunar.

## Özellikler

- YouTube videolarının otomatik olarak özetlenmesi
- Çoklu dil desteği (Türkçe, İngilizce, Rusça, Arapça)
- Özeti TXT veya PDF formatında indirme imkanı
- Basit ve kullanıcı dostu arayüz

## Nasıl Çalışır?

1. Kullanıcı bir YouTube video URL'si girer
2. Sistem video altyazılarını çeker
3. OpenRouter API aracılığıyla altyazı metni yapay zekaya gönderilir
4. Yapay zeka (Claude 3 Haiku) videoyu özetler
5. Özet kullanıcıya gösterilir

## Teknik Altyapı

- **Backend**: Go (Gin framework)
- **Frontend**: HTML, CSS, JavaScript
- **Altyazı Çekme**: yt-dlp
- **Yapay Zeka**: OpenRouter API (Claude 3 Haiku)

## Kurulum

### Yerel Geliştirme Ortamı

1. Repoyu klonlayın
   ```bash
   git clone https://github.com/Firudin085/-summarify.online.git
   cd -summarify.online
   ```

2. Bağımlılıkları yükleyin
   ```bash
   go mod download
   ```

3. `.env` dosyasını oluşturun
   ```
   OPENROUTER_API_KEY=your_api_key_here
   ```

4. Uygulamayı çalıştırın
   ```bash
   go run main.go
   ```

5. Tarayıcınızda http://localhost:8080 adresine gidin

### VPS'e Deploy Etme

Otomatik deployment için [GITHUB_SECRETS.md](GITHUB_SECRETS.md) dosyasını inceleyebilirsiniz.

## Katkıda Bulunma

1. Bu repoyu fork edin
2. Yeni bir özellik dalı oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'Yeni özellik: Harika özellik'`)
4. Dalınızı ana repoya push edin (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun

## İletişim

Firudin Mustafayev - firudinmustafayev00@gmail.com 