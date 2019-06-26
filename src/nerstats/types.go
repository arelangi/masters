package main

type RTTweet struct {
	OriginalText   string `json:"original_text"`
	NormalizedText string `json:"normalized_text"`
}

type Predictions struct {
	FtPrediction      string   `json:"ft_prediction"`
	LstmPrediction    string   `json:"lstm_prediction"`
	NbPrediction      string   `json:"nb_prediction"`
	SvmPrediction     string   `json:"svm_prediction"`
	TranslatedText    string   `json:"translated_text"`
	NonEnglishPercent float64  `json:"non_english_percent"`
	SpacyEntities     []string `json:"spacy_entities"`
	NLTKEntities      []string `json:"nltk_entities"`
}

type Response struct {
	Tweet       RTTweet             `json:"tweet"`
	Predictions Predictions         `json:"predictions"`
	Entities    map[string][]string `json:"entities"`
}
