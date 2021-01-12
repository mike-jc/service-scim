package controllers

type IndexController struct {
	AbstractController
}

// @Title Home page
// @router / [get]
func (c *IndexController) Home() {
	c.SuccessResponse()
}
