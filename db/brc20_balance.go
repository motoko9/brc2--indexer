package db

type brc20BalanceDao struct {
	dao *Dao
}

func (d *Dao) Brc20Balance() brc20BalanceDao {
	return brc20BalanceDao{
		dao: d,
	}
}

func (d brc20BalanceDao) Find(address string) (*Brc20Balance, error) {
	var brc20Balance Brc20Balance
	err := d.dao.db.Model(&Brc20Balance{}).
		Where("address = ?", address).
		First(&brc20Balance).Error
	return &brc20Balance, err
}

func (d brc20BalanceDao) Save(info []*Brc20Balance) {
	d.dao.db.Save(info)
}
