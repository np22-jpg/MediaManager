ARG VERSION
ARG BASE_URL=""
FROM node:24-alpine AS build
WORKDIR /app
ARG VERSION
ARG BASE_URL

# Copy package files first for better layer caching
COPY package*.json ./
RUN npm ci && npm cache clean --force

# Copy source code after dependencies are installed
COPY . .
RUN env PUBLIC_VERSION=${VERSION} BASE_URL=${BASE_URL} npm run build

FROM node:24-alpine AS frontend
ARG VERSION

USER node
EXPOSE 3000
WORKDIR /app

LABEL version=${VERSION}
LABEL description="Docker image for the web frontend of MediaManager"


ENV PUBLIC_VERSION=${VERSION}
ENV PUBLIC_SSR_WEB=false

# Copy built application and package files
COPY --from=build /app/build ./build
COPY --from=build /app/package*.json ./
COPY --chmod=755 entrypoint.sh .

# Install only production dependencies needed for the Node adapter
RUN npm ci --only=production && npm cache clean --force

CMD ["/app/entrypoint.sh"]

