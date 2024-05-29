# GoDir3num

GoDirEnum, bir web sitesindeki gizli dizinleri ve dosyaları bulmak için kullanılan çok iş parçacıklı bir Go programıdır. Kullanıcı tarafından sağlanan bir URL ve wordlist dosyasını kullanarak HTTP istekleri yapar ve yanıt kodlarına göre sonuçları filtreler.

## Gereksinimler

- Go 1.16 veya üstü
- Git
- İnternet bağlantısı

## Kurulum

### 1. Depoyu Klonlayın

Öncelikle, projeyi klonlayın:

```sh
git clone https://github.com/Muhammetcansimsek/GoDir3num.git
cd GoDir3num
```

### 2. Bağımlılıkları Yükleyin

Bağımlılıkları yüklemek için aşağıdaki komutu çalıştırın:

```sh
go mod tidy
```

## Kullanım

### 1. Wordlist Dosyasını Hazırlayın

Farklı dizin ve dosya adlarını içeren bir wordlist dosyasına ihtiyacınız olacak. Örneğin:

```
admin
login
dashboard
uploads
```

Bu dosyayı `wordlist.txt` olarak kaydedin.

### 2. Programı Derleyin

Programı derlemek için aşağıdaki komutu kullanın:

```sh
go build -o webfuzzer
```

### 3. Programı Çalıştırın

Programı çalıştırmak için aşağıdaki komutu kullanın:

```sh
./webfuzzer -url http://hedef-site.com -wordlist wordlist.txt -threads 150
```

### Komut Satırı Seçenekleri

- `-url`: Tarama yapılacak temel URL.
- `-wordlist`: Wordlist dosyasının yolu.
- `-threads`: Eşzamanlı iş parçacığı sayısı (varsayılan: 150).
- `-verbose`: Ayrıntılı çıktıyı etkinleştirir.

### 4. Örnek Kullanım

```sh
./webfuzzer -url http://example.com -wordlist common.txt -threads 100 -verbose true
```

## Geliştirme

Proje üzerinde çalışırken kullanışlı olabilecek bazı komutlar:

### Kodları Derlemeden Çalıştırmak

```sh
go run main.go result.go -url http://example.com -wordlist common.txt -threads 100
```

## Katkıda Bulunma

Katkıda bulunmak isterseniz, lütfen bir pull request gönderin veya bir issue açın.
