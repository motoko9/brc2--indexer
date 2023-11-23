package db

type receiptDao struct {
	dao *Dao
}

func (d *Dao) Receipt() receiptDao {
	return receiptDao{
		dao: d,
	}
}

func (d receiptDao) Save(infos []*Receipt) error {
	return d.dao.db.Save(infos).Error
}
