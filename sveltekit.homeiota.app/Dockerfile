# Stage 1: Build the app
FROM node:20-alpine AS builder
WORKDIR /app
# Install build dependencies for canvas and other native modules
RUN apk add --no-cache python3 make g++ cairo-dev pango-dev jpeg-dev giflib-dev
COPY package.json package-lock.json* pnpm-lock.yaml* yarn.lock* ./
RUN npm install
COPY . .
# Copy prisma schema and migrations
COPY prisma ./prisma
# Generate Prisma client
RUN npx prisma generate
RUN npm run build

# Stage 2: Run the app
FROM node:20-alpine
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/package.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/build ./build
COPY --from=builder /app/.svelte-kit ./.svelte-kit
COPY --from=builder /app/src ./src
COPY --from=builder /app/static ./static
COPY --from=builder /app/svelte.config.js ./svelte.config.js
COPY --from=builder /app/vite.config.ts ./vite.config.ts
COPY --from=builder /app/prisma ./prisma
# Optional: run migrations in production
# RUN npx prisma migrate deploy
EXPOSE 3000
CMD ["node", "build"]
