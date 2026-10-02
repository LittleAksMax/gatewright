import datetime
import os
from collections.abc import Generator
from dataclasses import dataclass
from io import TextIOWrapper
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
class AWSConfig:
    def write(self, file: TextIOWrapper) -> None:
        raise NotImplementedError


@dataclass
class GithubConfig:
    access_token: str
    access_token_expires_at: datetime.datetime
    refresh_token: str
    refresh_token_expires_at: datetime.datetime

    def write(self, file: TextIOWrapper) -> None:
        file.writelines(
            [
                f"accessToken={self.access_token}",
                f"accessTokenExpiresAt={self.access_token_expires_at.strftime(DATE_FORMAT)}",
                f"refreshToken={self.refresh_token}",
                f"refreshTokenExpiresAt={self.refresh_token_expires_at.strftime(DATE_FORMAT)}",
            ]
        )


@dataclass
class Config:
    aws: AWSConfig | None
    gh: GithubConfig | None


class ConfigFileManager:
    def __init__(self, config: Config) -> None:
        self._config = config

    def __enter__(self) -> Self:
        config_root = Path.home() / Path(".gatewright")
        config_root.mkdir(mode=0o600)

        if self._config.gh is not None:
            github_file = config_root / "github"
            github_file.touch(mode=0o600, exist_ok=False)
            with open(github_file, "w") as f:
                self._config.gh.write(f)

        if self._config.aws is not None:
            aws_file = config_root / "aws"
            aws_file.touch(mode=0o600, exist_ok=False)
            with open(aws_file, "w") as f:
                self._config.aws.write(f)
        return self

    def __exit__(
        self,
        _exc_type: type[BaseException] | None,
        _exc_value: BaseException | None,
        _exc_traceback: TracebackType | None,
    ) -> None:
        config_dir = Path("~/.gatewright").expanduser()
        rmtree(config_dir, ignore_errors=True)


@pytest.fixture
def empty_config() -> Generator[Config]:
    config = Config(aws=None, gh=None)
    with ConfigFileManager(config):
        yield config


@pytest.fixture
def github_config_context(request):
    return getattr(request, "gh", None)


@pytest.fixture
def aws_config_context(request):
    return getattr(request, "aws", None)
