import datetime
import os
import shlex
import subprocess
from collections.abc import Generator
from dataclasses import dataclass
from pathlib import Path
from shutil import rmtree
from types import TracebackType
from typing import Self, TypedDict, Unpack

import pytest


def pytest_sessionstart(session: pytest.Session) -> None:
    if os.environ.get("GATEWRIGHT_COMP_TESTS") != "docker":
        pytest.exit(
            "comp_tests must be run through Docker Compose",
            returncode=2,
        )


DATE_FORMAT = "%Y-%m-%d %H:%M:%S"


class _AwsConfigOverrides(TypedDict, total=False):
    pass


@dataclass
class AwsConfig:
    @classmethod
    def new(cls) -> AwsConfig:
        return cls()

    
    def write(self, path: Path) -> None:
        path.write_text("placeholder")


class _GithubConfigOverrides(TypedDict, total=False):
    access_token: str
    access_token_expires_at: datetime.datetime
    refresh_token: str
    refresh_token_expires_at: datetime.datetime


@dataclass
class GithubConfig:
    access_token: str
    access_token_expires_at: datetime.datetime
    refresh_token: str
    refresh_token_expires_at: datetime.datetime

    @classmethod
    def new(cls, **overrides: Unpack[_GithubConfigOverrides]) -> Self:
        now = datetime.datetime.now()
        return cls(
            access_token=overrides.get("access_token", "ACCESS_TOKEN"),
            access_token_expires_at=overrides.get(
                "access_token_expires_at", now + datetime.timedelta(minutes=15)
            ),
            refresh_token=overrides.get("refresh_token", "REFRESH_TOKEN"),
            refresh_token_expires_at=overrides.get(
                "refresh_token_expires_at", now + datetime.timedelta(days=7)
            ),
        )

    def write(self, path: Path) -> None:
        path.write_text(
            f"accessToken={self.access_token}\n"
            f"accessTokenExpiresAt={self.access_token_expires_at.strftime(DATE_FORMAT)}\n"
            f"refreshToken={self.refresh_token}\n"
            f"refreshTokenExpiresAt={self.refresh_token_expires_at.strftime(DATE_FORMAT)}\n"
        )


@dataclass
class Config:
    aws: AwsConfig | None
    github: GithubConfig | None


class ConfigFileManager:
    def __init__(self, config: Config) -> None:
        self._config = config
        self._root = Path.home() / ".gatewright"

    def __enter__(self) -> Self:
        rmtree(self._root, ignore_errors=True)
        self._root.mkdir(mode=0o700)
        if self._config.github is not None:
            self._config.github.write(self._root / "github")
        if self._config.aws is not None:
            self._config.aws.write(self._root / "aws")
        return self

    def __exit__(
        self,
        _exc_type: type[BaseException] | None,
        _exc_value: BaseException | None,
        _exc_traceback: TracebackType | None,
    ) -> None:
        rmtree(self._root, ignore_errors=True)


@pytest.fixture
def _github_config_fixture(request: pytest.FixtureRequest) -> GithubConfig | None:
    return getattr(request, "param", None)


@pytest.fixture
def _aws_config_fixture(request: pytest.FixtureRequest) -> AwsConfig | None:
    return getattr(request, "param", None)


def with_github_config(github_config: GithubConfig):
    return pytest.mark.parametrize(_github_config_fixture.__name__, (github_config,), indirect=True)


def with_aws_config(aws_config: AwsConfig):
    return pytest.mark.parametrize(_aws_config_fixture.__name__, (aws_config,), indirect=True)


@pytest.fixture
def config_fixture(
    _github_config_fixture: GithubConfig | None, _aws_config_fixture: AwsConfig | None
) -> Generator[Config]:
    config = Config(github=_github_config_fixture, aws=_aws_config_fixture)
    with ConfigFileManager(config):
        yield config


def run_command(command: str, stdin: str = "") -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        shlex.split(command),
        input=stdin,
        capture_output=True,
        text=True,
        timeout=5,
    )
