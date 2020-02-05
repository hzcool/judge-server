FROM registry.cn-shanghai.aliyuncs.com/aliyun_hzcool/lang_env


WORKDIR /src
COPY src/judgeServer ./judgeServer
COPY src/runner ./runner

CMD ["/src/judgeServer/judger"]

