package repo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// DictionaryValueDBEntry is a dictionary key/value contract which is serialized to BSON and saved in MongoDB.
type DictionaryValueDBEntry struct {
	Key string
	Value string
}

// DictionaryCacheDBEntry is a contract which is serialized to BSON and saved in MongoDB.
type DictionaryCacheDBEntry struct {
	Key         string
	Values      []DictionaryValueDBEntry
	ExpireAfter int64
	Added       int64
	Updated     int64
}

// IDictionaryCacheRepository is an interface for DictionaryCacheRepository.
type IDictionaryCacheRepository interface {
	GetAll() []DictionaryCacheDBEntry
	SaveAll(newEntries []DictionaryCacheDBEntry, updatedEntries []DictionaryCacheDBEntry)
}

// DictionaryCacheRepository for persisting cache entries to MongoDB.
type DictionaryCacheRepository struct {
	Host    string
	DBName  string
	ColName string
}

// GetAll returns cache snapshot from DB.
func (r DictionaryCacheRepository) GetAll() []DictionaryCacheDBEntry {
	session, err := mgo.Dial(r.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB(r.DBName).C(r.ColName)
	var result []DictionaryCacheDBEntry
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		log.Fatal(err)
	}
	if result != nil {
		log.Printf("[DictionaryCacheDBEntry] Read snapshot from %s.%s successfully.", r.DBName, r.ColName)
	}
	return result
}

// SaveAll saves cache snapshot to DB.
func (r DictionaryCacheRepository) SaveAll(newEntries []DictionaryCacheDBEntry, updatedEntries []DictionaryCacheDBEntry) {
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
	log.Printf("[DictionaryCacheDBEntry] Persisted data to %s.%s successfully.", r.DBName, r.ColName)
}

// EmptyDictionaryCacheRepository for testing only.
type EmptyDictionaryCacheRepository struct {

}

// GetAll returns empty cache.
func (r EmptyDictionaryCacheRepository) GetAll() []DictionaryCacheDBEntry {
	return make([]DictionaryCacheDBEntry, 0)
}

// SaveAll saves nothing.
func (r EmptyDictionaryCacheRepository) SaveAll(newEntries []DictionaryCacheDBEntry, updatedEntries []DictionaryCacheDBEntry) {
}