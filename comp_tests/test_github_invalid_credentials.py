import datetime

from conftest import (
    AwsConfig,
    Config,
    GithubConfig,
    config_fixture,
    run_command,
    with_aws_config,
    with_github_config,
)


class TestGithubNoConfig:
    GITHUB_NOT_FOUND_ERROR = "github not linked, use `gwt link gh` to link"
    
    def test_no_config_errors_linking_and_exits_with_1(self, config_fixture: Config) -> None:
        command_result = run_command("gwt init new_project")
        assert command_result.stderr.strip() == TestGithubNoConfig.GITHUB_NOT_FOUND_ERROR
        assert command_result.returncode == 1


class TestGithubExpiredRefreshToken:
    GITHUB_EXPIRED_REFRESH_TOKEN_ERROR = "github credentials expired, use `gwt link gh` to link"

    days_ago = datetime.datetime.now() - datetime.timedelta(days=5)

    @with_aws_config(aws_config=AwsConfig.new())
    @with_github_config(
        github_config=GithubConfig.new(
            refresh_token_expires_at=datetime.datetime.now() - datetime.timedelta(days=5)
        )
    )
    def test_expired_config_errors_linking_and_exits_with_1(self, config_fixture: Config) -> None:
        command_result = run_command("gwt init new_project")
        assert (
            command_result.stderr.strip()
            == TestGithubExpiredRefreshToken.GITHUB_EXPIRED_REFRESH_TOKEN_ERROR
        )
        assert command_result.returncode == 1
