# dapr cron issue

Container Apps の dapr で cron binding を利用しているが、 binding が正しく動作しない問題を再現するためのリポジトリである。

## サマリ

cron binding が完了するまでに数秒の遅延があると、以下のログを出力して cron binding が失敗する。

```log
daprd-1              | time="2024-06-20T07:28:25.876900176Z" level=info msg="app has not subscribed to binding h2-binding." app_id=moris-cron-sample instance=bd67e8399137 scope=dapr.runtime.processor.binding type=log ver=1.13.4
daprd-1              | time="2024-06-20T07:28:25.87691298Z" level=info msg="dapr initialized. Status: Running. Init Elapsed 4643ms" app_id=moris-cron-sample instance=bd67e8399137 scope=dapr.runtime type=log ver=1.13.4
```


## 背景

Container Apps の本番環境で dapr が更新された後、cron binding されない問題が発生したが、これらをローカル環境で再現できるところまで落とし込めた。結論から言うと、アプリケーションの起動から dapr の handler がバインドされるまでの間に数秒の遅延があると、binding に失敗する。これは dapr のバージョンアップによる挙動の変更と思われる。再現アプリではこの遅延を`sleep`によって再現しているが、実際は Managed identity のトークン取得や Key valut のアクセスで数秒の実行時間を要している。

## dapr バージョン

| version | description                                                               |
| ------- | ------------------------------------------------------------------------- |
| 1.11.6  | 5月にアップデートされる前に container apps で利用されていたバージョン     |
| 1.12.5  | 5月にアップデートされたバージョン。現時点のContainer Appsで利用されている |
| 1.13.4  | https://github.com/dapr/dapr/releases による最新版                        |

## 再現手順

コンテにメージをビルドする

```sh
go install github.com/google/ko@latest
make build
```

### 正常に動作する

`1.11.6` では、数秒の遅延があっても問題なく cron binding は成功し動作する。

```sh
# success
DAPR_VERSION=1.11.6 ZZZ_SLEEP=0s make up 
# success
DAPR_VERSION=1.11.6 ZZZ_SLEEP=5s make up
```

### 失敗する場合

`1.12.5` もしくは、最新版の `1.13.4` では、数秒の遅延があると cron binding が失敗する。ただし、遅延を `0s` にすると問題なく動作する。


実行例
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

## ログ

詳細なログは、log folder を参照のこと。


以上
