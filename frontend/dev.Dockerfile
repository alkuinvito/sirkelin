FROM node:19-alpine

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY next.config.js postcss.config.js jsconfig.json tailwind.config.js ./

CMD ["npm", "run", "dev"]