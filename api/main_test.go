package api

import (
	"os"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.IStore) *Server {
	config := util.Config{
		DBDriver:            "",
		DBSource:            "",
		ServerAddress:       "",
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server
}
