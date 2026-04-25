package mealdomain

import "fmt"

type MealType string
type Category string

const (
	MealTypeAlmuerzo MealType = "almuerzo"
	MealTypeCena     MealType = "cena"
)

const (
	CategoryComida          Category = "comida"
	CategoryPastas          Category = "pastas"
	CategoryMilanesas       Category = "milanesas"
	CategoryEnsaladas       Category = "ensaladas"
	CategorySandwichesWraps Category = "sandwiches_y_wraps"
	CategoryPollo           Category = "pollo"
	CategoryCarne           Category = "carne"
)

var validMealTypes = map[MealType]bool{
	MealTypeAlmuerzo: true,
	MealTypeCena:     true,
}

var validCategories = map[Category]bool{
	CategoryComida:          true,
	CategoryPastas:          true,
	CategoryMilanesas:       true,
	CategoryEnsaladas:       true,
	CategorySandwichesWraps: true,
	CategoryPollo:           true,
	CategoryCarne:           true,
}

func NewMealType(value string) (MealType, error) {
	mt := MealType(value)
	if !validMealTypes[mt] {
		return "", fmt.Errorf("%w: %s", ErrInvalidMealType, value)
	}
	return mt, nil
}

func NewCategory(value string) (Category, error) {
	c := Category(value)
	if !validCategories[c] {
		return "", fmt.Errorf("%w: %s", ErrInvalidCategory, value)
	}
	return c, nil
}

func (mt MealType) String() string {
	return string(mt)
}

func (c Category) String() string {
	return string(c)
}
