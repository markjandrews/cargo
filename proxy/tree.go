package proxy

import (
	"errors"
	"github.com/rs/zerolog/log"
	"strings"
)

// Public Errors
var (
	ErrExists   = errors.New("item already exists")
	ErrNotFound = errors.New("item not found")
)

type Context int

type Tree struct {
	Context *Context
	nodes   map[string]*Tree
}

func NewTree() *Tree {
	return &Tree{nodes: map[string]*Tree{}}
}

func (t *Tree) Add(path string, context *Context) error {
	if len(path) == 0 {
		err := ErrInvalidArg
		log.Error().Err(err).Msg("Path must not be empty")
		return err
	}

	parts := strings.Split(path, ".")
	ReverseStringList(parts)

	if t.nodes == nil {
		t.nodes = map[string]*Tree{}
	}

	return t.add(parts, context)
}

func (t *Tree) add(path []string, context *Context) error {
	if path == nil {
		err := ErrInvalidArg
		log.Error().Err(err).Msg("Path must not be nil")
		return err
	}

	pathLen := len(path)

	if pathLen == 0 {
		err := ErrInvalidArg
		log.Error().Err(err).Msg("Path must not be empty")
		return err
	}

	if t.nodes == nil {
		t.nodes = map[string]*Tree{}
	}

	root := t
	for i, part := range path {
		if len(part) == 0 {
			err := ErrInvalidArg
			log.Error().Err(err).Msg("Path must not contain empty part")
			return err
		}

		child, ok := root.nodes[part]
		if !ok {
			child, ok = root.nodes["*"]
		}

		if !ok {
			child = NewTree()
			root.nodes[part] = child

			if pathLen-i == 1 {
				child.Context = context
				return nil
			}
		}

		if pathLen-i == 1 {
			err := ErrExists
			log.Error().Err(err).Msgf("Child Node: %s already exists with value: %v", path[0], child.Context)
			return err
		}

		root = child
	}

	return nil
}

func (t *Tree) Get(path string) (*Context, error) {
	if len(path) == 0 {
		err := ErrInvalidArg
		log.Error().Err(err).Msg("Path must not be empty")
		return nil, err
	}

	parts := strings.Split(path, ".")
	ReverseStringList(parts)

	ctx, err := t.get(parts)
	if err != nil {
		log.Error().Err(err).Msgf("Item for path: %s not found", path)
		return nil, err
	}

	return ctx, nil
}

func (t *Tree) get(path []string) (*Context, error) {
	if t.nodes == nil {
		err := ErrInvalidData
		log.Error().Err(err).Msg("Tree is not initialized correctly. Use NewTree() not &Tree{}")
		return nil, err
	}

	if path == nil {
		err := ErrInvalidArg
		log.Error().Err(err).Msg("Path must not be nil")
		return nil, err
	}

	pathLen := len(path)

	if pathLen == 0 {
		err := ErrInvalidArg
		log.Error().Err(err).Msg("Path must not be empty")
		return nil, err
	}

	var isWild = false

	root := t
	for i, part := range path {
		if len(part) == 0 {
			err := ErrInvalidArg
			log.Error().Err(err).Msg("Path must not contain empty part")
			return nil, err
		}

		child, ok := root.nodes[part]
		if ok {
			isWild = false
		} else {
			child, ok = root.nodes["*"]
			if ok {
				isWild = true
			}
		}

		if !ok {
			if isWild && root.Context != nil {
				return root.Context, nil
			}

			return nil, ErrNotFound
		}

		if pathLen-i == 1 {
			return child.Context, nil
		}

		root = child
	}

	ctx := Context(0)
	return &ctx, nil
}
