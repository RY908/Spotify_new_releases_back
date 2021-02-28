#使うDockerイメージ
FROM mysql:8.0
#ポートを開ける
EXPOSE 3306
#MySQL設定ファイルをイメージ内にコピー
ADD ./my.cnf /etc/mysql/conf.d/my.cnf
COPY ./init/* /docker-entrypoint-initdb.d/
RUN chown mysql:mysql /docker-entrypoint-initdb.d/*.sql
#docker runに実行される
CMD ["mysqld"]
