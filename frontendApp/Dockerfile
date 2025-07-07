# Etapa 1: Compilación
FROM node:20.16.0 AS builder

WORKDIR /app

# Copiar archivos de definición de dependencias
COPY package*.json ./

# Instalar dependencias con resolución flexible
RUN npm install --legacy-peer-deps

# Copiar el resto de la app
COPY . .

# Construir la app en modo producción
RUN npm run build -- --configuration=production

# Etapa 2: Servidor Nginx para servir el frontend
FROM nginx:alpine

# Copiar los archivos compilados a nginx
COPY --from=builder /app/dist/frontend-app/browser /usr/share/nginx/html

# (Opcional) Reemplazar la config por defecto de Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

# Comando de inicio del contenedor
CMD ["nginx", "-g", "daemon off;"]
