package data

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type FilaEntrada struct {
	Departamento string `json:"departamento"`
	Hora         int    `json:"hora"`
	Modalidad    string `json:"modalidad"`
	Fallecidos   int    `json:"fallecidos"`
	Heridos      int    `json:"heridos"`
	Riesgo       string `json:"riesgo"`
}

func CargarCSV(path string) ([]FilaEntrada, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = ','
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var datos []FilaEntrada
	for i, row := range records {
		if i == 0 || len(row) < 5 {
			continue // encabezado o línea inválida
		}
		horaStr := strings.Split(row[1], ":")[0]
		hora, _ := strconv.Atoi(horaStr)
		fallecidos, _ := strconv.Atoi(row[3])
		heridos, _ := strconv.Atoi(row[4])
		riesgo := CalcularRiesgo(fallecidos, heridos)
		datos = append(datos, FilaEntrada{
			Departamento: row[0],
			Hora:         hora,
			Modalidad:    row[2],
			Fallecidos:   fallecidos,
			Heridos:      heridos,
			Riesgo:       riesgo,
		})
	}
	return datos, nil
}

func CalcularRiesgo(f, h int) string {
	if f > 0 {
		return "ALTO"
	} else if h >= 2 {
		return "MEDIO"
	}
	return "BAJO"
}