package reaction

//go:generate go run github.com/dmarkham/enumer -type=Reaction -transform=snake -trimprefix=Reaction
type Reaction int

const (
	ReactionNone Reaction = iota
	ReactionLike
	ReactionDislike
)
