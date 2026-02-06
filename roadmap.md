
# 📘 **BeanSpect — Fullstack System Development Plan**

> **Scope:**
> Sistem fullstack untuk klasifikasi spesies biji kopi berbasis AI, dilengkapi inference service dan visualisasi GIS.
> **ML modelling & training:** *COMPLETED (out of scope)*

---

## 🧱 **System Overview**

BeanSpect dibangun sebagai **multi-service architecture** dengan pemisahan concern yang jelas:

* **AI Inference Service** → prediksi spesies kopi
* **Application Backend** → orchestration & GIS logic
* **Frontend Web App** → user interface & visualization

---

## 🏗️ **High-Level Architecture**

```
[ Web Frontend ]
      |
      v
[ Application Backend ]
      |
      +--> [ AI Inference Service ]
      |
      +--> [ GIS / Origin Mapping ]
```

---

## 🔹 **Technology Stack (Recommended)**

### Backend

* **AI Inference:** FastAPI (Python)
* **App Backend:** Fiber (Go) /  (non-ML routes)
* **AI Framework:** TensorFlow + Keras 3
* **Model Format:** TensorFlow SavedModel (inference-only)

### Frontend

* React / Next.js
* Leaflet / Mapbox (GIS visualization)

### Infrastructure

* Docker (local deployment)
* Local filesystem (no cloud dependency)

---

# 🚀 **EPIC 1 — AI Inference Service (FastAPI)**

**Goal:**
Expose trained BeanSpect model as a stateless HTTP inference API.

---

### 1.1 Service Setup ✅

* [x] Initialize FastAPI project
* [x] Setup virtual environment & dependencies
* [x] Define project structure (`app/`, `models/`, `schemas/`)
* [x] Configure CORS for frontend access
* [x] Health check endpoint (`GET /health`)

---

### 1.2 Model Loading & Initialization

* [x] Load `SavedModel` via `keras.layers.TFSMLayer`
* [x] Verify `serving_default` signature
* [x] Load class label mapping
* [x] Warm-up model on startup
* [x] Log model metadata (version, size)

---

### 1.3 Inference API

* [x] Endpoint `POST /predict`
* [x] Accept image upload (multipart/form-data)
* [x] Image preprocessing (resize, normalize)
* [x] Run inference using TFSMLayer
* [x] Return JSON response

---

### 1.4 Error Handling & Validation

* [x] Validate file type (image only)
* [x] Handle corrupted images
* [x] Graceful inference failure handling
* [x] Standardized error response format

---

# 🌍 **EPIC 2 — Application Backend & GIS Logic**

**Goal:**
Handle non-ML logic, GIS mapping, and orchestration.

---

### 2.1 Backend Core Setup

* [x] Initialize backend service (Fiber)
* [x] Environment configuration
* [x] Service-to-service communication with inference API
* [x] Centralized logging

---

### 2.2 Coffee Origin Mapping (GIS)

* [x] Define species → origin dataset (static JSON or DB)
* [x] Map species to:
  * country
  * region
  * coordinates (lat/lng)
  * image
  * description
  * taste profile
  * caffeine level
  * etc.
* [x] Endpoint `GET /origin/{species}`

---

### 2.3 Orchestration Flow

* [x] Forward image to inference service
* [x] Receive prediction result
* [x] Fetch GIS origin data
* [x] Return combined response to frontend

---

# 🖥️ **EPIC 3 — Frontend Web Application**

**Goal:**
Provide intuitive UI for image upload, prediction display, and GIS visualization.

---

### 3.1 UI Foundation

* [x] Initialize frontend project
* [x] Layout & navigation
* [x] API client setup
* [x] Environment configuration

---

### 3.2 Image Upload & Prediction View

* [x] Image upload component
* [x] Preview uploaded image
* [x] Call backend prediction endpoint
* [x] Display:

  * [x] predicted species
  * [x] confidence score

---

### 3.3 GIS Visualization

* [x] Integrate map library (Leaflet / Mapbox)
* [x] Plot origin location based on species
* [x] Marker with species metadata
* [x] Responsive map rendering

---

### 3.4 Explainability (Optional UI)

* [ ] Display Grad-CAM overlays (static artifacts)
* [ ] Toggle explainability view
* [ ] Contextual explanation text

---

# 📦 **EPIC 4 — Deployment & Integration**

**Goal:**
Run full system locally in a reproducible manner.

---

### 4.1 Containerization

* [ ] Dockerfile for inference service
* [ ] Dockerfile for backend
* [ ] Dockerfile for frontend
* [ ] Docker Compose orchestration

---

### 4.2 Environment Validation

* [ ] Local startup instructions
* [ ] Service dependency checks
* [ ] Health verification for all services

---

# 📄 **EPIC 5 — Documentation & Handover**

**Goal:**
Make the project understandable, reproducible, and portfolio-ready.

---

### 5.1 Technical Documentation

* [ ] System architecture diagram
* [ ] API documentation
* [ ] Model artifact description
* [ ] GIS data source explanation

---

### 5.2 Final Deliverables

* [ ] `README.md`
* [ ] `run_manifest.json`
* [ ] Model artifacts
* [ ] Screenshots / demo GIFs
