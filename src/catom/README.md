# CAT\_OM 設定ファイル

git で管理されているドキュメントの追加・更新情報から、Atom Feed を生成します。  
設定ファイル([env.toml](./etc/env.toml.tmpl))で指定された、`root`のディレクトリ配下に配置された Markdown ファイルのみを対象とする。  
Markdown ファイルの拡張子は、\*.md, \*.markdown とする。  

## env.toml
- **dotgit**: Feed を生成する対象となる`.git`への絶対パス
- **diff**: Feed を生成する対象とする commit の時間指定
  - 実行時から指定した時間[h]が対象となる
- **root**: Feed を生成する対象のディレクトリを指定
  - **root** で指定されたディレクトリ群は **dotgit** で指定されたリポジトリで管理されていなければならない

```
root=[
"www_md",
"www_tmpl"
]
```

- **proto**: Feed の URI で使用するスキーム
  - http or https
- **host**: Feed の URI で使用する host 名
- **urlroot**: Feed の Entry で記述する、link タグの URL ルートのパス
  - `<link href>`のリンクは、以下のように構成される
    - **proto**://**host**/**urlroot**/**root** 配下のディレクトリとファイル`
- **feedurl**: Feed のリンクとなる URL
- **outpath**: Feed を出力する絶対パス
- **outfile**: Feed を出力するファイル名

### [feed]
- **feedid**: `<feed>`で指定する`<id>`の値
- **title**: `<feed>`で指定する`<title>`の値
- **subtitle**: `<feed>`で指定する`<subtitle>`の値

### [author]
- **name**: `<feed>`で指定する`<name>`の値
- **email**: `<feed>`で指定する`<email>`の値
- **url**: `<feed>`で指定する`<uri>`の値
