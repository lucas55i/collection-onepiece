FROM nginx:1.25-alpine

# Remove config padrão do nginx
RUN rm /etc/nginx/conf.d/default.conf

# Config customizada
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copia arquivos estáticos
COPY index.html style.css app.js utils.js /usr/share/nginx/html/

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
