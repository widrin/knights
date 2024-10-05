package enginelog

type EngineLog int64

const (
	S2SDistributionRuleServerReq EngineLog = 1 << 0
	S2SActorHeart                EngineLog = 1 << 1
	DistributionRuleServer       EngineLog = 1 << 2
	c                            EngineLog = 1 << 3
	d                            EngineLog = 1 << 4
	e                            EngineLog = 1 << 5
)
