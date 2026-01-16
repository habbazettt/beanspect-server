from pydantic import BaseModel, Field
from typing import List


class ClassPrediction(BaseModel):
    """Single class prediction with confidence score."""
    
    class_name: str = Field(..., alias="class", description="Predicted class name")
    confidence: float = Field(..., ge=0.0, le=1.0, description="Confidence score (0-1)")
    
    class Config:
        populate_by_name = True


class PredictionResponse(BaseModel):
    """Response schema for prediction endpoint."""
    
    predicted_class: str = Field(..., description="Top predicted class name")
    confidence: float = Field(..., ge=0.0, le=1.0, description="Confidence of top prediction")
    all_predictions: List[ClassPrediction] = Field(
        ..., description="All class predictions sorted by confidence"
    )
    
    class Config:
        json_schema_extra = {
            "example": {
                "predicted_class": "arabica",
                "confidence": 0.85,
                "all_predictions": [
                    {"class": "arabica", "confidence": 0.85},
                    {"class": "robusta", "confidence": 0.08},
                    {"class": "liberica", "confidence": 0.04},
                    {"class": "excelsa", "confidence": 0.03}
                ]
            }
        }


class ErrorResponse(BaseModel):
    """Standard error response schema."""
    
    error: bool = Field(default=True, description="Error indicator")
    code: str = Field(..., description="Error code")
    message: str = Field(..., description="Human-readable error message")
    
    class Config:
        json_schema_extra = {
            "example": {
                "error": True,
                "code": "INVALID_IMAGE",
                "message": "Uploaded file is not a valid image"
            }
        }


class HealthResponse(BaseModel):
    """Health check response schema."""
    
    status: str = Field(default="healthy", description="Service status")
    service: str = Field(..., description="Service name")
    version: str = Field(..., description="Service version")
    model_loaded: bool = Field(..., description="Whether model is loaded successfully")
