from __future__ import annotations

from typing import Any

import httpx


def assert_json_response(response: httpx.Response) -> Any:
    # Centralize content-type and JSON parsing checks so individual tests only
    # describe business assertions.
    content_type = response.headers.get("content-type", "")
    assert "application/json" in content_type, (
        f"expected JSON response, got content-type={content_type!r} "
        f"and body={response.text!r}"
    )

    try:
        return response.json()
    except ValueError as exc:
        raise AssertionError(
            f"response body is not valid JSON: status={response.status_code}, body={response.text!r}"
        ) from exc
