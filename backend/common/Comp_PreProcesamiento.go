package common

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Registro struct {
	Hora         int
	Departamento int
	CodigoVia    int
	Kilometro    float64
	Modalidad    int
	Fallecidos   int
	Heridos      int
}

// Codificadores y decodificadores
var (
	DeptoCod    = map[string]int{}
	ViaCod      = map[string]int{}
	ModoCod     = map[string]int{}
	DeptoDecod  = map[int]string{}
	ViaDecod    = map[int]string{}
	ModoDecod   = map[int]string{}
	nextCod     = 1
)

func codificar(valor string, cod map[string]int, decod map[int]string) int {
	if v, ok := cod[valor]; ok {
		return v
	}
	cod[valor] = nextCod
	decod[nextCod] = valor
	nextCod++
	return cod[valor]
}

// Cargar y convertir registros
func CargarDatosCSV(path string) ([]Registro, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // saltar cabecera
	registros := []Registro{}

	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		hora, _ := strconv.Atoi(row[1])
		kilometro, _ := strconv.ParseFloat(row[4], 64)
		fallecidos, _ := strconv.Atoi(row[6])
		heridos, _ := strconv.Atoi(row[7])

		departamento := codificar(row[2], DeptoCod, DeptoDecod)
		codigoVia := codificar(row[3], ViaCod, ViaDecod)
		modalidad := codificar(row[5], ModoCod, ModoDecod)

		reg := Registro{
			Hora:         hora,
			Departamento: departamento,
			CodigoVia:    codigoVia,
			Kilometro:    kilometro,
			Modalidad:    modalidad,
			Fallecidos:   fallecidos,
			Heridos:      heridos,
		}
		registros = append(registros, reg)
	}
	fmt.Printf("✅ Codificadores generados. Departamentos: %d, Vías: %d, Modalidades: %d\n",
		len(DeptoCod), len(ViaCod), len(ModoCod))
	return registros, nil
}
