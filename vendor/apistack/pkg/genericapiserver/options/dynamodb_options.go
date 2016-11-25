package options

import "github.com/spf13/pflag"

// AddMysqlStorageFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet
func (s *ServerRunOptions) AddDynamoDBStorageFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.StorageConfig.AWSDynamoDB.Region, "aws-region", s.StorageConfig.AWSDynamoDB.Region, ""+
		"specify the region where the session to connection.")

	fs.StringVar(&s.StorageConfig.AWSDynamoDB.Table, "aws-table", s.StorageConfig.AWSDynamoDB.Table, ""+
		"specify the table name. default(program name)")

	fs.StringVar(&s.StorageConfig.AWSDynamoDB.Token, "aws-cred-token", s.StorageConfig.AWSDynamoDB.Table, ""+
		"specify the token for credentials.")

	fs.StringVar(&s.StorageConfig.AWSDynamoDB.AccessID, "aws-cred-accessid", s.StorageConfig.AWSDynamoDB.AccessID, ""+
		"specify the access id for credentials.")

	fs.StringVar(&s.StorageConfig.AWSDynamoDB.AccessKey, "aws-cred-accesskey", s.StorageConfig.AWSDynamoDB.AccessKey, ""+
		"specify the access key for credentials.")
}
