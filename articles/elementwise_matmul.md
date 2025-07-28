title: Elementwise and reduce matmulとは何か
timestamp: 2025-07-28 09:43:01
lang: ja
---


また行列計算に関するブログで軽いやつ。
Elementwise and reduce Matmulについて。よく分かってなかったので調べてまとめる。「Elementwise and reduce Matmul」という言葉は存在しない可能性が高く、なんと言えばいいかわからないのでこのように記している。

まず、Elementwise-productという言葉も多分あり、これは単に行列を要素ごとに積を取る (いわゆる「アダマール積」) ものだ。
今から書くのはこれとは全く違う話で、普通に行列積 (dot-productとか言われる方) の話である。

Elementwise and reduce Matmulというのは何かというと、elementwise操作と、reduce操作だけで行列積が計算できるという話だ。

elementwiseオペレーションとは、テンソルに対する「要素ごと」の操作である。例えばベクトルの加算をするとき、

```python
a = array([1, 2, 3])
b = array([4, 5, 6])
a + b # [5, 7, 9]
```

これはaとbの要素ごとに足し算をしている。これがelementwiseな操作である。

reduceとは、要素を「減らす」オペレーションで、代表的なのはmaxやsumである。

例えば、

```python
a = array([1, 2, 3])
a.sum() # 6
```

これは次元が減っているのでreduce操作である。

行列積というのは、通常は3重ループが必要になる。

```c
for (int i = 0; i < ROWS; i++) {
    for (int j = 0; j < COLS; j++) {
        result[i][j] = 0;
        for (int k = 0; k < COLS; k++) {
            result[i][j] += firstMatrix[i][k] * secondMatrix[k][j];
        }
    }
}
```

これを、elementwise操作とreduce操作だけで (ループせずに) 代替可能か？というのが問いになる。
結論から言うとできる。次のようにやる。

実際にやってみると、例えば

```python
# 2 x 4
t1 = [
  [a, b, c, d],
  [e, f, g, h],
]

# 4 x 3
t2 = [
  [i, j, k],
  [l, m, n],
  [o, p, q],
  [r, s, t],
]

t1.dot(t2)
```

例えばこう言うことがしたいとする (結果のshapeは2 x 3になる) 。
まず普通にやると、結果は次のようになる。

```python
[
  [ai+bl+co+dr, aj+am+ap+as, ak+bn+cq+dt],
  [ei+fl+go+hr, ej+fm+gp+hs, ek+fn+gq+ht],
]

```

こうなる。

elementwise操作でやってみると、

まず、t1を変形して 2 x 3 x 4にする。reshape -> ブロードキャストになるので、

```python
# reshapeして2 x 1 x 4
[
  [
    [a, b, c, d],
  ],
  [
    [e, f, g, h],
  ],
]

```

こうして、


```python
# ブロードキャストして2 x 3 x 4
[
  [
    [a, b, c, d],
    [a, b, c, d],
    [a, b, c, d],
  ],
  [
    [e, f, g, h],
    [e, f, g, h],
    [e, f, g, h],
  ],
]
```

こうする。これはnumpy的には次元 (ストライド) 操作のみなので、メタデータの変更だけであり、内部のデータはいじる必要がない (という話 (numpyの内部構造) もいつか書きたい) 。

t2についても同じshapeにしたいので、これはreshape -> 転置 -> ブロードキャストでできる。すなわち、

```python
# reshapeして1 x 4 x 3
[
  [
    [i, j, k],
    [l, m, n],
    [o, p, q],
    [r, s, t],
  ]
]
```

こうして、

```python
# 転置して1 x 3 x 4
[
  [
    [i, l, o, r],
    [j, m, p, s],
    [k, n, q, t],
  ]
]
```

こう転置し、

```python
# ブロードキャストして2 x 3 x 4
[
  [
    [i, l, o, r],
    [j, m, p, s],
    [k, n, q, t],
  ],
  [
    [i, l, o, r],
    [j, m, p, s],
    [k, n, q, t],
  ],
]
```

こうなる。これも内部のデータはいじられておらず、メタデータだけ変えることで実現可能。

さて、ここまで来たら、この(2, 3, 4)テンソルをelementwiseに掛け算する。するとこうなる。

```python
[
  [
    [ai, bl, co, dr],
    [aj, bm, cp, ds],
    [ak, bn, cq, dt],
  ],
  [
    [ei, bl, co, dr],
    [ej, bm, cp, ds],
    [ek, bn, cq, dt],
  ],
]
```

で、これを一番内側の次元でsumを取る ( `np.sum(axis=2)` する感じ) 。すると (2, 3, 4) が (2, 3) になるので、こうなる。

```python
[
  [ai+bl+co+dr, aj+bm+cp+ds, ak+bn+cq+dt],
  [ai+bl+co+dr, aj+bm+cp+ds, ak+bn+cq+dt],
],
```

このように、上の方で計算した結果と一致していることがわかる。

これは何故こうなるかというと、

変形後のt1', t2'について

```python
t1'[i][j] = t1のi行 (4つの値)
t2'[i][j] = t1のj列 (4つの値)
```

である、これをelementwiseに掛け算すると、

```python
result[i][j][k] = t1'[i][j][k] * t2'[i][j][k]
```

これは、以下と同じであり、

```python
result[i][j][k] = t1[i][k] * t2[k][j]
```

これのk次元でのsumを取ると、まさに行列積の定義と一致するから、ということらしい。頭いい感。

## メリット

これの何が嬉しいのかというと、DNNコンパイラを作る人の視点からすると嬉しさがあるのである。
DNNコンパイラは高レベルのテンソル操作を低レベルのGPU上の命令 (あるいはGPUのカーネルコード) に変換するソフトウェアだ。

ニューラルネットワークは誤差逆伝播を計算しないといけない関係で、一連のテンソル操作を計算グラフとして定義する必要がある。
define-by-runと呼ばれるモデルでは、ランタイムでテンソル操作から計算グラフを構築するのだが、DNNコンパイラでは、計算グラフをASTとみなし、
ASTからカーネルコードを生成する。
この時、ASTからカーネルコードを生成するにあたり、matmul操作をelementwiseとreduceに置き換えることができれば、サポートすべき命令の数を減らせるので、この生成器の実装を簡単にすることができる。これがメリットである。
