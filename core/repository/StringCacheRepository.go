package repo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// StringCacheDBEntry is a contract which is serialized to BSON and saved in MongoDB.
type StringCacheDBEntry struct {
	Key         string
	Value       string
	ExpireAfter int64
	Added       int64
	Updated     int64
}

// StringCacheRepository for persisting cache entries to MongoDB.
type StringCacheRepository struct {
	Host    string
	DBName  string
	ColName string
}

// GetAll returns cache snapshot from DB.
func (r *StringCacheRepository) GetAll() []StringCacheDBEntry {
	session, err := mgo.Dial(r.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(r.DBName).C(r.ColName)
	var result []StringCacheDBEntry
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// SaveAll saves cache snapshot to DB.
func (r *StringCacheRepository) SaveAll(newEntries []StringCacheDBEntry, updatedEntries []StringCacheDBEntry) {
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
			e := c.Update(bson.M{"key": entry.Key}, bson.M{"$set": bson.M{"value": entry.Value, "updated": entry.Updated}})
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
}
