FROM golang:1.19 AS build

COPY . /building
WORKDIR /building/client

RUN go build -o trav 


FROM debian:11

COPY --from=build /building/client/trav /usr/bin/trav
COPY --from=build /building/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN apt update && apt install -y iptables

ENTRYPOINT ["/entrypoint.sh"]
