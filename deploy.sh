#!/bin/bash

# Summarify deployment script

echo "🚀 Başlatılıyor: Summarify Deployment"

# Uygulama dizini
APP_DIR="$DEPLOY_PATH"
cd "$APP_DIR" || exit 1

# Go uygulamasını build et
echo "🔨 Go uygulaması derleniyor..."
go build -o summarify main.go

# Çalıştırma izinlerini ayarla
chmod +x summarify
chmod +x yt-dlp

# Systemd servis dosyasını kopyala ve etkinleştir
if [ -f "summarify.service" ]; then
    echo "🔄 Systemd servisini yapılandırıyor..."
    cp summarify.service /etc/systemd/system/
    systemctl daemon-reload
    systemctl enable summarify
    systemctl restart summarify
    echo "✅ Servis başlatıldı: summarify"
else
    echo "❌ summarify.service dosyası bulunamadı"
fi

# Servis durumunu kontrol et
systemctl status summarify

echo "✅ Deployment tamamlandı!" 
