FROM 976d9160e30942228bbcb9a5fb4b2ac82d7aa63c034e073205a776692167937c


WORKDIR /src
COPY src/judgeServer ./judgeServer
COPY src/runner ./runner

CMD ["/src/judgeServer/judger"]

