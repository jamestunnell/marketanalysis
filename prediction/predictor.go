package prediction

// import (
// 	"github.com/jamestunnell/marketanalysis/models"
// 	"github.com/jamestunnell/marketanalysis/processing"
// 	"github.com/jamestunnell/marketanalysis/provision"
// )

// //go:generate mockgen -destination=mock_prediction/mocks.go . Predictor

// type EachPositionFunc func(*models.Position)

// type Predictor interface {
// 	Train(chain *processing.Chain,
// 		seqs provision.BarSequences) error
// 	Test(chain *processing.Chain,
// 		seqs provision.BarSequences) (models.Positions, error)
// 	Predict(chain *processing.Chain,
// 		bars provision.BarSequence,
// 		each EachPositionFunc) error
// }
