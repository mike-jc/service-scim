package modelsSdkJwt

import (
	"strconv"
)

const DomainUser = "user"
const DomainUserId = "1"

const ServiceKeeper = "keeper"

type ClaimsUserPayload struct {
	Domain   string `json:"domain"`
	Id       string `json:"id"`
	Instance string `json:"instance"`
}

func (o *ClaimsUserPayload) GetIdInt64() (id int64) {
	id, _ = strconv.ParseInt(o.Id, 10, 0)
	return
}
