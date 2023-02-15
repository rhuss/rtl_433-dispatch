FROM alpine AS build
RUN apk add --no-cache alpine-sdk gcc linux-headers ncurses-dev librtlsdr-dev cmake libusb-dev bash
RUN git clone -b add-hash-0x7c0e ttps://github.com/rhuss/wmbusmeters.git\
    git clone https://github.com/weetmuts/wmbusmeters.git && \
    git clone https://github.com/weetmuts/rtl-wmbus.git && \
    git clone https://github.com/merbanan/rtl_433.git && \
    git clone
WORKDIR /wmbusmeters
RUN make
WORKDIR /rtl-wmbus
RUN make release && chmod 755 build/rtl_wmbus
WORKDIR /rtl_433
RUN mkdir build && cd build && cmake ../ && make

FROM alpine as scratch
RUN apk add --no-cache mosquitto-clients libstdc++ curl libusb ncurses rtl-sdr netcat-openbsd
WORKDIR /wmbusmeters
COPY --from=build /wmbusmeters/build/wmbusmeters /wmbusmeters/wmbusmeters
COPY --from=build /rtl-wmbus/build/rtl_wmbus /usr/bin/rtl_wmbus
COPY --from=build /rtl_433/build/src/rtl_433 /usr/bin/rtl_433
COPY --from=build /wmbusmeters/docker/docker-entrypoint.sh /wmbusmeters/docker-entrypoint.sh
CMD ["sh", "/wmbusmeters/docker-entrypoint.sh"]
