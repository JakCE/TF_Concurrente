package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"app/common"
	"github.com/go-redis/redis/v8"
)

type Entrada struct {
	Hora         int     `json:"hora"`
	Departamento int     `json:"departamento"`
	CodigoVia    int     `json:"codigo_via"`
	Kilometro    float64 `json:"kilometro"`
	Modalidad    int     `json:"modalidad"`
}

type Salida struct {
	Riesgo string `json:"riesgo"`
}

type RandomForest struct {
	Arboles []*common.Nodo
}

func (rf *RandomForest) Entrenar(datos []common.Ejemplo, numArboles, profundidad int) {
	for i := 0; i < numArboles; i++ {
		if i%10 == 0 {
			fmt.Printf("%s - Entrenando árbol %d/%d...\n", time.Now().Format(time.RFC3339), i+1, numArboles)
		}
		muestra := bootstrapping(datos)
		arbol := common.EntrenarArbol(muestra, profundidad)
		rf.Arboles = append(rf.Arboles, arbol)
		fmt.Println("✅ Entrenamiento completado")
	}
}

func (rf *RandomForest) Predecir(entrada []float64) string {
	votos := map[string]int{}
	for _, arbol := range rf.Arboles {
		c := common.Predecir(arbol, entrada)
		votos[c]++
	}
	var ganador string
	max := 0
	for clase, count := range votos {
		if count > max {
			max = count
			ganador = clase
		}
	}
	return ganador
}

func bootstrapping(datos []common.Ejemplo) []common.Ejemplo {
	var out []common.Ejemplo
	for i := 0; i < len(datos); i++ {
		out = append(out, datos[rand.Intn(len(datos))])
	}
	return out
}

var (
	rf         RandomForest
	bitacora   *os.File
	ctx        = context.Background()
	redisCli   *redis.Client
	nodoNombre = "nodo2"
)

func registrarBitacora(msg string) {
	logLine := fmt.Sprintf("%s - %s\n", time.Now().Format(time.RFC3339), msg)
	fmt.Print(logLine)
	if bitacora != nil {
		bitacora.WriteString(logLine)
	}
}

func handlerPrediccion(w http.ResponseWriter, r *http.Request) {
	var entrada Entrada
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	atributos := []float64{
		float64(entrada.Hora),
		float64(entrada.Departamento),
		float64(entrada.CodigoVia),
		entrada.Kilometro,
		float64(entrada.Modalidad),
	}

	clase := rf.Predecir(atributos)
	registrarBitacora(fmt.Sprintf("Entrada recibida: %+v → Predicción: %s", entrada, clase))

	inputData, _ := json.Marshal(entrada)
	predData := map[string]interface{}{
		"entrada":   string(inputData),
		"resultado": clase,
		"hora":      time.Now().Format(time.RFC3339),
	}
	predJson, _ := json.Marshal(predData)
	clave := fmt.Sprintf("predicciones:%s:%d", nodoNombre, time.Now().UnixNano())

	if err := redisCli.Set(ctx, clave, predJson, 0).Err(); err != nil {
		registrarBitacora(fmt.Sprintf("❌ Error guardando en Redis: %v", err))
	}

	json.NewEncoder(w).Encode(Salida{Riesgo: clase})
}

func main() {
	var err error

	registrarBitacora("Iniciando nodo...")

	// Abrir bitácora
	bitacora, err = os.OpenFile("bitacora.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("❌ Error abriendo bitácora: %v", err)
	}
	defer bitacora.Close()

	// Conectar a Redis
	registrarBitacora("Conectando a Redis...")
	redisCli = redis.NewClient(&redis.Options{
		Addr: common.RedisAddr,
	})
	if pong, err := redisCli.Ping(ctx).Result(); err != nil {
		log.Fatalf("❌ Error conectando a Redis: %v", err)
	} else {
		registrarBitacora("✅ Conexión a Redis exitosa: " + pong)
	}

	// Cargar dataset
	registrarBitacora("Cargando dataset desde CSV...")
	registros, err := common.CargarDatosCSV(common.DatasetPath)
	if err != nil {
		log.Fatalf("❌ Error cargando dataset: %v", err)
	}
	rand.Shuffle(len(registros), func(i, j int) {
		registros[i], registros[j] = registros[j], registros[i]
	})
	registros = registros[:5000]

	registrarBitacora(fmt.Sprintf("✅ Registros cargados: %d", len(registros)))

	ejemplos := []common.Ejemplo{}
	for _, r := range registros {
		clase := "Leve"
		gravedad := r.Fallecidos + r.Heridos
		if gravedad >= 1 {
			clase = "Moderado"
		}
		if gravedad >= 3 {
			clase = "Grave"
		}

		atributos := []float64{
			float64(r.Hora),
			float64(r.Departamento),
			float64(r.CodigoVia),
			r.Kilometro,
			float64(r.Modalidad),
		}

		ejemplos = append(ejemplos, common.Ejemplo{
			Atributos: atributos,
			Clase:     clase,
		})
	}

	registrarBitacora("Entrenando Random Forest...")
	rf.Entrenar(ejemplos, common.ArbolesPorNodo, 5)
	registrarBitacora("✅ Modelo entrenado correctamente")

	registrarBitacora("Servidor iniciado en :8002")
	http.HandleFunc("/predecir", handlerPrediccion)
	log.Fatal(http.ListenAndServe(":8002", nil))
}
