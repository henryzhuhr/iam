from __future__ import annotations

import os
import signal
import subprocess
import time
from pathlib import Path

import httpx

from tests.support.config import ManagedAppConfig


class ManagedGoApp:
    # Own the lifecycle of a single `go run` process for the pytest session.
    def __init__(self, config: ManagedAppConfig) -> None:
        self._config = config
        self._repo_root = Path(__file__).resolve().parents[2]
        self._process: subprocess.Popen[str] | None = None

    def start(self, *, base_url: str) -> None:
        if self._process is not None:
            return

        entrypoint = self._repo_root / self._config.entrypoint
        if not entrypoint.exists():
            raise AssertionError(f"managed Go app entrypoint not found: {entrypoint}")

        config_file = self._repo_root / self._config.config_file
        if not config_file.exists():
            raise AssertionError(f"managed Go app config file not found: {config_file}")

        try:
            # Capture combined output so startup failures can surface a single
            # readable error message in pytest.
            self._process = subprocess.Popen(
                self.command,
                cwd=self._repo_root,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                start_new_session=True,
            )
        except OSError as exc:
            raise AssertionError(
                f"failed to start managed Go app with command `{self.command_display}`: {exc}"
            ) from exc

        self._wait_until_ready(base_url)

    def stop(self) -> None:
        if self._process is None:
            return

        process = self._process
        self._process = None

        # Terminate the whole process group so child processes created by
        # `go run` do not survive the pytest session.
        if process.poll() is None:
            try:
                os.killpg(process.pid, signal.SIGTERM)
            except ProcessLookupError:
                pass

            try:
                process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                try:
                    os.killpg(process.pid, signal.SIGKILL)
                except ProcessLookupError:
                    pass
                process.wait(timeout=5)

        if process.stdout is not None:
            process.stdout.close()

    @property
    def command(self) -> list[str]:
        return [
            "go",
            "run",
            self._config.entrypoint,
            "-f",
            self._config.config_file,
            *self._config.extra_args,
        ]

    @property
    def command_display(self) -> str:
        return " ".join(self.command)

    def _wait_until_ready(self, base_url: str) -> None:
        if self._process is None:
            raise AssertionError("managed Go app process was not created")

        health_url = f"{base_url.rstrip('/')}/health"
        deadline = time.monotonic() + self._config.startup_timeout_seconds
        last_error = "health check did not run yet"

        # The process can exist before the HTTP routes are ready, so wait for
        # the real health endpoint instead of only checking that the PID exists.
        while time.monotonic() < deadline:
            return_code = self._process.poll()
            if return_code is not None:
                output = self._collect_output()
                raise AssertionError(
                    "managed Go app exited before becoming healthy.\n"
                    f"command: {self.command_display}\n"
                    f"exit_code: {return_code}\n"
                    f"output:\n{output}"
                )

            try:
                response = httpx.get(health_url, timeout=1.0)
                if response.status_code == 200:
                    return
                last_error = (
                    f"health check returned status {response.status_code}: {response.text!r}"
                )
            except httpx.HTTPError as exc:
                last_error = str(exc)

            time.sleep(0.2)

        self._terminate_process_group()
        output = self._collect_output()
        raise AssertionError(
            "managed Go app did not become healthy in time.\n"
            f"command: {self.command_display}\n"
            f"health_url: {health_url}\n"
            f"last_error: {last_error}\n"
            f"output:\n{output}"
        )

    def _terminate_process_group(self) -> None:
        if self._process is None or self._process.poll() is not None:
            return

        # Reuse the same termination policy during startup failures/timeouts so
        # cleanup behavior is consistent with the normal session teardown.
        try:
            os.killpg(self._process.pid, signal.SIGTERM)
        except ProcessLookupError:
            return

        try:
            self._process.wait(timeout=5)
        except subprocess.TimeoutExpired:
            try:
                os.killpg(self._process.pid, signal.SIGKILL)
            except ProcessLookupError:
                pass
            self._process.wait(timeout=5)

    def _collect_output(self) -> str:
        if self._process is None:
            return ""

        # Only collect buffered output after the process has stopped; on a live
        # process this would block and hide the real startup problem.
        try:
            output, _ = self._process.communicate(timeout=1)
        except subprocess.TimeoutExpired:
            return ""

        return output.strip()
