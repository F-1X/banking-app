package banking

type CustomError struct {
	describtion string
}

func (c *CustomError) Error() string {
	return c.describtion
}

var (
	NotEnough = &CustomError{describtion: "not enough money on balance"}
)
