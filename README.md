# Sistema distribuido de predicción de accidentes - TF CC65

Este proyecto implementa nodos distribuidos en Go que leen un dataset de accidentes, donde predicen el nivel de riesgo y almacenan los resultados en Redis. Se comunican entre ellos usando sockets TCP.

## Cómo ejecutar

1. Asegúrate de tener Docker y Docker Compose instalados.
2. Coloca el archivo `accidentes_original_convertido.csv` en la raíz.
3. Ejecuta:
```bash
docker-compose up --build
docker-compose up --build

## descripcion data set
Se utilizó el dataset titulado "Accidentes de Tránsito en Carreteras", proporcionado por la Superintendencia de Transporte Terrestre de Personas, Carga y Mercancías (SUTRAN) y disponible en el portal de datos abiertos del gobierno del Perú. Este conjunto de datos ofrece información detallada sobre incidentes vehiculares ocurridos en las carreteras nacionales, incluyendo variables como la fecha, hora, ubicación, tipo de accidente, así como el número de fallecidos y heridos.

