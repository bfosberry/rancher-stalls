FROM haproxy
WORKDIR /
ENTRYPOINT ["./stalls"]
ADD dist/stalls .


