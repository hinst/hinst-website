FROM golang AS backend
ADD backend /app
WORKDIR /app
ENV GOARCH=arm64
RUN go mod download
RUN go build

FROM node:22 AS frontend
ADD frontend /app
WORKDIR /app
RUN rm -rf node_modules
RUN rm -rf dist
RUN rm -rf .parcel-cache
RUN npm install
RUN npm run build

FROM debian:bookworm
COPY --from=backend /app/hinst-website /app/hinst-website
COPY --from=frontend /app/dist /app/www
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/app/hinst-website", "--allowOrigin=https://hinst.github.io"]