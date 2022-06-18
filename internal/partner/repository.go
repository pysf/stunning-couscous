package partner

import "context"

type Repository interface {
	createPartner(context.Context, Partner) error
}
