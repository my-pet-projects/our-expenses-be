package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCatergory_ValidArgs_ReturnsCategory(t *testing.T) {
	t.Parallel()
	// Arrange
	id := "id1"
	parentID := "parentID1"
	name := "name1"
	icon := "icon1"
	lvl := 1
	path := "path1"

	// Act
	res, resErr := NewCategory(id, &parentID, name, &icon, lvl, path)

	// Assert
	assert.NotNil(t, res)
	assert.Nil(t, resErr)
	assert.Equal(t, id, res.ID())
	assert.Equal(t, name, res.Name())
	assert.Equal(t, &icon, res.Icon())
	assert.Equal(t, lvl, res.Level())
	assert.Equal(t, path, res.Path())
}

func TestNewCatergory_InvalidArgs_ThrowsError(t *testing.T) {
	t.Parallel()
	// Arrange
	id := "id2"
	parentID := "parentID2"
	name := ""
	icon := "icon2"
	lvl := 1
	path := "path2"

	// Act
	res, resErr := NewCategory(id, &parentID, name, &icon, lvl, path)

	// Assert
	assert.Nil(t, res)
	assert.NotNil(t, resErr)
}

func TestSetParents_SetsCategoryParents(t *testing.T) {
	t.Parallel()
	// Arrange
	id := "id3"
	parentID := "parentID3"
	name := "name3"
	icon := "icon3"
	lvl := 1
	path := "path3"
	parent1 := Category{id: "par1"}
	parent2 := Category{id: "par2"}

	// SUT
	sut, _ := NewCategory(id, &parentID, name, &icon, lvl, path)

	// Act
	sut.SetParents(&[]Category{parent1, parent2})

	// Assert
	assert.NotNil(t, sut.Parents())
	assert.Contains(t, *sut.Parents(), parent1)
	assert.Contains(t, *sut.Parents(), parent2)
}

func TestIRootCategory_ReturnsIsRoot(t *testing.T) {
	t.Parallel()
	// Arrange
	id := "id4"
	name := "name4"
	icon := "icon4"
	lvl := 1
	path := "path4"
	emptyParentID := ""
	parentID := "parentID4"
	type test struct {
		parentID *string
		isRoot   bool
	}
	tests := []test{
		{parentID: nil, isRoot: true},
		{parentID: &emptyParentID, isRoot: true},
		{parentID: &parentID, isRoot: false},
	}

	// SUT
	for _, tc := range tests {
		// SUT
		sut, _ := NewCategory(id, tc.parentID, name, &icon, lvl, path)

		// Act
		res := sut.IsRoot()

		// Assert
		assert.Equal(t, tc.isRoot, res)
	}
}
