FROM alpine AS build
RUN apk add --no-cache alpine-sdk gcc linux-headers ncurses-dev librtlsdr-dev cmake libusb-dev bash
RUN git clone -b add-hash-0x7c0e https://github.com/rhuss/wmbusmeters.git \
 && git clone https://github.com/merbanan/rtl_433.git \
 && git clone https://github.com/rhuss/rtl_433-dispatch.git
WORKDIR /wmbusmeters
RUN make
WORKDIR /rtl_433
RUN mkdir build && cd build && cmake ../ && make
WORKDIR /rtl_433-dispatch
RUN go build .

FROM alpine as scratch
RUN apk add --no-cache mosquitto-clients libstdc++ curl libusb ncurses rtl-sdr netcat-openbsd
WORKDIR /rtl_433-dispatch
COPY --from=build /wmbusmeters/build/wmbusmeters wmbusmeters
COPY --from=build /rtl_433/build/src/rtl_433 rtl_433
COPY --from=build /rtl_433-dispatch/rtl_433-dispatch rtl_433-dispatch

