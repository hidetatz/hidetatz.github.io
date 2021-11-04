SAMLについて調べた---2017-10-24 19:03:25

### SAMLとはなにか？

SAML とは **仕様** のことである。
シングルサインオンを実現するためのものがSAMLだ。

用語
* IdP(Identity Provider): 認証情報を提供する側。OneLoginのようなサービスもあるが、OSSのOpenAMなどを使って自分たちで構築することもできる。
* SP(Service Provider): 認証情報を利用してシングルサインオンさせてあげる側。Dropbx、Google、AWS IAMなどがあたる。

具体的な流れは、
1. 最初にIdPとSP間でトラストサークル(信頼関係)を構築
2. IdP・Spそれぞれで設定(連携する相手を教える)
3. ログイン時には一旦SP側にアクセスする(このへんの動きはSPによって異なる。例えば、SalesforceはsalesforceのサブドメインにアクセスするとOpenAMにリダイレクトされる、Dropboxの場合はメールアドレスのみ入力してログインするとOpenAMにリダイレクトされる、など。)

1.の内部的な動きは以下。
* IdP(OpenAMというOSSが動いているサーバ)が発行した証明書をSP(例: Dropbox)に渡す
* シングルサインオンしたいSP(=Dropbox)からもらった証明書をIdPに登録する
* つまりお互いがその時発行した証明書を交換する

これにより、証明書の有無でシングルサインオン時に正しくユーザがIdPにログインしていることを認可できる。

### OpenIDとはなにが違うのか？

OpenIDもできることは同じだが、仕様としての細かい部分が違う。

SAML仕様は、運用前に、信頼関係（トラストサークル）を構築する必要がある（SPとIdP間でのメタデータの交換、証明書の交換）。OpenID仕様では、信頼関係を構築する必要がない。
SAML仕様では、アサーションという認証情報（XML形式）を規定している。OpenID仕様は、アサーションという形では規定されておらず、キーと値のペアで認証情報をOPとRP間で交換する。
など。基本的にSAMLのほうがセキュアとのこと。