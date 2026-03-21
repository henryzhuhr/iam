from __future__ import annotations

from collections.abc import Iterator

import pytest

from tests.support.app import ManagedGoApp
from tests.support.client import APIClient, AuthSession, RequestContext
from tests.support.config import TestConfig, load_test_config


def pytest_addoption(parser: pytest.Parser) -> None:
    # Keep runtime overrides in one place so tests, local shells, and CI all
    # use the same parameter surface.
    group = parser.getgroup("iam-api")
    group.addoption(
        "--base-url",
        action="store",
        default=None,
        help="Base URL for the IAM API, for example http://127.0.0.1:8080/api",
    )
    group.addoption(
        "--use-existing-service",
        action="store_true",
        default=False,
        help="Do not start the Go service inside pytest; target an already-running service.",
    )
    group.addoption(
        "--app-entry",
        action="store",
        default=None,
        help="Go entrypoint to run when pytest manages the app process. Defaults to app/main.go.",
    )
    group.addoption(
        "--app-config",
        action="store",
        default=None,
        help="Config file passed to -f when pytest manages the app process. Defaults to etc/dev.yaml.",
    )
    group.addoption(
        "--go-run-extra-args",
        action="store",
        default=None,
        help="Extra arguments appended after `go run <entry> -f <config>`.",
    )
    group.addoption(
        "--app-startup-timeout",
        action="store",
        default=None,
        help="Seconds to wait for the managed Go app to become healthy.",
    )
    group.addoption(
        "--env",
        action="store",
        default=None,
        help="Logical environment name reserved for future test env mapping.",
    )


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
) -> APIClient:
    # `managed_go_app` is injected for its startup side effect; the client only
    # needs the resolved request configuration after the service is ready.
    with APIClient(
        base_url=test_config.base_url,
        default_headers=request_context.build_headers(),
        timeout_seconds=test_config.timeout_seconds,
    ) as client:
        yield client
