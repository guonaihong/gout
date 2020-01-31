package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestGetNoPortExists(t *testing.T) {
	port := GetNoPortExists()
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	assert.NoError(t, err)
	if l != nil {
		l.Close()
	}

}
