# ベースとなるDockerイメージ指定
FROM golang:1.14

ENV GO111MODULE=on

# アップデートとgitのインストール！！
RUN apt-get update && apt-get install --yes --no-install-recommends git 

# コンテナ内に作業ディレクトリを作成
RUN mkdir /go/src/Spotify_new_releases

# コンテナログイン時のディレクトリ指定
WORKDIR /go/src/Spotify_new_releases

# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . /go/src/Spotify_new_releases

RUN go mod download

ENV GO111MODULE=off
# RUN go get gopkg.in/urfave/cli.v2@master
RUN go get github.com/oxequa/realize 

ENV GO111MODULE=on

# CMD ["realize", "start"]