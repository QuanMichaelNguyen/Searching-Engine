# Stage 1: Build the application assets
FROM node:lts-alpine AS build

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install

# Copy source code and build
COPY . .
# Replace 'build' with the correct command for your framework (e.g., 'npm run build')
RUN npm run build

# Stage 2: Serve the production assets using a lightweight web server (like Nginx)
FROM nginx:alpine
# Remove default nginx config
RUN rm -rf /etc/nginx/conf.d

# Copy the build output from the previous stage into the Nginx static folder
COPY --from=build /app/dist /usr/share/nginx/html

# Copy a custom Nginx configuration file (if needed for routing)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]