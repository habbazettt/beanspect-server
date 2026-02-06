from fastapi import APIRouter
from .health import router as health_router
from .predict import router as predict_router

# Create main API router
api_router = APIRouter()

# Include routes
api_router.include_router(health_router)
api_router.include_router(predict_router)
