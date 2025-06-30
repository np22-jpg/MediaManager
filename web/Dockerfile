ARG VERSION
ARG BASE_URL=""
FROM node:24-alpine AS deps
WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

FROM node:24-alpine AS build
WORKDIR /app
ARG VERSION
ARG BASE_URL

COPY --from=deps /app/node_modules ./node_modules
COPY package*.json ./
RUN npm ci

COPY . .
RUN env PUBLIC_VERSION=${VERSION} BASE_URL=${BASE_URL} npm run build

FROM node:24-alpine AS frontend
ARG VERSION
WORKDIR /app
LABEL version=${VERSION}
LABEL description="Docker image for the web frontend of MediaManager"

ENV PUBLIC_VERSION=${VERSION}
ENV PUBLIC_SSR_WEB=false

COPY --from=build /app/build/ ./build/
COPY package*.json ./
RUN npm ci --only=production

EXPOSE 3000
USER node
CMD ["node","build/index.js"]