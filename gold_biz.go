package gold

import (
	"github.com/MarkLux/GOLD/serving/common"
	"log"
)

type InfoModel struct {
	Name string `bson:"name"`
	Sex  string `bson:"sex"`
	Age  string `bson:"age"`
	Mobile string `bson:"mobile"`
	Address string `bson:"address"`
}

func (s *GoldService) OnInit() {
}

func (s *GoldService) OnHandle(req *common.GoldRequest, rsp *common.GoldResponse) error {
	infoModel := InfoModel{}
	infoModel.Name = req.Data["name"].(string)
	infoModel.Sex = req.Data["sex"].(string)
	infoModel.Age = req.Data["age"].(string)
	infoModel.Mobile = req.Data["mobile"].(string)
	infoModel.Address = req.Data["address"].(string)

	rsp.Data = make(map[string]interface{})

	cacheKey := "prefix_" + infoModel.Mobile
	res, err := s.CacheClient.Get(cacheKey)
	if err == nil && res != nil {
		rsp.Data["sucess"] = false
		rsp.Data["message"] = "do not retry in 5 minutes."
		return nil
	} else {
		s.CacheClient.Set(cacheKey, true, 300 * 1000)
	}

	dbSession, err := s.DbFactory.NewDataBaseSession("test", "userInfo", "tst", "123")
	if err != nil {
			log.Println("create db session failed, ", err)
			return err
	}
	defer dbSession.Close()

	err = dbSession.Insert(&infoModel)
	if err != nil {
		log.Println("fail to delete info model, ", err)
		return err
	}

	greetingService := s.RpcFactory.NewRemoteServiceConsumer("greeting-service", 3000)
	rpcReq := make(map[string]interface{})
	rpcReq["name"] = infoModel.Name
	greeting, err := greetingService.Request(rpcReq)
	if err != nil {
		log.Println("fali to invoke rpc, ", err)
		return err
	}
	rsp.Data["message"] = greeting

	return nil
}

func (s *GoldService) OnError(err error) bool {
	log.Println(err)
	return false
}
