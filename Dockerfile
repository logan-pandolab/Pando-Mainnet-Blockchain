FROM golang:latest
LABEL maintainer="Bhupesh-Yadav"
LABEL build_date="2022-04-01"
WORKDIR /usr/local/go/
RUN mkdir -p src/github.com/pandotoken/pando && chmod -R 777 /usr/local/go
WORKDIR /usr/local/go/src/github.com/pandotoken/pando
COPY . /usr/local/go/src/github.com/pandotoken/pando/
WORKDIR /usr/local/go/src/github.com/pandotoken/pando/
RUN mkdir ../pandonet
WORKDIR /usr/local/go/src/github.com/pandotoken/pando/integration/pandonet/
RUN cp -ivr * /usr/local/go/src/github.com/pandotoken/pandonet/
WORKDIR /usr/local/go/src/github.com/pandotoken/pando
RUN make install
EXPOSE 16888
EXPOSE 12000
EXPOSE 16889
CMD ["/usr/local/go/bin/pando start"]