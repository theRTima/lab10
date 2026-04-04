import os

import httpx
from fastapi import FastAPI
from fastapi.responses import Response
from pydantic import BaseModel, EmailStr, Field

GO_PROFILE_URL = os.environ.get(
    "GO_PROFILE_URL",
    "http://localhost:8080/profile",
)

app = FastAPI(title="Profile proxy", version="1.0.0")


class Profile(BaseModel):
    """Matches profileBody in mid/3/3.go (JSON tags and validator rules)."""

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
