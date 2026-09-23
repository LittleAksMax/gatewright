import pytest
from dataclasses import dataclass
from types import TracebackType
from typing import Generator, Self
from shutil import rmtree
from pathlib import Path
from os import environ as os_environ


def pytest_sessionstart(session: pytest.Session) -> None:
    if os_environ.get("GATEWRIGHT_COMP_TESTS") != "docker":
        pytest.exit(
            "comp_tests must be run through Docker Compose",
            returncode=2,
        )


@dataclass
class AWSConfig:
    pass

@dataclass
class GithubConfig:
    pass

@dataclass
class Config:
    aws: AWSConfig | None
    gh: GithubConfig | None

class ConfigFileManager:
    def __init__(self, config: Config) -> None:
        self._config = config

    def __enter__(self) -> Self:
        # TODO: create relevant files
        return self

    def __exit__(
        self,
        _exc_type: type[BaseException] | None,
        _exc_value: BaseException | None,
        _exc_traceback: TracebackType | None,
    ) -> None:
        config_dir = Path("~/.gatewright").expanduser()
        rmtree(config_dir, ignore_errors=True)

@pytest.fixture()
def empty_config() -> Generator[Config, None, None]:
    config = Config(aws=None, gh=None)
    with ConfigFileManager(config):
        yield config
