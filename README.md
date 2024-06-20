# dapr cron issue

This repository is for reproducing an issue where the cron binding does not work properly in Container Apps' dapr.
This issue is reproduced in the local environment, but of course it is also reproduced in the Container App Environment.

## Summary

If there is a delay of a few seconds until the cron binding is completed, the following log is output and the cron binding fails.

```log
daprd-1              | time="2024-06-20T07:58:51.319209384Z" level=info msg="app has not subscribed to binding h1-binding." app_id=moris-cron-sample instance=1f5e885d3494 scope=dapr.runtime.processor.binding type=log ver=1.13.4
daprd-1              | time="2024-06-20T07:58:51.319232029Z" level=info msg="app has not subscribed to binding h2-binding." app_id=moris-cron-sample instance=1f5e885d3494 scope=dapr.runtime.processor.binding type=log ver=1.13.4
```

## Background

After dapr was updated in the production environment of Container Apps, a problem occurred where cron was not subscribed. To sum up, if there is a delay of a few seconds between the launch of the application and the binding of the dapr handler, the binding fails. This seems to be a change in behavior due to the version upgrade of dapr. In the reproduction application, this delay is reproduced by `sleep`, but in reality, it takes a few seconds to execute to obtain the token of Managed identity and access to Key valut.

## Dapr version

| version | description                                                           |
| ------- | --------------------------------------------------------------------- |
| 1.11.6  | The version that was used in Container Apps before the update in May  |
| 1.12.5  | The version that was updated in May. Currently used in Container Apps |
| 1.13.4  | The latest version according to https://github.com/dapr/dapr/releases |

## Reproduction procedure

Build the image to the container:

```sh
go install github.com/google/ko@latest
make build
```

### Successful case

With `1.11.6`, cron binding succeeds and operates normally even if there is a delay of a few seconds.

```sh
# success
DAPR_VERSION=1.11.6 ZZZ_SLEEP=0s make up
# success
DAPR_VERSION=1.11.6 ZZZ_SLEEP=5s make up
```

### Failure case

With `1.12.5` or the latest version `1.13.4`, if there is a delay of a few seconds, the cron binding fails. However, if you set the delay to `0s`, it works without any issues.

Example of execution:
```sh
# success
DAPR_VERSION=1.12.5 ZZZ_SLEEP=0s make up
# failure
DAPR_VERSION=1.12.5 ZZZ_SLEEP=5s make up
# success
DAPR_VERSION=1.13.4 ZZZ_SLEEP=0s make up
# failure
DAPR_VERSION=1.13.4 ZZZ_SLEEP=5s make up
```

## Logs

For detailed logs, please refer to the log folder.
