FROM node:alpine

# Create the working directory
WORKDIR /usr/src/app

# Copy package json only so that we can run npm install to setup the dependencies.
COPY package*.json ./

# Install all dependencies
RUN apk add --no-cache --virtual .gyp \
        python \
        make \
        g++ \
    && npm install \
    && apk del .gyp

# Copy all js files, avoid copying node_modules
COPY ./*.js ./

# Run the server.js file
CMD ["node", "server.js"]

