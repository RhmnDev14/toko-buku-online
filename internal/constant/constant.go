package constant

const (
	// role
	Admin = "admin"
	User  = "user"
	// keyCtx
	UserIDKey = "userID"
	RoleKey   = "role"
	// method
	GET     = "get"
	PUT     = "put"
	DELETE  = "delete"
	CREATE  = "create"
	GETBOOK = "get book"
	ORDER   = "order"
	// status pembayaran
	PENDING   = "PENDING"
	PAID      = "PAID"
	CANCELLED = "CANCELLED"
	//message grpc
	Succes = "Succes"
	// error db
	ErrorServerCreate     = "maaf, terjadi kesalahan saat menyimpan data. Silakan coba beberapa saat lagi"
	ErrorServerGet        = "maaf, data tidak dapat diambil saat ini. Silakan coba beberapa saat lagi"
	ErrorServerUpdate     = "maaf, terjadi kesalahan saat memperbarui data. Silakan coba beberapa saat lagi"
	ErrorDataNotFound     = "data tidak ditemukan"
	ErrorEmailHasBeenUsed = "email sudah dipakai"
	// error be
	ErrorInternalSystem = "maaf, terjadi kesalahan pada sistem. Silahkan coba beberapa saat lagi"
	ErrorLogin          = "email atau password tidak sesuai"
	ErrorDontPermission = "maaf, kamu tidak memiliki hak akses"
)
