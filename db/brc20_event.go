package db

type brc20EventDao struct {
	dao *Dao
}

func (d *Dao) Brc20Event() brc20EventDao {
	return brc20EventDao{
		dao: d,
	}
}

func (d brc20EventDao) Find(name string) (*Brc20Event, error) {
	var brc20Event Brc20Event
	err := d.dao.db.Model(&Brc20Event{}).
		Where("name = ?", name).
		First(&brc20Event).Error
	return &brc20Event, err
}

func (d brc20EventDao) Save(infos []*Brc20Event) {
	d.dao.db.Save(infos)
}
