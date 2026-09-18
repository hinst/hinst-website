docker run --rm --privileged multiarch/qemu-user-static --reset -p yes &&^
docker build -t 192.168.178.29:5000/hinst/hinst-website --progress=plain --platform=linux/arm64 .
