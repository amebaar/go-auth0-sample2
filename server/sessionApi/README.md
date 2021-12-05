# session-api
[Resource Owner Password Flow](https://auth0.com/docs/authorization/flows/resource-owner-password-flow) を用いてauth0でログインするサンプル。

## 動作確認
1. セッションを発行
```shell
$ curl -i -c cookie.txt http://localhost/state
{"state":"wIzwtxAQxte3dCtr5hVUDUfkaqzkrdgkv288dytrJ2Q="}
```
2. まだ認証前なのでsession（ユーザ情報）は取得できない
```shell
$ curl -i -b cookie.txt http://localhost/session
```
3. あらかじめauth0に作成しておいたユーザでログイン（stateは1で取得したものを指定）
```shell
$ curl -i -X POST -H "Content-Type: application/json" -d '{
    "username": "YOUR_NAME (must be email address)",
    "password": "YOUR_PASSWORD",
    "state": "wIzwtxAQxte3dCtr5hVUDUfkaqzkrdgkv288dytrJ2Q=",
    "redirect_to": "http://localhost/session"
}' -b cookie.txt http://localhost/login
```
3. 303がサーバから返っていればOK（パスワードをあえて間違えた場合は401が返る）
4. ログイン後にsessionを取得すると、Claimが正しくとれていることが確認できる。
```shell
$ curl -i -b cookie.txt http://localhost/session
```