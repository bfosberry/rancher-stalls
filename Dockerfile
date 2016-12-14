FROM haproxy
WORKDIR /
ENTRYPOINT ["./stalls"]
ADD stalls .


