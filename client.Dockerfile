FROM nginx:alpine

# Remove default nginx config & html
RUN rm -rf /etc/nginx/conf.d/default.conf /usr/share/nginx/html/*

# Copy our custom config
COPY client.nginx.conf /etc/nginx/conf.d/default.conf

# Copy frontend files
COPY client/ /usr/share/nginx/html/

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
