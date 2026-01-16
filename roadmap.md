
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

* [ ] Load `SavedModel` via `keras.layers.TFSMLayer`
* [ ] Verify `serving_default` signature
* [ ] Load class label mapping
* [ ] Warm-up model on startup
* [ ] Log model metadata (version, size)

---

### 1.3 Inference API

* [ ] Endpoint `POST /predict`
* [ ] Accept image upload (multipart/form-data)
* [ ] Image preprocessing (resize, normalize)
* [ ] Run inference using TFSMLayer
* [ ] Return JSON response:

  ```json
  {
    "class": "arabica",
    "confidence": 0.93
  }
  ```

---

### 1.4 Error Handling & Validation

* [ ] Validate file type (image only)
* [ ] Handle corrupted images
* [ ] Graceful inference failure handling
* [ ] Standardized error response format

---

# 🌍 **EPIC 2 — Application Backend & GIS Logic**

**Goal:**
Handle non-ML logic, GIS mapping, and orchestration.

---

### 2.1 Backend Core Setup

* [ ] Initialize backend service (Fiber / Node / FastAPI)
* [ ] Environment configuration
* [ ] Service-to-service communication with inference API
* [ ] Centralized logging

---

### 2.2 Coffee Origin Mapping (GIS)

* [ ] Define species → origin dataset (static JSON or DB)
* [ ] Map species to:

  * country
  * region
  * coordinates (lat/lng)
* [ ] Endpoint `GET /origin/{species}`

---

### 2.3 Orchestration Flow

* [ ] Receive image from frontend
* [ ] Forward image to inference service
* [ ] Receive prediction result
* [ ] Fetch GIS origin data
* [ ] Return combined response to frontend

---

# 🖥️ **EPIC 3 — Frontend Web Application**

**Goal:**
Provide intuitive UI for image upload, prediction display, and GIS visualization.

---

### 3.1 UI Foundation

* [ ] Initialize frontend project
* [ ] Layout & navigation
* [ ] API client setup
* [ ] Environment configuration

---

### 3.2 Image Upload & Prediction View

* [ ] Image upload component
* [ ] Preview uploaded image
* [ ] Call backend prediction endpoint
* [ ] Display:

  * predicted species
  * confidence score

---

### 3.3 GIS Visualization

* [ ] Integrate map library (Leaflet / Mapbox)
* [ ] Plot origin location based on species
* [ ] Marker with species metadata
* [ ] Responsive map rendering

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

---

## 🔒 **Out of Scope (Locked)**

* ❌ Model retraining
* ❌ Fine-tuning experiments
* ❌ Dataset changes
* ❌ Performance benchmarking at scale
