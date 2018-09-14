package outputadapter

import "TechPlat/datapipe/config"

type OutputAdapter func(conf config.OutputAdapter, appid string, logstr string)