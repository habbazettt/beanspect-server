import logging
from io import BytesIO

from fastapi import APIRouter, File, UploadFile, HTTPException, status
from PIL import Image

from app.core.config import get_settings
from app.models.inference import get_model_service
from app.schemas.prediction import PredictionResponse, ErrorResponse

logger = logging.getLogger(__name__)
router = APIRouter()
settings = get_settings()


def validate_file(file: UploadFile) -> None:
    """Validate uploaded file type and size."""
    # Check file extension
    if file.filename:
        extension = file.filename.split(".")[-1].lower()
        allowed = [ext.strip('"') for ext in settings.ALLOWED_EXTENSIONS]
        if extension not in allowed:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail={
                    "error": True,
                    "code": "INVALID_FILE_TYPE",
                    "message": f"File type '.{extension}' not allowed. Allowed: {allowed}"
                }
            )
    
    # Check content type
    if file.content_type and not file.content_type.startswith("image/"):
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail={
                "error": True,
                "code": "INVALID_CONTENT_TYPE",
                "message": "Uploaded file must be an image"
            }
        )


@router.post(
    "/predict",
    response_model=PredictionResponse,
    responses={
        400: {"model": ErrorResponse, "description": "Invalid image or request"},
        500: {"model": ErrorResponse, "description": "Inference error"},
        503: {"model": ErrorResponse, "description": "Model not loaded"}
    },
    summary="Predict coffee bean species",
    description="Upload an image of a coffee bean to classify its species."
)
async def predict(
    file: UploadFile = File(..., description="Image file (JPG, JPEG, PNG)")
):
    """
    Classify a coffee bean image into one of four species:
    - **arabica**: Premium coffee with acidic and complex flavor
    - **excelsa**: Rare coffee with unique taste profile  
    - **liberica**: Large beans with strong aroma
    - **robusta**: High caffeine content beans
    """
    model_service = get_model_service()
    
    # Check if model is loaded
    if not model_service.is_loaded:
        raise HTTPException(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            detail={
                "error": True,
                "code": "MODEL_NOT_LOADED",
                "message": "Model is not loaded. Please try again later."
            }
        )
    
    # Validate file
    validate_file(file)
    
    try:
        # Read file content
        contents = await file.read()
        
        # Check file size
        if len(contents) > settings.MAX_FILE_SIZE:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail={
                    "error": True,
                    "code": "FILE_TOO_LARGE",
                    "message": f"File size exceeds maximum of {settings.MAX_FILE_SIZE // (1024*1024)}MB"
                }
            )
        
        # Open image with PIL
        try:
            image = Image.open(BytesIO(contents))
            # Verify image is not corrupted
            image.verify()
            # Re-open after verify (verify() closes the file)
            image = Image.open(BytesIO(contents))
            # Force load to catch truncated images
            image.load()
        except Exception as e:
            logger.error(f"Failed to open/validate image: {e}")
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail={
                    "error": True,
                    "code": "CORRUPTED_IMAGE",
                    "message": "Uploaded image is corrupted or invalid"
                }
            )
        
        # Run inference
        logger.info(f"Running inference on image: {file.filename}")
        result = model_service.predict(image)
        
        logger.info(f"Prediction: {result['predicted_class']} ({result['confidence']:.2%})")
        
        return PredictionResponse(
            predicted_class=result["predicted_class"],
            confidence=result["confidence"],
            all_predictions=result["all_predictions"]
        )
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Inference error: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail={
                "error": True,
                "code": "INFERENCE_ERROR",
                "message": f"Failed to process image: {str(e)}"
            }
        )
    finally:
        await file.close()
