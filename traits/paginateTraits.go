package traits

//TODO ADD CHECKING OF NUMBER OF TOTAL PAGES IF THE PAGE LIMIT IS SET
import "gorm.io/gorm"

type Paginate struct {
	Limit int
	Page  int
}

func NewPaginate(limit int, page int) *Paginate {
	return &Paginate{Limit: limit, Page: page}
}

func (p *Paginate) PagintedResult(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.Limit

	return db.Offset(offset).Limit(p.Limit)
}
