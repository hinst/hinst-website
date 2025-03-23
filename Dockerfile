FROM --platform=linux/arm64 golang AS backend
ADD backend /app
WORKDIR /app
ENV CGO_ENABLED=1
RUN go mod download
RUN go build

FROM node:22 AS frontend
ADD frontend /app
WORKDIR /app
RUN rm -rf node_modules
RUN rm -rf dist
RUN rm -rf .parcel-cache
RUN npm install
RUN npm run build -- --public-url=/hinst-website

FROM --platform=linux/arm64 debian:bookworm
COPY --from=backend /app/hinst-website /app/hinst-website
COPY --from=frontend /app/dist /app/www
WORKDIR /app
EXPOSE 8080
ENV GOGC=50
ENTRYPOINT ["/app/hinst-website"]