package entity

type Sales struct {
	Omset            float64 `json:"omset"`
	TotalBukuTerjual int     `json:"total_buku_terjual"`
}

type BestBook struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type BestSeller struct {
	BestBook []BestBook `json:"buku"`
}

type PriceBook struct {
	Max float64 `json:"max"`
	Min float64 `json:"min"`
	Avg float64 `json:"avg"`
}
