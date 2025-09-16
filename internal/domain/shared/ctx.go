package shared

type RequestCtx struct {
	User UserCtx
}

type UserCtx struct {
	UserID IDType
	Role   UserRole
}

func NewRequestCtx(userID IDType, role UserRole) RequestCtx {
	return RequestCtx{
		User: UserCtx{
			UserID: userID,
			Role:   role,
		},
	}
}
