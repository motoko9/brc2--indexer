package db

type inscriptionTransactionDao struct {
	dao *Dao
}

func (d *Dao) InscriptionTransaction() inscriptionTransactionDao {
	return inscriptionTransactionDao{
		dao: d,
	}
}

func (d inscriptionTransactionDao) Find(inscriptionId string) (*Transaction, error) {
	var transaction Transaction
	err := d.dao.db.Model(&Transaction{}).
		Where("inscription_id = ?", inscriptionId).
		First(&transaction).Error
	return &transaction, err
}

func (d inscriptionTransactionDao) FindByHeight(height int64) ([]*Transaction, error) {
	var transaction []*Transaction
	err := d.dao.db.Model(&Transaction{}).
		Preload("Receipts").
		Where("height = ?", height).
		Find(&transaction).Error
	return transaction, err
}

func (d inscriptionTransactionDao) Save(infos []*Transaction) error {
	return d.dao.db.Save(infos).Error
}
