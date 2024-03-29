package usecases

import (
	"github.com/rs/zerolog/log"

	"github.com/berops/claudie/proto/pb"
)

// GetConfigBuilder is a gRPC service: function returns oldest config from the queueBuilder
func (u *Usecases) GetConfigBuilder(request *pb.GetConfigRequest) (*pb.GetConfigResponse, error) {
	configInfo := u.builderQueue.Dequeue()
	if configInfo != nil {
		log.Info().Msgf("Sending config %s to Builder", configInfo.GetName())
		config, err := u.DB.GetConfig(configInfo.GetName(), pb.IdType_NAME)
		if err != nil {
			return nil, err
		}
		return &pb.GetConfigResponse{Config: config}, nil
	}
	return &pb.GetConfigResponse{Config: nil}, nil
}
