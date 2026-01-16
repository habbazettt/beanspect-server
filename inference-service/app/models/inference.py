import json
import logging
from pathlib import Path
from typing import Optional, List, Dict, Any

import numpy as np
import tensorflow as tf
from PIL import Image

from app.core.config import get_settings

logger = logging.getLogger(__name__)


class ModelService:
    """Service for loading and running inference on the BeanSpect model."""
    
    _instance: Optional["ModelService"] = None
    
    def __init__(self):
        self.model = None
        self.class_names: List[str] = []
        self.is_loaded: bool = False
        self.settings = get_settings()
    
    @classmethod
    def get_instance(cls) -> "ModelService":
        """Get singleton instance of ModelService."""
        if cls._instance is None:
            cls._instance = cls()
        return cls._instance
    
    def load_model(self) -> bool:
        """
        Load the TensorFlow SavedModel and class names.
        
        Returns:
            bool: True if model loaded successfully, False otherwise.
        """
        try:
            model_path = Path(self.settings.MODEL_PATH)
            
            if not model_path.exists():
                logger.error(f"Model path does not exist: {model_path}")
                return False
            
            logger.info(f"Loading model from: {model_path}")
            
            # Load the SavedModel
            self.model = tf.saved_model.load(str(model_path))
            
            # Verify model signature
            if hasattr(self.model, "signatures"):
                signatures = list(self.model.signatures.keys())
                logger.info(f"Model signatures available: {signatures}")
                
                if "serving_default" in signatures:
                    logger.info("Found 'serving_default' signature")
            
            # Load class names
            self._load_class_names()
            
            # Warm up the model
            self._warm_up()
            
            self.is_loaded = True
            logger.info("Model loaded successfully!")
            logger.info(f"Classes: {self.class_names}")
            
            return True
            
        except Exception as e:
            logger.error(f"Failed to load model: {e}")
            self.is_loaded = False
            return False
    
    def _load_class_names(self) -> None:
        """Load class names from JSON file."""
        class_names_path = Path(self.settings.CLASS_NAMES_PATH)
        
        if class_names_path.exists():
            with open(class_names_path, "r") as f:
                data = json.load(f)
                if isinstance(data, list):
                    self.class_names = data
                elif isinstance(data, dict) and "class_names" in data:
                    self.class_names = data["class_names"]
                logger.info(f"Loaded {len(self.class_names)} class names")
        else:
            # Default class names for coffee species
            self.class_names = ["arabica", "excelsa", "liberica", "robusta"]
            logger.warning(f"Class names file not found, using defaults: {self.class_names}")
    
    def _warm_up(self) -> None:
        """Warm up the model with a dummy prediction to prevent cold start latency."""
        try:
            logger.info("Warming up model...")
            dummy_input = np.zeros((1, self.settings.IMAGE_SIZE, self.settings.IMAGE_SIZE, 3), dtype=np.float32)
            dummy_tensor = tf.constant(dummy_input)
            
            # Run inference with the model
            infer = self.model.signatures["serving_default"]
            _ = infer(dummy_tensor)
            
            logger.info("Model warm-up complete")
        except Exception as e:
            logger.warning(f"Model warm-up failed (this may be okay): {e}")
    
    def preprocess_image(self, image: Image.Image) -> np.ndarray:
        """
        Preprocess image for model inference.
        
        Args:
            image: PIL Image object
            
        Returns:
            Preprocessed numpy array ready for inference
        """
        # Convert to RGB if necessary
        if image.mode != "RGB":
            image = image.convert("RGB")
        
        # Resize to model input size
        image = image.resize((self.settings.IMAGE_SIZE, self.settings.IMAGE_SIZE))
        
        # Convert to numpy array
        img_array = np.array(image, dtype=np.float32)
        
        # Normalize pixel values (0-1)
        img_array = img_array / 255.0
        
        # Add batch dimension
        img_array = np.expand_dims(img_array, axis=0)
        
        return img_array
    
    def predict(self, image: Image.Image) -> Dict[str, Any]:
        """
        Run inference on an image.
        
        Args:
            image: PIL Image object
            
        Returns:
            Dictionary containing prediction results
        """
        if not self.is_loaded:
            raise RuntimeError("Model is not loaded")
        
        # Preprocess image
        input_array = self.preprocess_image(image)
        input_tensor = tf.constant(input_array)
        
        # Run inference
        infer = self.model.signatures["serving_default"]
        predictions = infer(input_tensor)
        
        # Get output tensor (the key may vary)
        output_key = list(predictions.keys())[0]
        probs = predictions[output_key].numpy()[0]
        
        # Build prediction results
        all_predictions = []
        for i, prob in enumerate(probs):
            class_name = self.class_names[i] if i < len(self.class_names) else f"class_{i}"
            all_predictions.append({
                "class": class_name,
                "confidence": float(prob)
            })
        
        # Sort by confidence descending
        all_predictions.sort(key=lambda x: x["confidence"], reverse=True)
        
        # Get top prediction
        top_prediction = all_predictions[0]
        
        return {
            "predicted_class": top_prediction["class"],
            "confidence": top_prediction["confidence"],
            "all_predictions": all_predictions
        }


def get_model_service() -> ModelService:
    """Get the singleton ModelService instance."""
    return ModelService.get_instance()
