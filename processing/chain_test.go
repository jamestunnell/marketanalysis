package processing_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/jamestunnell/marketanalysis/processing"
	"github.com/jamestunnell/marketanalysis/processing/mock_processing"
)

func TestChain_NoProcs(t *testing.T) {
	c := processing.NewChain()

	assert.False(t, c.Warm())

	assert.NoError(t, c.Initialize())

	assert.NoError(t, c.WarmUp([]float64{}))

	assert.True(t, c.Warm())
	assert.Equal(t, 0.0, c.Output())

	assert.NoError(t, c.WarmUp([]float64{1.0, 2.0}))

	assert.True(t, c.Warm())
	assert.Equal(t, 2.0, c.Output())
}

func TestChain_OneProc(t *testing.T) {
	ctrl := gomock.NewController(t)
	proc := mock_processing.NewMockProcessor(ctrl)
	c := processing.NewChain(proc)

	assert.False(t, c.Warm())

	proc.EXPECT().Initialize().Return(nil)
	proc.EXPECT().WarmupPeriod().AnyTimes().Return(2)

	assert.NoError(t, c.Initialize())

	assert.Equal(t, 2, c.WarmupPeriod())

	assert.Error(t, c.WarmUp([]float64{}))
	assert.Error(t, c.WarmUp([]float64{1.5}))

	inputVals := []float64{1.5, 2.5}

	proc.EXPECT().WarmUp(inputVals).Return(nil)
	proc.EXPECT().Output().Return(3.5)

	assert.NoError(t, c.WarmUp(inputVals))

	assert.True(t, c.Warm())
	assert.Equal(t, 3.5, c.Output())

	inputVals = []float64{1.5, 2.5, 3.5}

	proc.EXPECT().WarmUp(inputVals).Return(nil)
	proc.EXPECT().Output().Return(4.5)

	assert.NoError(t, c.WarmUp(inputVals))

	assert.True(t, c.Warm())
	assert.Equal(t, 4.5, c.Output())
}

func TestChain_TwoProcs(t *testing.T) {
	const (
		warmup1 = 3
		warmup2 = 2
	)

	ctrl := gomock.NewController(t)
	proc1 := mock_processing.NewMockProcessor(ctrl)
	proc2 := mock_processing.NewMockProcessor(ctrl)
	c := processing.NewChain(proc1, proc2)

	assert.False(t, c.Warm())

	gomock.InOrder(
		proc1.EXPECT().Initialize().Return(nil),
		proc1.EXPECT().WarmupPeriod().Return(warmup1),

		proc2.EXPECT().Initialize().Return(nil),
		proc2.EXPECT().WarmupPeriod().Return(warmup2),
	)

	assert.NoError(t, c.Initialize())

	assert.Equal(t, 4, c.WarmupPeriod())

	assert.Error(t, c.WarmUp([]float64{}))
	assert.Error(t, c.WarmUp([]float64{1.5}))
	assert.Error(t, c.WarmUp([]float64{1.5, 2.5}))
	assert.Error(t, c.WarmUp([]float64{1.5, 2.5, 3.5}))

	inputVals := []float64{1.5, 2.5, 3.5, 4.5}

	gomock.InOrder(
		proc1.EXPECT().WarmupPeriod().Return(warmup1),
		proc1.EXPECT().WarmUp([]float64{1.5, 2.5, 3.5}).Return(nil),
		proc1.EXPECT().Output().Return(5.5),

		proc2.EXPECT().WarmupPeriod().Return(warmup2),

		proc1.EXPECT().Update(4.5),
		proc1.EXPECT().Output().Return(6.5),

		proc2.EXPECT().WarmUp([]float64{5.5, 6.5}).Return(nil),
		proc2.EXPECT().Output().Return(7.5),
	)

	assert.NoError(t, c.WarmUp(inputVals))

	assert.True(t, c.Warm())
	assert.Equal(t, 7.5, c.Output())
}
