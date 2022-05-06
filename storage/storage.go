package storage

type Storage interface {
	// check if data with given id exists
	CheckExist(key string) (bool, error)

	// Fetch data from storage
	GetData(key string) (data StoredData, err error)

	// Save data to storage
	SaveData(data StoredData) error

	// Get all data
	GetAllData() (dataList []StoredData, err error)

}

// Abstraction of stored data in storage
type StoredData interface {
	GetId() (string)
	GetRawData() (string)
	GetSignature() (string)
	GetTimeStamp() (string)
}