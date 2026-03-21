from __future__ import annotations

import os
import shlex
from dataclasses import dataclass
from typing import Any, Protocol
from urllib.parse import urlparse

import pytest

DEFAULT_BASE_URL = "http://127.0.0.1:8080/api"
DEFAULT_TIMEOUT_SECONDS = 5.0
DEFAULT_APP_ENTRY = "app/main.go"
DEFAULT_APP_CONFIG = "etc/dev.yaml"
DEFAULT_APP_STARTUP_TIMEOUT_SECONDS = 30.0


@dataclass(frozen=True, slots=True)
class ManagedAppConfig:
    # Separate app lifecycle knobs from request config so pytest can either
    # launch the service itself or point at an existing deployment.
    enabled: bool
    entrypoint: str
    config_file: str
    extra_args: tuple[str, ...]
    startup_timeout_seconds: float


@dataclass(frozen=True, slots=True)
class TestConfig:
    base_url: str
    environment: str | None
    timeout_seconds: float
    token: str | None
    tenant_id: str | None
    managed_app: ManagedAppConfig


class PytestConfigLike(Protocol):
    # `load_test_config` only needs pytest's option lookup surface, so accept
    # any object that provides a compatible `getoption` method.
    def getoption(self, name: str) -> Any: ...


def _env_flag(name: str) -> bool:
    value = os.getenv(name)
    if value is None:
        return False
    return value.strip().lower() in {"1", "true", "yes", "on"}


def load_test_config(pytestconfig: PytestConfigLike) -> TestConfig:
    # Command-line options win over env vars, which win over baked-in defaults.
    # That keeps local debugging flexible without hiding the project defaults.
    raw_base_url = (
        pytestconfig.getoption("base_url")
        or os.getenv("IAM_BASE_URL")
        or DEFAULT_BASE_URL
    )
    base_url = raw_base_url.rstrip("/")
    parsed = urlparse(base_url)
    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        raise pytest.UsageError(
            f"invalid IAM API base URL {raw_base_url!r}; expected http(s)://host[:port][/prefix]"
        )

    # These flags decide whether pytest owns the Go process lifecycle or simply
    # points at an already-running environment.
    use_existing_service = pytestconfig.getoption("use_existing_service") or _env_flag(
        "IAM_USE_EXISTING_SERVICE"
    )
    raw_startup_timeout = (
        pytestconfig.getoption("app_startup_timeout")
        or os.getenv("IAM_APP_STARTUP_TIMEOUT")
        or DEFAULT_APP_STARTUP_TIMEOUT_SECONDS
    )
    try:
        startup_timeout_seconds = float(raw_startup_timeout)
    except ValueError as exc:
        raise pytest.UsageError(
            f"invalid app startup timeout {raw_startup_timeout!r}; expected a number of seconds"
        ) from exc
    if startup_timeout_seconds <= 0:
        raise pytest.UsageError("app startup timeout must be greater than 0 seconds")

    environment = pytestconfig.getoption("env") or os.getenv("IAM_ENV")
    return TestConfig(
        base_url=base_url,
        environment=environment,
        timeout_seconds=DEFAULT_TIMEOUT_SECONDS,
        token=os.getenv("IAM_TOKEN"),
        tenant_id=os.getenv("IAM_TENANT_ID"),
        managed_app=ManagedAppConfig(
            # Default to a pytest-managed Go process so service boot is covered
            # by the test run; callers can opt out with --use-existing-service.
            enabled=not use_existing_service,
            entrypoint=(
                pytestconfig.getoption("app_entry")
                or os.getenv("IAM_APP_ENTRY")
                or DEFAULT_APP_ENTRY
            ),
            config_file=(
                pytestconfig.getoption("app_config")
                or os.getenv("IAM_APP_CONFIG")
                or DEFAULT_APP_CONFIG
            ),
            extra_args=tuple(
                shlex.split(
                    pytestconfig.getoption("go_run_extra_args")
                    or os.getenv("IAM_GO_RUN_EXTRA_ARGS")
                    or ""
                )
            ),
            startup_timeout_seconds=startup_timeout_seconds,
        ),
    )
