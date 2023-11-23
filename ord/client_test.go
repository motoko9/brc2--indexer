package ord

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_BlockHeight(t *testing.T) {
	c := New("https://ordinals.com")
	height, err := c.BlockHeight()
	assert.NoError(t, err)
	fmt.Printf("latest blockheight: %d\n", height)
}

func TestClient_BlockHash(t *testing.T) {
	c := New("http://localhost:80")
	hash, err := c.BlockHash()
	assert.NoError(t, err)
	fmt.Printf("latest hash: %s\n", hash)
}

func TestClient_BlockTime(t *testing.T) {
	c := New("http://localhost:80")
	time, err := c.BlockTime()
	assert.NoError(t, err)
	fmt.Printf("latest blocktime: %d\n", time)
}

func TestClient_InscriptionById(t *testing.T) {
	c := New("http://localhost:80")
	inscription, err := c.InscriptionById("d6a05a03ea525fe08e3f2db71ed3adf276d657970964076761de1b42a7caede1i0")
	assert.NoError(t, err)
	fmt.Printf("inscription: %v\n", inscription)
}

func TestClient_Inscriptions(t *testing.T) {
	c := New("http://localhost:80")
	inscriptions, err := c.Inscriptions()
	assert.NoError(t, err)
	fmt.Printf("inscriptions: %v\n", inscriptions)
}

func TestClient_InscriptionsByBlock(t *testing.T) {
	c := New("https://ordiscan.com")
	inscriptions, err := c.InscriptionsByBlock(779630)
	assert.NoError(t, err)
	fmt.Printf("inscriptions: %v\n", inscriptions)
}

func TestClient_InscriptionContent(t *testing.T) {
	c := New("http://localhost:80")
	content, err := c.InscriptionContent("d6a05a03ea525fe08e3f2db71ed3adf276d657970964076761de1b42a7caede1i0")
	assert.NoError(t, err)
	fmt.Printf("content: %v\n", string(content))
}
