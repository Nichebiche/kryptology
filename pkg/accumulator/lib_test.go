package accumulator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/coinbase/kryptology/pkg/core/curves"
)

func TestEvaluatePolyG1(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly := polynomialPoint{
		curve.PointG1.Generator().Mul(curve.Scalar.New(3)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(2)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
	}
	output1, err := poly.evaluate(curve.Scalar.New(1))
	require.NoError(t, err)
	require.NotNil(t, output1)
	result1 := curve.PointG1.Generator().Mul(curve.Scalar.New(6))
	require.Equal(t, output1.ToAffineCompressed(), result1.ToAffineCompressed())

	output2, err := poly.evaluate(curve.Scalar.New(2))
	require.NoError(t, err)
	require.NotNil(t, output2)
	result2 := curve.PointG1.Generator().Mul(curve.Scalar.New(11))
	require.Equal(t, output2.ToAffineCompressed(), result2.ToAffineCompressed())
}

func TestEvaluatePolyG1Error(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly := polynomialPoint{
		nil,
		curve.PointG1.Generator().Mul(curve.Scalar.New(2)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
	}
	_, err := poly.evaluate(curve.Scalar.New(1))
	require.Error(t, err)
}

func TestAddAssignPolyG1(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	// Test polynomial with equal length
	poly1 := polynomialPoint{
		curve.PointG1.Generator().Mul(curve.Scalar.New(3)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(2)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
	}
	poly2 := polynomialPoint{
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(2)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(3)),
	}

	output, err := poly1.Add(poly2)
	require.NoError(t, err)
	require.NotNil(t, output)
	result := polynomialPoint{
		curve.PointG1.Generator().Mul(curve.Scalar.New(4)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(4)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(4)),
	}
	for i := 0; i < len(output); i++ {
		require.Equal(t, output[i].ToAffineCompressed(), result[i].ToAffineCompressed())
	}

	// Test polynomials with unequal length
	poly3 := polynomialPoint{
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(2)),
	}
	output2, err := poly1.Add(poly3)
	require.NoError(t, err)
	require.NotNil(t, output2)
	result2 := polynomialPoint{
		curve.PointG1.Generator().Mul(curve.Scalar.New(4)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(4)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
	}
	require.Equal(t, len(output2), len(result2))
	for i := 0; i < len(output2); i++ {
		require.Equal(t, output2[i].ToAffineCompressed(), result2[i].ToAffineCompressed())
	}

	// Test polynomial with Capacity
	poly4 := make(polynomialPoint, 0, 3)
	poly5, err := poly4.Add(poly1)
	require.NoError(t, err)
	require.Equal(t, len(poly5), len(poly1))
	for i := 0; i < len(poly5); i++ {
		require.Equal(t, poly5[i].ToAffineCompressed(), poly1[i].ToAffineCompressed())
	}

}

func TestAddAssignPolyG1Error(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly1 := polynomialPoint{
		nil,
		curve.PointG1.Generator().Mul(curve.Scalar.New(2)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
	}
	poly2 := polynomialPoint{
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(2)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(3)),
	}
	output, err := poly1.Add(poly2)
	require.Error(t, err)
	require.Nil(t, output)
}

func TestMulAssignPolyG1(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly := polynomialPoint{
		curve.PointG1.Generator().Mul(curve.Scalar.New(3)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(2)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
	}
	rhs := curve.Scalar.New(3)
	output, err := poly.Mul(rhs)
	require.NoError(t, err)
	require.NotNil(t, output)
	poly2 := polynomialPoint{
		curve.PointG1.Generator().Mul(curve.Scalar.New(9)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(6)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(3)),
	}
	for i := 0; i < len(poly2); i++ {
		require.Equal(t, output[i].ToAffineCompressed(), poly2[i].ToAffineCompressed())
	}
}

func TestMulAssignPolyG1Error(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly := polynomialPoint{
		nil,
		curve.PointG1.Generator().Mul(curve.Scalar.New(2)),
		curve.PointG1.Generator().Mul(curve.Scalar.New(1)),
	}
	rhs := curve.Scalar.New(3)
	output, err := poly.Mul(rhs)
	require.Error(t, err)
	require.Nil(t, output)
}

func TestPushPoly(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly := polynomial{
		curve.Scalar.New(3),
		curve.Scalar.New(2),
		curve.Scalar.New(1),
	}
	scalar := curve.Scalar.New(4)
	result := append(poly, scalar)
	require.Equal(t, result[3], scalar)

	// Push one more
	scalar2 := curve.Scalar.New(5)
	result2 := append(result, scalar2)
	require.Equal(t, result2[4], scalar2)

	// Push to a new polynomial
	newPoly := polynomial{}
	newPoly = append(newPoly, scalar)
	require.Equal(t, newPoly[0], scalar)
	newPoly = append(newPoly, scalar2)
	require.Equal(t, newPoly[1], scalar2)
}

func TestAddAssignPoly(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	// Test polynomial with equal length
	poly1 := polynomial{
		curve.Scalar.New(3),
		curve.Scalar.New(2),
		curve.Scalar.New(1),
	}
	poly2 := polynomial{
		curve.Scalar.New(1),
		curve.Scalar.New(2),
		curve.Scalar.New(3),
	}

	output, err := poly1.Add(poly2)
	require.NoError(t, err)
	require.NotNil(t, output)
	result := []curves.Scalar{
		curve.Scalar.New(4),
		curve.Scalar.New(4),
		curve.Scalar.New(4),
	}
	for i := 0; i < len(output); i++ {
		require.Equal(t, output[i], result[i])
	}

	// Test polynomials with unequal length
	poly3 := polynomial{
		curve.Scalar.New(1),
		curve.Scalar.New(2),
	}
	output2, err := poly1.Add(poly3)
	require.NoError(t, err)
	require.NotNil(t, output2)
	result2 := []curves.Scalar{
		curve.Scalar.New(4),
		curve.Scalar.New(4),
		curve.Scalar.New(1),
	}
	require.Equal(t, len(output2), len(result2))
	for i := 0; i < len(output2); i++ {
		require.Equal(t, output2[i], result2[i])
	}
}

func TestAddAssignPolyError(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	// Test polynomial with equal length
	poly1 := polynomial{
		nil,
		curve.Scalar.New(2),
		curve.Scalar.New(1),
	}
	poly2 := polynomial{
		curve.Scalar.New(1),
		curve.Scalar.New(2),
		curve.Scalar.New(3),
	}

	output, err := poly1.Add(poly2)
	require.Error(t, err)
	require.Nil(t, output)
}

func TestSubAssignPoly(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	// Test polynomial with equal length
	poly1 := polynomial{
		curve.Scalar.New(3),
		curve.Scalar.New(2),
		curve.Scalar.New(1),
	}
	poly2 := polynomial{
		curve.Scalar.New(1),
		curve.Scalar.New(2),
		curve.Scalar.New(3),
	}

	output, err := poly1.Sub(poly2)
	require.NoError(t, err)
	require.NotNil(t, output)
	result := []curves.Scalar{
		curve.Scalar.New(2),
		curve.Scalar.New(0),
		curve.Scalar.New(-2),
	}
	for i := 0; i < len(output); i++ {
		require.Equal(t, output[i].Bytes(), result[i].Bytes())
	}

	// Test polynomials with unequal length
	poly3 := polynomial{
		curve.Scalar.New(1),
		curve.Scalar.New(2),
		curve.Scalar.New(3),
		curve.Scalar.New(4),
	}
	output2, err := poly1.Sub(poly3)
	require.NoError(t, err)
	require.NotNil(t, output2)
	result2 := []curves.Scalar{
		curve.Scalar.New(2),
		curve.Scalar.New(0),
		curve.Scalar.New(-2),
		curve.Scalar.New(-4),
	}
	require.Equal(t, len(output2), len(result2))
	for i := 0; i < len(output2); i++ {
		require.Equal(t, output2[i].Bytes(), result2[i].Bytes())
	}
}

func TestSubAssignPolyError(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly1 := polynomial{
		nil,
		curve.Scalar.New(2),
		curve.Scalar.New(1),
	}
	poly2 := polynomial{
		curve.Scalar.New(1),
		curve.Scalar.New(2),
		curve.Scalar.New(3),
	}

	output, err := poly1.Sub(poly2)
	require.Error(t, err)
	require.Nil(t, output)
}

func TestMulAssignPoly(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	// Test polynomial with equal length
	poly1 := polynomial{
		curve.Scalar.New(3),
		curve.Scalar.New(2),
		curve.Scalar.New(1),
	}
	poly2 := polynomial{
		curve.Scalar.New(1),
		curve.Scalar.New(2),
		curve.Scalar.New(3),
	}

	output, err := poly1.Mul(poly2)
	require.NoError(t, err)
	require.NotNil(t, output)
	result := []curves.Scalar{
		curve.Scalar.New(3),
		curve.Scalar.New(8),
		curve.Scalar.New(14),
		curve.Scalar.New(8),
		curve.Scalar.New(3),
	}
	for i := 0; i < len(result); i++ {
		require.Equal(t, output[i].Bytes(), result[i].Bytes())
	}

	// Test polynomials with unequal length
	poly3 := polynomial{
		curve.Scalar.New(1),
		curve.Scalar.New(2),
	}
	output2, err := poly1.Mul(poly3)
	require.NoError(t, err)
	require.NotNil(t, output2)
	result2 := []curves.Scalar{
		curve.Scalar.New(3),
		curve.Scalar.New(8),
		curve.Scalar.New(5),
		curve.Scalar.New(2),
	}
	require.Equal(t, len(output2), 4)
	for i := 0; i < len(output2); i++ {
		require.Equal(t, output2[i].Bytes(), result2[i].Bytes())
	}
}

func TestMulAssignPolyError(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly1 := polynomial{
		nil,
		curve.Scalar.New(2),
		curve.Scalar.New(1),
	}
	poly2 := polynomial{
		curve.Scalar.New(1),
		curve.Scalar.New(2),
		curve.Scalar.New(3),
	}
	output, err := poly1.Mul(poly2)
	require.Error(t, err)
	require.Nil(t, output)
}

func TestMulValueAssignPoly(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly := polynomial{
		curve.Scalar.New(3),
		curve.Scalar.New(2),
		curve.Scalar.New(1),
	}
	rhs := curve.Scalar.New(3)
	output, err := poly.MulScalar(rhs)
	require.NoError(t, err)
	require.NotNil(t, output)
	coefficients2 := []curves.Scalar{
		curve.Scalar.New(9),
		curve.Scalar.New(6),
		curve.Scalar.New(3),
	}
	for i := 0; i < len(coefficients2); i++ {
		require.Equal(t, output[i].Bytes(), coefficients2[i].Bytes())
	}
}

func TestMulValueAssignPolyError(t *testing.T) {
	curve := curves.BLS12381(&curves.PointBls12381G1{})
	poly := polynomial{
		nil,
		curve.Scalar.New(2),
		curve.Scalar.New(1),
	}
	rhs := curve.Scalar.New(3)
	output, err := poly.MulScalar(rhs)
	require.Error(t, err)
	require.Nil(t, output)
}

func TestPerformanceDeviation(t *testing.T) {
	goodComputationBench := "Benchmark_ABC             842688              1396 ns/op"
	badComputationBench := "Benchmark_ABC             842688              13960 ns/op"
	if err := Compare(strings.NewReader(goodComputationBench), strings.NewReader(badComputationBench)); err == nil {
		t.Errorf("Expected performance deviation: [%v], but did not detect any!", err)
	}
}

func TestParsing(t *testing.T) {
	// TODO: the current parser ignores the 3rd and 4th column of data (e.g., the custom benchmarks)
	output := `
garbage data
BenchmarkSigning/Secp256k1_-_5_of_9-16                         1        5794642205 ns/op
BenchmarkSign2p-16                                             2         685590314 ns/op             29319 bytes/sign           16.00 msgs/sign
garbage data
`
	o, _, err := parseBenchmarks(strings.NewReader(output), strings.NewReader(output))
	if err != nil {
		t.Errorf("Failed to parse test input %v", err)
	}

	if len(o) != 2 {
		t.Errorf("Incorrect output length. Expected 2, got %#v", o)
	}
	if _, ok := o["BenchmarkSigning/Secp256k1_-_5_of_9-16"]; !ok {
		t.Errorf("Did not find BenchmarkSigning/Secp256k1_-_5_of_9-16 in the parsed output")
	}
	if _, ok := o["BenchmarkSign2p-16"]; !ok {
		t.Errorf("Did not find BenchmarkSign2p-16 in the parsed output")
	}
}
