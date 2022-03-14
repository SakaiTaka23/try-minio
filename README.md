# try-minio

- [try-minio](#try-minio)
  - [操作](#操作)
  - [参考](#参考)

## 操作
**API**
- 画像一覧表示
  - http://127.0.0.1:5000/static?bucket=static
**minio**
- 管理画面
  - http://127.0.0.1:9001
  - パスワードはdocker-composeに書いたもの
- 画像表示
  - http://127.0.0.1:9000/バケット名/ファイル名
  - **管理者用のAPIとはポートが違う**

## 参考
- minio公式
https://min.io/
- dockerでの環境構築・操作
https://qiita.com/reflet/items/3e0f07bc9d64314515c1
- goの実装
https://penguin-note.tech/go-docker-minio/
