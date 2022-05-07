package utils

import "gitlab.com/alphaticks/xchanger/models"

func GetProtocolAssetID(asset *models.Asset, protocol *models.Protocol, chain *models.Chain) uint64 {
	ID := uint64(asset.ID)<<32 + uint64(protocol.ID)<<16 + uint64(chain.ID)
	return ID
}
