FROM node:24-alpine AS build
USER node:node
WORKDIR /app

COPY --chown=node:node . .

RUN npm ci
RUN npm run build

FROM node:24-alpine AS frontend
ARG VERSION
USER node:node
LABEL version=${VERSION}
LABEL description="Docker image for the web frontend of MediaManager"

ENV PUBLIC_VERSION=${VERSION}
ENV PUBLIC_SSR_WEB=false
WORKDIR /app

COPY --chown=node:node package*.json ./
COPY --chown=node:node --from=build /app/build/ ./build/

RUN npm install

EXPOSE 3000

CMD ["node","build/index.js"]
