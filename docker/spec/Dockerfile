FROM node:22-alpine
WORKDIR /spec

COPY spec/package.json ./
RUN npm install --ignore-scripts

COPY spec/ .
CMD ["npm", "run", "compile"]
