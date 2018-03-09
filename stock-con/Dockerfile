FROM mhart/alpine-node:8
MAINTAINER Michael Hausenblas <mhausenb@redhat.com>
WORKDIR /app
COPY . .
RUN npm install express && \
    npm install
USER nobody:nobody
CMD ["node", "service.js"]