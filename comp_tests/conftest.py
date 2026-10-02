import datetime
import os
from collections.abc import Generator
from dataclasses import dataclass
from pathlib import Path
from shutil import rmtree
from types import TracebackType
from typing import Self

import pytest


def pytest_sessionstart(session: pytest.Session) -> None:
    if os.environ.get("GATEWRIGHT_COMP_TESTS") != "docker":
        pytest.exit(
            "comp_tests must be run through Docker Compose",
            returncode=2,
        )


DATE_FORMAT = "%Y-%m-%d %H:%M:%S"


@dataclass
class AwsConfig:
    def write(self, path: Path) -> None:
        raise NotImplementedError


@dataclass
class GithubConfig:
    access_token: str
    access_token_expires_at: datetime.datetime
    refresh_token: str
    refresh_token_expires_at: datetime.datetime

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
def github_config_fixture(request: pytest.FixtureRequest) -> GithubConfig | None:
    return getattr(request, "param", None)


@pytest.fixture
def aws_config_fixture(request: pytest.FixtureRequest) -> AwsConfig | None:
    return getattr(request, "param", None)


def with_github_config(github_config: GithubConfig):
    return pytest.mark.parametrize(github_config_fixture.__name__, (github_config,), indirect=True)


def with_aws_config(aws_config: AwsConfig):
    return pytest.mark.parametrize(aws_config_fixture.__name__, (aws_config,), indirect=True)


@pytest.fixture
def config_fixture(
    github_config_fixture: GithubConfig | None, aws_config_fixture: AwsConfig | None
) -> Generator[Config]:
    config = Config(github=github_config_fixture, aws=aws_config_fixture)
    with ConfigFileManager(config):
        yield config
