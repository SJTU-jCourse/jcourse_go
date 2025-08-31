package shared

type RequestCtx struct {
	User UserCtx
}

type UserCtx struct {
	UserID IDType
	Role   UserRole
}
