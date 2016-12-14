FROM haproxy

ADD stalls .

ENTRYPOINT ["./stalls"]

