package db

type inscriptionDao struct {
	dao *Dao
}

func (d *Dao) Inscription() inscriptionDao {
	return inscriptionDao{
		dao: d,
	}
}

func (d inscriptionDao) Find(inscriptionId string) (*Inscription, error) {
	var inscription Inscription
	err := d.dao.db.Model(&Inscription{}).
		Where("inscription_id = ?", inscriptionId).
		First(&inscription).Error
	return &inscription, err
}

func (d inscriptionDao) Save(inscriptions []*Inscription) error {
	return d.dao.db.Save(inscriptions).Error
}
