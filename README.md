# Sistema distribuido de predicción de accidentes - TF CC65

Este proyecto implementa nodos distribuidos en Go que leen un dataset de accidentes, donde predicen el nivel de riesgo y almacenan los resultados en Redis. Se comunican entre ellos usando sockets TCP.

## Dataset

Se utilizó el dataset titulado **"Accidentes de Tránsito en Carreteras"**, proporcionado por la Superintendencia de Transporte Terrestre de Personas, Carga y Mercancías (SUTRAN), disponible en el [portal de datos abiertos del gobierno del Perú](https://www.datosabiertos.gob.pe/dataset/accidentes-de-tr%C3%A1nsito-en-carreteras).

Este conjunto de datos incluye variables como:

- Fecha y hora del accidente
- Ubicación geográfica (departamento, kilómetro)
- Tipo de accidente (modalidad)
- Número de fallecidos y heridos

## Cómo ejecutar

1. Asegúrate de tener Docker y Docker Compose instalados.
2. Coloca el archivo `accidentes_original_convertido.csv` en la raíz.
3. Ejecuta:
```bash
docker-compose up --build
docker-compose up --build

