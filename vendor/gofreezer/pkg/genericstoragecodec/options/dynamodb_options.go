package options

import "github.com/spf13/pflag"

// AddMysqlStorageFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet
func (s *StorageOptions) AddDynamoDBStorageFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.StorageConfig.AWSDynamoDB.Region, "aws-region", s.StorageConfig.AWSDynamoDB.Region, ""+
		"specify the region where the session to connection.")

	fs.StringVar(&s.StorageConfig.AWSDynamoDB.Table, "table", s.StorageConfig.AWSDynamoDB.Table, ""+
		"specify the table name. default(runtime)")
}
