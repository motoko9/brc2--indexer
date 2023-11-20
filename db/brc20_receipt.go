package db

type brc20ReceiptDao struct {
	dao *Dao
}

func (d *Dao) Brc20Receipt() brc20ReceiptDao {
	return brc20ReceiptDao{
		dao: d,
	}
}

func (d brc20ReceiptDao) Find(name string) (*Brc20Receipt, error) {
	var brc20Receipt Brc20Receipt
	err := d.dao.db.Model(&Brc20Receipt{}).
		Where("name = ?", name).
		First(&brc20Receipt).Error
	return &brc20Receipt, err
}

func (d brc20ReceiptDao) Save(infos []*Brc20Receipt) {
	d.dao.db.Save(infos)
}
