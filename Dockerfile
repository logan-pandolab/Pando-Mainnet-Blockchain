FROM golang:latest
LABEL maintainer="Bhupesh-Yadav"
LABEL build_date="2022-04-01"
RUN mkdir -p /usr/local/go/src/github.com/pandotoken/pando
RUN chmod -R 777 /usr/local/go
WORKDIR /usr/local/go/src/github.com/pandotoken/pando
COPY . /usr/local/go/src/github.com/pandotoken/pando/
WORKDIR /usr/local/go/src/github.com/pandotoken/pando/
WORKDIR /usr/local/go/src/github.com/pandotoken/pando/integration
RUN cp -r /usr/local/go/src/github.com/pandotoken/pando/integration/pandonet /usr/local/go/src/github.com/pandotoken/
WORKDIR /usr/local/go/src/github.com/pandotoken/pando
RUN chmod -R 777 /usr/local/go
RUN make install
EXPOSE 16888
EXPOSE 12000
EXPOSE 16889
CMD ["/usr/local/go/bin/pando start"]
