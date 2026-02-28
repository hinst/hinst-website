FROM golang AS backend
ADD backend /app
WORKDIR /app
RUN bash build/go_mod_no_replace.sh
RUN go mod tidy
ENV CGO_ENABLED=1
RUN go mod download
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build

FROM node:22 AS frontend
ADD frontend /app
WORKDIR /app
RUN rm -rf node_modules
RUN rm -rf dist
RUN rm -rf .parcel-cache
RUN npm install
RUN npm run build -- --public-url=/hinst-website

FROM debian:trixie
RUN apt-get update
RUN apt-get -y install git
COPY --from=backend /app/hinst-website /app/hinst-website
COPY --from=backend /app/primeNumbers.json /app/primeNumbers.json
COPY --from=backend /app/pages /app/pages
COPY --from=frontend /app/dist /app/www
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/app/hinst-website"]