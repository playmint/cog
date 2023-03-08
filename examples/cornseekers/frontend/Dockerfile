FROM node:17.1-alpine as builder
WORKDIR /tmp

# dependencies
RUN mkdir frontend
COPY ./package*.json ./frontend/
RUN	(cd frontend && npm install)

# build
COPY . ./frontend
RUN (cd frontend && npm run build)

# final image
FROM nginx:1.21.5
RUN rm -rf /usr/share/nginx/html
COPY --from=builder /tmp/frontend/dist /usr/share/nginx/html

# nginx config for nextjs apps
RUN echo '                                               \
	server_tokens off;                                   \
	server {                                             \
		listen       80;                                 \
		root   /usr/share/nginx/html;                    \
		location /assets {                               \
			add_header Cache-Control "max-age=31536000"; \
		}                                                \
		location / {                                     \
			try_files $uri.html  $uri $uri/ /index.html; \
			add_header Cache-Control "no-cache";         \
		}                                                \
	}                                                    \
' > /etc/nginx/conf.d/default.conf

