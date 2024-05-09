package main

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RAGRequest struct {
	Input       string    `json:"input"`
	Messages    []Message `json:"messages,omitempty"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature"`
	K           int       `json:"k"`
	Relevance   float64   `json:"relevance"`
}

type GPTRequest struct {
	Messages    []Message `json:"messages"`
	MaxTokens   *int      `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Choice struct {
	Index                int     `json:"index"`
	FinishReason         string  `json:"finish_reason"`
	ContentFilterResults any     `json:"content_filter_results"`
	Logprobs             any     `json:"logprobs"`
	Message              Message `json:"message"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type GPTResponse struct {
	Id      string   `json:"id"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type AISearchResult struct {
	SearchScore   float64 `json:"@search.score"`
	RerankerScore float64 `json:"@search.rerankerScore"`
	ChunkId       string  `json:"chunk_id"`
	ParentId      string  `json:"parent_id"`
	Chunk         string  `json:"chunk"`
	Title         string  `json:"title"`
}

type VectorQuery struct {
	Text   string `json:"text"`
	Kind   string `json:"kind"`
	K      int    `json:"k"`
	Fields string `json:"fields"`
}

type AISearchRequest struct {
	Search                string        `json:"search"`
	VectorQueries         []VectorQuery `json:"vectorQueries"`
	SemanticConfiguration string        `json:"semanticConfiguration"`
	Top                   int           `json:"top"`
	QueryType             string        `json:"queryType"`
	Select                string        `json:"select"`
	QueryLanguage         string        `json:"queryLanguage"`
}

// {
// 	"@odata.context": "https://alemorsearch.search.windows.net/indexes('adventureworksai-20240421')/$metadata#docs(*)",
// 	"value": [
// 	  {
// 		"@search.score": 0.1,
// 		"@search.rerankerScore": 1.0,
// 		"chunk_id": "",
// 		"chunk": "",
// 		"title": ""
// 	  }
// }

type AISearchResponse struct {
	Context string           `json:"@odata.context"`
	Value   []AISearchResult `json:"value"`
}
