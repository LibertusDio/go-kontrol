# pull official base image
FROM node:13.12.0-alpine

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH
ARG DOMAIN_API_GATEWAY
ENV DOMAIN_API_GATEWAY = $DOMAIN_API_GATEWAY
# install app dependencies
COPY fe/package.json ./
COPY fe/package-lock.json ./
RUN npm install --silent
#RUN npm install react-scripts@3.4.1 -g --silent

# add app
COPY fe ./

# start app
CMD ["npm", "start"]