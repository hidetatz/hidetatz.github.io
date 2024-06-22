import re
import urllib.request

content = """![PXL_20240620_102157852](https://github.com/hidetatz/hidetatz.github.io/assets/60682957/9c9bfd85-c5e2-4bc9-b8ca-e81783bd9ce8)

はま寿司。

![PXL_20240620_102326807](https://github.com/hidetatz/hidetatz.github.io/assets/60682957/67944768-0a5f-4f8c-824b-296efcf1a637)

寿司に囲まれてる。
私は寿司に囲まれるのが好きだ。

店で待っている時、偶然出くわしたような家族同士が話していた。全然喋らない子供に対して、父親のような人が「どうして喋らないんだ、恥ずかしいのか？」というようなことを聞いていた。

私もあまり喋らない子どもだったので、同じようなことを言われたことがあるなあ、と思い出した。

こっちからすると、恥ずかしいというより、ただ単に話したいと思っていないのだ。聞きたいこともないし、伝えたいこともないので、特に口を開かないのは合理的な選択の結果だ。でも、恥ずかしがって喋らないと思われている。

私は、人は他人からどう思われているかによって振る舞いが決まると思っている。他人から恥ずかしがりだと思われている人は、その人を驚かせないために、わざと恥ずかしがりのふりをする。それが楽だから。

自分のネイティブの言語ではない言葉を話しているときに明るくなる人というのがいるけど、あれは話し相手が自分に何の仮定も置いてないから、自分のありたい自分として自由に振る舞えるようになった結果なのだと考えている。

今日のはま寿司も、相変わらず炭治郎がうるさかったな。"""

images = re.finditer("!\[\S*\]\(\S+\)", content)
for i, image in enumerate(images):
    alt, url = image.group().lstrip("![").rstrip(")").split("](")
    urllib.request.urlretrieve(url, f"./{i}.jpg")
    replace = f"![{alt}](./{i}.jpg)"
    content = content.replace(image.group(), replace)
