package db

type inscriptionDao struct {
	dao *Dao
}

func (d *Dao) Inscription() inscriptionDao {
	return inscriptionDao{
		dao: d,
	}
}

func (d inscriptionDao) Find(name string) (*Inscription, error) {
	var inscription Inscription
	err := d.dao.db.Model(&Inscription{}).
		Where("name = ?", name).
		First(&inscription).Error
	return &inscription, err
}
