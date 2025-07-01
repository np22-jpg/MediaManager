ARG VERSION
ARG BASE_URL=""
FROM node:24-alpine AS build
WORKDIR /app
ARG VERSION
ARG BASE_URL

COPY package*.json ./
RUN npm ci && npm cache clean --force

COPY . .
RUN env PUBLIC_VERSION=${VERSION} BASE_URL=${BASE_URL} npm run build

FROM node:24-alpine AS frontend
ARG VERSION

WORKDIR /app

LABEL version=${VERSION}
LABEL description="Docker image for the web frontend of MediaManager"

ENV PUBLIC_VERSION=${VERSION}
ENV PUBLIC_SSR_WEB=false

COPY --from=build /app/build ./build
COPY --from=build /app/package*.json ./

RUN npm ci --omit=dev

EXPOSE 3000
USER node
CMD ["node","build/index.js"]