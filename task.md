# Fixing Docker build

## Log below:

C:\Dev\hinst-website>build && push

C:\Dev\hinst-website>docker build -t 192.168.0.23:5000/hinst/hinst-website --progress=plain --platform=linux/arm64 .
#0 building with "desktop-linux" instance using docker driver

#1 [internal] load build definition from Dockerfile
#1 transferring dockerfile: 888B done
#1 DONE 0.0s

#2 [internal] load metadata for docker.io/library/debian:trixie
#2 DONE 0.5s

#3 [internal] load metadata for docker.io/library/node:26
#3 DONE 0.5s

#4 [internal] load metadata for docker.io/library/golang:latest
#4 DONE 0.5s

#5 [internal] load .dockerignore
#5 transferring context: 2B done
#5 DONE 0.0s

#6 [frontend  1/10] FROM docker.io/library/node:26@sha256:e961046fec20896e8904f2b4a8b4c7e5ca91826d84d8d33d83dbaa61f942069e
#6 resolve docker.io/library/node:26@sha256:e961046fec20896e8904f2b4a8b4c7e5ca91826d84d8d33d83dbaa61f942069e 0.0s done
#6 DONE 0.0s

#7 [backend 1/6] FROM docker.io/library/golang:latest@sha256:f44f6e88636cfb311f9ebace870ded69d943f227bb3cb27d32ffd84ea18c43ea
#7 resolve docker.io/library/golang:latest@sha256:f44f6e88636cfb311f9ebace870ded69d943f227bb3cb27d32ffd84ea18c43ea 0.0s done
#7 DONE 0.0s

#8 [stage-2 1/8] FROM docker.io/library/debian:trixie@sha256:f324c7ff54321e8d9c588493a20244965938ce0aa50bbd1022d38010e9ffc4b1
#8 resolve docker.io/library/debian:trixie@sha256:f324c7ff54321e8d9c588493a20244965938ce0aa50bbd1022d38010e9ffc4b1 0.0s done
#8 DONE 0.0s

#9 [internal] load build context
#9 transferring context: 97.83MB 5.0s
#9 transferring context: 216.59MB 10.1s
#9 transferring context: 338.34MB 15.2s
#9 transferring context: 461.54MB 20.3s
#9 transferring context: 576.77MB 25.3s
#9 transferring context: 663.50MB 28.9s done
#9 DONE 29.0s

#10 [backend 2/6] ADD backend /app
#10 CACHED

#11 [backend 3/6] WORKDIR /app
#11 CACHED

#12 [frontend  2/10] ADD frontend /app
#12 CACHED

#13 [frontend  3/10] WORKDIR /app
#13 DONE 0.8s

#14 [backend 4/6] RUN go mod download -x
#14 0.247 exec /bin/sh: exec format error
#14 ERROR: process "/bin/sh -c go mod download -x" did not complete successfully: exit code: 255
------
 > [backend 4/6] RUN go mod download -x:
0.247 exec /bin/sh: exec format error
------
Dockerfile:5
--------------------
   3 |     WORKDIR /app
   4 |     ENV CGO_ENABLED=1
   5 | >>> RUN go mod download -x
   6 |     ENV GOCACHE=/root/.cache/go-build
   7 |     RUN --mount=type=cache,target="/root/.cache/go-build" go build
--------------------
ERROR: failed to build: failed to solve: process "/bin/sh -c go mod download -x" did not complete successfully: exit code: 255

What's next:
    Debug this build failure with Gordon → docker ai "help me fix this build failure"
