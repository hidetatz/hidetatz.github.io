映画のホームページのURLについて語る---2018-01-09 22:25:23

本エントリでは、映画が公開された際に作られる特設サイトにおけるURLについて書いていく。

### パターン分け
URLについて注意深く見ていくと、いくつかのパターンがあることがわかる。

## 配給会社のドメイン + 映画名のパス

けっこう多いパターン。
例としては

* ドクター・ストレンジ
 - http://marvel.disney.co.jp/movie/dr-strange.html
* レゴバットマン
 - http://wwws.warnerbros.co.jp/legobatmanmovie/
* サラエヴォの銃声
 - http://www.bitters.co.jp/tanovic/sarajevo.html

など。
ディズニー、ワーナー、gagaなど、大手の会社がリリースした映画は、このパターンになることが多いと思われる。

このパターンのメリットは、ドメインの下に映画名がぶら下がる構造となるため、ドメインの管理がしやすいこと。
削除時はパスを削除してしまえばいいし、新たに登録したい際はパスを切ればいい。

デメリットは、パス名の切り方を決めにくいこと。例えば、同じ名前の映画を公開したいときはどうするか？
映画のタイトルが非常に長い場合、どのようなパス名にすればいいか？日本語タイトルは、英訳するか？など、考えることがたくさんある。
(これはこの形式に限らないが。)

## シリーズ名のドメイン + 年度のパス

あまり見ない形式だが、筆者はこれが好き。例としては、

* ドラえもん のび太の南極カチコチ大冒険
 - http://doraeiga.com/2017/

がある。
このパターンは、パスの変更に強い。問題が発生するとすれば、同一年度に映画を2本作る場合どうするか？だが、
ドラえもんはそういうことないんだと思う。

## 映画名 + movie

これもめちゃくちゃ多い。

* スノーデン
 - http://www.snowden-movie.jp/
* セル
 - http://cell-movie.jp/
* SING
 - http://sing-movie.jp/
* スプリット
 - http://split-movie.jp/

など。
このパターンは、映画のタイトルが一般名詞の場合につくのだと考えている。 `split.jp` だと、スプリットという言葉があるだけに、
映画のURLであることが伝わりづらい、という意図ではないだろうか。

このパターンは一見わかりやすいが、ドメインがどんどん増えていって管理が煩雑では？と筆者は考えている。新規登録もまあまあめんどうだし…

## 映画名をダイレクトにURLにするケース


* コクソン
 - http://kokuson.com/
* バーニングオーシャン
 - http://burningocean.jp/
* ハクソーリッジ
 - http://hacksawridge.jp/
* 光をくれた人
 - http://hikariwokuretahito.com/

などがある。日本語でもローマ字表記でそのままURLにしているのがポイント。
筆者としては、 `-movie` をつけるくらいならこれでよくないか？とは思っている。

## 映画名を省略するケース

意外と省略パターンは多い。


* ゴーストインザシェル
 - http://ghostshell.jp/
* 歓びのトスカーナ
 - http://yorokobino.com/
* あさがくるまえに
 - http://www.reallylikefilms.com/asakuru
* 愛を綴る女
 - http://aiotsuzuru.com/

など。タイトルが非常に長いなら仕方ない気もする。Twitterでつぶやくときのこととかも考慮しているのだと思う。

## 原題を用いているケース

けっこうレア。

* マリアンヌ
 - http://www.alliedmovie.com/
* 女神の見えざる手
 - http://miss-sloane.jp/
* メッセージ
 - http://www.bd-dvd.sonypictures.jp/arrival/

原題をそのまま使っているケース。
`.jp` ドメインのように、TLDでリージョンを分けているパターンと、完全に多言語対応していないケースがある。

## サブドメインにしているケース

相当レア。

* ジグソウ: ソウ・レガシー
 - http://jigsaw.asmik-ace.co.jp/
* 米軍が最も恐れた男 その名は、カメジロー
 - http://www.kamejiro.ayapro.ne.jp/

## 独自のURLのケース

多い。省略に近いが、省略とも言い難いものとも言える。

* 僕とカミンスキーの旅
 - http://meandkaminski.com/
* 不都合な真実2
 - http://futsugou2.jp/
* ヒトラーへの285枚の手紙
 - http://hitler-hagaki-movie.com/
* ジュリーと恋と靴工場
 - http://julie-kutsu.com/
* メアリと魔女の花
 - http://www.maryflower.jp/

映画のタイトルを若干変えて、読みやすく覚えやすくしているケースと考えている。
個人的には、 `メアリと魔女の花` がこのパターンなのはけっこう意外だが、よく考えたらジブリじゃないしそんなもんか、と思った。
個人的には早めにスタジオポノック配下に持っていくほうがいいと思っている。

## シリーズの考慮がないケース

アニメにありがち。


* 怪盗グルーのミニオン大脱走
 - http://minions.jp/
* クレヨンしんちゃん 襲来!! 宇宙人シリリ
 - http://www.shinchan-movie.com/
* 名探偵コナン から紅のラブレター
 - http://www.conan-movie.jp/

この問題点は、常に最新のもののみトップに表示され、古いものについては上書きされてしまうこと。
たぶん常に最新のものさえあればいいって考え方だと思う。

## その他

### \<title\> が 絶賛公開中！になってる

* パトリオットデイ
 - http://www.patriotsday.jp/

なおしたほうが良いと思う。

### カーソルがかっこいいもの

* フリーファイヤー
 - http://freefire.jp/

昔のホームページっぽい。クリックすると発砲音が鳴るのが最高。

### ドメインは残っていても、コンテンツが削除されているもの

* トンネル
 - http://tunnel-movie.net/

Apacheのエラーが出る。サーバーの容量不足とかなのかな。

### トップレベルドメイン

`.jp` `.com` がほとんど。
まれに `.info` `.asia` がある。

* バンコクナイツ
 - http://www.bangkok-nites.asia/

アジアっぽくて良い。

### https対応
1件もなし。(静的コンテンツしかないし、ユーザーがなにか入力することも基本ないので、OKだと思う。)

### あまり見ない文字が入っている

ドット

* ル・コルビジェとアイリーン
 - http://www.transformer.co.jp/m/lecorbusier.eileen/

エクスクラメーション

* リュミエール
 - http://gaga.ne.jp/lumiere!/

### 所感

URL設計は難しい。  
シリーズ物なら前述のドラえもんパターンが良いと思うが、それ以外だといい命名はけっこう難しそう。  
加えて、映画の公開が終わってしばらくしたらサイト自体閉鎖して良いような性質のものなので、URLなんてなんでもいい、というのも少なからずあるとは思う。  
映画ファンとしては、 `配給会社名ドメイン + 映画名のパス` のほうが、基本ずっと残るはずなので嬉しさはある。

## URL一覧

以下に、[筆者が昨年映画館で見た映画](https://yagi5.com/2018/01/03/2017-movies/)について、URL一覧を記す。
なお、昨年見たとはいえ、昨年公開されたとは限らないため、特設サイトが存在しないものはスルーする。

* 沈黙 -サイレンス-
 - http://chinmoku.jp/
* スノーデン
 - http://www.snowden-movie.jp/
* ドクター・ストレンジ
 - http://marvel.disney.co.jp/movie/dr-strange.html
* マリアンヌ
 - http://www.alliedmovie.com/
* セル
 - http://cell-movie.jp/
* ララランド
 - http://gaga.ne.jp/lalaland/
* 素晴らしきかな、人生
 - http://wwws.warnerbros.co.jp/subarashiki-movie/
* We Are X
 - http://www.wearexfilm.jp/
* ドラえもん のび太の南極カチコチ大冒険
 - http://doraeiga.com/2017/
* モアナと伝説の海
 - http://www.disney.co.jp/movie/moana/about.html
* コクソン
 - http://kokuson.com/
* SING
 - http://sing-movie.jp/
* パッセンジャー
 - http://www.bd-dvd.sonypictures.jp/passenger/
* キングコング 髑髏島の巨人
 - http://wwws.warnerbros.co.jp/kingkong/
* ムーンライト
 - http://moonlight-movie.jp/
* ジャッキー/ファーストレディ 最後の使命
 - http://jackie-movie.jp/
* レゴバットマン
 - http://wwws.warnerbros.co.jp/legobatmanmovie/
* ハードコア
 - http://hardcore-eiga.com/
* ライオン
 - http://gaga.ne.jp/lion/
* ゴーストインザシェル
 - http://ghostshell.jp/
* グレートウォール
 - http://greatwall-movie.jp/ (日)
* パージ
 - http://purge-movie.com/
* 美女と野獣
 - http://www.disney.co.jp/movie/beautyandbeast.html
* バーニングオーシャン
 - http://burningocean.jp/
* ワイルド・スピード ICE BREAK
 - http://wildspeed-official.jp/
* メットガラ ドレスをまとった美術館
 - http://metgala-movie.com/
* 汚れたミルク
 - http://www.bitters.co.jp/tanovic/milk.html
* サラエヴォの銃声
 - http://www.bitters.co.jp/tanovic/sarajevo.html
* 午後8時の訪問者
 - http://www.bitters.co.jp/pm8/
* マイビューティフルガーデン
 - http://my-beautiful-garden.com/
* スウィート17モンスター
 - http://www.sweet17monster.com/
* 僕とカミンスキーの旅
 - http://meandkaminski.com/
* フリーファイヤー
 - http://freefire.jp/
* カフェソサエティ
 - http://www.longride.jp/movie-cafesociety/
* スプリット
 - http://split-movie.jp/
* パーソナルショッパー
 - http://personalshopper-movie.com/
* マンチェスターバイザシー
 - http://manchesterbythesea.jp/
* メッセージ
 - http://www.bd-dvd.sonypictures.jp/arrival/
* パトリオットデイ
 - http://www.patriotsday.jp/
* セールスマン
 - http://www.thesalesman.jp/
* レイルロードタイガー
 - http://railroadtiger-movie.jp/
* 22年目の告白
 - http://wwws.warnerbros.co.jp/22-kokuhaku/
* 映画 山田孝之3D
 - http://www.tv-tokyo.co.jp/yamada_cannes/smp/takayuki_yamada_the_movie_3d/
* TAP THE LAST SHOW
 - http://www.tap-movie.jp/
* キングアーサー
 - http://wwws.warnerbros.co.jp/king-arthur/
* フィフティシェイズダーカー
 - http://fiftyshadesmovie.jp/
* ハクソーリッジ
 - http://hacksawridge.jp/
* パイレーツ・オブ・カリビアン
 - http://www.disney.co.jp/movie/pirates.html
* ジョンウィック
 - http://johnwick.jp/
* メアリと魔女の花
 - http://www.maryflower.jp/
* ライフ
 - http://www.bd-dvd.sonypictures.jp/life/
* カーズ
 - http://www.disney.co.jp/movie/cars-crossroad.html
* パワーレンジャー
 - http://www.power-rangers.jp/
* 怪盗グルーのミニオン大脱走
 - http://minions.jp/
* ビニー
 - http://vinny-movie.com/
* ザ・マミー
 - http://themummy.jp/
* トンネル
 - http://tunnel-movie.net/
* ラストプリンセス
 - http://lastprincess.info/
* ディストピア パンドラの少女
 - http://pandora-movie.jp/
* オリーブの樹は呼んでいる
 - http://olive-tree-jp.com/
* 光をくれた人
 - http://hikariwokuretahito.com/
* クレヨンしんちゃん 襲来!! 宇宙人シリリ
 - http://www.shinchan-movie.com/
* 名探偵コナン から紅のラブレター
 - http://www.conan-movie.jp/
* ヒトラーへの285枚の手紙
 - http://hitler-hagaki-movie.com/
* 歓びのトスカーナ
 - http://yorokobino.com/
* 幸せな人生の選択
 - http://www.finefilms.co.jp/shiawase/
* 残像
 - http://zanzou-movie.com/
* ボンジュール、アン
 - http://bonjour-anne.jp/
* わすれな草
 - http://www.gnome15.com/wasurenagusa/
* バンコクナイツ
 - http://www.bangkok-nites.asia/
* パターソン
 - http://paterson-movie.com/
* アイヒマンの後継者
 - http://next-eichmann.com/
* ファウンダー ハンバーガー帝国の秘密
 - http://thefounder.jp/
* ブランカとギター弾き
 - http://www.transformer.co.jp/m/blanka/
* トランスフォーマー
 - http://tf-movie.jp/
* ハイドリヒを撃て
 - http://shoot-heydrich.com/
* エル ELLE
 - http://gaga.ne.jp/elle/
* ザ・ウォール
 - http://thewall-movie.jp/
* スキップトレース
 - http://skiptrace-movie.jp/
* 新感染
 - http://shin-kansen.com/
* ダンケルク
 - http://wwws.warnerbros.co.jp/dunkirk/
* エイリアン: コヴェナント
 - http://www.foxmovies-jp.com/alien/
* スイス・アーミーマン
 - http://sam-movie.jp/
* アトミックブロンド
 - http://atomic-blonde.jp/
* 女神の見えざる手
 - http://miss-sloane.jp/
* ブレードランナー
 - http://www.bladerunner2049.jp/
* ゲットアウト
 - http://getout.jp/
* IT
 - http://wwws.warnerbros.co.jp/itthemovie/
* マイティ・ソー
 - http://marvel.disney.co.jp/movie/thor-br.html
* ジグソウ: ソウ・レガシー
 - http://jigsaw.asmik-ace.co.jp/
* ローガンラッキー
 - http://www.logan-lucky.jp/
* ジャスティスリーグ
 - http://wwws.warnerbros.co.jp/justiceleaguejp/news/news.html
* ドクターエクソシスト
 - http://dr-exorcist.com/
* オリエント急行殺人事件
 - http://www.foxmovies-jp.com/orient-movie/
* 不都合な真実2
 - http://futsugou2.jp/
* ジュリーと恋と靴工場
 - http://julie-kutsu.com/
* パッションフラメンコ
 - http://passion-flamenco.net/
* 夜明けの祈り
 - http://yoake-inori.com/
* アメイジングジャーニー
 - http://amazing-journey.jp/
* 静かなる情熱
 - http://dickinson-film.jp/
* 少女ファニーと運命の旅
 - http://shojo-fanny-movie.jp/
* 新世紀 パリ・オペラ座
 - http://gaga.ne.jp/parisopera/
* 愛を綴る女
 - http://aiotsuzuru.com/
* 米軍が最も恐れた男 その名は、カメジロー
 - http://www.kamejiro.ayapro.ne.jp/
* ナインイレヴン
 - http://nineeleven.jp/
* 素敵な遺産相続
 - http://www.finefilms.co.jp/isan/
* ギフテッド
 - http://www.foxmovies-jp.com/gifted/
* はじまりの街
 - http://www.crest-inter.co.jp/hajimarinomachi/
* リュミエール
 - http://gaga.ne.jp/lumiere!/
* わたしたち
 - http://www.watashitachi-movie.com/
* ネルーダ
 - http://neruda-movie.jp/
* セザンヌと過ごした時間
 - http://www.cetera.co.jp/cezanne/
* ハートストーン
 - http://www.magichour.co.jp/heartstone/
* 皆はこう呼んだ、鋼鉄ジーグ
 - http://www.zaziefilms.com/jeegmovie/
* ギミーデンジャー
 - http://movie-gimmedanger.com/
* 破裏拳ポリマー
 - http://polimar.jp/
* ル・コルビジェとアイリーン
 - http://www.transformer.co.jp/m/lecorbusier.eileen/
* あさがくるまえに
 - http://www.reallylikefilms.com/asakuru
* 笑う故郷
 - http://www.waraukokyo.com/
* Avicii TRUE STORIES
 - http://liveviewing.jp/contents/avicii/
* 婚約者の友人
 - http://frantz-movie.com/
* スターウォーズ
 - http://starwars.disney.co.jp/movie/lastjedi.html
* フラットライナーズ
 - http://www.flatliners.jp/
