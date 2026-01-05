package internal

type Hurtable interface {
	DoDamage(int)
	Slow(int)
	KnockOut(int)
}
