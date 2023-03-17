FROM node
ADD . /src
#RUN rm -rf /src/node_modules
WORKDIR /src
RUN npm install
CMD ["npm", "run", "dev"]
