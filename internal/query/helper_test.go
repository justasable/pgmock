package query_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/justasable/pgmock/internal/connect"
	"github.com/stretchr/testify/suite"
)

type DefaultSuite struct {
	suite.Suite
	conn *pgx.Conn
}

func TestDefaultSuite(t *testing.T) {
	suite.Run(t, new(DefaultSuite))
}

func (s *DefaultSuite) SetupTest() {
	err := connect.SetupDBWithScript("test_setup.sql")
	s.NoError(err)

	conn, err := connect.Connect()
	s.NoError(err)

	s.conn = conn
}

func (s *DefaultSuite) TearDownTest() {
	s.conn.Close(context.Background())
}
