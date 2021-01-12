package modelsSdkKeeper

const NotFoundCode = 404

type Response struct {
	Code         int
	Body         Error
	ParsingError error
}
