from __future__ import annotations

import pytest

from tests.helpers.assertions import assert_json_response

pytestmark = [pytest.mark.api, pytest.mark.smoke]


def test_health_returns_ok(api_client) -> None:
    # Pin the first smoke test to the real health contract returned by the
    # managed Go service, not the outdated Swagger example.
    response = api_client.get("/health")

    assert response.status_code == 200
    payload = assert_json_response(response)
    assert payload == {"message": "ok"}
