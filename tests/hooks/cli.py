from __future__ import annotations

import pytest


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
