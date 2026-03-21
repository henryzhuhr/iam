from __future__ import annotations

from collections.abc import Iterator

import pytest

from tests.helpers.app import ManagedGoApp
from tests.helpers.client import APIClient, AuthSession, RequestContext
from tests.helpers.config import TestConfig, load_test_config


@pytest.fixture(scope="session")
def test_config(pytestconfig: pytest.Config) -> TestConfig:
    return load_test_config(pytestconfig)


@pytest.fixture(scope="session")
def auth_session(test_config: TestConfig) -> AuthSession:
    return AuthSession(token=test_config.token, tenant_id=test_config.tenant_id)


@pytest.fixture(scope="session")
def request_context(
    test_config: TestConfig,
    auth_session: AuthSession,
) -> RequestContext:
    return RequestContext(
        environment=test_config.environment,
        auth=auth_session,
        extra_headers={},
    )


@pytest.fixture(scope="session")
def managed_go_app(test_config: TestConfig) -> Iterator[ManagedGoApp | None]:
    if not test_config.managed_app.enabled:
        yield None
        return

    # One managed process per session keeps startup cost low while still making
    # service boot a first-class part of the API test flow.
    app = ManagedGoApp(test_config.managed_app)
    app.start(base_url=test_config.base_url)
    try:
        yield app
    finally:
        app.stop()


@pytest.fixture()
def api_client(
    test_config: TestConfig,
    request_context: RequestContext,
    managed_go_app: ManagedGoApp | None,
) -> Iterator[APIClient]:
    # `managed_go_app` is injected for its startup side effect; the client only
    # needs the resolved request configuration after the service is ready.
    with APIClient(
        base_url=test_config.base_url,
        default_headers=request_context.build_headers(),
        timeout_seconds=test_config.timeout_seconds,
    ) as client:
        yield client
