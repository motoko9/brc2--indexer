package db

type syncDao struct {
	dao *Dao
}

func (d *Dao) Sync() syncDao {
	return syncDao{
		dao: d,
	}
}

func (d syncDao) Find() (*Sync, error) {
	var sync Sync
	err := d.dao.db.Model(&Sync{}).First(&sync).Error
	return &sync, err
}

func (d syncDao) Save(sync *Sync) {
	d.dao.db.Save(sync)
}

func (d syncDao) UpdateSyncHeight(height int64) error {
	return d.dao.db.Model(&Sync{}).
		Where("1 = 1").
		Update("sync_height", height).Error
}

func (d syncDao) UpdateCommitHeight(height int64) error {
	return d.dao.db.Model(&Sync{}).
		Where("1 = 1").
		Update("commit_height", height).Error
}
