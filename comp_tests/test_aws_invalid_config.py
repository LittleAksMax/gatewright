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


class TestAwsNoConfig:
    AWS_NOT_FOUND_ERROR = "aws not linked, use `gwt link aws` to link"

    @with_github_config(github_config=GithubConfig.new())
    @with_aws_config(aws_config=AwsConfig.new())
    def test_no_config_errors_linking_and_exits_with_1(self, config_fixture: Config) -> None:
        command_result = run_command("gwt init new_project")
        assert command_result.stderr.strip() == TestAwsNoConfig.AWS_NOT_FOUND_ERROR
        assert command_result.returncode == 1


class TestAwsExpiredRefreshToken:
    AWS_EXPIRED_REFRESH_TOKEN_ERROR = "aws credentials expired, use `gwt link aws` to link"

    @with_github_config(github_config=GithubConfig.new())
    @with_aws_config(aws_config=AwsConfig.new())
    def test_expired_config_errors_linking_and_exits_with_1(self, config_fixture: Config) -> None:
        command_result = run_command("gwt init new_project")
        assert (
            command_result.stderr.strip()
            == TestAwsExpiredRefreshToken.AWS_EXPIRED_REFRESH_TOKEN_ERROR
        )
        assert command_result.returncode == 1