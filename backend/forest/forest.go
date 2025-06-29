package forest

import (
	"backend/data"
	"encoding/json"
	"sync"
)

type Arbol struct {
	MayorHeridos int    `json:"mayor_heridos"`
	MayorFalle   bool   `json:"mayor_falle"`
	LabelIfHigh  string `json:"label_high"`
	LabelIfMed   string `json:"label_med"`
	LabelIfLow   string `json:"label_low"`
}

type Bosque struct {
	Arboles []Arbol `json:"arboles"`
}

// Esta función procesa un subset y lo envía al canal
func procesarSubset(subset []data.FilaEntrada, canal chan<- Arbol, wg *sync.WaitGroup) {
	defer wg.Done()
	canal <- entrenar(subset)
}

func EntrenamientoDistribuido(data []data.FilaEntrada, n int) Bosque {
	var wg sync.WaitGroup
	canal := make(chan Arbol, n)
	dist := dividir(data, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go procesarSubset(dist[i], canal, &wg)
	}

	wg.Wait()
	close(canal)

	var bosque Bosque
	for a := range canal {
		bosque.Arboles = append(bosque.Arboles, a)
	}
	return bosque
}

func dividir(datos []data.FilaEntrada, partes int) [][]data.FilaEntrada {
	size := len(datos) / partes
	var chunks [][]data.FilaEntrada
	for i := 0; i < partes; i++ {
		start := i * size
		end := start + size
		if i == partes-1 {
			end = len(datos)
		}
		chunks = append(chunks, datos[start:end])
	}
	return chunks
}

func entrenar(data []data.FilaEntrada) Arbol {
	return Arbol{
		MayorHeridos: 2,
		MayorFalle:   true,
		LabelIfHigh:  "ALTO",
		LabelIfMed:   "MEDIO",
		LabelIfLow:   "BAJO",
	}
}

func (b Bosque) Predict(f data.FilaEntrada) string {
	cuenta := map[string]int{}
	for _, a := range b.Arboles {
		if f.Fallecidos > 0 && a.MayorFalle {
			cuenta[a.LabelIfHigh]++
		} else if f.Heridos >= a.MayorHeridos {
			cuenta[a.LabelIfMed]++
		} else {
			cuenta[a.LabelIfLow]++
		}
	}
	max := "BAJO"
	for k, v := range cuenta {
		if v > cuenta[max] {
			max = k
		}
	}
	return max
}

func (b Bosque) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

func FromJSON(data []byte) (Bosque, error) {
	var b Bosque
	err := json.Unmarshal(data, &b)
	return b, err
}