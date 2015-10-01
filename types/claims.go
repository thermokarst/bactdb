package types

type Claims struct {
	Name string
	Iss  string
	Sub  int64
	Role string
	Iat  int64
	Exp  int64
	Ref  string
}
