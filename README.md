# BeanSpect Server 🌱☕

Backend service untuk **BeanSpect** - Aplikasi identifikasi spesies biji kopi menggunakan AI/Machine Learning.

## 📋 Overview

BeanSpect Server terdiri dari layanan AI Inference yang menggunakan TensorFlow SavedModel untuk mengklasifikasikan gambar biji kopi ke dalam 4 spesies:

| Spesies | Deskripsi |
|---------|-----------|
| **Arabica** | Biji kopi premium dengan rasa asam dan kompleks |
| **Excelsa** | Biji kopi langka dengan profil rasa unik |
| **Liberica** | Biji kopi besar dengan aroma kuat |
| **Robusta** | Biji kopi dengan kadar kafein tinggi |

## 🏗️ Arsitektur

```
beanspect-server/
├── docker-compose.yml          # Container orchestration
├── .gitignore                  # Git ignore rules
├── beanspect_savedmodel/       # TensorFlow SavedModel
│   ├── saved_model.pb
│   ├── variables/
│   ├── assets/
│   └── class_names.json
└── inference-service/          # FastAPI Inference Service
    ├── Dockerfile
    ├── requirements.txt
    ├── .env
    └── app/
        ├── main.py
        ├── api/routes/
        ├── core/
        ├── models/
        └── schemas/
```

## 🚀 Quick Start

### Prerequisites

- Docker Desktop
- Git

### Menjalankan dengan Docker

```bash
# Clone repository
git clone <repository-url>
cd beanspect-server

# Buat file .env dari template
cp inference-service/.env.example inference-service/.env

# Build dan jalankan
docker-compose up --build -d

# Cek status
docker-compose ps

# Lihat logs
docker-compose logs -f
```

### Verifikasi Service

```bash
# Health check
curl http://localhost:8001/health

# Response:
# {"status":"healthy","service":"BeanSpect AI Inference Service","version":"1.0.0","model_loaded":true}
```

## 📡 API Endpoints

### Inference Service (Port 8001)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/` | Service info |
| GET | `/health` | Health check |
| GET | `/docs` | Swagger UI |
| GET | `/redoc` | ReDoc API Docs |
| POST | `/predict` | Klasifikasi gambar biji kopi |

### Contoh Request Predict

```bash
curl -X POST "http://localhost:8001/predict" \
  -H "Content-Type: multipart/form-data" \
  -F "file=@coffee_bean_image.jpg"
```

### Contoh Response

```json
{
  "predictions": [
    {"class_name": "arabica", "confidence": 0.89},
    {"class_name": "robusta", "confidence": 0.08},
    {"class_name": "liberica", "confidence": 0.02},
    {"class_name": "excelsa", "confidence": 0.01}
  ],
  "top_prediction": {
    "class_name": "arabica",
    "confidence": 0.89
  }
}
```

## ⚙️ Environment Variables

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `APP_NAME` | BeanSpect AI Inference Service | Nama aplikasi |
| `APP_VERSION` | 1.0.0 | Versi aplikasi |
| `DEBUG` | false | Mode debug |
| `HOST` | 0.0.0.0 | Host address |
| `PORT` | 8000 | Port internal |
| `MODEL_PATH` | /app/models/beanspect_savedmodel | Path model |
| `CLASS_NAMES_PATH` | /app/models/class_names.json | Path class names |
| `IMAGE_SIZE` | 224 | Ukuran input gambar |
| `MAX_FILE_SIZE` | 10485760 | Max file size (10MB) |
| `ALLOWED_EXTENSIONS` | ["jpg","jpeg","png"] | Ekstensi yang diizinkan |

## 🧪 Development

### Menjalankan tanpa Docker (Development)

```bash
cd inference-service

# Buat virtual environment
python -m venv venv
source venv/bin/activate  # Linux/Mac
venv\Scripts\activate     # Windows

# Install dependencies
pip install -r requirements.txt

# Jalankan server
uvicorn app.main:app --reload --port 8000
```

### Menghentikan Service

```bash
# Stop containers
docker-compose down

# Stop dan hapus volumes
docker-compose down -v
```

## 📁 Model

Model TensorFlow SavedModel terletak di `beanspect_savedmodel/`. Model ini dilatih untuk mengklasifikasikan gambar biji kopi dengan input size 224x224 pixels.

**Input**: Gambar RGB 224x224  
**Output**: Probabilitas untuk 4 kelas (arabica, excelsa, liberica, robusta)

## 🔧 Tech Stack

- **Framework**: FastAPI
- **ML**: TensorFlow 2.15.0
- **Image Processing**: Pillow
- **Server**: Uvicorn
- **Containerization**: Docker

## 📄 License

MIT License
