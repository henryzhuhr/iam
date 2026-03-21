from __future__ import annotations

from collections.abc import Mapping
from dataclasses import dataclass, field
from typing import Any

import httpx


@dataclass(frozen=True, slots=True)
class AuthSession:
    # Store auth material separately so tests can opt into headers without
    # rebuilding client configuration by hand.
    token: str | None = None
    tenant_id: str | None = None

    def as_headers(self) -> dict[str, str]:
        headers: dict[str, str] = {}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        if self.tenant_id:
            headers["X-Tenant-Id"] = self.tenant_id
        return headers


@dataclass(frozen=True, slots=True)
class RequestContext:
    # Layer request-scoped metadata on top of auth state before it reaches the
    # shared API client.
    environment: str | None
    auth: AuthSession
    extra_headers: Mapping[str, str] = field(default_factory=dict)

    def build_headers(self) -> dict[str, str]:
        headers = {
            "Accept": "application/json",
        }
        headers.update(self.auth.as_headers())
        headers.update(self.extra_headers)
        return headers


class APIClient:
    # Keep tests focused on endpoint behavior instead of repetitive httpx setup
    # and transport-level error handling.
    def __init__(
        self,
        *,
        base_url: str,
        default_headers: Mapping[str, str],
        timeout_seconds: float,
    ) -> None:
        self._client = httpx.Client(
            base_url=base_url,
            headers=dict(default_headers),
            timeout=timeout_seconds,
        )

    def __enter__(self) -> "APIClient":
        return self

    def __exit__(self, exc_type: object, exc: object, exc_tb: object) -> None:
        self.close()

    def close(self) -> None:
        self._client.close()

    def request(
        self,
        method: str,
        path: str,
        *,
        headers: Mapping[str, str] | None = None,
        **kwargs: Any,
    ) -> httpx.Response:
        request = self._client.build_request(method, path, headers=headers, **kwargs)
        try:
            return self._client.send(request)
        except httpx.HTTPError as exc:
            # Surface transport failures as test assertions so the failure reads
            # like a broken test environment, not a raw client traceback.
            raise AssertionError(
                f"API request failed: {request.method} {request.url} -> {exc}"
            ) from exc

    def get(
        self,
        path: str,
        *,
        headers: Mapping[str, str] | None = None,
        **kwargs: Any,
    ) -> httpx.Response:
        return self.request("GET", path, headers=headers, **kwargs)

    def post(
        self,
        path: str,
        *,
        headers: Mapping[str, str] | None = None,
        **kwargs: Any,
    ) -> httpx.Response:
        return self.request("POST", path, headers=headers, **kwargs)
