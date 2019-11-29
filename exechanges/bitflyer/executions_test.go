package bitflyer

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestAppend(t *testing.T) {
	v := make([]float64, 12)
	for i := 0; i < 12; i++ {
		v[i] = float64(i)
	}
	// Create a new matrix
	A := mat.NewDense(3, 4, v)
	r, c := A.Dims()
	B := mat.NewDense(r+1, c, nil)
	add := mat.NewDense(1, c, []float64{11, 11, 11, 11})
	B.Stack(A, add)
	println("B:")
	fmt.Printf("%+v\n", mat.Formatted(B))

	A = B
	println("A:")
	fmt.Printf("%+v\n", mat.Formatted(B))
}
