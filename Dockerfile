FROM golang:1.17
WORKDIR /go/src
RUN mkdir -p github.com/pandotoken && chmod -R 777 /go 
WORKDIR /go/src/github.com/pandotoken/pando
COPY . /go/src/github.com/pandotoken/pando/
WORKDIR /go/src/github.com/pandotoken/pando
RUN mkdir ../pandonet
WORKDIR /go/src/github.com/pandotoken/pando/integration/pandonet/
RUN cp -ivr * /go/src/github.com/pandotoken/pandonet/
WORKDIR /go/src/github.com/pandotoken/pando
RUN make install
EXPOSE 16888
EXPOSE 16889
EXPOSE 12000
CMD ["ls"]
