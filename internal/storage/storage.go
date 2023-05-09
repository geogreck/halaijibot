package storage

import (
	"fmt"
	"strconv"

	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const defaultRaiting = 1000

type Storage interface {
	ChangeRaiting(username string, inc int) (int, error)
}

type storage struct {
	db     *bolt.DB
	logger *zap.Logger
}

func New(logger *zap.Logger) (Storage, error) {
	db, err := bolt.Open("raitings.db", 0666, nil)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}

	_, err = tx.CreateBucketIfNotExists([]byte("raitings"))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &storage{
		db:     db,
		logger: logger,
	}, nil
}

func (st *storage) ChangeRaiting(username string, inc int) (int, error) {
	var i int
	err := st.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("raitings"))
		v := b.Get([]byte(username))
		if v == nil {
			b.Put([]byte(username), []byte(strconv.Itoa(defaultRaiting+inc)))
			i = defaultRaiting + inc
			return nil
		}
		num, err := strconv.Atoi(string(v))
		if err != nil {
			st.logger.Error("Corrupted raiting value", zap.String("user", username), zap.ByteString("data", v))
			return err
		}
		fmt.Println(num, inc)
		i = num + inc
		return b.Put([]byte(username), []byte([]byte(strconv.Itoa(num+inc))))
	})

	if err != nil {
		return 0, err
	}

	st.logger.Info("Successfull raiting change", zap.String("user", username), zap.Int("new raiting", i))
	return i, err
}
