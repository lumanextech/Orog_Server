package types

import "github.com/gagliardetto/solana-go"

var (
	SerumOpenBookProgramID = solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")

	RaydiumAmmMarket = &MarketProgram{
		MarketTypeName:     "RaydiumAMM",
		ProgramID:          solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"),
		AuthorityProgramID: solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1"),
	}

	//RaydiumCLMMMarket https://solscan.io/account/7ZSthSUfB4aw7tzuVDY5kzvt77qFYqmeRV99yW4NpM1c#data
	RaydiumCLMMMarket = &MarketProgram{
		MarketTypeName:     "RaydiumCLMM",
		ProgramID:          solana.MustPublicKeyFromBase58("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK"),
		AuthorityProgramID: solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1"),
	}

	JupterV6Market = &MarketProgram{
		MarketTypeName:     "JupterV6",
		ProgramID:          solana.MustPublicKeyFromBase58("JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4"),
		AuthorityProgramID: solana.MustPublicKeyFromBase58("69yhtoJR4JYPPABZcSNkzuqbaFbwHsCkja1sP1Q2aVT5"),
	}

	RaydiumCPMMMarket = &MarketProgram{
		MarketTypeName:     "RaydiumCPMM",
		ProgramID:          solana.MustPublicKeyFromBase58("CLmmVkC6jXjR7j5AXzL5m7K1s7VTwXJEbXkXUQJ6q6Mk"),
		AuthorityProgramID: solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1"),
	}
)
