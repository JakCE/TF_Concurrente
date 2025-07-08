# Sistema distribuido de predicción de accidentes - PC4 CC65

Este proyecto implementa nodos distribuidos en Go que leen un dataset de accidentes, donde predicen el nivel de riesgo y almacenan los resultados en Redis. Se comunican entre ellos usando sockets TCP.

## Cómo ejecutar

1. Asegúrate de tener Docker y Docker Compose instalados.
2. Coloca el archivo `accidentes_original_convertido.csv` en la raíz.
3. Ejecuta:
```bash
docker-compose up --build
docker-compose up --build
