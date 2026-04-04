import json
from unittest.mock import AsyncMock, MagicMock, patch

from fastapi.testclient import TestClient

from main import GO_PROFILE_URL, app

client = TestClient(app)


def test_post_profile_invalid_body_returns_422() -> None:
    """Pydantic rejects invalid email before any outbound httpx call."""
    resp = client.post(
        "/profile",
        json={
            "display_name": "Timur",
            "email": "not-an-email",
            "age": 20,
        },
    )
    assert resp.status_code == 422


def test_post_profile_valid_forwards_to_go_and_returns_upstream() -> None:
    payload = {
        "display_name": "Timur",
        "email": "t@example.com",
        "age": 20,
    }
    upstream_json = json.dumps(payload).encode("utf-8")

    mock_response = MagicMock()
    mock_response.status_code = 201
    mock_response.content = upstream_json
    mock_response.headers = {"content-type": "application/json"}

    mock_http = MagicMock()
    mock_http.post = AsyncMock(return_value=mock_response)
    mock_client_cm = MagicMock()
    mock_client_cm.__aenter__ = AsyncMock(return_value=mock_http)
    mock_client_cm.__aexit__ = AsyncMock(return_value=None)

    with patch("main.httpx.AsyncClient", return_value=mock_client_cm):
        resp = client.post("/profile", json=payload)

    assert resp.status_code == 201
    assert resp.json() == payload
    mock_http.post.assert_awaited_once()
    args, kwargs = mock_http.post.await_args
    assert args[0] == GO_PROFILE_URL
    assert kwargs["json"] == payload
