from __future__ import annotations

import httpx
import pytest

from tests.support.assertions import assert_json_response
from tests.support.client import APIClient, AuthSession, RequestContext


def test_auth_session_as_headers_omits_empty_values() -> None:
    assert AuthSession().as_headers() == {}
    assert AuthSession(token="token-123", tenant_id="tenant-1").as_headers() == {
        "Authorization": "Bearer token-123",
        "X-Tenant-Id": "tenant-1",
    }


def test_request_context_build_headers_merges_auth_and_extra_headers() -> None:
    context = RequestContext(
        environment="ci",
        auth=AuthSession(token="token-123"),
        extra_headers={
            "Accept": "application/problem+json",
            "X-Trace-Id": "trace-1",
        },
    )

    assert context.build_headers() == {
        "Accept": "application/problem+json",
        "Authorization": "Bearer token-123",
        "X-Trace-Id": "trace-1",
    }


def test_assert_json_response_parses_json_payload() -> None:
    response = httpx.Response(
        200,
        headers={"content-type": "application/json; charset=utf-8"},
        json={"message": "ok"},
    )

    assert assert_json_response(response) == {"message": "ok"}


def test_assert_json_response_rejects_non_json_content_type() -> None:
    response = httpx.Response(
        200,
        headers={"content-type": "text/plain"},
        text="ok",
    )

    with pytest.raises(AssertionError, match="expected JSON response"):
        assert_json_response(response)


def test_assert_json_response_rejects_invalid_json_body() -> None:
    response = httpx.Response(
        200,
        headers={"content-type": "application/json"},
        content=b"not-json",
        request=httpx.Request("GET", "http://example.test/api/health"),
    )

    with pytest.raises(AssertionError, match="response body is not valid JSON"):
        assert_json_response(response)


def test_api_client_wraps_transport_errors(
    monkeypatch: pytest.MonkeyPatch,
) -> None:
    with APIClient(
        base_url="http://example.test/api",
        default_headers={"Accept": "application/json"},
        timeout_seconds=1.0,
    ) as client:
        def fake_send(request: httpx.Request) -> httpx.Response:
            raise httpx.ConnectError("connection refused", request=request)

        monkeypatch.setattr(client._client, "send", fake_send)

        with pytest.raises(
            AssertionError,
            match=r"API request failed: GET http://example\.test/api/health",
        ):
            client.get("/health")
