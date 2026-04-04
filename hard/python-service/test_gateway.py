import json
from datetime import UTC, datetime, timedelta
from unittest.mock import AsyncMock, MagicMock, patch

import jwt
from fastapi.testclient import TestClient

import main

client = TestClient(main.app)


def _valid_profile() -> dict:
    return {
        "display_name": "Timur",
        "email": "t@example.com",
        "age": 20,
    }


def _make_token(secret: str | None = None) -> str:
    secret = secret or main.JWT_SECRET
    return jwt.encode(
        {
            "sub": "test-user",
            "exp": datetime.now(UTC) + timedelta(hours=1),
            "iat": datetime.now(UTC),
        },
        secret,
        algorithm=main.JWT_ALGORITHM,
    )


def test_invalid_body_rejected_before_upstream() -> None:
    """Pydantic 422; httpx must not run."""
    called: list[bool] = []

    async def boom(*_a, **_kw):
        called.append(True)
        raise AssertionError("httpx should not be called")

    mock_http = MagicMock()
    mock_http.post = boom
    mock_cm = MagicMock()
    mock_cm.__aenter__ = AsyncMock(return_value=mock_http)
    mock_cm.__aexit__ = AsyncMock(return_value=None)

    with patch.object(main.httpx, "AsyncClient", return_value=mock_cm):
        resp = client.post(
            "/profile",
            json={
                "display_name": "Timur",
                "email": "not-an-email",
                "age": 20,
            },
            headers={"Authorization": "Bearer " + _make_token()},
        )

    assert resp.status_code == 422
    assert called == []


def test_missing_bearer_401_no_upstream() -> None:
    called: list[bool] = []

    async def boom(*_a, **_kw):
        called.append(True)
        raise AssertionError("httpx should not be called")

    mock_http = MagicMock()
    mock_http.post = boom
    mock_cm = MagicMock()
    mock_cm.__aenter__ = AsyncMock(return_value=mock_http)
    mock_cm.__aexit__ = AsyncMock(return_value=None)

    with patch.object(main.httpx, "AsyncClient", return_value=mock_cm):
        resp = client.post("/profile", json=_valid_profile())

    assert resp.status_code == 401
    assert called == []


def test_forwards_authorization_header_to_go() -> None:
    token = _make_token()
    profile = _valid_profile()
    upstream_json = json.dumps(profile).encode("utf-8")

    mock_response = MagicMock()
    mock_response.status_code = 201
    mock_response.content = upstream_json
    mock_response.headers = {"content-type": "application/json"}

    mock_http = MagicMock()
    mock_http.post = AsyncMock(return_value=mock_response)
    mock_cm = MagicMock()
    mock_cm.__aenter__ = AsyncMock(return_value=mock_http)
    mock_cm.__aexit__ = AsyncMock(return_value=None)

    auth_value = "Bearer " + token
    with patch.object(main.httpx, "AsyncClient", return_value=mock_cm):
        resp = client.post(
            "/profile",
            json=profile,
            headers={"Authorization": auth_value},
        )

    assert resp.status_code == 201
    mock_http.post.assert_awaited_once()
    _args, kwargs = mock_http.post.await_args
    assert kwargs["headers"]["Authorization"] == auth_value
    assert kwargs["json"] == profile


def test_upstream_status_passed_through_to_client() -> None:
    """HTTP status from Go (mocked) must match the gateway response status."""
    token = _make_token()
    profile = _valid_profile()

    cases = [
        (201, {"display_name": "Timur", "email": "t@example.com", "age": 20}),
        (400, {"error": "Key: 'profileBody.Email' Error:Field validation for 'Email' failed on the 'email' tag"}),
    ]

    for upstream_status, upstream_body in cases:
        raw = json.dumps(upstream_body).encode("utf-8")
        mock_response = MagicMock()
        mock_response.status_code = upstream_status
        mock_response.content = raw
        mock_response.headers = {"content-type": "application/json"}

        mock_http = MagicMock()
        mock_http.post = AsyncMock(return_value=mock_response)
        mock_cm = MagicMock()
        mock_cm.__aenter__ = AsyncMock(return_value=mock_http)
        mock_cm.__aexit__ = AsyncMock(return_value=None)

        with patch.object(main.httpx, "AsyncClient", return_value=mock_cm):
            resp = client.post(
                "/profile",
                json=profile,
                headers={"Authorization": "Bearer " + token},
            )

        assert resp.status_code == upstream_status, upstream_status
        assert resp.content == raw
