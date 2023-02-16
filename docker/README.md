# How to build for ARM and AMD with Docker Desktop


```
# Prepare
docker buildx create --name xbuilder
docker buildx use xbuilder
docker buildx inspect --bootstrap

# Build
docker buildx build --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --progress plain --push -t rhuss/rtl_433-dispatch .
```
