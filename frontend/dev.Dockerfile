FROM node:20-alpine3.17

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn --frozen-lockfile

COPY .env.local ./
COPY .eslintrc.json jsconfig.json next.config.js postcss.config.js tailwind.config.js ./

CMD ["yarn", "dev"]