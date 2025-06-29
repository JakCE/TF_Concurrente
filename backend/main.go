package main

import (
	"context"
	"fmt"
	"os"
	"net"
	"backend/data"
	"backend/forest"
	"backend/redis"
)

func main() {
	fmt.Println("Nodo iniciado - Entrenamiento Random Forest")

	dataset, err := data.CargarCSV("accidentes_procesado.csv")
	if err != nil {
		fmt.Println("Error cargando dataset:", err)
		os.Exit(1)
	}

	rf := forest.EntrenamientoDistribuido(dataset, 5) // 5 árboles concurrentes
	fmt.Println("Entrenamiento completado.")

	// Guardar modelo en Redis
	rdb := redis.NewClient()
	err = redis.GuardarModelo(context.Background(), rdb, "modelo-nodo", rf)
	if err != nil {
		fmt.Println("Error guardando modelo en Redis:", err)
	} else {
		fmt.Println("Modelo guardado en Redis correctamente")
	}

	// Leer modelo de Redis
	rf2, err := redis.LeerModelo(context.Background(), rdb, "modelo-nodo")
	if err != nil {
		fmt.Println("Error leyendo modelo desde Redis:", err)
		return
	}

	fmt.Println("Ejemplo de predicción desde modelo almacenado:")
	entrada := data.FilaEntrada{
		Departamento: "LIMA",
		Hora:         13,
		Modalidad:    "CHOQUE",
		Fallecidos:   0,
		Heridos:      3,
	}
	riesgo := rf2.Predict(entrada)
	fmt.Println("Predicción:", riesgo)

	// Iniciar servidor TCP por nodo
	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = "6000"
	}
	fmt.Println("Servidor P2P escuchando en :" + puerto)

	ln, err := net.Listen("tcp", ":"+puerto)
	if err != nil {
		fmt.Println("Error iniciando servidor TCP:", err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error de conexión:", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			fmt.Fprintf(c, "Nodo activo - listo para predicción en puerto %s\n", puerto)
		}(conn)
	}
}