package entity

// grpc auth validation
type AuthRequest struct {
	AccessToken string
}

// example:
// "accessToken": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0ODc5NzU0NTgsImlhdCI6MTQ4N
// zk3MTg1OCwiaXNzIjoiYWNtZS5jb20iLCJzdWIiOiI4NThhNGIwMS02MmM4LTRjMmYtYmZhNy02ZDAxODgzM2JlYTc
// iLCJhcHBsaWNhdGlvbklkIjoiM2MyMTllNTgtZWQwZS00YjE4LWFkNDgtZjRmOTI3OTNhZTMyIiwicm9sZXMiOlsiY
// WRtaW4iXX0.O29_m_NDa8Cj7kcpV7zw5BfFmVGsK1n3EolCj5u1M9hZ09EnkaOl5n68OLsIcpCrX0Ue58qsabag3MC
// NS6H4ldt6kMnH6k4bVg4TvIjoR8WE-yGcu_xDUObYKZYaHWiNeuDL1EuQQI_8HajQLND-c9juy5ILuz6Fhx8CLfHCz
// iEHX_aQPt7jQ2IIasVzprKkgvWS07Hiv2Oskryx49wqCesl46b-30c6nfttHUDEQrVq9gaepca3Nhjj_cPtC400JgL
// CN9DOYIbtd69zvD8vDUOvVzMr2HGdWtKthqa35NF-3xMZKD8CShe8ZT74fNd9YZ0WRE-YeIf3T_Hv5p5V2w",

// "refreshToken": "YRjxLpsjRqL7zYuKstXogqioA_P3Z4fiEuga0NCVRcDSc8cy_9msxg"

type AuthResponse struct {
	Username string // Username - email
	Error    string
}

type RefreshTokenCredits struct { // refresh
	RefreshToken string
}

type UserCredits struct {
	Username     string
	AccessToken  string
	RefreshToken string
	Error        string
}
