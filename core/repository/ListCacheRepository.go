package repo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// ListCacheDBEntry is a contract which is serialized to BSON and saved in MongoDB.
type ListCacheDBEntry struct {
	Key         string
	Values      []string
	ExpireAfter int64
	Added       int64
	Updated     int64
}

// IListCacheRepository is an interface for ListCacheRepository.
type IListCacheRepository interface {
	GetAll() []ListCacheDBEntry
	SaveAll(newEntries []ListCacheDBEntry, updatedEntries []ListCacheDBEntry)
}

// ListCacheRepository for persisting cache entries to MongoDB.
type ListCacheRepository struct {
	Host    string
	DBName  string
	ColName string
}

// GetAll returns cache snapshot from DB.
func (r ListCacheRepository) GetAll() []ListCacheDBEntry {
	session, err := mgo.Dial(r.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(r.DBName).C(r.ColName)
	var result []ListCacheDBEntry
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		log.Fatal(err)
	}
	if result != nil {
		log.Printf("[ListCacheDBEntry] Read snapshot from %s.%s successfully.", r.DBName, r.ColName)
	}
	return result
}

// SaveAll saves cache snapshot to DB.
func (r ListCacheRepository) SaveAll(newEntries []ListCacheDBEntry, updatedEntries []ListCacheDBEntry) {
	session, err := mgo.Dial(r.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(r.DBName).C(r.ColName)
	c.EnsureIndexKey("key")
	var existingKeys []string

	for _, entry := range newEntries {
		existingKeys = append(existingKeys, entry.Key)
		e := c.Insert(entry)
		if e != nil {
			log.Fatal(err)
		}
	}

	for _, entry := range updatedEntries {
		existingKeys = append(existingKeys, entry.Key)
		if entry.Updated > entry.Added {
			e := c.Update(bson.M{"key": entry.Key}, bson.M{"$set": bson.M{"values": entry.Values, "updated": entry.Updated}})
			if e != nil {
				log.Fatal(err)
			}
		}
	}
	if existingKeys == nil {
		existingKeys = make([]string, 0)
	}
	_, e := c.RemoveAll(bson.M{"key": bson.M{"$nin": existingKeys}})
	if e != nil {
		log.Fatal(err)
	}
	log.Printf("[ListCacheDBEntry] Persisted data to %s.%s successfully.", r.DBName, r.ColName)
}

// EmptyListCacheRepository for testing only.
type EmptyListCacheRepository struct {
}

// GetAll returns empty cache.
func (r EmptyListCacheRepository) GetAll() []ListCacheDBEntry {
	return make([]ListCacheDBEntry, 0)
}

// SaveAll saves nothing.
func (r EmptyListCacheRepository) SaveAll(newEntries []ListCacheDBEntry, updatedEntries []ListCacheDBEntry) {
}