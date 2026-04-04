import os
from typing import Annotated

import httpx
import jwt
from fastapi import FastAPI, Header, HTTPException, Response, status
from pydantic import BaseModel, EmailStr, Field

_go_base = os.environ.get("GO_BACKEND_URL", "http://127.0.0.1:8080").rstrip("/")
GO_PROFILE_URL = os.environ.get("GO_PROFILE_URL", f"{_go_base}/profile")

JWT_SECRET = os.environ.get("JWT_SECRET", "dev-secret-change-me")
JWT_ALGORITHM = "HS256"

app = FastAPI(
    title="Gateway",
    version="1.0.0",
    description="Validates JWT, then proxies /profile to the Go service.",
)


class Profile(BaseModel):
    """Same shape as profileBody in hard/go-service (JSON + validation)."""

    display_name: str = Field(min_length=2, max_length=80)
    email: EmailStr
    age: int = Field(ge=1, le=150)


def verify_jwt_bearer(authorization: Annotated[str | None, Header()] = None) -> str:
    if not authorization or not authorization.startswith("Bearer "):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Missing Authorization: Bearer token",
        )
    token = authorization.removeprefix("Bearer ").strip()
    if not token:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Empty bearer token",
        )
    try:
        jwt.decode(token, JWT_SECRET, algorithms=[JWT_ALGORITHM])
    except jwt.PyJWTError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid or expired token",
        ) from None
    return authorization


@app.post("/profile")
async def forward_profile(
    profile: Profile,
    authorization: Annotated[str | None, Header()] = None,
) -> Response:
    auth_header = verify_jwt_bearer(authorization)
    async with httpx.AsyncClient(timeout=30.0) as client:
        upstream = await client.post(
            GO_PROFILE_URL,
            json=profile.model_dump(mode="json"),
            headers={"Authorization": auth_header},
        )
    ct = upstream.headers.get("content-type", "application/json")
    return Response(
        content=upstream.content,
        status_code=upstream.status_code,
        media_type=ct.split(";")[0].strip() or "application/json",
    )
