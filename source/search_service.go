package main

import (
	"math"
	"sort"
	"strings"
)

const IS_DEBUG = true

type search_service struct {
	DocFreq        metric
	InverseDocFreq metric
	TermIDF        map[string]metric
}

type match struct {
	Result string
	Score  float64
}

type doc_tokens map[string][]string
type metric map[string]float64

func GetSearchService(_doc_tokens doc_tokens) *search_service {
	service := search_service{}
	service.InverseDocFreq = *CalcIDF(_doc_tokens)
	service.TermIDF = *GetTermFreqIDF(_doc_tokens, service.InverseDocFreq)
	return &service
}

func (sv search_service) Search(input_str string, is_asc bool) []string {
	input_tokens := Tokenize(input_str)
	input_tfidf := CalcTF_IDF(input_tokens, sv.InverseDocFreq)
	matches := sv.FindMatches(*input_tfidf)
	sorted_matches := SortMatches(matches, is_asc)
	return sorted_matches
}

func (sv search_service) FindMatches(input_tfidf metric) []match {
	matches := []match{}
	for name, tf_idf := range sv.TermIDF {
		cosine_similarity := CalcCosineSimilarity(input_tfidf, tf_idf)
		if cosine_similarity > 0 {
			matches = append(matches, match{
				Result: name,
				Score:  cosine_similarity,
			})
		}
	}
	return matches
}

func SortMatches(matches []match, is_asc bool) []string {
	results := make([]string, 0, len(matches))
	for i := range matches {
		results = append(results, matches[i].Result)
	}
	sort.SliceStable(results, func(i, j int) bool {
		if is_asc {
			return matches[i].Score > matches[j].Score
		} else {
			return matches[i].Score < matches[j].Score
		}
	})
	return results
}

func GetTermFreqIDF(_doc_tokens doc_tokens, idf metric) *map[string]metric {
	tf_idf := map[string]metric{}
	for doc, tokens := range _doc_tokens {
		tf_idf[doc] = *CalcTF_IDF(tokens, idf)
	}
	return &tf_idf
}

func CalcTF_IDF(tokens []string, idf metric) *metric {
	tf := metric{}
	for _, token := range tokens {
		tf[token] = tf[token] + 1
	}

	tf_idf := metric{}
	for key := range tf {
		tf[key] = tf[key] / float64(len(tokens))
		tf_idf[key] = tf[key] * idf[key]
	}
	return &tf_idf
}

func CalcIDF(_docs_tokens doc_tokens) *metric {
	tokens := []string{}
	for _, _tokens := range _docs_tokens {
		for _, token := range _tokens {
			tokens = append(tokens, token)
		}
	}

	df := metric{}
	for _, token := range tokens {
		df[token] = df[token] + 1
	}

	idf := metric{}
	for _, token := range tokens {
		val := float64(len(_docs_tokens)) / float64(df[token]+1)
		idf[token] = math.Log(val)
	}
	return &idf
}

func Standardlize(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	return s
}

func Tokenize(s string) []string {
	delimiters := []string{",", ".", "|"}
	for _, delimiter := range delimiters {
		s = strings.ReplaceAll(s, delimiter, "")
	}
	tokens := strings.Split(s, " ")
	for key := range tokens {
		tokens[key] = Standardlize(tokens[key])
	}
	return tokens
}
func CalcCosineSimilarity(vec1 map[string]float64, vec2 map[string]float64) float64 {
	all_keys := map[string]int{}
	for key := range vec1 {
		all_keys[key] = all_keys[key] + 1
	}
	for key := range vec2 {
		all_keys[key] = all_keys[key] + 1
	}

	dot_product := 0.0
	mag1 := 0.0
	mag2 := 0.0
	for key := range all_keys {
		dot_product += vec1[key] * vec2[key]
		mag1 += vec1[key] * vec1[key]
		mag2 += vec2[key] * vec2[key]
	}

	mag1 = math.Sqrt(mag1)
	mag2 = math.Sqrt(mag2)
	if mag1 == 0 || mag2 == 0 {
		return 0.0
	}
	return dot_product / (mag1 * mag2)
}
