"""Pull a small image model once. Nothing at runtime depends on this; portraits
are generated offline and the PNGs are what the game ships."""
from huggingface_hub import snapshot_download
p = snapshot_download("stabilityai/sd-turbo", allow_patterns=[
    "*.json", "*.txt", "**/*.safetensors", "*.safetensors"])
print("MODEL_AT", p)
