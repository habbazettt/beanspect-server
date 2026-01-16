from fastapi import APIRouter, Depends
from app.core.config import Settings, get_settings
from app.schemas import HealthResponse

router = APIRouter(tags=["Health"])


@router.get(
    "/health",
    response_model=HealthResponse,
    summary="Health Check",
    description="Check the health status of the inference service"
)
async def health_check(settings: Settings = Depends(get_settings)) -> HealthResponse:
    """
    Health check endpoint.
    
    Returns the service status, version, and whether the model is loaded.
    """
    # Import here to avoid circular dependency
    from app.models.inference import get_model_service
    
    model_service = get_model_service()
    
    return HealthResponse(
        status="healthy",
        service=settings.APP_NAME,
        version=settings.APP_VERSION,
        model_loaded=model_service.is_loaded
    )
