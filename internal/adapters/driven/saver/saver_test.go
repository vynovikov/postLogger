package saver

import (
	"errors"
	"postLogger/internal/logger"
	"testing"

	"github.com/stretchr/testify/suite"
)

type saverSuite struct {
	suite.Suite
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(saverSuite))
}

func (s *saverSuite) TestNewSaver() {
	rootFolder := "../../../../logs"
	saver, err := NewSaver(rootFolder, int64(100))
	logger.L.Infof("in saver.TestNewSaver saver %v, err %v\n", saver, err)

}

var testStore *SaverStruct

func (s *saverSuite) TestSave() {
	testStore, _ = NewSaver("../../../../logs", int64(25))

	tt := []struct {
		name      string
		id        int
		logString string
		wantError error
	}{
		{
			name:      "no file",
			id:        0,
			logString: "000000000000000",
			wantError: errors.New(""),
		},
		{
			name:      "file exists",
			id:        1,
			logString: "111111111111111",
		},
		{
			name:      "file len exceeded",
			id:        2,
			logString: "222222222222222",
		},
	}
	for _, v := range tt {
		s.Run(v.name, func() {
			err := testStore.Save(v.logString)
			if err != nil {
				s.Equal(v.wantError, err)
			}

		})
	}
}
