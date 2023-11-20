package db

type brc20InfoDao struct {
	dao *Dao
}

func (d *Dao) Brc20Info() brc20InfoDao {
	return brc20InfoDao{
		dao: d,
	}
}

func (d brc20InfoDao) Find(name string) (*Brc20Info, error) {
	var rrc20Info Brc20Info
	err := d.dao.db.Model(&Brc20Info{}).
		Where("name = ?", name).
		First(&rrc20Info).Error
	return &rrc20Info, err
}

func (d brc20InfoDao) Save(info *Brc20Info) {
	d.dao.db.Save(info)
}
