import os

import httpx
from fastapi import FastAPI
from fastapi.responses import Response
from pydantic import BaseModel, EmailStr, Field

# Base URL of the Gin backend (task 3). Profile POST is forwarded to {base}/profile.
_go_base = os.environ.get("GO_BACKEND_URL", "http://127.0.0.1:8080").rstrip("/")
GO_PROFILE_URL = os.environ.get("GO_PROFILE_URL", f"{_go_base}/profile")

app = FastAPI(title="Gateway", version="1.0.0", description="Proxies /profile to the Go service.")


class Profile(BaseModel):
    """Same shape as profileBody in hard/go-service (JSON + validation)."""

    display_name: str = Field(min_length=2, max_length=80)
    email: EmailStr
    age: int = Field(ge=1, le=150)


@app.post("/profile")
async def forward_profile(profile: Profile) -> Response:
    async with httpx.AsyncClient(timeout=30.0) as client:
        upstream = await client.post(
            GO_PROFILE_URL,
            json=profile.model_dump(mode="json"),
        )
    ct = upstream.headers.get("content-type", "application/json")
    return Response(
        content=upstream.content,
        status_code=upstream.status_code,
        media_type=ct.split(";")[0].strip() or "application/json",
    )
