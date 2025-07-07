package common

import (
	"math"
)

type Nodo struct {
	Campo     int
	Valor     float64
	Clase     string
	Izquierda *Nodo
	Derecha   *Nodo
}

type Ejemplo struct {
	Atributos []float64
	Clase     string
}

// Entrenamiento del árbol
func EntrenarArbol(datos []Ejemplo, profundidad int) *Nodo {
	if profundidad == 0 || len(datos) <= 1 {
		return &Nodo{Clase: claseMayoritaria(datos)}
	}

	campo, valor := mejorDivision(datos)
	if campo == -1 {
		return &Nodo{Clase: claseMayoritaria(datos)}
	}

	izq, der := dividir(datos, campo, valor)
	return &Nodo{
		Campo:     campo,
		Valor:     valor,
		Izquierda: EntrenarArbol(izq, profundidad-1),
		Derecha:   EntrenarArbol(der, profundidad-1),
	}
}

func mejorDivision(datos []Ejemplo) (int, float64) {
	n := len(datos[0].Atributos)
	bestCampo := -1
	bestValor := 0.0
	bestGanancia := 0.0

	for i := 0; i < n; i++ {
		for _, ejemplo := range datos {
			valor := ejemplo.Atributos[i]
			izq, der := dividir(datos, i, valor)
			ganancia := informacionGanada(datos, izq, der)
			if ganancia > bestGanancia {
				bestGanancia = ganancia
				bestCampo = i
				bestValor = valor
			}
		}
	}
	return bestCampo, bestValor
}

func dividir(datos []Ejemplo, campo int, valor float64) ([]Ejemplo, []Ejemplo) {
	izq := []Ejemplo{}
	der := []Ejemplo{}
	for _, e := range datos {
		if e.Atributos[campo] < valor {
			izq = append(izq, e)
		} else {
			der = append(der, e)
		}
	}
	return izq, der
}

func informacionGanada(total, izq, der []Ejemplo) float64 {
	return entropia(total) - (float64(len(izq))/float64(len(total))*entropia(izq) +
		float64(len(der))/float64(len(total))*entropia(der))
}

func entropia(datos []Ejemplo) float64 {
	frecuencias := map[string]int{}
	for _, d := range datos {
		frecuencias[d.Clase]++
	}
	ent := 0.0
	for _, f := range frecuencias {
		p := float64(f) / float64(len(datos))
		ent -= p * log2(p)
	}
	return ent
}

func log2(x float64) float64 {
	return (math.Log(x) / math.Log(2))
}

func claseMayoritaria(datos []Ejemplo) string {
	c := map[string]int{}
	for _, e := range datos {
		c[e.Clase]++
	}
	var mayor string
	max := 0
	for clase, f := range c {
		if f > max {
			max = f
			mayor = clase
		}
	}
	return mayor
}

// Predicción con árbol
func Predecir(n *Nodo, atributos []float64) string {
	if n.Clase != "" {
		return n.Clase
	}
	if atributos[n.Campo] < n.Valor {
		return Predecir(n.Izquierda, atributos)
	}
	return Predecir(n.Derecha, atributos)
}