# session-api
auth0でログインするサンプル。
- [Resource Owner Password Flow](https://auth0.com/docs/authorization/flows/resource-owner-password-flow) 
- [Social Login](https://auth0.com/docs/authorization/flows/call-your-api-using-the-authorization-code-flow)


## 動作確認
### Resource Owner Password Flow
0. `.env.sample`を、`.env`にリネームして適切に設定、サーバを起動
1. セッションを発行
```shell
$ curl -i -c cookie.txt http://localhost:8080/state
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
    "redirect_to": "http://localhost:8080/session"
}' -b cookie.txt http://localhost:8080/login
```
3. 303がサーバから返っていればOK（パスワードをあえて間違えた場合は401が返る）
4. ログイン後にsessionを取得すると、Claimが正しくとれていることが確認できる。
```shell
$ curl -i -b cookie.txt http://localhost:8080/session
```


### Social Login
ブラウザで実施

0. `.env.sample`を、`.env`にリネームして適切に設定し、サーバを起動
1. ブラウザで http://localhost:8080/state にリクエスト。返却された`state`をコピー。
2. 次にアクセス。http://localhost:8080/login?connection=google-oauth2&state=<stateをURLエンコードした文字列>&redirect=http://localhost:8080/session
3. Googleのログインフォームが表示されるのでログイン（初回の場合はサインアップとなる）
4. ログイン後、 http://localhost:8080/session に遷移すればOK