FROM node:22-alpine AS frontend
WORKDIR /app
COPY frontend/package.json ./
RUN npm install --ignore-scripts
COPY frontend/ .
RUN npm run build

FROM golang:1.23-alpine AS backend
WORKDIR /build
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o /app/server .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=backend /app/server /server
COPY --from=frontend /app/dist /frontend
EXPOSE 8080
CMD ["/server"]
