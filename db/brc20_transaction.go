package db

type brc20TransactionDao struct {
	dao *Dao
}

func (d *Dao) Brc20Transaction() brc20TransactionDao {
	return brc20TransactionDao{
		dao: d,
	}
}

func (d brc20TransactionDao) Find(name string) (*Brc20Transaction, error) {
	var brc20Transaction Brc20Transaction
	err := d.dao.db.Model(&Brc20Transaction{}).
		Where("name = ?", name).
		First(&brc20Transaction).Error
	return &brc20Transaction, err
}

func (d brc20TransactionDao) Save(infos []*Brc20Transaction) {
	d.dao.db.Save(infos)
}
