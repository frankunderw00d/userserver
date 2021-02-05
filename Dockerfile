FROM golang:1.14-alpine as builder
WORKDIR /usr/src/userserver
COPY ./userserver ./
RUN apk add --no-cache tzdata upx
RUN upx --best userserver -o _upx_userserver && \
mv -f _upx_userserver userserver

FROM scratch
WORKDIR /opt/userserver
COPY --from=builder /usr/src/userserver/userserver ./
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/
ENV TZ=Asia/Shanghai
CMD ["./userserver"]