package utils

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sony/sonyflake"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSonyflake(t *testing.T) {
	st := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := st.NextID()
	assert.NoError(t, err)
	t.Log(id)
}
