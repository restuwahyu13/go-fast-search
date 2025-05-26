FROM node:22-alpine

LABEL name="Restu Wahyu Saputra"
LABEL email="restuwahyu13@gmail.com"
LABEL linkedin="https://www.linkedin.com/in/restuwahyu13"
LABEL gitHub="https://github.com/restuwahyu13"

WORKDIR /app
COPY ./package*.json .
COPY . .

RUN rm -rf node_modules .next .cache .turbo .~/.npm; \
    npm cache clean -f; \
    npm i -g npm@latest pm2 --loglevel verbose; \
    npm i --legacy-peer-deps --loglevel verbose; \
    npm run build

RUN apk update; \
    apk -u list; \
    apk upgrade; \
    apk add --no-cache tzdata ca-certificates

EXPOSE 3000
CMD npm run start:pm2