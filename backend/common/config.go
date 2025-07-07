package common

var (
	Nodos = map[string]string{
		"nodo1": "localhost:8001",
		"nodo2": "localhost:8002",
		"nodo3": "localhost:8003",
	}

	RedisAddr = "redis:6379"

	// Ruta del dataset (ya cargado y preprocesado)
	DatasetPath = "./common/accidentes_completo.csv"

	// Número de árboles por nodo (para Random Forest distribuido)
	ArbolesPorNodo = 5
)