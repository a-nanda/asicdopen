//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package rpc

import (
	"asicd/pluginManager"
	"asicd/publisher"
	"asicdServices"
	"git.apache.org/thrift.git/lib/go/thrift"
	"utils/logging"
)

// Struct that will provide service interfaces
type AsicDaemonServiceHandler struct {
	pluginMgr  *pluginManager.PluginManager
	logger     *logging.Writer
	notifyChan *publisher.PubChannels
}

func NewAsicDaemonServiceHandler(pluginMgr *pluginManager.PluginManager, logger *logging.Writer, notifyChan *publisher.PubChannels) *AsicDaemonServiceHandler {
	return &AsicDaemonServiceHandler{
		pluginMgr:  pluginMgr,
		logger:     logger,
		notifyChan: notifyChan,
	}
}

type AsicDaemonServerInfo struct {
	Socket           string
	Handler          *AsicDaemonServiceHandler
	Processor        *asicdServices.ASICDServicesProcessor
	Transport        *thrift.TServerSocket
	TransportFactory thrift.TTransportFactory
	ProtocolFactory  *thrift.TBinaryProtocolFactory
	Server           *thrift.TSimpleServer
}

func NewAsicdServer(socket string, pluginMgr *pluginManager.PluginManager, logger *logging.Writer, notifyChan *publisher.PubChannels) *AsicDaemonServerInfo {
	transport, err := thrift.NewTServerSocket(socket)
	if err != nil {
		panic(err)
	}
	handler := NewAsicDaemonServiceHandler(pluginMgr, logger, notifyChan)
	processor := asicdServices.NewASICDServicesProcessor(handler)
	transportFactory := thrift.NewTBufferedTransportFactory(16384)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	return &AsicDaemonServerInfo{
		Socket:           socket,
		Handler:          handler,
		Processor:        processor,
		Transport:        transport,
		TransportFactory: transportFactory,
		ProtocolFactory:  protocolFactory,
		Server:           server,
	}
}
