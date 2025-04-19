# GitHub Secrets Ayarlama Rehberi

GitHub Actions ile VPS deployment'ı için aşağıdaki secret'ları GitHub reponuza eklemeniz gerekmektedir:

## Gerekli Secret'lar

1. **SSH_PRIVATE_KEY**
   - SSH private key (id_rsa dosyasının içeriği)
   - VPS sunucusuna bağlanmak için kullanılır

2. **SSH_USER**
   - VPS sunucunuza bağlanmak için kullanıcı adı (genellikle "root")

3. **VPS_HOST**
   - VPS sunucunuzun IP adresi veya hostname (örn: 69.62.112.232)

4. **DEPLOY_PATH**
   - VPS sunucunuzda uygulamanın deploy edileceği dizin (örn: /root/summarify)

## Secret'ları Nasıl Eklersiniz?

1. GitHub reponuzda "Settings" sekmesine tıklayın
2. Sol menüden "Secrets and variables" > "Actions" seçeneğine tıklayın
3. "New repository secret" butonuna tıklayın
4. Her bir secret için:
   - İsim kısmına yukarıdaki isimlerden birini girin (örn: SSH_PRIVATE_KEY)
   - Value kısmına ilgili değeri girin
   - "Add secret" butonuna tıklayın

## SSH Private Key Nasıl Oluşturulur?

Eğer henüz bir SSH key'iniz yoksa:

```bash
# Lokal makinenizde çalıştırın
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"

# Private key'i görüntüleyin ve GitHub secret olarak ekleyin
cat ~/.ssh/id_rsa

# Public key'i VPS'e ekleyin
cat ~/.ssh/id_rsa.pub

# VPS'de ~/.ssh/authorized_keys dosyasına public key'i ekleyin
echo "BURAYA_PUBLIC_KEY_ICERIGI" >> ~/.ssh/authorized_keys
```

## VPS Tarafındaki Hazırlıklar

VPS sunucunuza SSH ile bağlanın:

```bash
ssh root@69.62.112.232
```

Uygulama dizinini oluşturun:

```bash
mkdir -p /root/summarify
chmod 755 /root/summarify
```

Go dilinin yüklü olduğundan emin olun:

```bash
go version || apt-get update && apt-get install -y golang-go
```

Deployment tamamlandıktan sonra servis durumunu kontrol edin:

```bash
systemctl status summarify
``` 
