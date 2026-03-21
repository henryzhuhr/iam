from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any

import pytest

from tests.support.config import (
    DEFAULT_APP_CONFIG,
    DEFAULT_APP_ENTRY,
    DEFAULT_APP_STARTUP_TIMEOUT_SECONDS,
    DEFAULT_BASE_URL,
    DEFAULT_TIMEOUT_SECONDS,
    load_test_config,
)


@dataclass(slots=True)
class DummyPytestConfig:
    options: dict[str, Any] = field(default_factory=dict)

    def getoption(self, name: str) -> Any:
        return self.options.get(name)


def test_load_test_config_uses_defaults(monkeypatch: pytest.MonkeyPatch) -> None:
    for name in (
        "IAM_BASE_URL",
        "IAM_USE_EXISTING_SERVICE",
        "IAM_APP_ENTRY",
        "IAM_APP_CONFIG",
        "IAM_GO_RUN_EXTRA_ARGS",
        "IAM_APP_STARTUP_TIMEOUT",
        "IAM_ENV",
        "IAM_TOKEN",
        "IAM_TENANT_ID",
    ):
        monkeypatch.delenv(name, raising=False)

    config = load_test_config(DummyPytestConfig())

    assert config.base_url == DEFAULT_BASE_URL
    assert config.environment is None
    assert config.timeout_seconds == DEFAULT_TIMEOUT_SECONDS
    assert config.token is None
    assert config.tenant_id is None
    assert config.managed_app.enabled is True
    assert config.managed_app.entrypoint == DEFAULT_APP_ENTRY
    assert config.managed_app.config_file == DEFAULT_APP_CONFIG
    assert config.managed_app.extra_args == ()
    assert (
        config.managed_app.startup_timeout_seconds
        == DEFAULT_APP_STARTUP_TIMEOUT_SECONDS
    )


def test_load_test_config_uses_env_when_cli_is_absent(
    monkeypatch: pytest.MonkeyPatch,
) -> None:
    monkeypatch.setenv("IAM_BASE_URL", "https://env.example.test/api/")
    monkeypatch.setenv("IAM_USE_EXISTING_SERVICE", "yes")
    monkeypatch.setenv("IAM_APP_ENTRY", "cmd/iam/main.go")
    monkeypatch.setenv("IAM_APP_CONFIG", "etc/test.yaml")
    monkeypatch.setenv("IAM_GO_RUN_EXTRA_ARGS", "--log-level debug --dry-run")
    monkeypatch.setenv("IAM_APP_STARTUP_TIMEOUT", "18")
    monkeypatch.setenv("IAM_ENV", "staging")
    monkeypatch.setenv("IAM_TOKEN", "token-from-env")
    monkeypatch.setenv("IAM_TENANT_ID", "tenant-from-env")

    config = load_test_config(DummyPytestConfig())

    assert config.base_url == "https://env.example.test/api"
    assert config.environment == "staging"
    assert config.token == "token-from-env"
    assert config.tenant_id == "tenant-from-env"
    assert config.managed_app.enabled is False
    assert config.managed_app.entrypoint == "cmd/iam/main.go"
    assert config.managed_app.config_file == "etc/test.yaml"
    assert config.managed_app.extra_args == ("--log-level", "debug", "--dry-run")
    assert config.managed_app.startup_timeout_seconds == 18.0


def test_load_test_config_cli_options_override_environment(
    monkeypatch: pytest.MonkeyPatch,
) -> None:
    monkeypatch.setenv("IAM_BASE_URL", "https://env.example.test/api/")
    monkeypatch.setenv("IAM_USE_EXISTING_SERVICE", "0")
    monkeypatch.setenv("IAM_APP_ENTRY", "env/main.go")
    monkeypatch.setenv("IAM_APP_CONFIG", "etc/env.yaml")
    monkeypatch.setenv("IAM_GO_RUN_EXTRA_ARGS", "--from-env")
    monkeypatch.setenv("IAM_APP_STARTUP_TIMEOUT", "20")
    monkeypatch.setenv("IAM_ENV", "staging")

    config = load_test_config(
        DummyPytestConfig(
            {
                "base_url": "https://cli.example.test/custom/",
                "use_existing_service": True,
                "app_entry": "app/custom_main.go",
                "app_config": "etc/ci.yaml",
                "go_run_extra_args": "--listen=:9090 --feature-flag",
                "app_startup_timeout": "12.5",
                "env": "ci",
            }
        )
    )

    assert config.base_url == "https://cli.example.test/custom"
    assert config.environment == "ci"
    assert config.managed_app.enabled is False
    assert config.managed_app.entrypoint == "app/custom_main.go"
    assert config.managed_app.config_file == "etc/ci.yaml"
    assert config.managed_app.extra_args == ("--listen=:9090", "--feature-flag")
    assert config.managed_app.startup_timeout_seconds == 12.5


@pytest.mark.parametrize("raw_base_url", ["localhost:8080/api", "ftp://example.test"])
def test_load_test_config_rejects_invalid_base_url(
    raw_base_url: str,
) -> None:
    with pytest.raises(pytest.UsageError, match="invalid IAM API base URL"):
        load_test_config(DummyPytestConfig({"base_url": raw_base_url}))


def test_load_test_config_rejects_non_numeric_startup_timeout() -> None:
    with pytest.raises(pytest.UsageError, match="invalid app startup timeout"):
        load_test_config(DummyPytestConfig({"app_startup_timeout": "soon"}))


@pytest.mark.parametrize("raw_timeout", ["0", "-1"])
def test_load_test_config_rejects_non_positive_startup_timeout(
    raw_timeout: str,
) -> None:
    with pytest.raises(
        pytest.UsageError,
        match="app startup timeout must be greater than 0 seconds",
    ):
        load_test_config(DummyPytestConfig({"app_startup_timeout": raw_timeout}))
