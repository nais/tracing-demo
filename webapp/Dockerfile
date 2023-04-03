FROM node AS builder
ADD webapp/ /src
WORKDIR /src
RUN npm install
RUN npm run build

FROM nginxinc/nginx-unprivileged
COPY --from=builder /src/dist /usr/share/nginx/html
