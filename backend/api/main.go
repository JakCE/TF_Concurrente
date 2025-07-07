package main

import (
	"app/common"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

type Entrada struct {
	Hora         int     `json:"hora"`
	Departamento int     `json:"departamento"`
	CodigoVia    int     `json:"codigo_via"`
	Kilometro    float64 `json:"kilometro"`
	Modalidad    int     `json:"modalidad"`
}

type RespuestaNodo struct {
	Riesgo string `json:"riesgo"`
}

type ResultadoFinal struct {
	RiesgoFinal string            `json:"riesgo_final"`
	Votos       map[string]int    `json:"votos"`
	Respuestas  map[string]string `json:"respuestas"`
}

var nodos = []string{
	"http://nodo1:8001/predecir",
	"http://nodo2:8002/predecir",
	"http://nodo3:8003/predecir",
}

func handlerPrediccion(w http.ResponseWriter, r *http.Request) {
	var entrada Entrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	votos := make(map[string]int)
	respuestas := make(map[string]string)

	for _, url := range nodos {
		wg.Add(1)
		go func(nodo string) {
			defer wg.Done()
			cuerpo, _ := json.Marshal(entrada)
			resp, err := http.Post(nodo, "application/json", bytes.NewBuffer(cuerpo))
			if err != nil {
				log.Printf("Error con %s: %v\n", nodo, err)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			var r RespuestaNodo
			if err := json.Unmarshal(body, &r); err != nil {
				log.Printf("Error al parsear respuesta de %s\n", nodo)
				return
			}

			mu.Lock()
			votos[r.Riesgo]++
			respuestas[nodo] = r.Riesgo
			mu.Unlock()
		}(url)
	}

	wg.Wait()

	// Votación mayoritaria
	ganador := ""
	maxVotos := 0
	for clase, count := range votos {
		if count > maxVotos {
			ganador = clase
			maxVotos = count
		}
	}

	resultado := ResultadoFinal{
		RiesgoFinal: ganador,
		Votos:       votos,
		Respuestas:  respuestas,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}

func handlerDiccionarios(w http.ResponseWriter, r *http.Request) {
	diccionario := map[string]map[int]string{
		"departamentos": common.DeptoDecod,
		"vias":          common.ViaDecod,
		"modalidades":   common.ModoDecod,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diccionario)
}

func handlerEstadisticas(w http.ResponseWriter, r *http.Request) {
	registros, err := common.CargarDatosCSV(common.DatasetPath)
	if err != nil {
		http.Error(w, "Error al cargar registros", http.StatusInternalServerError)
		return
	}

	deptoCount := make(map[string]int)
	horaCount := make(map[int]int)
	modalidadCount := make(map[string]int)

	for _, r := range registros {
		depto := common.DeptoDecod[r.Departamento]
		deptoCount[depto]++

		horaCount[r.Hora]++

		modalidad := common.ModoDecod[r.Modalidad]
		modalidadCount[modalidad]++
	}

	respuesta := map[string]interface{}{
		"por_departamento": deptoCount,
		"por_hora":         horaCount,
		"por_modalidad":    modalidadCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}

// Middleware CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	_, err := common.CargarDatosCSV(common.DatasetPath)
	if err != nil {
		log.Fatalf("Error cargando dataset: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/predecir", handlerPrediccion)
	mux.HandleFunc("/diccionarios", handlerDiccionarios)
    mux.HandleFunc("/estadisticas", handlerEstadisticas)

	log.Println("API escuchando en puerto 8080...")
	err = http.ListenAndServe(":8080", corsMiddleware(mux))
	if err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}