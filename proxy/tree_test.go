package proxy

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func init() {
	log.Logger = log.With().Caller().Logger()
}

func TestNewTree(t *testing.T) {
	tr := Tree{}
	assert.Nil(t, tr.nodes)

	tr2 := NewTree()
	assert.NotNil(t, tr2.nodes)
}

func TestTree_Add(t *testing.T) {
	host := "www.example.com.au"
	hostParts := strings.Split(host, ".")
	ReverseStringList(hostParts)

	tr := Tree{}
	ctx := Context(6)
	assert.NotPanics(t, func() {tr.add(hostParts, &ctx)})

	tr2 := Tree{}
	ctx2 := Context(10)
	err := tr2.add(hostParts, &ctx2)
	assert.Nil(t, err)
	assert.True(t, *tr2.nodes["au"].nodes["com"].nodes["example"].nodes["www"].Context == 10)
}

func TestTree_AddWildcard(t *testing.T) {
	host := "www.example.*.au"
	hostParts := strings.Split(host, ".")
	ReverseStringList(hostParts)

	tr := Tree{}
	ctx := Context(6)
	assert.NotPanics(t, func() {tr.add(hostParts, &ctx)})

	tr2 := Tree{}
	ctx2 := Context(10)
	err := tr2.add(hostParts, &ctx2)
	assert.Nil(t, err)
	assert.True(t, *tr2.nodes["au"].nodes["*"].nodes["example"].nodes["www"].Context == 10)
}


func TestTree_AddNil(t *testing.T) {
	tr := Tree{}
	ctx := Context(8)
	err := tr.add(nil, &ctx)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrInvalidArg)
}

func TestTree_AddEmpty(t *testing.T) {
	tr := Tree{}
	ctx := Context(8)
	err := tr.Add("", &ctx)
	assert.NotNil(t, err)
	assert.Equal(t, err, ErrInvalidArg)

	tr2 := Tree{}
	ctx2 := Context(8)
	err2 := tr2.add([]string{""}, &ctx2)
	assert.NotNil(t, err2)
	assert.Equal(t, err2, ErrInvalidArg)
}

func TestTree_AddRawStruct(t *testing.T) {
	tr := Tree{}
	assert.Nil(t, tr.nodes)

	ctx := Context(8)
	tr.add([]string{"test"}, &ctx)
	assert.NotNil(t, tr.nodes)
}

func TestTree_AddNewStruct(t *testing.T) {
	tr := NewTree()
	assert.NotNil(t, tr.nodes)
}

func TestTree_Get(t *testing.T) {
	tr := NewTree()
	ctx := Context(10)
	err := tr.Add("www.example.com", &ctx)
	assert.Nil(t, err)

	ctxResult, err := tr.Get("www.example.com")
	assert.Nil(t, err)
	assert.Equal(t, *ctxResult, Context(10))

	tr2 := NewTree()
	ctx2 := Context(20)
	err = tr2.Add("www.fred.com", &ctx2)
	assert.Nil(t, err)

	_, err = tr2.Get("www.example.com")
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestTree_GetWildCatchAll(t *testing.T) {
	tr := NewTree()
	ctx := Context(10)
	err := tr.Add("*", &ctx)
	assert.Nil(t, err)

	ctxResult, err := tr.Get("www.example.com")
	require.Nil(t, err)
	require.Equal(t, *ctxResult, Context(10))
}

func TestTree_GetWildCatchSpecific(t *testing.T) {
	tr := NewTree()
	ctx := Context(10)
	err := tr.Add("www.*.com", &ctx)
	assert.Nil(t, err)

	ctxResult, err := tr.Get("www.example.com")
	require.Nil(t, err)
	require.Equal(t, *ctxResult, Context(10))

	ctxResult2, err := tr.Get("www.fred.com")
	require.Nil(t, err)
	require.Equal(t, *ctxResult2, Context(10))

	_, err = tr.Get("www.fred.com.bz")
	require.NotNil(t, err)
	require.Equal(t, err, ErrNotFound)

	_, err = tr.Get("bill.example.com")
	require.NotNil(t, err)
	require.Equal(t, err, ErrNotFound)
}

func TestTree_GetWildCatchTrailing(t *testing.T) {
	tr := NewTree()
	ctx := Context(10)
	err := tr.Add("www.example.*", &ctx)
	assert.Nil(t, err)

	ctxResult, err := tr.Get("www.example.com")
	require.Nil(t, err)
	require.Equal(t, *ctxResult, Context(10))

	ctxResult2, err := tr.Get("www.example.biz")
	require.Nil(t, err)
	require.Equal(t, *ctxResult2, Context(10))

	_, err = tr.Get("www.example.com.biz")
	require.NotNil(t, err)
	require.Equal(t, err, ErrNotFound)
}

func TestTree_GetMultiLevel(t *testing.T) {
	tr := NewTree()
	ctx := Context(10)
	err := tr.Add("www.example.com", &ctx)
	assert.Nil(t, err)

	ctx2 := Context(20)
	err = tr.Add("www.example.com.biz", &ctx2)
	assert.Nil(t, err)

	ctxResult, err := tr.Get("www.example.com")
	require.Nil(t, err)
	require.Equal(t, *ctxResult, Context(10))

	ctxResult2, err := tr.Get("www.example.com.biz")
	require.Nil(t, err)
	require.Equal(t, *ctxResult2, Context(20))

	_, err = tr.Get("www.example.biz")
	require.NotNil(t, err)
	require.Equal(t, err, ErrNotFound)
}